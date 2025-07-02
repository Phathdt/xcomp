-- name: GetCustomer :one
SELECT id, username, email, created_at, updated_at
FROM customers
WHERE id = $1;

-- name: GetCustomerByUsername :one
SELECT id, username, email, created_at, updated_at
FROM customers
WHERE username = $1;

-- name: GetCustomerByEmail :one
SELECT id, username, email, created_at, updated_at
FROM customers
WHERE email = $1;

-- name: ListCustomers :many
SELECT id, username, email, created_at, updated_at
FROM customers
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: SearchCustomers :many
SELECT id, username, email, created_at, updated_at
FROM customers
WHERE (username ILIKE '%' || $1 || '%' OR email ILIKE '%' || $1 || '%')
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: CreateCustomer :one
INSERT INTO customers (username, email)
VALUES ($1, $2)
RETURNING id, username, email, created_at, updated_at;

-- name: UpdateCustomer :one
UPDATE customers
SET username = $2, email = $3, updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING id, username, email, created_at, updated_at;

-- name: DeleteCustomer :exec
DELETE FROM customers
WHERE id = $1;

-- name: CountCustomers :one
SELECT COUNT(*) FROM customers;
