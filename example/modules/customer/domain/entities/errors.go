package entities

import (
	"errors"
)

var (
	ErrCustomerNotFound         = errors.New("customer not found")
	ErrCustomerUsernameRequired = errors.New("customer username is required")
	ErrCustomerEmailRequired    = errors.New("customer email is required")
	ErrCustomerUsernameExists   = errors.New("customer username already exists")
	ErrCustomerEmailExists      = errors.New("customer email already exists")
)
