// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0

package gen

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type Customer struct {
	ID        pgtype.UUID        `db:"id"`
	Username  string             `db:"username"`
	Email     string             `db:"email"`
	CreatedAt pgtype.Timestamptz `db:"created_at"`
	UpdatedAt pgtype.Timestamptz `db:"updated_at"`
}
