package db

import (
	"context"
	"github.com/google/uuid"
)

type UpdateProductTxParams struct {
	ID          uuid.UUID `json:"id"`
	Name        *string   `json:"name,omitempty"`
	Description *string   `json:"description,omitempty"`
	Price       *float64  `json:"price,omitempty"`
	Stock       *int32    `json:"stock,omitempty"`
}

type UpdateProductTxResult Product

func (store *SQLStore) UpdateProductTx(ctx context.Context, arg UpdateProductTxParams) (Product, error, error) {
	var result Product
	execErr, txErr := store.execTx(ctx, func(q *Queries) error {
		var err error
		product, err := q.GetOneProduct(ctx, arg.ID)
		if err != nil {
			return err
		}
		if arg.Name == nil {
			arg.Name = &product.Name
		}
		if arg.Description == nil {
			arg.Description = &product.Description
		}
		if arg.Price == nil {
			arg.Price = &product.Price
		}
		if arg.Stock == nil {
			arg.Stock = &product.Stock
		}
		updatedProduct, err := q.UpdateOneProduct(ctx, UpdateOneProductParams{
			ID:          arg.ID,
			Name:        *arg.Name,
			Description: *arg.Description,
			Price:       *arg.Price,
			Stock:       *arg.Stock,
		})
		if err != nil {
			return err
		}
		result = updatedProduct
		return nil
	})
	return result, execErr, txErr
}
