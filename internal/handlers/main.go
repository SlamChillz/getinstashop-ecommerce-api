package handlers

import (
	"github.com/gin-gonic/gin"
	db "github.com/slamchillz/getinstashop-ecommerce-api/internal/db/sqlc"
	"github.com/slamchillz/getinstashop-ecommerce-api/pkg/token"
)

type AllHandler struct {
	*UserHandler
	*ProductHandler
	*OrderHandler
}

type Handler interface {
	CreateUser(ctx *gin.Context)
	LoginUser(ctx *gin.Context)
}

func RegisterHandlers(store db.Store, jwtToken *token.JWT) *AllHandler {
	return &AllHandler{
		UserHandler:    NewUserHandler(store, jwtToken),
		ProductHandler: NewProductHandler(store),
		OrderHandler:   NewOrderHandler(store),
	}
}
