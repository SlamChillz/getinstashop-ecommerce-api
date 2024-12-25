package types

import (
	db "github.com/slamchillz/getinstashop-ecommerce-api/internal/db/sqlc"
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
