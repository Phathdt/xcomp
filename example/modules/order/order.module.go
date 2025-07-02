package order

import (
	"example/modules/order/application/services"
	"example/modules/order/infrastructure/repositories"
	"xcomp"
)

func NewOrderModule() xcomp.Module {
	return xcomp.NewModule().
		AddFactory("OrderService", func(c *xcomp.Container) any {
			service := services.NewOrderService()
			if err := c.Inject(service); err != nil {
				if logger, ok := c.Get("Logger").(xcomp.Logger); ok {
					logger.Error("Failed to inject OrderService dependencies",
						xcomp.Field("error", err))
				}
				panic("Failed to inject OrderService dependencies: " + err.Error())
			}
			return service
		}).
		AddFactory("OrderRepository", func(c *xcomp.Container) any {
			repo := &repositories.OrderRepositoryImpl{}
			if err := c.Inject(repo); err != nil {
				if logger, ok := c.Get("Logger").(xcomp.Logger); ok {
					logger.Error("Failed to inject OrderRepository dependencies",
						xcomp.Field("error", err))
				}
				panic("Failed to inject OrderRepository dependencies: " + err.Error())
			}
			return repo
		}).
		AddFactory("OrderItemRepository", func(c *xcomp.Container) any {
			repo := &repositories.OrderItemRepositoryImpl{}
			if err := c.Inject(repo); err != nil {
				if logger, ok := c.Get("Logger").(xcomp.Logger); ok {
					logger.Error("Failed to inject OrderItemRepository dependencies",
						xcomp.Field("error", err))
				}
				panic("Failed to inject OrderItemRepository dependencies: " + err.Error())
			}
			return repo
		}).
		AddFactory("OrderCacheRepository", func(c *xcomp.Container) any {
			cacheRepo := &repositories.OrderCacheRepositoryImpl{}
			if err := c.Inject(cacheRepo); err != nil {
				if logger, ok := c.Get("Logger").(xcomp.Logger); ok {
					logger.Error("Failed to inject OrderCacheRepository dependencies",
						xcomp.Field("error", err))
				}
				panic("Failed to inject OrderCacheRepository dependencies: " + err.Error())
			}
			return cacheRepo
		}).
		Build()
}
