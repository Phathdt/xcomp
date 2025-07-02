package repositories

import (
	"context"
	"time"

	"example/modules/product/domain/entities"

	"github.com/google/uuid"
)

type ProductCacheRepository interface {
	Get(ctx context.Context, id uuid.UUID) (*entities.Product, error)
	Set(ctx context.Context, product *entities.Product, ttl time.Duration) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetServiceName() string
}
