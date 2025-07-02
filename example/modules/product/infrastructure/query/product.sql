-- name: GetProduct :one
SELECT id, name, description, price, stock_quantity, category, is_active, created_at, updated_at
FROM products
WHERE id = $1 AND is_active = true;

-- name: ListProducts :many
SELECT id, name, description, price, stock_quantity, category, is_active, created_at, updated_at
FROM products
WHERE is_active = true
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: ListProductsByCategory :many
SELECT id, name, description, price, stock_quantity, category, is_active, created_at, updated_at
FROM products
WHERE category = $1 AND is_active = true
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: SearchProducts :many
SELECT id, name, description, price, stock_quantity, category, is_active, created_at, updated_at
FROM products
WHERE (name ILIKE '%' || $1 || '%' OR description ILIKE '%' || $1 || '%')
  AND is_active = true
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: CreateProduct :one
INSERT INTO products (name, description, price, stock_quantity, category)
VALUES ($1, $2, $3, $4, $5)
RETURNING id, name, description, price, stock_quantity, category, is_active, created_at, updated_at;

-- name: UpdateProduct :one
UPDATE products
SET name = $2, description = $3, price = $4, stock_quantity = $5, category = $6, updated_at = CURRENT_TIMESTAMP
WHERE id = $1 AND is_active = true
RETURNING id, name, description, price, stock_quantity, category, is_active, created_at, updated_at;

-- name: UpdateProductStock :one
UPDATE products
SET stock_quantity = $2, updated_at = CURRENT_TIMESTAMP
WHERE id = $1 AND is_active = true
RETURNING id, name, description, price, stock_quantity, category, is_active, created_at, updated_at;

-- name: DeleteProduct :exec
UPDATE products
SET is_active = false, updated_at = CURRENT_TIMESTAMP
WHERE id = $1;

-- name: CountProducts :one
SELECT COUNT(*) FROM products WHERE is_active = true;

-- name: CountProductsByCategory :one
SELECT COUNT(*) FROM products WHERE category = $1 AND is_active = true;
