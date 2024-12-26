package services

import (
	"context"
	"database/sql"
	"github.com/google/uuid"
	"github.com/slamchillz/getinstashop-ecommerce-api/internal/constants"
	db "github.com/slamchillz/getinstashop-ecommerce-api/internal/db/sqlc"
	"github.com/slamchillz/getinstashop-ecommerce-api/internal/types"
	"github.com/slamchillz/getinstashop-ecommerce-api/internal/utils"
	"log"
	"net/http"
	"strings"
)

var (
	OrderCancelled = "CANCELLED"
	OrderCompleted = "COMPLETED"
	OrderPending   = "PENDING"
	orderStatus    = map[string]int8{OrderCancelled: 1, OrderPending: 1, OrderCompleted: 1}
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
	log.Printf("Items: %+v", items)
	if len(orderReq.Items) <= 0 {
		errMessage.Items = map[string]string{"productId": "must be a valid product id", "quantity": "must be greater than zero"}
		return db.Order{}, errMessage, http.StatusBadRequest, nil
	}
	for _, item := range orderReq.Items {
		if item.Quantity <= 0 {
			errMessage.Items = map[string]string{"productId": "must be a product id", "quantity": "must be greater than zero"}
			return db.Order{}, errMessage, http.StatusBadRequest, nil
		}
	}
	userId, _ := ctx.Value(constants.ContextUserIdKey).(uuid.UUID)
	for _, item := range orderReq.Items {
		productId := utils.ParseStringToUUID(item.ProductId)
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
	log.Printf("Order Id: %v", orderId)
	userId, _ := ctx.Value(constants.ContextUserIdKey).(uuid.UUID)
	order, err := s.store.UpdateOrderTx(ctx, db.UpdateOrderTxParams{
		ID:     orderId,
		UserId: userId,
		Admin:  false,
		Status: db.OrderStatusCANCELLED,
	})
	if err != nil {
		errMessage.ID = "order not found or has already been cancelled"
		return order, errMessage, http.StatusBadRequest, err
	}
	return order, errMessage, http.StatusOK, nil
}

func (s *OrderService) UpdateOrderStatus(ctx context.Context, orderId uuid.UUID, status string) (db.Order, types.OrderErrMessage, int, error) {
	var order db.Order
	var errMessage types.OrderErrMessage
	status = strings.ToUpper(status)
	userId, _ := ctx.Value(constants.ContextUserIdKey).(uuid.UUID)
	if orderStatus[status] != 1 {
		errMessage.Status = "Unknown order status"
		return db.Order{}, errMessage, http.StatusBadRequest, nil
	}
	order, err := s.store.UpdateOrderTx(ctx, db.UpdateOrderTxParams{
		ID:     orderId,
		Status: db.OrderStatus(status),
		Admin:  true,
		UserId: userId,
	})
	log.Printf("Error: %v", err)
	if err != nil {
		if strings.Replace(sql.ErrNoRows.Error(), "sql: ", "", 1) == err.Error() {
			errMessage.ID = "order id not found"
			return order, errMessage, http.StatusNotFound, err
		}
		return order, errMessage, http.StatusInternalServerError, err
	}
	return order, errMessage, http.StatusOK, nil
}
