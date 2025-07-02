package customer

import (
	"example/modules/customer/application/services"
	"example/modules/customer/domain/interfaces"
	"example/modules/customer/infrastructure/repositories"
	"xcomp"
)

func CreateCustomerModule() xcomp.Module {
	return xcomp.NewModule().
		AddFactory("CustomerService", func(c *xcomp.Container) any {
			service := services.NewCustomerService()

			// Manual inject lowercase fields via method
			customerRepo := c.Get("CustomerRepository").(interfaces.CustomerRepository)
			customerCacheRepo := c.Get("CustomerCacheRepository").(interfaces.CustomerCacheRepository)

			service.SetDependencies(customerRepo, customerCacheRepo)

			return service
		}).
		AddFactory("CustomerRepository", func(c *xcomp.Container) any {
			repo := &repositories.CustomerRepositoryImpl{}
			c.Inject(repo)
			return repo
		}).
		AddFactory("CustomerCacheRepository", func(c *xcomp.Container) any {
			cacheRepo := &repositories.CustomerCacheRepositoryImpl{}
			c.Inject(cacheRepo)
			return cacheRepo
		}).
		Build()
}
