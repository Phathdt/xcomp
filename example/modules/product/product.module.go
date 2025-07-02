package product

import (
	"example/modules/product/application/services"
	"example/modules/product/infrastructure/repositories"
	"xcomp"
)

func CreateProductModule() xcomp.Module {
	return xcomp.NewModule().
		AddFactory("ProductService", func(c *xcomp.Container) any {
			service := &services.ProductService{}
			c.Inject(service)
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
