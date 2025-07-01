package routes

import (
	"example/modules/product/infrastructure/http/controllers"

	"github.com/gofiber/fiber/v2"
)

type ProductRoutes struct {
	ProductController *controllers.ProductController `inject:"ProductController"`
}

func (pr *ProductRoutes) GetServiceName() string {
	return "ProductRoutes"
}

func (pr *ProductRoutes) SetupRoutes(app *fiber.App) {
	api := app.Group("/api/v1")
	products := api.Group("/products")

	products.Get("/", pr.ProductController.ListProducts)
	products.Get("/search", pr.ProductController.SearchProducts)
	products.Get("/:id", pr.ProductController.GetProduct)
	products.Post("/", pr.ProductController.CreateProduct)
	products.Put("/:id", pr.ProductController.UpdateProduct)
	products.Patch("/:id/stock", pr.ProductController.UpdateProductStock)
	products.Delete("/:id", pr.ProductController.DeleteProduct)
}
