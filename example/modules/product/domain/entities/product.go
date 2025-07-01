package entities

import (
	"time"

	"github.com/google/uuid"
)

type Product struct {
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

func (p *Product) Validate() error {
	if p.Name == "" {
		return ErrProductNameRequired
	}
	if p.Price < 0 {
		return ErrProductPriceInvalid
	}
	if p.StockQuantity < 0 {
		return ErrProductStockInvalid
	}
	return nil
}

func (p *Product) UpdateStock(quantity int32) error {
	if quantity < 0 {
		return ErrProductStockInvalid
	}
	p.StockQuantity = quantity
	return nil
}
