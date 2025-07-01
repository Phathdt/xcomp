package product

import (
	"example/modules/product/application/services"
	"example/modules/product/infrastructure/http/controllers"
	"example/modules/product/infrastructure/http/routes"
	"example/modules/product/infrastructure/persistence"
	"xcomp"
)

func CreateProductModule() xcomp.Module {
	return xcomp.NewModule().
		AddFactory("ProductRepository", func(container *xcomp.Container) any {
			repo := &persistence.ProductRepositoryImpl{}
			container.Inject(repo)
			repo.Initialize()
			return repo
		}).
		AddFactory("ProductService", func(container *xcomp.Container) any {
			service := &services.ProductService{}
			container.Inject(service)
			return service
		}).
		AddFactory("ProductController", func(container *xcomp.Container) any {
			controller := &controllers.ProductController{}
			container.Inject(controller)
			return controller
		}).
		AddFactory("ProductRoutes", func(container *xcomp.Container) any {
			routes := &routes.ProductRoutes{}
			container.Inject(routes)
			return routes
		}).
		Build()
}
