package types

import (
	"github.com/google/uuid"
	db "github.com/slamchillz/getinstashop-ecommerce-api/internal/db/sqlc"
	"time"
)

type OrderStatus db.OrderStatus

type Order struct {
	ID        uuid.UUID   `json:"id"`
	UserId    uuid.UUID   `json:"userId"`
	Total     float64     `json:"total"`
	Status    OrderStatus `json:"status"`
	CreatedAt time.Time   `json:"createdAt"`
	UpdatedAt time.Time   `json:"updatedAt"`
}

type Item struct {
	ProductId string `json:"productId"`
	Quantity  int32  `json:"quantity"`
}

type CreateOrderInput struct {
	Items []Item `json:"items"`
}

type OrderErrMessage struct {
	Items  map[string]string `json:"items,omitempty"`
	ID     string            `json:"id,omitempty"`
	Status string            `json:"status,omitempty"`
}

type ItemError struct {
	ProductId string `json:"productId"`
	Quantity  string `json:"quantity"`
}

type OrderError struct {
	Status  string      `json:"status,omitempty"`
	Message string      `json:"message,omitempty"`
	Error   []ItemError `json:"error,omitempty"`
}

type OrderCancelError struct {
	Status  string `json:"status,omitempty"`
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
}

type UpdateOrderStatusInput struct {
	Status db.OrderStatus `json:"status"`
}
