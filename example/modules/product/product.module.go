package product

import (
	"example/modules/product/application/services"
	"example/modules/product/domain/interfaces"
	"example/modules/product/infrastructure/repositories"
	"xcomp"
)

func CreateProductModule() xcomp.Module {
	return xcomp.NewModule().
		AddFactory("ProductService", func(c *xcomp.Container) any {
			service := &services.ProductService{}

			// Auto inject Logger (uppercase field with inject tag)
			if err := c.Inject(service); err != nil {
				if logger, ok := c.Get("Logger").(xcomp.Logger); ok {
					logger.Error("Failed to inject ProductService Logger",
						xcomp.Field("error", err))
				}
				panic("Failed to inject ProductService Logger: " + err.Error())
			}

			// Manual inject lowercase fields via method
			productRepo := c.Get("ProductRepository").(interfaces.ProductRepository)
			productCacheRepo := c.Get("ProductCacheRepository").(interfaces.ProductCacheRepository)

			service.SetDependencies(productRepo, productCacheRepo)

			return service
		}).
		AddFactory("ProductRepository", func(c *xcomp.Container) any {
			repo := &repositories.ProductRepositoryImpl{}
			c.Inject(repo)
			return repo
		}).
		AddFactory("ProductCacheRepository", func(c *xcomp.Container) any {
			cacheRepo := &repositories.ProductCacheRepositoryImpl{}
			c.Inject(cacheRepo)
			return cacheRepo
		}).
		Build()
}
