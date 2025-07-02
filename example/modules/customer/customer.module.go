package customer

import (
	"example/infrastructure/database"
	"example/modules/customer/application/services"
	"example/modules/customer/infrastructure/http/controllers"
	"example/modules/customer/infrastructure/http/routes"
	"example/modules/customer/infrastructure/persistence"
	"xcomp"
)

func CreateCustomerModule() xcomp.Module {
	return xcomp.NewModule().
		AddFactory("RedisClient", func(container *xcomp.Container) any {
			redisService := &database.RedisService{}
			container.Inject(redisService)
			redisService.Initialize()
			return redisService.GetClient()
		}).
		AddFactory("CustomerRepository", func(container *xcomp.Container) any {
			repo := &persistence.CustomerRepositoryImpl{}
			container.Inject(repo)
			repo.Initialize()
			return repo
		}).
		AddFactory("CustomerCacheRepository", func(container *xcomp.Container) any {
			cacheRepo := &persistence.CustomerCacheRepositoryImpl{}
			container.Inject(cacheRepo)
			return cacheRepo
		}).
		AddFactory("CustomerService", func(container *xcomp.Container) any {
			service := &services.CustomerService{}
			container.Inject(service)
			return service
		}).
		AddFactory("CustomerController", func(container *xcomp.Container) any {
			controller := &controllers.CustomerController{}
			container.Inject(controller)
			return controller
		}).
		AddFactory("CustomerRoutes", func(container *xcomp.Container) any {
			routes := &routes.CustomerRoutes{}
			container.Inject(routes)
			return routes
		}).
		Build()
}
