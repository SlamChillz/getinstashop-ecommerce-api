package services

import (
	"context"
	"database/sql"
	"errors"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/slamchillz/getinstashop-ecommerce-api/internal/constants"
	db "github.com/slamchillz/getinstashop-ecommerce-api/internal/db/sqlc"
	"github.com/slamchillz/getinstashop-ecommerce-api/internal/types"
	"github.com/slamchillz/getinstashop-ecommerce-api/internal/validators"
	"log"
	"math"
	"net/http"
	"strings"
)

// ProductService provides business logic for product operations.
type ProductService struct {
	store db.Store
}

// NewProductService creates a new ProductService instance.
func NewProductService(store db.Store) *ProductService {
	return &ProductService{
		store: store,
	}
}

func (s *ProductService) CreateProduct(ctx context.Context, product types.CreateProductInput) (types.ProductOutput, types.ProductErrMessage, int, error) {
	errMessage, err := validators.ValidateProduct(product)
	if err != nil {
		return types.ProductOutput{}, errMessage, http.StatusBadRequest, err
	}
	userId, _ := ctx.Value(constants.ContextUserIdKey).(uuid.UUID)
	newProduct, err := s.store.CreateProduct(ctx, db.CreateProductParams{
		ID:          uuid.New(),
		Name:        product.Name,
		Description: product.Description,
		Price:       math.Round(product.Price*100) / 100,
		Stock:       int32(product.Stock),
		CreatedBy:   userId,
	})
	log.Print(newProduct, userId, err)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				errMessage.Name = "product already exists"
				return types.ProductOutput{}, errMessage, http.StatusBadRequest, err
			}
		}
		return types.ProductOutput{}, errMessage, http.StatusInternalServerError, err
	}
	return types.ProductOutput{
		ID:          newProduct.ID,
		Name:        newProduct.Name,
		Description: newProduct.Description,
		Price:       newProduct.Price,
		Stock:       newProduct.Stock,
		CreatedBy:   newProduct.CreatedBy,
		CreatedAt:   newProduct.CreatedAt,
		UpdatedAt:   newProduct.UpdatedAt,
	}, errMessage, http.StatusCreated, nil
}

func (s *ProductService) GetAllProduct(ctx context.Context) ([]types.ProductOutput, types.ProductErrMessage, int, error) {
	var errMessage types.ProductErrMessage
	allProduct, err := s.store.GetAllProduct(ctx)
	if err != nil {
		return nil, errMessage, http.StatusInternalServerError, err
	}
	// Convert GetAllProductRow slice to CreateProductOutput slice
	var allProductOutput []types.ProductOutput
	for _, product := range allProduct {
		productOutput := types.ProductOutput{
			ID:          product.ID,
			Name:        product.Name,
			Description: product.Description,
			Price:       product.Price,
			Stock:       product.Stock,
			CreatedAt:   product.CreatedAt,
			UpdatedAt:   product.UpdatedAt,
			CreatedBy:   product.CreatedBy,
		}
		allProductOutput = append(allProductOutput, productOutput)
	}
	// Return the converted slice along with the status code
	return allProductOutput, errMessage, http.StatusOK, nil
}

func (s *ProductService) GetOneProduct(ctx context.Context, productId uuid.UUID) (types.ProductOutput, types.ProductErrMessage, int, error) {
	var errMessage types.ProductErrMessage
	product, err := s.store.GetOneProduct(ctx, productId)
	if err != nil {
		if strings.Replace(sql.ErrNoRows.Error(), "sql: ", "", 1) == err.Error() {
			errMessage.ID = "product not found"
			return types.ProductOutput(product), errMessage, http.StatusNotFound, err
		}
		return types.ProductOutput(product), errMessage, http.StatusInternalServerError, err
	}
	return types.ProductOutput(product), errMessage, http.StatusOK, nil
}

func (s *ProductService) DeleteOneProduct(ctx context.Context, productId uuid.UUID) (types.ProductOutput, types.ProductErrMessage, int, error) {
	var product types.ProductOutput
	var errMessage types.ProductErrMessage
	err := s.store.DeleteOneProduct(ctx, productId)
	if err != nil {
		if strings.Replace(sql.ErrNoRows.Error(), "sql: ", "", 1) == err.Error() {
			errMessage.ID = "product not found"
			return product, errMessage, http.StatusNotFound, err
		}
		return product, errMessage, http.StatusInternalServerError, err
	}
	return product, errMessage, http.StatusNoContent, nil
}

func (s *ProductService) UpdateOneProduct(ctx context.Context, productId uuid.UUID, product types.ProductUpdateInput) (db.Product, types.ProductErrMessage, int, error) {
	var updatedProduct db.Product
	errMessage, err := validators.ValidateProductUpdateInput(product)
	if err != nil {
		return updatedProduct, errMessage, http.StatusBadRequest, err
	}
	updatedProduct, execErr, txErr := s.store.UpdateProductTx(ctx, db.UpdateProductTxParams{
		ID:          productId,
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		Stock:       product.Stock,
	})
	if execErr != nil || txErr != nil {
		if execErr != nil {
			if strings.Replace(sql.ErrNoRows.Error(), "sql: ", "", 1) == execErr.Error() {
				return updatedProduct, types.ProductErrMessage{
					ID: "product not found",
				}, http.StatusNotFound, err
			}
		}
		return updatedProduct, errMessage, http.StatusInternalServerError, err
	}
	return updatedProduct, errMessage, http.StatusOK, nil
}
