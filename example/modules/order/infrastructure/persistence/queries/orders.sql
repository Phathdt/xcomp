-- Order queries
-- name: CreateOrder :one
INSERT INTO orders (
    id, customer_id, status, total_amount, shipping_cost, tax_amount,
    discount_amount, notes, shipping_address, billing_address, created_at, updated_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12
) RETURNING *;

-- name: GetOrderByID :one
SELECT * FROM orders WHERE id = $1;

-- name: GetOrdersByCustomerID :many
SELECT * FROM orders
WHERE customer_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: GetOrdersByStatus :many
SELECT * FROM orders
WHERE status = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: GetAllOrders :many
SELECT * FROM orders
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: UpdateOrder :one
UPDATE orders
SET status = $2, total_amount = $3, shipping_cost = $4, tax_amount = $5,
    discount_amount = $6, notes = $7, shipping_address = $8, billing_address = $9,
    updated_at = $10
WHERE id = $1
RETURNING *;

-- name: DeleteOrder :exec
DELETE FROM orders WHERE id = $1;

-- name: CountOrders :one
SELECT COUNT(*) FROM orders;

-- name: CountOrdersByCustomerID :one
SELECT COUNT(*) FROM orders WHERE customer_id = $1;

-- Order Item queries
-- name: CreateOrderItem :one
INSERT INTO order_items (
    id, order_id, product_id, product_name, quantity, unit_price, total_price
) VALUES (
    $1, $2, $3, $4, $5, $6, $7
) RETURNING *;

-- name: GetOrderItemByID :one
SELECT * FROM order_items WHERE id = $1;

-- name: GetOrderItemsByOrderID :many
SELECT * FROM order_items WHERE order_id = $1 ORDER BY id;

-- name: UpdateOrderItem :one
UPDATE order_items
SET quantity = $2, unit_price = $3, total_price = $4
WHERE id = $1
RETURNING *;

-- name: DeleteOrderItem :exec
DELETE FROM order_items WHERE id = $1;

-- name: DeleteOrderItemsByOrderID :exec
DELETE FROM order_items WHERE order_id = $1;
