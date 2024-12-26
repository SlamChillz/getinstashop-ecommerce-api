package types

import (
	"github.com/google/uuid"
	db "github.com/slamchillz/getinstashop-ecommerce-api/internal/db/sqlc"
)

type CreateOrderInput struct {
	Items map[uuid.UUID]int32 `json:"items"`
}

type OrderErrMessage struct {
	Items  map[string]string `json:"items,omitempty"`
	ID     string            `json:"id,omitempty"`
	Status string            `json:"status,omitempty"`
}

type UpdateOrderStatusInput struct {
	Status db.OrderStatus `json:"status"`
}
