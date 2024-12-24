package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/slamchillz/getinstashop-ecommerce-api/internal/handlers"
)

func InitRouters(handler *handlers.AllHandler) *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// Health check
	router.GET("/health", handlers.Health)

	// Register authentication routes
	router.POST("/auth/register", handler.UserHandler.CreateUser)

	return router
}
