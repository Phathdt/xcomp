package interfaces

import (
	"context"

	"example/modules/order/application/dto"
	"example/modules/order/domain/entities"

	"github.com/google/uuid"
)

type OrderService interface {
	CreateOrder(ctx context.Context, req dto.CreateOrderRequest) (*dto.OrderResponse, error)
	GetOrderByID(ctx context.Context, id uuid.UUID) (*dto.OrderResponse, error)
	GetOrdersByCustomerID(ctx context.Context, customerID uuid.UUID, page, pageSize int32) (*dto.OrderListResponse, error)
	GetAllOrders(ctx context.Context, page, pageSize int32) (*dto.OrderListResponse, error)
	GetOrdersByStatus(ctx context.Context, status entities.OrderStatus, page, pageSize int32) (*dto.OrderListResponse, error)
	UpdateOrder(ctx context.Context, id uuid.UUID, req dto.UpdateOrderRequest) (*dto.OrderResponse, error)
	ConfirmOrder(ctx context.Context, id uuid.UUID) (*dto.OrderResponse, error)
	ShipOrder(ctx context.Context, id uuid.UUID) (*dto.OrderResponse, error)
	DeliverOrder(ctx context.Context, id uuid.UUID) (*dto.OrderResponse, error)
	CancelOrder(ctx context.Context, id uuid.UUID) (*dto.OrderResponse, error)
	AddOrderItem(ctx context.Context, orderID uuid.UUID, req dto.AddOrderItemRequest) (*dto.OrderResponse, error)
	UpdateOrderItemQuantity(ctx context.Context, orderID, productID uuid.UUID, req dto.UpdateOrderItemQuantityRequest) (*dto.OrderResponse, error)
	RemoveOrderItem(ctx context.Context, orderID, productID uuid.UUID) (*dto.OrderResponse, error)
	DeleteOrder(ctx context.Context, id uuid.UUID) error
}
