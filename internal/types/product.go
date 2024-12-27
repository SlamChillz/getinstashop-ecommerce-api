package types

import (
	"github.com/google/uuid"
	db "github.com/slamchillz/getinstashop-ecommerce-api/internal/db/sqlc"
	"time"
)

type CreateProductInput struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Stock       int     `json:"stock"`
}

type ProductErrMessage struct {
	ID          string `json:"id,omitempty"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	Price       string `json:"price,omitempty"`
	Stock       string `json:"stock,omitempty"`
}

type CreateProductOutput db.GetAllProductRow

type ProductOutput db.GetAllProductRow

type ProductUpdateInput struct {
	Name        *string  `json:"name,omitempty"`
	Description *string  `json:"description,omitempty"`
	Price       *float64 `json:"price,omitempty"`
	Stock       *int32   `json:"stock,omitempty"`
}

type Product struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	Stock       int32     `json:"stock"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	CreatedBy   uuid.UUID `json:"createdBy"`
}

type ProductError struct {
	Status  string            `json:"status"`
	Message string            `json:"message"`
	Error   ProductErrMessage `json:"error"`
}
