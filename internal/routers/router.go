package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/slamchillz/getinstashop-ecommerce-api/internal/handlers"
	"github.com/slamchillz/getinstashop-ecommerce-api/internal/middlewares"
	"github.com/slamchillz/getinstashop-ecommerce-api/pkg/token"
)

func InitRouters(handler *handlers.AllHandler, token *token.JWT) *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// Health check
	router.GET("/health", handlers.Health)

	// Register authentication routes
	router.POST("/auth/register", handler.UserHandler.CreateUser)
	router.POST("/auth/login", handler.UserHandler.LoginUser)

	// Authenticated Endpoints
	apiV1Router := router.Group("/api/v1/").Use(middlewares.AuthMiddy(token))

	// Authenticated User Endpoints
	apiV1Router.POST("/orders", handler.OrderHandler.CreateOrder)
	apiV1Router.GET("/orders", handler.OrderHandler.GetUserOrders)
	apiV1Router.PATCH("/orders/:id", handler.OrderHandler.CancelOrder)

	// Admin routes
	adminRouter := apiV1Router.Use(middlewares.AdminMiddy)
	// Create a product
	adminRouter.POST("/admin/products", handler.ProductHandler.CreateProduct)
	// Fetch all Product
	adminRouter.GET("/admin/products", handler.ProductHandler.GetAllProduct)
	// Fetch a single Product
	adminRouter.GET("/admin/products/:id", handler.ProductHandler.GetOneProduct)
	// Delete a single Product
	adminRouter.DELETE("/admin/products/:id", handler.ProductHandler.DeleteOneProduct)
	// Update a Product
	adminRouter.PUT("/admin/products/:id", handler.ProductHandler.UpdateOneProduct)
	// Update Order status
	adminRouter.PATCH("/admin/orders/:id", handler.OrderHandler.UpdateOrderStatus)

	return router
}
