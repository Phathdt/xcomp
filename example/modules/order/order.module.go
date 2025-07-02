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
			service := &services.OrderService{}
			c.Inject(service)
			return service
		}).
		AddFactory("OrderRepository", func(c *xcomp.Container) any {
			repo := &persistence.OrderRepositoryImpl{}
			c.Inject(repo)
			return repo
		}).
		AddFactory("OrderItemRepository", func(c *xcomp.Container) any {
			repo := &persistence.OrderItemRepositoryImpl{}
			c.Inject(repo)
			return repo
		}).
		AddFactory("OrderCacheRepository", func(c *xcomp.Container) any {
			cacheRepo := &persistence.OrderCacheRepositoryImpl{}
			c.Inject(cacheRepo)
			return cacheRepo
		}).
		AddFactory("OrderController", func(c *xcomp.Container) any {
			controller := &controllers.OrderController{}
			c.Inject(controller)
			return controller
		}).
		AddFactory("OrderRoutes", func(c *xcomp.Container) any {
			routes := &routes.OrderRoutes{}
			c.Inject(routes)
			return routes
		}).
		Build()
}
