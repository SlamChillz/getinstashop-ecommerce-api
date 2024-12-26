package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	db "github.com/slamchillz/getinstashop-ecommerce-api/internal/db/sqlc"
	"github.com/slamchillz/getinstashop-ecommerce-api/internal/services"
	"github.com/slamchillz/getinstashop-ecommerce-api/internal/types"
	"log"
	"net/http"
)

// ProductHandler handles user-related operations.
type ProductHandler struct {
	productService *services.ProductService
}

// NewProductHandler creates a new UserHandler instance.
func NewProductHandler(store db.Store) *ProductHandler {
	return &ProductHandler{productService: services.NewProductService(store)}
}

func (h *ProductHandler) CreateProduct(ctx *gin.Context) {
	var err error
	var req types.CreateProductInput
	if err = ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "failed",
			"message": "Invalid JSON payload",
		})
		return
	}
	response, errMessage, statusCode, err := h.productService.CreateProduct(ctx, req)
	if err != nil {
		ctx.JSON(statusCode, gin.H{
			"status":  "failed",
			"message": "Product not created",
			"error":   errMessage,
		})
		log.Printf("Error while creating product: %v", err)
		return
	}
	ctx.JSON(statusCode, gin.H{
		"status":  "success",
		"message": "Product created",
		"data":    response,
	})
}

func (h *ProductHandler) GetAllProduct(ctx *gin.Context) {
	var err error
	response, errMessage, statusCode, err := h.productService.GetAllProduct(ctx)
	if err != nil {
		ctx.JSON(statusCode, gin.H{
			"status":  "failed",
			"message": "Unable to fetch products",
			"error":   errMessage,
		})
		log.Printf("Error while fetching product: %v", err)
		return
	}
	ctx.JSON(statusCode, gin.H{
		"status":  "success",
		"message": "Products retrieved",
		"data":    response,
	})
}

func (h *ProductHandler) GetOneProduct(ctx *gin.Context) {
	var err error
	var productId uuid.UUID = uuid.UUID([]byte(ctx.Param("id")))
	response, errMessage, statusCode, err := h.productService.GetOneProduct(ctx, productId)
	if err != nil {
		ctx.JSON(statusCode, gin.H{
			"status":  "failed",
			"message": "Unable to fetch product",
			"error":   errMessage,
		})
		log.Printf("Error while fetching product: %v", err)
		return
	}
	ctx.JSON(statusCode, gin.H{
		"status":  "success",
		"message": "Product retrieved",
		"data":    response,
	})
}

func (h *ProductHandler) DeleteOneProduct(ctx *gin.Context) {
	var err error
	var productId uuid.UUID = uuid.UUID([]byte(ctx.Param("id")))
	response, errMessage, statusCode, err := h.productService.DeleteOneProduct(ctx, productId)
	if err != nil {
		ctx.JSON(statusCode, gin.H{
			"status":  "failed",
			"message": "Unable to delete product",
			"error":   errMessage,
		})
		log.Printf("Error while deleting product: %v", err)
		return
	}
	ctx.JSON(statusCode, gin.H{
		"status":  "success",
		"message": "Product deleted",
		"data":    response,
	})
}

func (h *ProductHandler) UpdateOneProduct(ctx *gin.Context) {
	var err error
	var req types.ProductUpdateInput
	var productId uuid.UUID = uuid.UUID([]byte(ctx.Param("id")))
	if err = ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "failed",
			"message": "Invalid JSON payload",
		})
		return
	}
	response, errMessage, statusCode, err := h.productService.UpdateOneProduct(ctx, productId, req)
	if err != nil {
		ctx.JSON(statusCode, gin.H{
			"status":  "failed",
			"message": "Product not updated",
			"error":   errMessage,
		})
		log.Printf("Error while updating product: %v", err)
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "Product updated",
		"data":    response,
	})
}
