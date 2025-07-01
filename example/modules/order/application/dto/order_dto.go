package dto

import (
	"time"

	"example/modules/order/domain/entities"

	"github.com/google/uuid"
)

type CreateOrderRequest struct {
	CustomerID      uuid.UUID                `json:"customer_id" validate:"required"`
	ShippingAddress *string                  `json:"shipping_address"`
	BillingAddress  *string                  `json:"billing_address"`
	Notes           *string                  `json:"notes"`
	Items           []CreateOrderItemRequest `json:"items" validate:"required,min=1"`
}

type CreateOrderItemRequest struct {
	ProductID   uuid.UUID `json:"product_id" validate:"required"`
	ProductName string    `json:"product_name" validate:"required"`
	Quantity    int32     `json:"quantity" validate:"required,min=1"`
	UnitPrice   float64   `json:"unit_price" validate:"required,min=0.01"`
}

type UpdateOrderRequest struct {
	Status          *entities.OrderStatus `json:"status"`
	ShippingCost    *float64              `json:"shipping_cost"`
	TaxAmount       *float64              `json:"tax_amount"`
	DiscountAmount  *float64              `json:"discount_amount"`
	ShippingAddress *string               `json:"shipping_address"`
	BillingAddress  *string               `json:"billing_address"`
	Notes           *string               `json:"notes"`
}

type AddOrderItemRequest struct {
	ProductID   uuid.UUID `json:"product_id" validate:"required"`
	ProductName string    `json:"product_name" validate:"required"`
	Quantity    int32     `json:"quantity" validate:"required,min=1"`
	UnitPrice   float64   `json:"unit_price" validate:"required,min=0.01"`
}

type UpdateOrderItemQuantityRequest struct {
	Quantity int32 `json:"quantity" validate:"required,min=1"`
}

type OrderResponse struct {
	ID              uuid.UUID            `json:"id"`
	CustomerID      uuid.UUID            `json:"customer_id"`
	Status          entities.OrderStatus `json:"status"`
	TotalAmount     float64              `json:"total_amount"`
	ShippingCost    float64              `json:"shipping_cost"`
	TaxAmount       float64              `json:"tax_amount"`
	DiscountAmount  float64              `json:"discount_amount"`
	Notes           *string              `json:"notes"`
	ShippingAddress *string              `json:"shipping_address"`
	BillingAddress  *string              `json:"billing_address"`
	OrderItems      []OrderItemResponse  `json:"order_items"`
	CreatedAt       time.Time            `json:"created_at"`
	UpdatedAt       time.Time            `json:"updated_at"`
}

type OrderItemResponse struct {
	ID          uuid.UUID `json:"id"`
	OrderID     uuid.UUID `json:"order_id"`
	ProductID   uuid.UUID `json:"product_id"`
	ProductName string    `json:"product_name"`
	Quantity    int32     `json:"quantity"`
	UnitPrice   float64   `json:"unit_price"`
	TotalPrice  float64   `json:"total_price"`
}

type OrderListResponse struct {
	Orders     []OrderResponse `json:"orders"`
	Total      int64           `json:"total"`
	Page       int32           `json:"page"`
	PageSize   int32           `json:"page_size"`
	TotalPages int32           `json:"total_pages"`
}

func ToOrderResponse(order *entities.Order) OrderResponse {
	items := make([]OrderItemResponse, len(order.OrderItems))
	for i, item := range order.OrderItems {
		items[i] = ToOrderItemResponse(item)
	}

	return OrderResponse{
		ID:              order.ID,
		CustomerID:      order.CustomerID,
		Status:          order.Status,
		TotalAmount:     order.TotalAmount,
		ShippingCost:    order.ShippingCost,
		TaxAmount:       order.TaxAmount,
		DiscountAmount:  order.DiscountAmount,
		Notes:           order.Notes,
		ShippingAddress: order.ShippingAddress,
		BillingAddress:  order.BillingAddress,
		OrderItems:      items,
		CreatedAt:       order.CreatedAt,
		UpdatedAt:       order.UpdatedAt,
	}
}

func ToOrderItemResponse(item *entities.OrderItem) OrderItemResponse {
	return OrderItemResponse{
		ID:          item.ID,
		OrderID:     item.OrderID,
		ProductID:   item.ProductID,
		ProductName: item.ProductName,
		Quantity:    item.Quantity,
		UnitPrice:   item.UnitPrice,
		TotalPrice:  item.TotalPrice,
	}
}

func ToOrderListResponse(orders []*entities.Order, total int64, page, pageSize int32) OrderListResponse {
	orderResponses := make([]OrderResponse, len(orders))
	for i, order := range orders {
		orderResponses[i] = ToOrderResponse(order)
	}

	totalPages := int32((total + int64(pageSize) - 1) / int64(pageSize))

	return OrderListResponse{
		Orders:     orderResponses,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}
}
