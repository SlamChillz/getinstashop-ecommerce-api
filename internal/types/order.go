package types

import "github.com/google/uuid"

type CreateOrderInput struct {
	Items map[uuid.UUID]int32 `json:"items"`
}

type OrderErrMessage struct {
	Items map[string]string `json:"items"`
}
