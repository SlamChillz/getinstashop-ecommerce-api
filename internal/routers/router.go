package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/slamchillz/getinstashop-ecommerce-api/internal/handlers"
)

func InitRouters() *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// Health check
	router.GET("/health", handlers.Health)
	return router
}
