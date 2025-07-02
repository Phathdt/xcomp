package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"example/infrastructure/async"
	"example/infrastructure/database"
	"example/modules/customer"
	customerInterfaces "example/modules/customer/domain/interfaces"
	"example/modules/order"
	orderInterfaces "example/modules/order/domain/interfaces"
	"example/modules/product"

	"xcomp"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/redis/go-redis/v9"
	"github.com/urfave/cli/v2"
)

var (
	Version   = "1.0.0"
	BuildTime = "unknown"
	GitCommit = "unknown"
)

func createInfrastructureModule(container *xcomp.Container) xcomp.Module {
	return xcomp.NewModule().
		AddFactory("ConfigService", func(container *xcomp.Container) any {
			configFile := os.Getenv("CONFIG_FILE")
			if configFile == "" {
				configFile = "config-dev.yaml"
			}

			return xcomp.NewConfigService(configFile)
		}).
		AddFactory("Logger", func(container *xcomp.Container) any {
			configService, _ := container.Get("ConfigService").(*xcomp.ConfigService)
			if configService != nil {
				return xcomp.NewLogger(configService)
			}
			return xcomp.NewDevelopmentLogger()
		}).
		AddFactory("RedisClient", func(container *xcomp.Container) any {
			redisService := &database.RedisService{}
			container.Inject(redisService)
			redisService.Initialize()
			return redisService.GetClient()
		}).
		AddFactory("DatabaseConnection", func(container *xcomp.Container) any {
			dbConn := &database.DatabaseConnection{}
			if err := container.Inject(dbConn); err != nil {
				if logger, ok := container.Get("Logger").(xcomp.Logger); ok {
					logger.Error("Failed to inject DatabaseConnection dependencies",
						xcomp.Field("error", err))
				}
				panic("Failed to inject DatabaseConnection dependencies: " + err.Error())
			}
			if err := dbConn.Initialize(); err != nil {
				if logger, ok := container.Get("Logger").(xcomp.Logger); ok {
					logger.Fatal("Failed to initialize database connection",
						xcomp.Field("error", err))
				} else {
					log.Fatalf("Failed to initialize database connection: %v", err)
				}
			}
			if logger, ok := container.Get("Logger").(xcomp.Logger); ok {
				logger.Info("Database connection initialized successfully")
			}
			return dbConn.GetDB()
		}).
		Build()
}

func createAppModule(container *xcomp.Container) xcomp.Module {
	infrastructureModule := createInfrastructureModule(container)
	productModule := product.CreateProductModule()
	orderModule := order.NewOrderModule()
	customerModule := customer.CreateCustomerModule()
	transportModule := CreateTransportModule()

	// Register all business modules first - do NOT include AsyncModule here
	return xcomp.NewModule().
		Import(infrastructureModule).
		Import(productModule).
		Import(orderModule).
		Import(customerModule).
		Import(transportModule).
		Build()
}

func setupFiberApp(configService *xcomp.ConfigService) *fiber.App {
	app := fiber.New(fiber.Config{
		ReadTimeout:  time.Duration(configService.GetInt("server.read_timeout_seconds", 30)) * time.Second,
		WriteTimeout: time.Duration(configService.GetInt("server.write_timeout_seconds", 30)) * time.Second,
		IdleTimeout:  time.Duration(configService.GetInt("server.timeout_seconds", 30)) * time.Second,
		Prefork:      configService.GetBool("server.prefork", false),
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}

			return c.Status(code).JSON(fiber.Map{
				"error":   "Request failed",
				"message": err.Error(),
			})
		},
	})

	app.Use(recover.New())
	app.Use(logger.New(logger.Config{
		Format: "${time} ${method} ${path} - ${status} - ${latency}\n",
	}))

	if configService.GetBool("server.cors.enabled", true) {
		allowedOrigins := configService.GetString("server.cors.allowed_origins", "*")
		allowedMethods := configService.GetString("server.cors.allowed_methods", "GET,POST,PUT,DELETE,OPTIONS,PATCH")
		allowedHeaders := configService.GetString("server.cors.allowed_headers", "Content-Type,Authorization")

		app.Use(cors.New(cors.Config{
			AllowOrigins: allowedOrigins,
			AllowMethods: allowedMethods,
			AllowHeaders: allowedHeaders,
		}))
	}

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":    "healthy",
			"timestamp": time.Now().Unix(),
			"service":   configService.GetString("app.name", "API Server"),
			"version":   Version,
		})
	})

	return app
}

func serveCommand(c *cli.Context) error {
	container := xcomp.NewContainer()

	appModule := createAppModule(container)
	if err := container.RegisterModule(appModule); err != nil {
		return fmt.Errorf("failed to register app module: %w", err)
	}

	configService, ok := container.Get("ConfigService").(*xcomp.ConfigService)
	if !ok {
		return fmt.Errorf("failed to get ConfigService from container")
	}

	logger, ok := container.Get("Logger").(xcomp.Logger)
	if !ok {
		return fmt.Errorf("failed to get Logger from container")
	}

	logger.Info("Starting API Server",
		xcomp.Field("version", Version),
		xcomp.Field("build_time", BuildTime),
		xcomp.Field("git_commit", GitCommit),
		xcomp.Field("environment", configService.GetString("app.environment", "development")),
		xcomp.Field("name", configService.GetString("app.name", "API Server")))

	services := container.ListServices()
	logger.Info("Dependency injection container initialized",
		xcomp.Field("registered_services_count", len(services)),
		xcomp.Field("services", services))

	app := setupFiberApp(configService)

	// Setup centralized routes
	setupRoutes(app, container)
	logger.Debug("All routes registered")

	// Create AsyncService AFTER all modules are registered and dependencies are available
	redisClient, ok := container.Get("RedisClient").(*redis.Client)
	if !ok || redisClient == nil {
		return fmt.Errorf("failed to get RedisClient from container")
	}

	orderService, ok := container.Get("OrderService").(orderInterfaces.OrderService)
	if !ok || orderService == nil {
		return fmt.Errorf("failed to get OrderService from container")
	}

	customerService, ok := container.Get("CustomerService").(customerInterfaces.CustomerService)
	if !ok || customerService == nil {
		return fmt.Errorf("failed to get CustomerService from container")
	}

	logger.Info("Creating AsyncService manually after all dependencies are available")
	asyncService := async.NewAsyncService(redisClient, orderService, customerService, logger)

	asyncCtx, asyncCancel := context.WithCancel(context.Background())
	defer asyncCancel()

	if err := asyncService.Start(asyncCtx); err != nil {
		return fmt.Errorf("failed to start async service: %w", err)
	}

	// Setup asynq monitoring endpoint
	monitorHandler := asyncService.GetMonitorHandler()
	go func() {
		monitorPort := configService.GetInt("async.monitor.port", 8080)
		logger.Info("Asynq monitor starting",
			xcomp.Field("port", monitorPort),
			xcomp.Field("path", "/monitoring"))

		if err := http.ListenAndServe(fmt.Sprintf(":%d", monitorPort), monitorHandler); err != nil {
			logger.Error("Asynq monitor failed to start",
				xcomp.Field("port", monitorPort),
				xcomp.Field("error", err))
		}
	}()

	port := c.Int("port")
	if port == 0 {
		port = configService.GetInt("app.port", 3000)
	}

	go func() {
		logger.Info("HTTP server starting",
			xcomp.Field("port", port),
			xcomp.Field("address", fmt.Sprintf(":%d", port)))
		if err := app.Listen(fmt.Sprintf(":%d", port)); err != nil {
			logger.Error("Server failed to start",
				xcomp.Field("port", port),
				xcomp.Field("error", err))
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	// Cancel async context first
	asyncCancel()

	// Shutdown Fiber server
	if err := app.ShutdownWithTimeout(30 * time.Second); err != nil {
		logger.Error("Server forced to shutdown", xcomp.Field("error", err))
		return err
	}

	logger.Info("Server exited successfully")
	return nil
}

func main() {
	app := &cli.App{
		Name:    "API Server",
		Usage:   "XComp-powered API server with dependency injection",
		Version: Version,
		Commands: []*cli.Command{
			{
				Name:    "serve",
				Aliases: []string{"s", "start"},
				Usage:   "Start the HTTP server",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "config",
						Aliases: []string{"c"},
						Usage:   "Configuration file path",
						EnvVars: []string{"CONFIG_FILE"},
						Value:   "config.yaml",
					},
					&cli.IntFlag{
						Name:    "port",
						Aliases: []string{"p"},
						Usage:   "Port to listen on",
						EnvVars: []string{"PORT"},
						Value:   0, // 0 means use config file value
					},
				},
				Action: func(c *cli.Context) error {
					if configFile := c.String("config"); configFile != "" {
						os.Setenv("CONFIG_FILE", configFile)
					}
					return serveCommand(c)
				},
			},
			{
				Name:    "version",
				Aliases: []string{"v"},
				Usage:   "Show version information",
				Action: func(c *cli.Context) error {
					fmt.Printf("Version: %s\n", Version)
					fmt.Printf("Build Time: %s\n", BuildTime)
					fmt.Printf("Git Commit: %s\n", GitCommit)
					return nil
				},
			},
			{
				Name:  "health",
				Usage: "Check application health",
				Action: func(c *cli.Context) error {
					fmt.Println("✅ Application is healthy")
					fmt.Printf("Version: %s\n", Version)
					fmt.Printf("Build Time: %s\n", BuildTime)
					return nil
				},
			},
		},
		DefaultCommand: "serve",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "verbose",
				Aliases: []string{"V"},
				Usage:   "Enable verbose logging",
				EnvVars: []string{"VERBOSE"},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
