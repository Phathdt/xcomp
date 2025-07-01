package entities

import (
	"time"

	"github.com/google/uuid"
)

type OrderStatus string

const (
	OrderStatusPending   OrderStatus = "pending"
	OrderStatusConfirmed OrderStatus = "confirmed"
	OrderStatusShipped   OrderStatus = "shipped"
	OrderStatusDelivered OrderStatus = "delivered"
	OrderStatusCancelled OrderStatus = "cancelled"
)

type Order struct {
	ID              uuid.UUID    `json:"id"`
	CustomerID      uuid.UUID    `json:"customer_id"`
	Status          OrderStatus  `json:"status"`
	TotalAmount     float64      `json:"total_amount"`
	ShippingCost    float64      `json:"shipping_cost"`
	TaxAmount       float64      `json:"tax_amount"`
	DiscountAmount  float64      `json:"discount_amount"`
	Notes           *string      `json:"notes"`
	ShippingAddress *string      `json:"shipping_address"`
	BillingAddress  *string      `json:"billing_address"`
	OrderItems      []*OrderItem `json:"order_items"`
	CreatedAt       time.Time    `json:"created_at"`
	UpdatedAt       time.Time    `json:"updated_at"`
}

func NewOrder(customerID uuid.UUID) *Order {
	return &Order{
		ID:         uuid.New(),
		CustomerID: customerID,
		Status:     OrderStatusPending,
		OrderItems: make([]*OrderItem, 0),
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
}

func (o *Order) Validate() error {
	if o.CustomerID == uuid.Nil {
		return ErrOrderNotFound
	}

	if len(o.OrderItems) == 0 {
		return ErrEmptyOrder
	}

	if !o.isValidStatus() {
		return ErrInvalidOrderStatus
	}

	calculatedTotal := o.calculateItemsTotal()
	expectedTotal := calculatedTotal + o.ShippingCost + o.TaxAmount - o.DiscountAmount

	if abs(o.TotalAmount-expectedTotal) > 0.01 {
		return ErrOrderTotalMismatch
	}

	for _, item := range o.OrderItems {
		if err := item.Validate(); err != nil {
			return err
		}
	}

	return nil
}

func (o *Order) AddItem(productID uuid.UUID, productName string, quantity int32, unitPrice float64) error {
	if quantity <= 0 {
		return ErrOrderItemQuantityInvalid
	}

	if unitPrice <= 0 {
		return ErrOrderItemPriceInvalid
	}

	if !o.canBeModified() {
		return ErrOrderCannotBeModified
	}

	for _, item := range o.OrderItems {
		if item.ProductID == productID {
			item.Quantity += quantity
			item.TotalPrice = float64(item.Quantity) * item.UnitPrice
			o.UpdatedAt = time.Now()
			return nil
		}
	}

	orderItem := &OrderItem{
		ID:          uuid.New(),
		OrderID:     o.ID,
		ProductID:   productID,
		ProductName: productName,
		Quantity:    quantity,
		UnitPrice:   unitPrice,
		TotalPrice:  float64(quantity) * unitPrice,
	}

	o.OrderItems = append(o.OrderItems, orderItem)
	o.UpdatedAt = time.Now()
	return nil
}

func (o *Order) RemoveItem(productID uuid.UUID) error {
	if !o.canBeModified() {
		return ErrOrderCannotBeModified
	}

	for i, item := range o.OrderItems {
		if item.ProductID == productID {
			o.OrderItems = append(o.OrderItems[:i], o.OrderItems[i+1:]...)
			o.UpdatedAt = time.Now()
			return nil
		}
	}

	return ErrOrderItemNotFound
}

func (o *Order) UpdateItemQuantity(productID uuid.UUID, newQuantity int32) error {
	if newQuantity <= 0 {
		return ErrOrderItemQuantityInvalid
	}

	if !o.canBeModified() {
		return ErrOrderCannotBeModified
	}

	for _, item := range o.OrderItems {
		if item.ProductID == productID {
			item.Quantity = newQuantity
			item.TotalPrice = float64(newQuantity) * item.UnitPrice
			o.UpdatedAt = time.Now()
			return nil
		}
	}

	return ErrOrderItemNotFound
}

func (o *Order) ConfirmOrder() error {
	if o.Status != OrderStatusPending {
		return ErrOrderCannotBeModified
	}

	o.Status = OrderStatusConfirmed
	o.UpdatedAt = time.Now()
	return nil
}

func (o *Order) ShipOrder() error {
	if o.Status != OrderStatusConfirmed {
		return ErrOrderCannotBeModified
	}

	o.Status = OrderStatusShipped
	o.UpdatedAt = time.Now()
	return nil
}

func (o *Order) DeliverOrder() error {
	if o.Status != OrderStatusShipped {
		return ErrOrderCannotBeModified
	}

	o.Status = OrderStatusDelivered
	o.UpdatedAt = time.Now()
	return nil
}

func (o *Order) CancelOrder() error {
	if o.Status == OrderStatusCancelled {
		return ErrOrderAlreadyCancelled
	}

	if o.Status == OrderStatusDelivered {
		return ErrOrderAlreadyCompleted
	}

	o.Status = OrderStatusCancelled
	o.UpdatedAt = time.Now()
	return nil
}

func (o *Order) CalculateTotal() {
	itemsTotal := o.calculateItemsTotal()
	o.TotalAmount = itemsTotal + o.ShippingCost + o.TaxAmount - o.DiscountAmount
	o.UpdatedAt = time.Now()
}

func (o *Order) calculateItemsTotal() float64 {
	total := 0.0
	for _, item := range o.OrderItems {
		total += item.TotalPrice
	}
	return total
}

func (o *Order) canBeModified() bool {
	return o.Status == OrderStatusPending
}

func (o *Order) isValidStatus() bool {
	validStatuses := []OrderStatus{
		OrderStatusPending,
		OrderStatusConfirmed,
		OrderStatusShipped,
		OrderStatusDelivered,
		OrderStatusCancelled,
	}

	for _, status := range validStatuses {
		if o.Status == status {
			return true
		}
	}

	return false
}

func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}
