package controllers

import (
	"strconv"

	"example/modules/product/application/dto"
	"example/modules/product/domain/entities"
	"example/modules/product/domain/interfaces"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type ProductController struct {
	ProductService interfaces.ProductService `inject:"ProductService"`
}

func (pc *ProductController) GetServiceName() string {
	return "ProductController"
}

func (pc *ProductController) GetProduct(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Invalid product ID",
			"message": "Product ID must be a valid UUID",
		})
	}

	product, err := pc.ProductService.GetProduct(c.Context(), id)
	if err != nil {
		if err == entities.ErrProductNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error":   "Product not found",
				"message": "The requested product does not exist",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Internal server error",
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    product,
	})
}

func (pc *ProductController) ListProducts(c *fiber.Ctx) error {
	page, _ := strconv.ParseInt(c.Query("page", "1"), 10, 32)
	pageSize, _ := strconv.ParseInt(c.Query("page_size", "10"), 10, 32)
	category := c.Query("category")

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	var products *dto.ProductListResponse
	var err error

	if category != "" {
		products, err = pc.ProductService.ListProductsByCategory(c.Context(), category, int32(page), int32(pageSize))
	} else {
		products, err = pc.ProductService.ListProducts(c.Context(), int32(page), int32(pageSize))
	}

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Internal server error",
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    products,
	})
}

func (pc *ProductController) SearchProducts(c *fiber.Ctx) error {
	query := c.Query("q")
	if query == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Missing search query",
			"message": "Search query parameter 'q' is required",
		})
	}

	page, _ := strconv.ParseInt(c.Query("page", "1"), 10, 32)
	pageSize, _ := strconv.ParseInt(c.Query("page_size", "10"), 10, 32)

	searchReq := &dto.ProductSearchRequest{
		Query:    query,
		Page:     int32(page),
		PageSize: int32(pageSize),
	}

	products, err := pc.ProductService.SearchProducts(c.Context(), searchReq)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Internal server error",
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    products,
	})
}

func (pc *ProductController) CreateProduct(c *fiber.Ctx) error {
	var req dto.CreateProductRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Invalid request body",
			"message": err.Error(),
		})
	}

	product, err := pc.ProductService.CreateProduct(c.Context(), &req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Failed to create product",
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"data":    product,
	})
}

func (pc *ProductController) UpdateProduct(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Invalid product ID",
			"message": "Product ID must be a valid UUID",
		})
	}

	var req dto.UpdateProductRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Invalid request body",
			"message": err.Error(),
		})
	}

	product, err := pc.ProductService.UpdateProduct(c.Context(), id, &req)
	if err != nil {
		if err == entities.ErrProductNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error":   "Product not found",
				"message": "The requested product does not exist",
			})
		}
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Failed to update product",
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    product,
	})
}

func (pc *ProductController) UpdateProductStock(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Invalid product ID",
			"message": "Product ID must be a valid UUID",
		})
	}

	var req dto.UpdateStockRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Invalid request body",
			"message": err.Error(),
		})
	}

	product, err := pc.ProductService.UpdateProductStock(c.Context(), id, &req)
	if err != nil {
		if err == entities.ErrProductNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error":   "Product not found",
				"message": "The requested product does not exist",
			})
		}
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Failed to update product stock",
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    product,
	})
}

func (pc *ProductController) DeleteProduct(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Invalid product ID",
			"message": "Product ID must be a valid UUID",
		})
	}

	err = pc.ProductService.DeleteProduct(c.Context(), id)
	if err != nil {
		if err == entities.ErrProductNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error":   "Product not found",
				"message": "The requested product does not exist",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to delete product",
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Product deleted successfully",
	})
}
