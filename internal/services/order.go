package services

import (
	"context"
	"github.com/google/uuid"
	"github.com/slamchillz/getinstashop-ecommerce-api/internal/constants"
	db "github.com/slamchillz/getinstashop-ecommerce-api/internal/db/sqlc"
	"github.com/slamchillz/getinstashop-ecommerce-api/internal/types"
	"github.com/slamchillz/getinstashop-ecommerce-api/internal/utils"
	"net/http"
)

// OrderService provides business logic for order operations.
type OrderService struct {
	store db.Store
}

// NewOrderService creates a new OrderService instance.
func NewOrderService(store db.Store) *OrderService {
	return &OrderService{
		store: store,
	}
}

func (s *OrderService) CreateOrder(ctx context.Context, orderReq types.CreateOrderInput) (*db.Order, *types.OrderErrMessage, int, error) {
	var productIds []uuid.UUID
	if len(orderReq.Items) <= 0 {
		return nil, &types.OrderErrMessage{
			Items: map[string]string{"key is productId": "value is quantity"},
		}, http.StatusBadRequest, nil
	}
	userId, _ := ctx.Value(constants.ContextUserIdKey).(uuid.UUID)
	for productId, _ := range orderReq.Items {
		productIds = append(productIds, productId)
	}
	order, orderErrMessage, execErr, txErr := s.store.CreateOrderTx(ctx, db.CreateOrderTxParams{
		ID:         uuid.New(),
		UserId:     userId,
		ProductIds: productIds,
		Items:      orderReq.Items,
	})
	if len(orderErrMessage) > 0 {
		return nil, &types.OrderErrMessage{
			Items: orderErrMessage,
		}, http.StatusBadRequest, nil
	}
	if execErr != nil || txErr != nil {
		return nil, nil, http.StatusInternalServerError, utils.ConcatenateErrors(execErr, txErr)
	}
	return &order, nil, http.StatusOK, nil
}
