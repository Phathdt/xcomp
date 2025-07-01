package entities

import (
	"github.com/google/uuid"
)

type OrderItem struct {
	ID          uuid.UUID `json:"id"`
	OrderID     uuid.UUID `json:"order_id"`
	ProductID   uuid.UUID `json:"product_id"`
	ProductName string    `json:"product_name"`
	Quantity    int32     `json:"quantity"`
	UnitPrice   float64   `json:"unit_price"`
	TotalPrice  float64   `json:"total_price"`
}

func NewOrderItem(orderID, productID uuid.UUID, productName string, quantity int32, unitPrice float64) *OrderItem {
	return &OrderItem{
		ID:          uuid.New(),
		OrderID:     orderID,
		ProductID:   productID,
		ProductName: productName,
		Quantity:    quantity,
		UnitPrice:   unitPrice,
		TotalPrice:  float64(quantity) * unitPrice,
	}
}

func (oi *OrderItem) Validate() error {
	if oi.ProductID == uuid.Nil {
		return ErrOrderItemNotFound
	}

	if oi.Quantity <= 0 {
		return ErrOrderItemQuantityInvalid
	}

	if oi.UnitPrice <= 0 {
		return ErrOrderItemPriceInvalid
	}

	expectedTotal := float64(oi.Quantity) * oi.UnitPrice
	if abs(oi.TotalPrice-expectedTotal) > 0.01 {
		return ErrOrderTotalMismatch
	}

	return nil
}

func (oi *OrderItem) UpdateQuantity(newQuantity int32) error {
	if newQuantity <= 0 {
		return ErrOrderItemQuantityInvalid
	}

	oi.Quantity = newQuantity
	oi.TotalPrice = float64(newQuantity) * oi.UnitPrice
	return nil
}

func (oi *OrderItem) UpdatePrice(newUnitPrice float64) error {
	if newUnitPrice <= 0 {
		return ErrOrderItemPriceInvalid
	}

	oi.UnitPrice = newUnitPrice
	oi.TotalPrice = float64(oi.Quantity) * newUnitPrice
	return nil
}
