package entities

import "errors"

var (
	ErrOrderNotFound            = errors.New("order not found")
	ErrOrderItemNotFound        = errors.New("order item not found")
	ErrOrderAlreadyCancelled    = errors.New("order is already cancelled")
	ErrOrderAlreadyCompleted    = errors.New("order is already completed")
	ErrOrderCannotBeModified    = errors.New("order cannot be modified in current status")
	ErrInvalidOrderStatus       = errors.New("invalid order status")
	ErrOrderItemQuantityInvalid = errors.New("order item quantity must be greater than 0")
	ErrOrderItemPriceInvalid    = errors.New("order item price must be greater than 0")
	ErrOrderTotalMismatch       = errors.New("order total does not match sum of items")
	ErrEmptyOrder               = errors.New("order must contain at least one item")
)
