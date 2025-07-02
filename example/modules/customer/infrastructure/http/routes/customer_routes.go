package routes

import (
	"example/modules/customer/infrastructure/http/controllers"

	"github.com/gofiber/fiber/v2"
)

type CustomerRoutes struct {
	CustomerController *controllers.CustomerController `inject:"CustomerController"`
}

func (cr *CustomerRoutes) GetServiceName() string {
	return "CustomerRoutes"
}

func (cr *CustomerRoutes) SetupRoutes(app *fiber.App) {
	api := app.Group("/api/v1")
	customers := api.Group("/customers")

	customers.Get("/", cr.CustomerController.ListCustomers)
	customers.Get("/search", cr.CustomerController.SearchCustomers)
	customers.Get("/by-email", cr.CustomerController.GetCustomerByEmail)
	customers.Get("/:id", cr.CustomerController.GetCustomer)
	customers.Get("/username/:username", cr.CustomerController.GetCustomerByUsername)
	customers.Post("/", cr.CustomerController.CreateCustomer)
	customers.Put("/:id", cr.CustomerController.UpdateCustomer)
	customers.Delete("/:id", cr.CustomerController.DeleteCustomer)
}
