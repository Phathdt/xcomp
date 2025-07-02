package dto

import (
	"time"

	"github.com/google/uuid"
)

type CreateCustomerRequest struct {
	Username string `json:"username" validate:"required,min=3,max=100"`
	Email    string `json:"email" validate:"required,email,max=255"`
}

type UpdateCustomerRequest struct {
	Username string `json:"username" validate:"required,min=3,max=100"`
	Email    string `json:"email" validate:"required,email,max=255"`
}

type CustomerResponse struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CustomerListResponse struct {
	Customers  []*CustomerResponse `json:"customers"`
	TotalCount int64               `json:"total_count"`
	Page       int32               `json:"page"`
	PageSize   int32               `json:"page_size"`
	TotalPages int32               `json:"total_pages"`
}

type CustomerSearchRequest struct {
	Query    string `json:"query" validate:"required,min=1"`
	Page     int32  `json:"page" validate:"gte=1"`
	PageSize int32  `json:"page_size" validate:"gte=1,lte=100"`
}
