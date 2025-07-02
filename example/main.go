package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"example/infrastructure/database"
	"example/modules/customer"
	"example/modules/order"
	"example/modules/product"

	"xcomp"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
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
				configFile = "config.yaml"
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
				return dbConn
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
			return dbConn
		}).
		Build()
}

func createAppModule(container *xcomp.Container) xcomp.Module {
	infrastructureModule := createInfrastructureModule(container)
	productModule := product.CreateProductModule()
	orderModule := order.NewOrderModule()
	customerModule := customer.CreateCustomerModule()
	transportModule := CreateTransportModule()

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
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
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
		app.Use(cors.New(cors.Config{
			AllowOrigins: "*",
			AllowMethods: "GET,POST,PUT,DELETE,OPTIONS,PATCH",
			AllowHeaders: "Content-Type,Authorization",
		}))
	}

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":    "healthy",
			"timestamp": time.Now().Unix(),
			"service":   "API Server",
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

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	logger.Info("Shutdown signal received, beginning graceful shutdown")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := app.ShutdownWithContext(ctx); err != nil {
		logger.Error("Server forced to shutdown",
			xcomp.Field("error", err),
			xcomp.Field("timeout", "30s"))
	} else {
		logger.Info("HTTP server shutdown completed")
	}

	if dbConn, ok := container.Get("DatabaseConnection").(*database.DatabaseConnection); ok {
		if err := dbConn.Close(); err != nil {
			logger.Error("Failed to close database connection",
				xcomp.Field("error", err))
		} else {
			logger.Info("Database connection closed successfully")
		}
	}

	logger.Info("Application shutdown completed")
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
