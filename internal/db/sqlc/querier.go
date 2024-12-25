// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package db

import (
	"context"

	"github.com/google/uuid"
)

type Querier interface {
	CreateOrder(ctx context.Context, arg CreateOrderParams) (Order, error)
	CreateProduct(ctx context.Context, arg CreateProductParams) (Product, error)
	CreateUser(ctx context.Context, arg CreateUserParams) (User, error)
	DeleteOneProduct(ctx context.Context, id uuid.UUID) error
	GetAllOrderByUserId(ctx context.Context, userid uuid.UUID) (Order, error)
	GetAllOrderItem(ctx context.Context, orderid uuid.UUID) ([]OrderItem, error)
	GetAllProduct(ctx context.Context) ([]GetAllProductRow, error)
	GetMultipleProductById(ctx context.Context, dollar_1 []uuid.UUID) ([]GetMultipleProductByIdRow, error)
	GetOneProduct(ctx context.Context, id uuid.UUID) (GetOneProductRow, error)
	GetOrderById(ctx context.Context, id uuid.UUID) (Order, error)
	GetUserById(ctx context.Context, email string) (GetUserByIdRow, error)
	UpdateOneProduct(ctx context.Context, arg UpdateOneProductParams) (Product, error)
	UpdateProductStock(ctx context.Context, arg UpdateProductStockParams) (Product, error)
}

var _ Querier = (*Queries)(nil)
