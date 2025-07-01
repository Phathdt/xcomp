package routes

import (
	"example/modules/order/infrastructure/http/controllers"

	"github.com/gofiber/fiber/v2"
)

type OrderRoutes struct {
	orderController *controllers.OrderController `inject:"OrderController"`
}

func NewOrderRoutes() *OrderRoutes {
	return &OrderRoutes{}
}

func (r *OrderRoutes) SetupRoutes(app *fiber.App) {
	orderGroup := app.Group("/api/orders")

	orderGroup.Post("/", r.orderController.CreateOrder)
	orderGroup.Get("/", r.orderController.GetOrders)
	orderGroup.Get("/:id", r.orderController.GetOrder)
	orderGroup.Put("/:id", r.orderController.UpdateOrder)
	orderGroup.Delete("/:id", r.orderController.DeleteOrder)

	orderGroup.Post("/:id/confirm", r.orderController.ConfirmOrder)
	orderGroup.Post("/:id/ship", r.orderController.ShipOrder)
	orderGroup.Post("/:id/deliver", r.orderController.DeliverOrder)
	orderGroup.Post("/:id/cancel", r.orderController.CancelOrder)

	orderGroup.Post("/:id/items", r.orderController.AddOrderItem)
	orderGroup.Put("/:id/items/:product_id", r.orderController.UpdateOrderItemQuantity)
	orderGroup.Delete("/:id/items/:product_id", r.orderController.RemoveOrderItem)
}
