package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	db "github.com/slamchillz/getinstashop-ecommerce-api/internal/db/sqlc"
	"github.com/slamchillz/getinstashop-ecommerce-api/internal/services"
	"github.com/slamchillz/getinstashop-ecommerce-api/internal/types"
	"log"
	"net/http"
)

// OrderHandler handles order related operations.
type OrderHandler struct {
	OrderService *services.OrderService
}

// NewOrderHandler creates a new OrderHandler instance.
func NewOrderHandler(store db.Store) *OrderHandler {
	return &OrderHandler{OrderService: services.NewOrderService(store)}
}

func (h *OrderHandler) CreateOrder(ctx *gin.Context) {
	var err error
	var req types.CreateOrderInput
	if err = ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "failed",
			"message": "Invalid JSON payload",
		})
		return
	}
	response, errMessage, statusCode, err := h.OrderService.CreateOrder(ctx, req)
	if err != nil {
		ctx.JSON(statusCode, gin.H{
			"status":  "failed",
			"message": "Order not created",
			"error":   errMessage,
		})
		log.Printf("Error while creating Order: %v", err)
		return
	}
	ctx.JSON(statusCode, gin.H{
		"status":  "success",
		"message": "Order created",
		"data":    response,
	})
}

func (h *OrderHandler) GetUserOrders(ctx *gin.Context) {
	var err error
	response, statusCode, err := h.OrderService.GetUserOrders(ctx)
	if err != nil {
		ctx.JSON(statusCode, gin.H{
			"status":  "failed",
			"message": "Unable to get user orders",
			"error":   gin.H{},
		})
		log.Printf("Error while fetching user orders: %v", err)
		return
	}
	ctx.JSON(statusCode, gin.H{
		"status":  "success",
		"message": "User orders fetched successfully",
		"data":    response,
	})
}

func (h *OrderHandler) CancelOrder(ctx *gin.Context) {
	var err error
	orderId := uuid.UUID([]byte(ctx.Param("id")))
	response, errMessage, statusCode, err := h.OrderService.CancelOrder(ctx, orderId)
	if err != nil {
		ctx.JSON(statusCode, gin.H{
			"status":  "failed",
			"message": "Order not cancelled",
			"error":   errMessage,
		})
		log.Printf("Error while fetching user orders: %v", err)
		return
	}
	ctx.JSON(statusCode, gin.H{
		"status":  "success",
		"message": "Order cancelled successfully",
		"data":    response,
	})
}

func (h *OrderHandler) UpdateOrderStatus(ctx *gin.Context) {
	var err error
	var req types.UpdateOrderStatusInput
	orderId := uuid.UUID([]byte(ctx.Param("id")))
	if err = ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "failed",
			"message": "Invalid JSON payload",
		})
		return
	}
	response, errMessage, statusCode, err := h.OrderService.UpdateOrderStatus(ctx, orderId, req)
	if err != nil {
		ctx.JSON(statusCode, gin.H{
			"status":  "failed",
			"message": "Unable to update order status",
			"error":   errMessage,
		})
		log.Printf("Error while fetching user orders: %v", err)
		return
	}
	ctx.JSON(statusCode, gin.H{
		"status":  "success",
		"message": "Order status updated successfully",
		"data":    response,
	})
}
