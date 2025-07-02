package interfaces

import (
	"context"

	"example/modules/customer/application/dto"

	"github.com/google/uuid"
)

type CustomerService interface {
	CreateCustomer(ctx context.Context, req *dto.CreateCustomerRequest) (*dto.CustomerResponse, error)
	UpdateCustomer(ctx context.Context, id uuid.UUID, req *dto.UpdateCustomerRequest) (*dto.CustomerResponse, error)
	DeleteCustomer(ctx context.Context, id uuid.UUID) error
	GetCustomer(ctx context.Context, id uuid.UUID) (*dto.CustomerResponse, error)
	GetCustomerByUsername(ctx context.Context, username string) (*dto.CustomerResponse, error)
	GetCustomerByEmail(ctx context.Context, email string) (*dto.CustomerResponse, error)
	ListCustomers(ctx context.Context, page, pageSize int32) (*dto.CustomerListResponse, error)
	SearchCustomers(ctx context.Context, req *dto.CustomerSearchRequest) (*dto.CustomerListResponse, error)
}
