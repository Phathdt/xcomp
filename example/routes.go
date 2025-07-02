package main

import (
	"example/controllers"
	"xcomp"

	"github.com/gofiber/fiber/v2"
)

func setupRoutes(app *fiber.App, container *xcomp.Container) {
	// Get controllers from container
	productController, ok := container.Get("ProductController").(*controllers.ProductController)
	if !ok {
		panic("Failed to get ProductController from container")
	}

	orderController, ok := container.Get("OrderController").(*controllers.OrderController)
	if !ok {
		panic("Failed to get OrderController from container")
	}

	customerController, ok := container.Get("CustomerController").(*controllers.CustomerController)
	if !ok {
		panic("Failed to get CustomerController from container")
	}

	// Setup API routes
	api := app.Group("/api/v1")

	// Product routes
	products := api.Group("/products")
	products.Get("/", productController.ListProducts)
	products.Get("/search", productController.SearchProducts)
	products.Get("/:id", productController.GetProduct)
	products.Post("/", productController.CreateProduct)
	products.Put("/:id", productController.UpdateProduct)
	products.Patch("/:id/stock", productController.UpdateProductStock)
	products.Delete("/:id", productController.DeleteProduct)

	// Order routes
	orders := api.Group("/orders")
	orders.Get("/", orderController.GetOrders)
	orders.Get("/:id", orderController.GetOrder)
	orders.Post("/", orderController.CreateOrder)
	orders.Put("/:id", orderController.UpdateOrder)
	orders.Patch("/:id/confirm", orderController.ConfirmOrder)
	orders.Patch("/:id/ship", orderController.ShipOrder)
	orders.Patch("/:id/deliver", orderController.DeliverOrder)
	orders.Patch("/:id/cancel", orderController.CancelOrder)
	orders.Post("/:id/items", orderController.AddOrderItem)
	orders.Put("/:order_id/items/:product_id", orderController.UpdateOrderItemQuantity)
	orders.Delete("/:order_id/items/:product_id", orderController.RemoveOrderItem)
	orders.Delete("/:id", orderController.DeleteOrder)

	// Customer routes
	customers := api.Group("/customers")
	customers.Get("/", customerController.ListCustomers)
	customers.Get("/search", customerController.SearchCustomers)
	customers.Get("/username/:username", customerController.GetCustomerByUsername)
	customers.Get("/by-email", customerController.GetCustomerByEmail)
	customers.Get("/:id", customerController.GetCustomer)
	customers.Post("/", customerController.CreateCustomer)
	customers.Put("/:id", customerController.UpdateCustomer)
	customers.Delete("/:id", customerController.DeleteCustomer)
}
