package persistence

import (
	"context"
	"log"
	"math/big"

	"example/modules/order/domain/entities"
	"example/modules/order/infrastructure/query/gen"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type OrderRepositoryImpl struct {
	db *pgxpool.Pool `inject:"DatabaseConnection"`
	q  *gen.Queries
}

type OrderItemRepositoryImpl struct {
	db *pgxpool.Pool `inject:"DatabaseConnection"`
	q  *gen.Queries
}

func (r *OrderRepositoryImpl) GetServiceName() string {
	return "OrderRepository"
}

func (r *OrderItemRepositoryImpl) GetServiceName() string {
	return "OrderItemRepository"
}

func (r *OrderRepositoryImpl) ensureQueries() {
	if r.q == nil {
		r.q = gen.New(r.db)
	}
}

func (r *OrderItemRepositoryImpl) ensureQueries() {
	if r.q == nil {
		r.q = gen.New(r.db)
	}
}

func (r *OrderRepositoryImpl) Create(ctx context.Context, order *entities.Order) error {
	r.ensureQueries()
	log.Printf("OrderRepository: Creating order %s", order.ID)

	params := gen.CreateOrderParams{
		ID:              uuidToPgUUID(order.ID),
		CustomerID:      uuidToPgUUID(order.CustomerID),
		Status:          string(order.Status),
		TotalAmount:     float64ToNumeric(order.TotalAmount),
		ShippingCost:    float64ToNumeric(order.ShippingCost),
		TaxAmount:       float64ToNumeric(order.TaxAmount),
		DiscountAmount:  float64ToNumeric(order.DiscountAmount),
		Notes:           order.Notes,
		ShippingAddress: order.ShippingAddress,
		BillingAddress:  order.BillingAddress,
		CreatedAt:       pgtype.Timestamptz{Time: order.CreatedAt, Valid: true},
		UpdatedAt:       pgtype.Timestamptz{Time: order.UpdatedAt, Valid: true},
	}

	_, err := r.q.CreateOrder(ctx, params)
	return err
}

func (r *OrderRepositoryImpl) GetByID(ctx context.Context, id uuid.UUID) (*entities.Order, error) {
	r.ensureQueries()
	log.Printf("OrderRepository: Getting order by ID %s", id)

	row, err := r.q.GetOrderByID(ctx, uuidToPgUUID(id))
	if err != nil {
		return nil, err
	}

	return convertOrderFromDB(*row), nil
}

func (r *OrderRepositoryImpl) GetByCustomerID(ctx context.Context, customerID uuid.UUID, limit, offset int32) ([]*entities.Order, error) {
	r.ensureQueries()
	log.Printf("OrderRepository: Getting orders for customer %s", customerID)

	params := gen.GetOrdersByCustomerIDParams{
		CustomerID: uuidToPgUUID(customerID),
		Limit:      limit,
		Offset:     offset,
	}

	rows, err := r.q.GetOrdersByCustomerID(ctx, params)
	if err != nil {
		return nil, err
	}

	orders := make([]*entities.Order, len(rows))
	for i, row := range rows {
		orders[i] = convertOrderFromDB(*row)
	}

	return orders, nil
}

func (r *OrderRepositoryImpl) Update(ctx context.Context, order *entities.Order) error {
	r.ensureQueries()
	log.Printf("OrderRepository: Updating order %s", order.ID)

	params := gen.UpdateOrderParams{
		ID:              uuidToPgUUID(order.ID),
		Status:          string(order.Status),
		TotalAmount:     float64ToNumeric(order.TotalAmount),
		ShippingCost:    float64ToNumeric(order.ShippingCost),
		TaxAmount:       float64ToNumeric(order.TaxAmount),
		DiscountAmount:  float64ToNumeric(order.DiscountAmount),
		Notes:           order.Notes,
		ShippingAddress: order.ShippingAddress,
		BillingAddress:  order.BillingAddress,
		UpdatedAt:       pgtype.Timestamptz{Time: order.UpdatedAt, Valid: true},
	}

	_, err := r.q.UpdateOrder(ctx, params)
	return err
}

func (r *OrderRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
	r.ensureQueries()
	log.Printf("OrderRepository: Deleting order %s", id)

	return r.q.DeleteOrder(ctx, uuidToPgUUID(id))
}

func (r *OrderRepositoryImpl) GetByStatus(ctx context.Context, status entities.OrderStatus, limit, offset int32) ([]*entities.Order, error) {
	r.ensureQueries()
	log.Printf("OrderRepository: Getting orders by status %s", status)

	params := gen.GetOrdersByStatusParams{
		Status: string(status),
		Limit:  limit,
		Offset: offset,
	}

	rows, err := r.q.GetOrdersByStatus(ctx, params)
	if err != nil {
		return nil, err
	}

	orders := make([]*entities.Order, len(rows))
	for i, row := range rows {
		orders[i] = convertOrderFromDB(*row)
	}

	return orders, nil
}

func (r *OrderRepositoryImpl) GetAll(ctx context.Context, limit, offset int32) ([]*entities.Order, error) {
	r.ensureQueries()
	log.Printf("OrderRepository: Getting all orders")

	params := gen.GetAllOrdersParams{
		Limit:  limit,
		Offset: offset,
	}

	rows, err := r.q.GetAllOrders(ctx, params)
	if err != nil {
		return nil, err
	}

	orders := make([]*entities.Order, len(rows))
	for i, row := range rows {
		orders[i] = convertOrderFromDB(*row)
	}

	return orders, nil
}

func (r *OrderRepositoryImpl) Count(ctx context.Context) (int64, error) {
	r.ensureQueries()
	log.Printf("OrderRepository: Counting orders")

	return r.q.CountOrders(ctx)
}

func (r *OrderRepositoryImpl) CountByCustomerID(ctx context.Context, customerID uuid.UUID) (int64, error) {
	r.ensureQueries()
	log.Printf("OrderRepository: Counting orders for customer %s", customerID)

	return r.q.CountOrdersByCustomerID(ctx, uuidToPgUUID(customerID))
}

func (r *OrderItemRepositoryImpl) Create(ctx context.Context, orderItem *entities.OrderItem) error {
	r.ensureQueries()
	log.Printf("OrderItemRepository: Creating order item %s", orderItem.ID)

	params := gen.CreateOrderItemParams{
		ID:          uuidToPgUUID(orderItem.ID),
		OrderID:     uuidToPgUUID(orderItem.OrderID),
		ProductID:   uuidToPgUUID(orderItem.ProductID),
		ProductName: orderItem.ProductName,
		Quantity:    orderItem.Quantity,
		UnitPrice:   float64ToNumeric(orderItem.UnitPrice),
		TotalPrice:  float64ToNumeric(orderItem.TotalPrice),
	}

	_, err := r.q.CreateOrderItem(ctx, params)
	return err
}

func (r *OrderItemRepositoryImpl) GetByID(ctx context.Context, id uuid.UUID) (*entities.OrderItem, error) {
	r.ensureQueries()
	log.Printf("OrderItemRepository: Getting order item by ID %s", id)

	row, err := r.q.GetOrderItemByID(ctx, uuidToPgUUID(id))
	if err != nil {
		return nil, err
	}

	return convertOrderItemFromDB(*row), nil
}

func (r *OrderItemRepositoryImpl) GetByOrderID(ctx context.Context, orderID uuid.UUID) ([]*entities.OrderItem, error) {
	r.ensureQueries()
	log.Printf("OrderItemRepository: Getting order items for order %s", orderID)

	rows, err := r.q.GetOrderItemsByOrderID(ctx, uuidToPgUUID(orderID))
	if err != nil {
		return nil, err
	}

	orderItems := make([]*entities.OrderItem, len(rows))
	for i, row := range rows {
		orderItems[i] = convertOrderItemFromDB(*row)
	}

	return orderItems, nil
}

func (r *OrderItemRepositoryImpl) Update(ctx context.Context, orderItem *entities.OrderItem) error {
	r.ensureQueries()
	log.Printf("OrderItemRepository: Updating order item %s", orderItem.ID)

	params := gen.UpdateOrderItemParams{
		ID:         uuidToPgUUID(orderItem.ID),
		Quantity:   orderItem.Quantity,
		UnitPrice:  float64ToNumeric(orderItem.UnitPrice),
		TotalPrice: float64ToNumeric(orderItem.TotalPrice),
	}

	_, err := r.q.UpdateOrderItem(ctx, params)
	return err
}

func (r *OrderItemRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
	r.ensureQueries()
	log.Printf("OrderItemRepository: Deleting order item %s", id)

	return r.q.DeleteOrderItem(ctx, uuidToPgUUID(id))
}

func (r *OrderItemRepositoryImpl) DeleteByOrderID(ctx context.Context, orderID uuid.UUID) error {
	r.ensureQueries()
	log.Printf("OrderItemRepository: Deleting order items for order %s", orderID)

	return r.q.DeleteOrderItemsByOrderID(ctx, uuidToPgUUID(orderID))
}

func (r *OrderItemRepositoryImpl) CreateBatch(ctx context.Context, orderItems []*entities.OrderItem) error {
	r.ensureQueries()
	log.Printf("OrderItemRepository: Creating batch of %d order items", len(orderItems))

	for _, orderItem := range orderItems {
		if err := r.Create(ctx, orderItem); err != nil {
			return err
		}
	}

	return nil
}

func convertOrderFromDB(row gen.Order) *entities.Order {
	order := &entities.Order{
		ID:              pgUUIDToUUID(row.ID),
		CustomerID:      pgUUIDToUUID(row.CustomerID),
		Status:          entities.OrderStatus(row.Status),
		TotalAmount:     numericToFloat64(row.TotalAmount),
		ShippingCost:    numericToFloat64(row.ShippingCost),
		TaxAmount:       numericToFloat64(row.TaxAmount),
		DiscountAmount:  numericToFloat64(row.DiscountAmount),
		Notes:           row.Notes,
		ShippingAddress: row.ShippingAddress,
		BillingAddress:  row.BillingAddress,
	}

	if row.CreatedAt.Valid {
		order.CreatedAt = row.CreatedAt.Time
	}

	if row.UpdatedAt.Valid {
		order.UpdatedAt = row.UpdatedAt.Time
	}

	return order
}

func convertOrderItemFromDB(row gen.OrderItem) *entities.OrderItem {
	return &entities.OrderItem{
		ID:          pgUUIDToUUID(row.ID),
		OrderID:     pgUUIDToUUID(row.OrderID),
		ProductID:   pgUUIDToUUID(row.ProductID),
		ProductName: row.ProductName,
		Quantity:    row.Quantity,
		UnitPrice:   numericToFloat64(row.UnitPrice),
		TotalPrice:  numericToFloat64(row.TotalPrice),
	}
}

func uuidToPgUUID(u uuid.UUID) pgtype.UUID {
	return pgtype.UUID{
		Bytes: u,
		Valid: true,
	}
}

func pgUUIDToUUID(u pgtype.UUID) uuid.UUID {
	if !u.Valid {
		return uuid.Nil
	}
	return u.Bytes
}

func float64ToNumeric(f float64) pgtype.Numeric {
	cents := int64(f * 100)
	return pgtype.Numeric{
		Int:   big.NewInt(cents),
		Valid: true,
	}
}

func numericToFloat64(n pgtype.Numeric) float64 {
	if !n.Valid {
		return 0.0
	}
	return float64(n.Int.Int64()) / 100.0
}
