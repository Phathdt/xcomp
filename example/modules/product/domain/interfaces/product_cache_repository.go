package interfaces

import (
	"context"
	"time"

	"example/modules/product/domain/entities"

	"github.com/google/uuid"
)

type ProductCacheRepository interface {
	Get(ctx context.Context, id uuid.UUID) (*entities.Product, error)
	Set(ctx context.Context, product *entities.Product, expiration time.Duration) error
	Delete(ctx context.Context, id uuid.UUID) error
	Clear(ctx context.Context) error
}
