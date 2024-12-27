package handlers

import (
	"github.com/gin-gonic/gin"
	db "github.com/slamchillz/getinstashop-ecommerce-api/internal/db/sqlc"
	"github.com/slamchillz/getinstashop-ecommerce-api/internal/services"
	"github.com/slamchillz/getinstashop-ecommerce-api/internal/types"
	"github.com/slamchillz/getinstashop-ecommerce-api/pkg/token"
	"log"
	"net/http"
)

// UserHandler handles user-related operations.
type UserHandler struct {
	userService *services.UserService
}

// NewUserHandler creates a new UserHandler instance.
func NewUserHandler(store db.Store, jwtToken *token.JWT) *UserHandler {
	return &UserHandler{userService: services.NewUserService(store, jwtToken)}
}

// CreateUser godoc
// @Summary      New user signup. Create a new user
// @Description  New user signup. Register a new user
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        payload   body	types.RegisterUserInput  true  "Register request body"
// @Success      201  {object}  types.CreateUserOk
// @Failure      400  {object}  types.CreateUserError
// @Failure      500  {object}  types.InterServerError
// @Router       /auth/register [post]
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

// LoginUser godoc
// @Summary      User Login. Generates an access token for a valid user.
// @Description  User Login. Generates an access token for a valid user.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        payload   body	types.LoginUserInput  true  "Login request body"
// @Success      200  {object}  types.LoginUserOutput
// @Failure      400  {object}  types.LoginUserError
// @Failure      404  {object}  types.LoginUserError
// @Failure      500  {object}  types.InterServerError
// @Router       /auth/login [post]
func (h *UserHandler) LoginUser(ctx *gin.Context) {
	var err error
	var req types.LoginUserInput
	if err = ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "failed",
			"message": "Invalid JSON payload",
		})
		return
	}
	response, errMessage, statusCode, err := h.userService.LoginUser(ctx, req)
	if err != nil {
		ctx.JSON(statusCode, gin.H{
			"status":  "failed",
			"message": "User not authenticated",
			"error":   errMessage,
		})
		log.Printf("Error while creating user: %v", err)
		return
	}
	ctx.JSON(statusCode, gin.H{
		"status":  "success",
		"message": "User authenticated",
		"data":    response,
	})
}
