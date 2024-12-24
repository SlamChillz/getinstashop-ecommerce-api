package handlers

import (
	"github.com/gin-gonic/gin"
	db "github.com/slamchillz/getinstashop-ecommerce-api/internal/db/sqlc"
	"github.com/slamchillz/getinstashop-ecommerce-api/internal/services"
	"github.com/slamchillz/getinstashop-ecommerce-api/internal/types"
	"log"
	"net/http"
)

// UserHandler handles user-related operations.
type UserHandler struct {
	userService *services.UserService
}

// NewUserHandler creates a new UserHandler instance.
func NewUserHandler(store db.Store) *UserHandler {
	return &UserHandler{userService: services.NewUserService(store)}
}

func (h *UserHandler) CreateUser(ctx *gin.Context) {
	var err error
	var req types.RegisterUserInput
	if err = ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "failed",
			"message": "Invalid JSON payload",
		})
		return
	}
	response, errMessage, statusCode, err := h.userService.CreateUser(ctx, req)
	if err != nil {
		ctx.JSON(statusCode, gin.H{
			"status":  "failed",
			"message": "User not created",
			"error":   errMessage,
		})
		log.Printf("Error while creating user: %v", err)
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "User created",
		"data":    response,
	})
}
