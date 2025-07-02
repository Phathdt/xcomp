package interfaces

import (
	"context"
	"time"

	"example/modules/customer/domain/entities"

	"github.com/google/uuid"
)

type CustomerCacheRepository interface {
	Set(ctx context.Context, key string, customer *entities.Customer, ttl time.Duration) error
	Get(ctx context.Context, key string) (*entities.Customer, error)
	Delete(ctx context.Context, key string) error
	GetCustomerCacheKey(id uuid.UUID) string
	GetCustomerUsernameCacheKey(username string) string
	GetCustomerEmailCacheKey(email string) string
}
