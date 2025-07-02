package order

import (
	"example/modules/order/application/services"
	"example/modules/order/infrastructure/repositories"
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
			repo := &repositories.OrderRepositoryImpl{}
			c.Inject(repo)
			return repo
		}).
		AddFactory("OrderItemRepository", func(c *xcomp.Container) any {
			repo := &repositories.OrderItemRepositoryImpl{}
			c.Inject(repo)
			return repo
		}).
		AddFactory("OrderCacheRepository", func(c *xcomp.Container) any {
			cacheRepo := &repositories.OrderCacheRepositoryImpl{}
			c.Inject(cacheRepo)
			return cacheRepo
		}).
		Build()
}
