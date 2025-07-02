package main

import (
	"example/controllers"
	"xcomp"
)

func CreateTransportModule() xcomp.Module {
	return xcomp.NewModule().
		AddFactory("ProductController", func(c *xcomp.Container) any {
			controller := &controllers.ProductController{}
			c.Inject(controller)
			return controller
		}).
		AddFactory("OrderController", func(c *xcomp.Container) any {
			controller := &controllers.OrderController{}
			c.Inject(controller)
			return controller
		}).
		AddFactory("CustomerController", func(c *xcomp.Container) any {
			controller := &controllers.CustomerController{}
			c.Inject(controller)
			return controller
		}).
		Build()
}
