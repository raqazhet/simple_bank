package storage

import (
	"context"
	"log"

	"bank/model"
)

type CreateEntryParams struct {
	AccountID int `json:"account_id"`
	Amount    int `json:"amount"`
}

func (q *Queries) CreateEntry(ctx context.Context, arg CreateEntryParams) (model.Entry, error) {
	query := `INSERT INTO entries 
	(account_id, amount)
	Values($1,$2)
	Returning id,account_id,amount,created_at`
	args := []any{arg.AccountID, arg.Amount}
	entry := model.Entry{}
	if err := q.db.QueryRowContext(ctx, query, args...).Scan(&entry.ID, &entry.AccountID, &entry.Amount, &entry.CreatedAt); err != nil {
		log.Printf("createEntry err: %v", err)
		return model.Entry{}, err
	}
	return entry, nil
}

func (q *Queries) GetEntry(ctx context.Context, entryId int) (model.Entry, error) {
	query := `SELECT  id,account_id,amount,created_at FROM entries
	WHERE id = $1`
	entry := model.Entry{}
	if err := q.db.QueryRowContext(ctx, query, entryId).Scan(&entry.ID, &entry.AccountID, &entry.Amount, &entry.CreatedAt); err != nil {
		log.Printf("getEntry err: %v", err)
		return model.Entry{}, err
	}
	return entry, nil
}
