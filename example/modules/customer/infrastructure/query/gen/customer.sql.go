// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: customer.sql

package gen

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const countCustomers = `-- name: CountCustomers :one
SELECT COUNT(*) FROM customers
`

func (q *Queries) CountCustomers(ctx context.Context) (int64, error) {
	row := q.db.QueryRow(ctx, countCustomers)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const createCustomer = `-- name: CreateCustomer :one
INSERT INTO customers (username, email)
VALUES ($1, $2)
RETURNING id, username, email, created_at, updated_at
`

type CreateCustomerParams struct {
	Username string `db:"username"`
	Email    string `db:"email"`
}

func (q *Queries) CreateCustomer(ctx context.Context, arg CreateCustomerParams) (*Customer, error) {
	row := q.db.QueryRow(ctx, createCustomer, arg.Username, arg.Email)
	var i Customer
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.Email,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return &i, err
}

const deleteCustomer = `-- name: DeleteCustomer :exec
DELETE FROM customers
WHERE id = $1
`

func (q *Queries) DeleteCustomer(ctx context.Context, id pgtype.UUID) error {
	_, err := q.db.Exec(ctx, deleteCustomer, id)
	return err
}

const getCustomer = `-- name: GetCustomer :one
SELECT id, username, email, created_at, updated_at
FROM customers
WHERE id = $1
`

func (q *Queries) GetCustomer(ctx context.Context, id pgtype.UUID) (*Customer, error) {
	row := q.db.QueryRow(ctx, getCustomer, id)
	var i Customer
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.Email,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return &i, err
}

const getCustomerByEmail = `-- name: GetCustomerByEmail :one
SELECT id, username, email, created_at, updated_at
FROM customers
WHERE email = $1
`

func (q *Queries) GetCustomerByEmail(ctx context.Context, email string) (*Customer, error) {
	row := q.db.QueryRow(ctx, getCustomerByEmail, email)
	var i Customer
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.Email,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return &i, err
}

const getCustomerByUsername = `-- name: GetCustomerByUsername :one
SELECT id, username, email, created_at, updated_at
FROM customers
WHERE username = $1
`

func (q *Queries) GetCustomerByUsername(ctx context.Context, username string) (*Customer, error) {
	row := q.db.QueryRow(ctx, getCustomerByUsername, username)
	var i Customer
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.Email,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return &i, err
}

const listCustomers = `-- name: ListCustomers :many
SELECT id, username, email, created_at, updated_at
FROM customers
ORDER BY created_at DESC
LIMIT $1 OFFSET $2
`

type ListCustomersParams struct {
	Limit  int32 `db:"limit"`
	Offset int32 `db:"offset"`
}

func (q *Queries) ListCustomers(ctx context.Context, arg ListCustomersParams) ([]*Customer, error) {
	rows, err := q.db.Query(ctx, listCustomers, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []*Customer
	for rows.Next() {
		var i Customer
		if err := rows.Scan(
			&i.ID,
			&i.Username,
			&i.Email,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, &i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const searchCustomers = `-- name: SearchCustomers :many
SELECT id, username, email, created_at, updated_at
FROM customers
WHERE (username ILIKE '%' || $1 || '%' OR email ILIKE '%' || $1 || '%')
ORDER BY created_at DESC
LIMIT $2 OFFSET $3
`

type SearchCustomersParams struct {
	Column1 *string `db:"column_1"`
	Limit   int32   `db:"limit"`
	Offset  int32   `db:"offset"`
}

func (q *Queries) SearchCustomers(ctx context.Context, arg SearchCustomersParams) ([]*Customer, error) {
	rows, err := q.db.Query(ctx, searchCustomers, arg.Column1, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []*Customer
	for rows.Next() {
		var i Customer
		if err := rows.Scan(
			&i.ID,
			&i.Username,
			&i.Email,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, &i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateCustomer = `-- name: UpdateCustomer :one
UPDATE customers
SET username = $2, email = $3, updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING id, username, email, created_at, updated_at
`

type UpdateCustomerParams struct {
	ID       pgtype.UUID `db:"id"`
	Username string      `db:"username"`
	Email    string      `db:"email"`
}

func (q *Queries) UpdateCustomer(ctx context.Context, arg UpdateCustomerParams) (*Customer, error) {
	row := q.db.QueryRow(ctx, updateCustomer, arg.ID, arg.Username, arg.Email)
	var i Customer
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.Email,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return &i, err
}
