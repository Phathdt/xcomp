package interfaces

import (
	"context"

	"example/modules/order/domain/entities"

	"github.com/google/uuid"
)

type OrderItemRepository interface {
	Create(ctx context.Context, item *entities.OrderItem) error
	Update(ctx context.Context, item *entities.OrderItem) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetByID(ctx context.Context, id uuid.UUID) (*entities.OrderItem, error)
	GetByOrderID(ctx context.Context, orderID uuid.UUID) ([]*entities.OrderItem, error)
	DeleteByOrderID(ctx context.Context, orderID uuid.UUID) error
}
