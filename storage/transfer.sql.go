package storage

import (
	"context"
	"log"

	"bank/model"
)

type CreateTransferParams struct {
	FromAccountID int `json:"from_account_id"`
	ToAccountID   int `json:"to_account_id"`
	Amount        int `json:"amount"`
}

func (q *Queries) CreateTransfer(ctx context.Context, arg CreateTransferParams) (model.Transfer, error) {
	query := `insert into transfers (from_account_id,to_account_id,amount)
	VALUES($1,$2,$3)
	returning id,from_account_id,to_account_id,amount,created_at`
	transfer := model.Transfer{}
	if err := q.db.QueryRowContext(ctx, query, arg.FromAccountID, arg.ToAccountID, arg.Amount).Scan(
		&transfer.ID, &transfer.FromAccountID,
		&transfer.ToAccountID, &transfer.Amount,
		&transfer.CreatedAt,
	); err != nil {
		log.Printf("createTransfer err: %v", err)
		return model.Transfer{}, err
	}
	return transfer, nil
}

func (q *Queries) GetTransfer(ctx context.Context, transferId int) (model.Transfer, error) {
	query := `SELECT id,from_account_id,to_account_id,amount,created_at FROM transfers
	WHERE id = $1`
	transfer := model.Transfer{}
	if err := q.db.QueryRowContext(ctx, query, transferId).Scan(
		&transfer.ID, &transfer.FromAccountID,
		&transfer.ToAccountID, &transfer.Amount,
		&transfer.CreatedAt,
	); err != nil {
		log.Printf("getTransfer err: %v", err)
		return model.Transfer{}, err
	}
	return transfer, nil
}
