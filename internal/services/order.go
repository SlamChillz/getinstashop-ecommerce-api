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

func (s *OrderService) CreateOrder(ctx context.Context, orderReq types.CreateOrderInput) (db.Order, types.OrderErrMessage, int, error) {
	var productIds []uuid.UUID
	var items = make(map[uuid.UUID]int32)
	var order db.Order
	var errMessage types.OrderErrMessage
	if len(orderReq.Items) <= 0 {
		errMessage.Items = map[string]string{"key is productId": "value is quantity"}
		return db.Order{}, errMessage, http.StatusBadRequest, nil
	}
	userId, _ := ctx.Value(constants.ContextUserIdKey).(uuid.UUID)
	for _, item := range orderReq.Items {
		productId := uuid.UUID([]byte(item.ProductId))
		productIds = append(productIds, productId)
		items[productId] = item.Quantity
	}
	order, orderErrMessage, execErr, txErr := s.store.CreateOrderTx(ctx, db.CreateOrderTxParams{
		ID:         uuid.New(),
		UserId:     userId,
		ProductIds: productIds,
		Items:      items,
	})
	if len(orderErrMessage) > 0 {
		errMessage.Items = orderErrMessage
		return order, errMessage, http.StatusBadRequest, nil
	}
	if execErr != nil || txErr != nil {
		return order, errMessage, http.StatusInternalServerError, utils.ConcatenateErrors(execErr, txErr)
	}
	return order, errMessage, http.StatusCreated, nil
}

func (s *OrderService) GetUserOrders(ctx context.Context) ([]db.Order, int, error) {
	userId, _ := ctx.Value(constants.ContextUserIdKey).(uuid.UUID)
	userOrders, err := s.store.GetAllOrderByUserId(ctx, userId)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	return userOrders, http.StatusOK, nil
}

func (s *OrderService) CancelOrder(ctx context.Context, orderId uuid.UUID) (db.Order, types.OrderErrMessage, int, error) {
	var order db.Order
	var errMessage types.OrderErrMessage
	userId, _ := ctx.Value(constants.ContextUserIdKey).(uuid.UUID)
	order, err := s.store.CancelOrder(ctx, db.CancelOrderParams{
		ID:     orderId,
		UserId: userId,
	})
	if err != nil {
		errMessage.ID = "order not found or has already been canceled"
		return order, errMessage, http.StatusBadRequest, nil
	}
	return order, errMessage, http.StatusOK, nil
}

func (s *OrderService) UpdateOrderStatus(ctx context.Context, orderId uuid.UUID, payload types.UpdateOrderStatusInput) (db.Order, types.OrderErrMessage, int, error) {
	var order db.Order
	var errMessage types.OrderErrMessage
	order, err := s.store.UpdateOrderStatus(ctx, db.UpdateOrderStatusParams{
		ID:     orderId,
		Status: payload.Status,
	})
	if err != nil {
		if strings.Replace(sql.ErrNoRows.Error(), "sql: ", "", 1) == err.Error() {
			errMessage.ID = "order not found"
			return order, errMessage, http.StatusBadRequest, err
		}
		return order, errMessage, http.StatusInternalServerError, err
	}
	return order, errMessage, http.StatusOK, nil
}
