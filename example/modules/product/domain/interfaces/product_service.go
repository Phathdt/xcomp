package interfaces

import (
	"context"

	"example/modules/product/application/dto"

	"github.com/google/uuid"
)

type ProductService interface {
	GetServiceName() string
	GetProduct(ctx context.Context, id uuid.UUID) (*dto.ProductResponse, error)
	ListProducts(ctx context.Context, page, pageSize int32) (*dto.ProductListResponse, error)
	ListProductsByCategory(ctx context.Context, category string, page, pageSize int32) (*dto.ProductListResponse, error)
	SearchProducts(ctx context.Context, searchReq *dto.ProductSearchRequest) (*dto.ProductListResponse, error)
	CreateProduct(ctx context.Context, req *dto.CreateProductRequest) (*dto.ProductResponse, error)
	UpdateProduct(ctx context.Context, id uuid.UUID, req *dto.UpdateProductRequest) (*dto.ProductResponse, error)
	UpdateProductStock(ctx context.Context, id uuid.UUID, req *dto.UpdateStockRequest) (*dto.ProductResponse, error)
	DeleteProduct(ctx context.Context, id uuid.UUID) error
}
