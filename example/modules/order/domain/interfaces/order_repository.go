package interfaces

import (
	"context"

	"example/modules/order/domain/entities"

	"github.com/google/uuid"
)

type OrderRepository interface {
	Create(ctx context.Context, order *entities.Order) error
	Update(ctx context.Context, order *entities.Order) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetByID(ctx context.Context, id uuid.UUID) (*entities.Order, error)
	GetByCustomerID(ctx context.Context, customerID uuid.UUID, limit, offset int32) ([]*entities.Order, error)
	GetAll(ctx context.Context, limit, offset int32) ([]*entities.Order, error)
	GetByStatus(ctx context.Context, status entities.OrderStatus, limit, offset int32) ([]*entities.Order, error)
	Count(ctx context.Context) (int64, error)
	CountByCustomerID(ctx context.Context, customerID uuid.UUID) (int64, error)
}
