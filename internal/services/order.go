package services

import (
	"context"
	"database/sql"
	"github.com/google/uuid"
	"github.com/slamchillz/getinstashop-ecommerce-api/internal/constants"
	db "github.com/slamchillz/getinstashop-ecommerce-api/internal/db/sqlc"
	"github.com/slamchillz/getinstashop-ecommerce-api/internal/types"
	"github.com/slamchillz/getinstashop-ecommerce-api/internal/utils"
	"net/http"
	"strings"
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

func (s *OrderService) GetUserOrders(ctx context.Context) ([]db.Order, int, error) {
	userId, _ := ctx.Value(constants.ContextUserIdKey).(uuid.UUID)
	userOrders, err := s.store.GetAllOrderByUserId(ctx, userId)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	return userOrders, http.StatusOK, nil
}

func (s *OrderService) CancelOrder(ctx context.Context, orderId uuid.UUID) (*db.Order, *types.OrderErrMessage, int, error) {
	userId, _ := ctx.Value(constants.ContextUserIdKey).(uuid.UUID)
	order, err := s.store.CancelOrder(ctx, db.CancelOrderParams{
		ID:     orderId,
		UserId: userId,
	})
	if err != nil {
		return nil, &types.OrderErrMessage{
			ID: "order not found or has already been canceled",
		}, http.StatusBadRequest, nil
	}
	return &order, nil, http.StatusOK, nil
}

func (s *OrderService) UpdateOrderStatus(ctx context.Context, orderId uuid.UUID, payload types.UpdateOrderStatusInput) (*db.Order, *types.OrderErrMessage, int, error) {
	order, err := s.store.UpdateOrderStatus(ctx, db.UpdateOrderStatusParams{
		ID:     orderId,
		Status: payload.Status,
	})
	if err != nil {
		if strings.Replace(sql.ErrNoRows.Error(), "sql: ", "", 1) == err.Error() {
			return nil, &types.OrderErrMessage{
				ID: "order does not exist",
			}, http.StatusBadRequest, err
		}
		return nil, nil, http.StatusInternalServerError, err
	}
	return &order, nil, http.StatusOK, nil
}
