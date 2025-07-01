package controllers

import (
	"strconv"

	"example/modules/order/application/dto"
	"example/modules/order/application/services"
	"example/modules/order/domain/entities"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type OrderController struct {
	orderService *services.OrderService `inject:"OrderService"`
}

func NewOrderController() *OrderController {
	return &OrderController{}
}

func (c *OrderController) CreateOrder(ctx *fiber.Ctx) error {
	var req dto.CreateOrderRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	order, err := c.orderService.CreateOrder(ctx.Context(), req)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusCreated).JSON(order)
}

func (c *OrderController) GetOrder(ctx *fiber.Ctx) error {
	idParam := ctx.Params("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid order ID",
		})
	}

	order, err := c.orderService.GetOrderByID(ctx.Context(), id)
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Order not found",
		})
	}

	return ctx.JSON(order)
}

func (c *OrderController) GetOrders(ctx *fiber.Ctx) error {
	page, _ := strconv.Atoi(ctx.Query("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.Query("page_size", "10"))
	customerIDParam := ctx.Query("customer_id")
	statusParam := ctx.Query("status")

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	var orders *dto.OrderListResponse
	var err error

	if customerIDParam != "" {
		customerID, parseErr := uuid.Parse(customerIDParam)
		if parseErr != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid customer ID",
			})
		}
		orders, err = c.orderService.GetOrdersByCustomerID(ctx.Context(), customerID, int32(page), int32(pageSize))
	} else if statusParam != "" {
		status := entities.OrderStatus(statusParam)
		orders, err = c.orderService.GetOrdersByStatus(ctx.Context(), status, int32(page), int32(pageSize))
	} else {
		orders, err = c.orderService.GetAllOrders(ctx.Context(), int32(page), int32(pageSize))
	}

	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.JSON(orders)
}

func (c *OrderController) UpdateOrder(ctx *fiber.Ctx) error {
	idParam := ctx.Params("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid order ID",
		})
	}

	var req dto.UpdateOrderRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	order, err := c.orderService.UpdateOrder(ctx.Context(), id, req)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.JSON(order)
}

func (c *OrderController) ConfirmOrder(ctx *fiber.Ctx) error {
	idParam := ctx.Params("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid order ID",
		})
	}

	order, err := c.orderService.ConfirmOrder(ctx.Context(), id)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.JSON(order)
}

func (c *OrderController) ShipOrder(ctx *fiber.Ctx) error {
	idParam := ctx.Params("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid order ID",
		})
	}

	order, err := c.orderService.ShipOrder(ctx.Context(), id)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.JSON(order)
}

func (c *OrderController) DeliverOrder(ctx *fiber.Ctx) error {
	idParam := ctx.Params("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid order ID",
		})
	}

	order, err := c.orderService.DeliverOrder(ctx.Context(), id)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.JSON(order)
}

func (c *OrderController) CancelOrder(ctx *fiber.Ctx) error {
	idParam := ctx.Params("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid order ID",
		})
	}

	order, err := c.orderService.CancelOrder(ctx.Context(), id)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.JSON(order)
}

func (c *OrderController) AddOrderItem(ctx *fiber.Ctx) error {
	idParam := ctx.Params("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid order ID",
		})
	}

	var req dto.AddOrderItemRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	order, err := c.orderService.AddOrderItem(ctx.Context(), id, req)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.JSON(order)
}

func (c *OrderController) UpdateOrderItemQuantity(ctx *fiber.Ctx) error {
	idParam := ctx.Params("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid order ID",
		})
	}

	productIDParam := ctx.Params("product_id")
	productID, err := uuid.Parse(productIDParam)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid product ID",
		})
	}

	var req dto.UpdateOrderItemQuantityRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	order, err := c.orderService.UpdateOrderItemQuantity(ctx.Context(), id, productID, req)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.JSON(order)
}

func (c *OrderController) RemoveOrderItem(ctx *fiber.Ctx) error {
	idParam := ctx.Params("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid order ID",
		})
	}

	productIDParam := ctx.Params("product_id")
	productID, err := uuid.Parse(productIDParam)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid product ID",
		})
	}

	order, err := c.orderService.RemoveOrderItem(ctx.Context(), id, productID)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.JSON(order)
}

func (c *OrderController) DeleteOrder(ctx *fiber.Ctx) error {
	idParam := ctx.Params("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid order ID",
		})
	}

	err = c.orderService.DeleteOrder(ctx.Context(), id)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusNoContent).Send(nil)
}
