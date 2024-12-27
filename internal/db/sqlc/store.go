package db

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Store defines all functions to execute db queries and transactions
type Store interface {
	Querier
	UpdateProductTx(ctx context.Context, arg UpdateProductTxParams) (Product, error, error)
	CreateOrderTx(ctx context.Context, arg CreateOrderTxParams) (Order, map[string]string, error, error)
	UpdateOrderTx(ctx context.Context, arg UpdateOrderTxParams) (Order, error)
}

// SQLStore provides all functions to execute SQL queries and transactions
type SQLStore struct {
	connPool *pgxpool.Pool
	*Queries
}

// NewStore creates a new store
func NewStore(connPool *pgxpool.Pool) Store {
	return &SQLStore{
		connPool: connPool,
		Queries:  New(connPool),
	}
}
