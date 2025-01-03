package db

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"math"
	"strings"
)

type CreateOrderTxParams struct {
	ID         uuid.UUID           `json:"id"`
	UserId     uuid.UUID           `json:"userId"`
	ProductIds []uuid.UUID         `json:"productIds"`
	Items      map[uuid.UUID]int32 `json:"items"`
}

func (store *SQLStore) CreateOrderTx(ctx context.Context, arg CreateOrderTxParams) (Order, map[string]string, error, error) {
	var order Order
	var invalidProducts = make(map[string]string)
	var orderTotal float64
	var prices = make(map[uuid.UUID]float64)
	var orderQuantities = make(map[uuid.UUID]int32)
	var values []interface{}
	var placeholders []string
	products, err := store.GetMultipleProductById(ctx, arg.ProductIds)
	if err != nil {
		return order, invalidProducts, err, nil
	}
	for i, product := range products {
		prices[product.ID] = product.Price
		quantity, ok := arg.Items[product.ID]
		if !ok {
			invalidProducts[product.ID.String()] = "product not found"
		}
		if quantity > product.Stock {
			invalidProducts[product.ID.String()] = "quantity less than available stock"
		}
		orderQuantities[product.ID] = quantity
		productPrice := math.Round((product.Price*float64(quantity))*100) / 100
		// Create a group of placeholders for each record
		placeholders = append(placeholders, fmt.Sprintf("($%d, $%d, $%d, $%d, $%d)",
			i*5+1, i*5+2, i*5+3, i*5+4, i*5+5))
		values = append(values, uuid.New(), arg.ID, product.ID, quantity, productPrice)
		orderTotal += productPrice
	}
	if len(invalidProducts) > 0 {
		return order, invalidProducts, nil, nil
	}
	// Join placeholders with commas and append to the query
	query := fmt.Sprint(`INSERT`, ` INTO`, ` "orderItem"`, ` ("id", "orderId", "productId", "quantity", "price")`, ` VALUES `, strings.Join(placeholders, ", "))
	execErr, txErr := store.execTx(ctx, func(q *Queries) error {
		order, err = q.CreateOrder(ctx, CreateOrderParams{
			ID:     arg.ID,
			UserId: arg.UserId,
			Total:  orderTotal,
		})
		if err != nil {
			return err
		}
		_, err = q.db.Exec(ctx, query, values...)
		if err != nil {
			return err
		}
		for productId, quantity := range orderQuantities {
			_, err = q.UpdateProductStock(ctx, UpdateProductStockParams{
				ID:    productId,
				Stock: quantity,
			})
			if err != nil {
				return err
			}
		}
		return nil
	})
	return order, invalidProducts, execErr, txErr
}

type UpdateOrderTxParams struct {
	ID     uuid.UUID   `json:"id"`
	UserId uuid.UUID   `json:"userId"`
	Admin  bool        `json:"admin"`
	Status OrderStatus `json:"status"`
}

func (store *SQLStore) UpdateOrderTx(ctx context.Context, arg UpdateOrderTxParams) (Order, error) {
	if arg.Status != OrderStatusCANCELLED {
		return store.UpdateOrderStatus(ctx, UpdateOrderStatusParams{
			ID:     arg.ID,
			Status: arg.Status,
		})
	}
	products, err := store.GetAllProductInOrder(ctx, arg.ID)
	if err != nil {
		return Order{}, err
	}
	var order Order
	execErr, _ := store.execTx(ctx, func(q *Queries) error {
		for _, product := range products {
			stock := product.Quantity * -1
			if arg.Status == OrderStatusPENDING {
				stock = product.Quantity
			}
			_, err := q.UpdateProductStock(ctx, UpdateProductStockParams{
				ID:    product.ProductId,
				Stock: stock,
			})
			if err != nil {
				return err
			}
		}
		if !arg.Admin {
			order, err = q.CancelOrder(ctx, CancelOrderParams{
				ID:     arg.ID,
				UserId: arg.UserId,
			})
			if err != nil {
				return err
			}
		} else {
			order, err = q.UpdateOrderStatus(ctx, UpdateOrderStatusParams{
				ID:     arg.ID,
				Status: arg.Status,
			})
			if err != nil {
				return err
			}
		}
		return nil
	})
	return order, execErr
}
