package interfaces

import (
	"context"

	"example/modules/customer/domain/entities"

	"github.com/google/uuid"
)

type CustomerRepository interface {
	Create(ctx context.Context, customer *entities.Customer) (*entities.Customer, error)
	Update(ctx context.Context, customer *entities.Customer) (*entities.Customer, error)
	Delete(ctx context.Context, id uuid.UUID) error
	GetByID(ctx context.Context, id uuid.UUID) (*entities.Customer, error)
	GetByUsername(ctx context.Context, username string) (*entities.Customer, error)
	GetByEmail(ctx context.Context, email string) (*entities.Customer, error)
	List(ctx context.Context, limit, offset int32) ([]*entities.Customer, error)
	Search(ctx context.Context, query string, limit, offset int32) ([]*entities.Customer, error)
	Count(ctx context.Context) (int64, error)
}
