package services

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/slamchillz/getinstashop-ecommerce-api/internal/constants"
	db "github.com/slamchillz/getinstashop-ecommerce-api/internal/db/sqlc"
	"github.com/slamchillz/getinstashop-ecommerce-api/internal/types"
	"github.com/slamchillz/getinstashop-ecommerce-api/internal/validators"
	"math"
	"net/http"
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

func (s *ProductService) CreateProduct(ctx context.Context, product types.CreateProductInput) (*types.CreateProductOutput, *types.CreateProductErrMessage, int, error) {
	errMessage, err := validators.ValidateProduct(product)
	if err != nil {
		return nil, &errMessage, http.StatusBadRequest, err
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
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				return nil, &types.CreateProductErrMessage{
					Name: "product already exists",
				}, http.StatusBadRequest, err
			}
		}
		return nil, nil, http.StatusInternalServerError, err
	}
	return &types.CreateProductOutput{
		ID:          newProduct.ID,
		Name:        newProduct.Name,
		Description: newProduct.Description,
		Price:       newProduct.Price,
		Stock:       newProduct.Stock,
		CreatedBy:   newProduct.CreatedBy,
	}, nil, http.StatusCreated, nil
}
