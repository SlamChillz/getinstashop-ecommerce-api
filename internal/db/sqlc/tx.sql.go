package db

import (
	"context"
)

// ExecTx executes a function within a database transaction
func (store *SQLStore) execTx(ctx context.Context, fn func(*Queries) error) (execErr error, txErr error) {
	tx, err := store.connPool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			return err, rbErr
		}
		return err, nil
	}
	return tx.Commit(ctx), nil
}
