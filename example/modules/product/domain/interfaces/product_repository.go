package interfaces

import (
	"context"

	"example/modules/product/domain/entities"

	"github.com/google/uuid"
)

type ProductRepository interface {
	Create(ctx context.Context, product *entities.Product) (*entities.Product, error)
	Update(ctx context.Context, product *entities.Product) (*entities.Product, error)
	UpdateStock(ctx context.Context, id uuid.UUID, stockQuantity int32) (*entities.Product, error)
	Delete(ctx context.Context, id uuid.UUID) error
	GetByID(ctx context.Context, id uuid.UUID) (*entities.Product, error)
	List(ctx context.Context, limit, offset int32) ([]*entities.Product, error)
	ListByCategory(ctx context.Context, category string, limit, offset int32) ([]*entities.Product, error)
	Search(ctx context.Context, query string, limit, offset int32) ([]*entities.Product, error)
	Count(ctx context.Context) (int64, error)
	CountByCategory(ctx context.Context, category string) (int64, error)
}
