package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	db "github.com/slamchillz/getinstashop-ecommerce-api/internal/db/sqlc"
	"github.com/slamchillz/getinstashop-ecommerce-api/internal/services"
	"github.com/slamchillz/getinstashop-ecommerce-api/internal/types"
	"github.com/slamchillz/getinstashop-ecommerce-api/internal/utils"
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

// CreateProduct godoc
// @Summary      Create a new product. Requires admin privilege
// @Description  Create a new product. Requires admin privilege
// @Tags         product
// @Accept       json
// @Produce      json
// @Param        payload   body	types.CreateProductInput  true  "Create Product request body"
// @Success      200  {object}  types.Product
// @Failure      400  {object}  types.ProductError
// @Failure      500  {object}  types.InterServerError
// @Security	 BearerAuth
// @Router       /admin/products [post]
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

// GetAllProduct godoc
// @Summary      List all products. None admin users should be able to see products before placing an order.
// @Description  List all products. None admin users should be able to see products before placing an order.
// @Tags         product
// @Accept       json
// @Produce      json
// @Success      200  {array}  types.Product
// @Failure      400  {object}  types.ProductError
// @Failure      500  {object}  types.InterServerError
// @Security	 BearerAuth
// @Router       /products [get]
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

// GetOneProduct godoc
// @Summary      Fetch One Product. Requires admin privilege
// @Description  Fetch One Product. Requires admin privilege
// @Tags         product
// @Accept       json
// @Produce      json
// @Param        productId   path	string  true  "Unique product id"
// @Success      200  {object}  types.Product
// @Failure      400  {object}  types.ProductError
// @Failure      500  {object}  types.InterServerError
// @Security	 BearerAuth
// @Router       /admin/products/{productId} [get]
func (h *ProductHandler) GetOneProduct(ctx *gin.Context) {
	var err error
	var productId uuid.UUID = utils.ParseStringToUUID(ctx.Param("id"))
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

// DeleteOneProduct godoc
// @Summary      Delete One Product. Requires admin privilege
// @Description  Delete One Product. Requires admin privilege
// @Tags         product
// @Accept       json
// @Produce      json
// @Param        productId   path	string  true  "Unique product id"
// @Success      204
// @Failure      400  {object}  types.ProductError
// @Failure      500  {object}  types.InterServerError
// @Security	 BearerAuth
// @Router       /admin/products/{productId} [delete]
func (h *ProductHandler) DeleteOneProduct(ctx *gin.Context) {
	var err error
	var productId uuid.UUID = utils.ParseStringToUUID(ctx.Param("id"))
	_, errMessage, statusCode, err := h.productService.DeleteOneProduct(ctx, productId)
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
		"data":    gin.H{},
	})
}

// UpdateOneProduct godoc
// @Summary      Update a single Product. Requires admin privilege
// @Description  Update a single Product. Requires admin privilege
// @Tags         product
// @Accept       json
// @Produce      json
// @Param        productId   path	string  true  "Unique product id"
// @Param        payload   	 body	types.CreateProductInput  true  "Update Product request body"
// @Success      200  {object}	types.Product
// @Failure      400  {object}  types.ProductError
// @Failure      500  {object}  types.InterServerError
// @Security	 BearerAuth
// @Router       /admin/products/{productId} [put]
func (h *ProductHandler) UpdateOneProduct(ctx *gin.Context) {
	var err error
	var req types.ProductUpdateInput
	var productId uuid.UUID = utils.ParseStringToUUID(ctx.Param("id"))
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
	ctx.JSON(statusCode, gin.H{
		"status":  "success",
		"message": "Product updated",
		"data":    response,
	})
}
