package controllers

import (
	"strconv"

	"example/modules/customer/application/dto"
	"example/modules/customer/domain/entities"
	"example/modules/customer/domain/interfaces"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type CustomerController struct {
	CustomerService interfaces.CustomerService `inject:"CustomerService"`
}

func (cc *CustomerController) GetServiceName() string {
	return "CustomerController"
}

func (cc *CustomerController) GetCustomer(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Invalid customer ID",
			"message": "Customer ID must be a valid UUID",
		})
	}

	customer, err := cc.CustomerService.GetCustomer(c.Context(), id)
	if err != nil {
		if err == entities.ErrCustomerNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error":   "Customer not found",
				"message": "The requested customer does not exist",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Internal server error",
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    customer,
	})
}

func (cc *CustomerController) GetCustomerByUsername(c *fiber.Ctx) error {
	username := c.Params("username")
	if username == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Missing username",
			"message": "Username parameter is required",
		})
	}

	customer, err := cc.CustomerService.GetCustomerByUsername(c.Context(), username)
	if err != nil {
		if err == entities.ErrCustomerNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error":   "Customer not found",
				"message": "The requested customer does not exist",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Internal server error",
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    customer,
	})
}

func (cc *CustomerController) GetCustomerByEmail(c *fiber.Ctx) error {
	email := c.Query("email")
	if email == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Missing email",
			"message": "Email query parameter is required",
		})
	}

	customer, err := cc.CustomerService.GetCustomerByEmail(c.Context(), email)
	if err != nil {
		if err == entities.ErrCustomerNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error":   "Customer not found",
				"message": "The requested customer does not exist",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Internal server error",
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    customer,
	})
}

func (cc *CustomerController) ListCustomers(c *fiber.Ctx) error {
	page, _ := strconv.ParseInt(c.Query("page", "1"), 10, 32)
	pageSize, _ := strconv.ParseInt(c.Query("page_size", "10"), 10, 32)

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	customers, err := cc.CustomerService.ListCustomers(c.Context(), int32(page), int32(pageSize))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Internal server error",
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    customers,
	})
}

func (cc *CustomerController) SearchCustomers(c *fiber.Ctx) error {
	query := c.Query("q")
	if query == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Missing search query",
			"message": "Search query parameter 'q' is required",
		})
	}

	page, _ := strconv.ParseInt(c.Query("page", "1"), 10, 32)
	pageSize, _ := strconv.ParseInt(c.Query("page_size", "10"), 10, 32)

	searchReq := &dto.CustomerSearchRequest{
		Query:    query,
		Page:     int32(page),
		PageSize: int32(pageSize),
	}

	customers, err := cc.CustomerService.SearchCustomers(c.Context(), searchReq)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Internal server error",
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    customers,
	})
}

func (cc *CustomerController) CreateCustomer(c *fiber.Ctx) error {
	var req dto.CreateCustomerRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Invalid request body",
			"message": err.Error(),
		})
	}

	customer, err := cc.CustomerService.CreateCustomer(c.Context(), &req)
	if err != nil {
		if err == entities.ErrCustomerUsernameExists || err == entities.ErrCustomerEmailExists {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"error":   "Conflict",
				"message": err.Error(),
			})
		}
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Failed to create customer",
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"data":    customer,
	})
}

func (cc *CustomerController) UpdateCustomer(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Invalid customer ID",
			"message": "Customer ID must be a valid UUID",
		})
	}

	var req dto.UpdateCustomerRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Invalid request body",
			"message": err.Error(),
		})
	}

	customer, err := cc.CustomerService.UpdateCustomer(c.Context(), id, &req)
	if err != nil {
		if err == entities.ErrCustomerNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error":   "Customer not found",
				"message": "The requested customer does not exist",
			})
		}
		if err == entities.ErrCustomerUsernameExists || err == entities.ErrCustomerEmailExists {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"error":   "Conflict",
				"message": err.Error(),
			})
		}
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Failed to update customer",
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    customer,
	})
}

func (cc *CustomerController) DeleteCustomer(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Invalid customer ID",
			"message": "Customer ID must be a valid UUID",
		})
	}

	err = cc.CustomerService.DeleteCustomer(c.Context(), id)
	if err != nil {
		if err == entities.ErrCustomerNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error":   "Customer not found",
				"message": "The requested customer does not exist",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to delete customer",
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Customer deleted successfully",
	})
}
