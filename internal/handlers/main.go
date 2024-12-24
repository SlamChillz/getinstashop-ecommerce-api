package handlers

import (
	"github.com/gin-gonic/gin"
	db "github.com/slamchillz/getinstashop-ecommerce-api/internal/db/sqlc"
)

type AllHandler struct {
	*UserHandler
}

type Handler interface {
	CreateUser(ctx *gin.Context)
}

func RegisterHandlers(store db.Store) *AllHandler {
	return &AllHandler{
		UserHandler: NewUserHandler(store),
	}
}
