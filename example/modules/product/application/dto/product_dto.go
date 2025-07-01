package dto

import (
	"time"

	"github.com/google/uuid"
)

type CreateProductRequest struct {
	Name          string  `json:"name" validate:"required,min=1,max=255"`
	Description   *string `json:"description" validate:"omitempty,max=1000"`
	Price         float64 `json:"price" validate:"required,gte=0"`
	StockQuantity int32   `json:"stock_quantity" validate:"gte=0"`
	Category      *string `json:"category" validate:"omitempty,max=100"`
}

type UpdateProductRequest struct {
	Name          string  `json:"name" validate:"required,min=1,max=255"`
	Description   *string `json:"description" validate:"omitempty,max=1000"`
	Price         float64 `json:"price" validate:"required,gte=0"`
	StockQuantity int32   `json:"stock_quantity" validate:"gte=0"`
	Category      *string `json:"category" validate:"omitempty,max=100"`
}

type UpdateStockRequest struct {
	StockQuantity int32 `json:"stock_quantity" validate:"gte=0"`
}

type ProductResponse struct {
	ID            uuid.UUID `json:"id"`
	Name          string    `json:"name"`
	Description   *string   `json:"description"`
	Price         float64   `json:"price"`
	StockQuantity int32     `json:"stock_quantity"`
	Category      *string   `json:"category"`
	IsActive      bool      `json:"is_active"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type ProductListResponse struct {
	Products   []*ProductResponse `json:"products"`
	TotalCount int64              `json:"total_count"`
	Page       int32              `json:"page"`
	PageSize   int32              `json:"page_size"`
	TotalPages int32              `json:"total_pages"`
}

type ProductSearchRequest struct {
	Query    string `json:"query" validate:"required,min=1"`
	Category string `json:"category" validate:"omitempty"`
	Page     int32  `json:"page" validate:"gte=1"`
	PageSize int32  `json:"page_size" validate:"gte=1,lte=100"`
}
