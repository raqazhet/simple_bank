package storage

import (
	"context"
	"database/sql"
	"fmt"
)

//	type SqlSrore interface {
//		Queries
//		TransferTx(ctx context.Context, arg TransferTxParams) (model.TransferTxResult, error)
//	}
type Storage struct {
	*Queries
	db *sql.DB
}

func NewStorage(db *sql.DB) *Storage {
	return &Storage{
		db:      db,
		Queries: New(db),
	}
}

type TransferTxParams struct {
	// is id of the account where money will be sent from
	FromAccountID int `json:"from_account_id"`
	// is the id of the account where money will be sent to
	ToAccountID int `json:"to_account_id"`
	// And the last field is the Amount of money to be sent
	Amount int `json:"amount"`
}

// The struct contains the result of the transfer transaction
// It has 5 fields

func (r *Storage) execTx(ctx context.Context, fn func(*Queries) error) error {
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
