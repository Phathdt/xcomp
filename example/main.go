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
	"example/modules/order"
	orderRoutes "example/modules/order/infrastructure/http/routes"
	"example/modules/product"
	"example/modules/product/infrastructure/http/routes"

	"xcomp"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func createInfrastructureModule(container *xcomp.Container) xcomp.Module {
	return xcomp.NewModule().
		AddFactory("ConfigService", func(container *xcomp.Container) any {
			return xcomp.NewConfigService("config.yaml")
		}).
		AddFactory("DatabaseConnection", func(container *xcomp.Container) any {
			dbConn := &database.DatabaseConnection{}
			if err := container.Inject(dbConn); err != nil {
				log.Printf("Failed to inject DatabaseConnection: %v", err)
				return dbConn
			}
			if err := dbConn.Initialize(); err != nil {
				log.Fatalf("Failed to initialize database connection: %v", err)
			}
			return dbConn
		}).
		Build()
}

func createAppModule(container *xcomp.Container) xcomp.Module {
	infrastructureModule := createInfrastructureModule(container)
	productModule := product.CreateProductModule()
	orderModule := order.NewOrderModule()

	return xcomp.NewModule().
		Import(infrastructureModule).
		Import(productModule).
		Import(orderModule).
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
		})
	})

	return app
}

func main() {
	log.Println("Starting API Server...")

	container := xcomp.NewContainer()

	appModule := createAppModule(container)
	if err := container.RegisterModule(appModule); err != nil {
		log.Fatalf("Failed to register app module: %v", err)
	}

	log.Println("Registered Services:")
	for _, serviceName := range container.ListServices() {
		log.Printf("- %s", serviceName)
	}

	configService, ok := container.Get("ConfigService").(*xcomp.ConfigService)
	if !ok {
		log.Fatal("Failed to get ConfigService")
	}

	app := setupFiberApp(configService)

	productRoutes, ok := container.Get("ProductRoutes").(*routes.ProductRoutes)
	if !ok {
		log.Fatal("Failed to get ProductRoutes")
	}
	productRoutes.SetupRoutes(app)

	orderRoutesInstance, ok := container.Get("OrderRoutes").(*orderRoutes.OrderRoutes)
	if !ok {
		log.Fatal("Failed to get OrderRoutes")
	}
	orderRoutesInstance.SetupRoutes(app)

	port := configService.GetInt("app.port", 3000)

	go func() {
		log.Printf("Server starting on port %d", port)
		if err := app.Listen(fmt.Sprintf(":%d", port)); err != nil {
			log.Printf("Server failed to start: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := app.ShutdownWithContext(ctx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	}

	if dbConn, ok := container.Get("DatabaseConnection").(*database.DatabaseConnection); ok {
		if err := dbConn.Close(); err != nil {
			log.Printf("Failed to close database connection: %v", err)
		}
	}

	log.Println("Server exited")
}
