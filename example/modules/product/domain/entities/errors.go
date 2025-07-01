package entities

import "errors"

var (
	ErrProductNotFound      = errors.New("product not found")
	ErrProductNameRequired  = errors.New("product name is required")
	ErrProductPriceInvalid  = errors.New("product price must be greater than or equal to 0")
	ErrProductStockInvalid  = errors.New("product stock quantity must be greater than or equal to 0")
	ErrProductAlreadyExists = errors.New("product already exists")
)
