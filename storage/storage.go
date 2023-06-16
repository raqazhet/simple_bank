package storage

import (
	"context"
	"database/sql"
	"fmt"

	"bank/model"
)

type (
	Store interface {
		Querier
		TransferTx(ctx context.Context, arg TransferTxParams) (model.TransferTxResult, error)
	}
	SqlStorage struct {
		*Queries
		db *sql.DB
	}
)

func NewStorage(db *sql.DB) Store {
	return &SqlStorage{
		db:      db,
		Queries: New(db),
	}
}

func (r *SqlStorage) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	q := New(tx)
	err = fn(q)
	if err != nil {
		if rberr := tx.Rollback(); rberr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rberr)
		}
		return err
	}
	return tx.Commit()
}
