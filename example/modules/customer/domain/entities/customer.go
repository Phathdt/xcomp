package entities

import (
	"time"

	"github.com/google/uuid"
)

type Customer struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (c *Customer) Validate() error {
	if c.Username == "" {
		return ErrCustomerUsernameRequired
	}
	if c.Email == "" {
		return ErrCustomerEmailRequired
	}
	return nil
}
