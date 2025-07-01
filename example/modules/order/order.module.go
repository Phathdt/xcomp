package order

import (
	"example/modules/order/application/services"
	"example/modules/order/infrastructure/http/controllers"
	"example/modules/order/infrastructure/http/routes"
	"example/modules/order/infrastructure/persistence"
	"xcomp"
)

func NewOrderModule() xcomp.Module {
	return xcomp.NewModule().
		AddFactory("OrderService", func(c *xcomp.Container) any {
			return services.NewOrderService()
		}).
		AddFactory("OrderRepository", func(c *xcomp.Container) any {
			return persistence.NewOrderRepository()
		}).
		AddFactory("OrderItemRepository", func(c *xcomp.Container) any {
			return persistence.NewOrderItemRepository()
		}).
		AddFactory("OrderController", func(c *xcomp.Container) any {
			return controllers.NewOrderController()
		}).
		AddFactory("OrderRoutes", func(c *xcomp.Container) any {
			return routes.NewOrderRoutes()
		}).
		Build()
}
