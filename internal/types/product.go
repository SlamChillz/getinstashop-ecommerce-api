package types

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type CreateProductInput struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Stock       int     `json:"stock"`
}

type CreateProductErrMessage struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	Price       string `json:"price,omitempty"`
	Stock       string `json:"stock,omitempty"`
}

type CreateProductOutput struct {
	ID          uuid.UUID        `json:"id"`
	Name        string           `json:"name"`
	Description string           `json:"description"`
	Price       float64          `json:"price"`
	Stock       int32            `json:"stock"`
	CreatedAt   pgtype.Timestamp `json:"createdAt"`
	UpdatedAt   pgtype.Timestamp `json:"updatedAt"`
	CreatedBy   uuid.UUID        `json:"createdBy"`
}
