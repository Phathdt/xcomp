package interfaces

import (
	"context"
	"time"

	"example/modules/order/domain/entities"

	"github.com/google/uuid"
)

type OrderCacheRepository interface {
	Get(ctx context.Context, id uuid.UUID) (*entities.Order, error)
	Set(ctx context.Context, order *entities.Order, expiration time.Duration) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetByCustomerID(ctx context.Context, customerID uuid.UUID) ([]*entities.Order, error)
	SetByCustomerID(ctx context.Context, customerID uuid.UUID, orders []*entities.Order, expiration time.Duration) error
	DeleteByCustomerID(ctx context.Context, customerID uuid.UUID) error
	Clear(ctx context.Context) error
}
