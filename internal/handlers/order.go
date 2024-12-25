package handlers

import (
	"github.com/gin-gonic/gin"
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
	ctx.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "Order created",
		"data":    response,
	})
}
