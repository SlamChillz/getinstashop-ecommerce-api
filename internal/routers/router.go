package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/slamchillz/getinstashop-ecommerce-api/internal/handlers"
	"github.com/slamchillz/getinstashop-ecommerce-api/internal/middlewares"
	"github.com/slamchillz/getinstashop-ecommerce-api/pkg/token"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func InitRouters(handler *handlers.AllHandler, token *token.JWT) *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// Health check
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.GET("/health", handlers.Health)
	v1 := router.Group("/api/v1")
	{
		auth := v1.Group("/auth")
		{
			auth.POST("/register", handler.UserHandler.CreateUser)
			auth.POST("/login", handler.UserHandler.LoginUser)
		}
		v1.Use(middlewares.AuthMiddy(token))
		// Orders routes
		orders := v1.Group("/orders")
		{
			orders.POST("", handler.CreateOrder)
			orders.GET("", handler.GetUserOrders)
			orders.PATCH("/:id", handler.CancelOrder)
		}
		// Admin routes
		admin := v1.Group("/admin")
		{
			admin.Use(middlewares.AdminMiddy)
			admin.POST("/products", handler.CreateProduct)
			admin.GET("/products", handler.GetAllProduct)
			admin.GET("/products/:id", handler.GetOneProduct)
			admin.DELETE("/products/:id", handler.DeleteOneProduct)
			admin.PUT("/products/:id", handler.UpdateOneProduct)
			admin.PATCH("/orders/:id", handler.OrderHandler.UpdateOrderStatus)
		}
	}
	//v1.GET("/docs", ginSwagger.WrapHandler(swaggerFiles.Handler))
	return router
}
