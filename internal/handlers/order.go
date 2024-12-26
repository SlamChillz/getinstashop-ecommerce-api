package handlers

import (
	"github.com/gin-gonic/gin"
	db "github.com/slamchillz/getinstashop-ecommerce-api/internal/db/sqlc"
	"github.com/slamchillz/getinstashop-ecommerce-api/internal/services"
	"github.com/slamchillz/getinstashop-ecommerce-api/internal/types"
	"github.com/slamchillz/getinstashop-ecommerce-api/internal/utils"
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

// CreateOrder godoc
// @Summary      Place an order for one or more Product
// @Description  Place an order for one or more Product
// @Tags         order
// @Accept       json
// @Produce      json
// @Param        payload   body	types.CreateOrderInput  true  "Create Order request body"
// @Success      200  {object}  types.Order
// @Failure      400  {object}  types.OrderError
// @Failure      500  {object}  types.InterServerError
// @Security	 BearerAuth
// @Router       /orders [post]
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
	if errMessage.Items != nil {
		ctx.JSON(statusCode, gin.H{
			"status":  "failed",
			"message": "Order not created",
			"error":   errMessage.Items,
		})
		log.Printf("Error while creating Order: %v", err)
		return
	}
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

// GetUserOrders godoc
// @Summary      Fetch all orders placed by a user
// @Description  Fetch all orders placed by a user
// @Tags         order
// @Accept       json
// @Produce      json
// @Success      200  {array}   types.Order
// @Failure      400  {object}  types.OrderError
// @Failure      500  {object}  types.InterServerError
// @Security	 BearerAuth
// @Router       /orders [get]
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

// CancelOrder godoc
// @Summary      Cancels and Order
// @Description  Cancels an order only if it is in PENDING state
// @Tags         order
// @Accept       json
// @Produce      json
// @Param        orderId   path		string  	true  "Unique uuid of the order whose status is to be cancelled"
// @Success      200  {object}  types.Order
// @Failure      400  {object}  types.OrderCancelError
// @Failure      500  {object}  types.InterServerError
// @Security	 BearerAuth
// @Router       /orders/{orderId} [patch]
func (h *OrderHandler) CancelOrder(ctx *gin.Context) {
	var err error
	orderId := utils.ParseStringToUUID(ctx.Param("id"))
	response, errMessage, statusCode, err := h.OrderService.CancelOrder(ctx, orderId)
	if errMessage.ID != "" {
		ctx.JSON(statusCode, gin.H{
			"status":  "failed",
			"message": "Order not cancelled",
			"error":   errMessage.ID,
		})
		log.Printf("Error while fetching user orders: %v", err)
		return
	}
	if err != nil {
		ctx.JSON(statusCode, gin.H{
			"status":  "failed",
			"message": "Order not cancelled",
			"error":   errMessage.ID,
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

// UpdateOrderStatus godoc
// @Summary      Updates the status of any order. Requires admin privilege
// @Description  Updates the status of any order. Requires admin privilege
// @Tags         order
// @Accept       json
// @Produce      json
// @Param        orderId   path		string  	true  "Unique uuid of the order whose status is to be updated"
// @Param        payload   body		types.UpdateOrderStatusInput  	true  "The new status of the order"
// @Success      200  {object}  types.Order
// @Failure      400  {object}  types.OrderCancelError
// @Failure      404  {object}  types.OrderCancelError
// @Failure      500  {object}  types.InterServerError
// @Security	 BearerAuth
// @Router       /admin/orders/{orderId} [patch]
func (h *OrderHandler) UpdateOrderStatus(ctx *gin.Context) {
	var err error
	var req types.UpdateOrderStatusInput
	orderId := utils.ParseStringToUUID(ctx.Param("id"))
	if err = ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "failed",
			"message": "Invalid JSON payload",
		})
		return
	}
	response, errMessage, statusCode, err := h.OrderService.UpdateOrderStatus(ctx, orderId, string(req.Status))
	if errMessage.Status != "" {
		ctx.JSON(statusCode, gin.H{
			"status":  "failed",
			"message": "Unable to update order status",
			"error":   errMessage,
		})
		log.Printf("Error while updating order status: %v", err)
		return
	}
	if err != nil {
		ctx.JSON(statusCode, gin.H{
			"status":  "failed",
			"message": "Unable to update order status",
			"error":   errMessage,
		})
		log.Printf("Error while updating order status: %v", err)
		return
	}
	ctx.JSON(statusCode, gin.H{
		"status":  "success",
		"message": "Order status updated successfully",
		"data":    response,
	})
}
