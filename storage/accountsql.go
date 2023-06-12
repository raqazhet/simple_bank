package storage

import (
	"context"
	"log"

	"bank/model"
)

type CreateAccountParams struct {
	Owner    string `json:"owner"`
	Balance  int    `json:"balance"`
	Currency string `json:"currency"`
}

func (q *Queries) CreateAccount(ctx context.Context, arg CreateAccountParams) (model.Account, error) {
	query := `INSERT INTO accounts (owner,balance,currency)
	VALUES($1,$2,$3)
	RETURNING id,owner,balance,currency,created_at`
	account := model.Account{}
	if err := q.db.QueryRowContext(ctx, query, arg.Owner, arg.Balance, arg.Currency).Scan(&account.ID, &account.Owner, &account.Balance, &account.Currency, &account.CreatedAt); err != nil {
		log.Printf("createAccount storage err: %v", err)
		return model.Account{}, err
	}
	return account, nil
}

func (q *Queries) GetAccountById(ctx context.Context, accountId int) (model.Account, error) {
	query := `SELECT * FROM accounts
				where id = $1
				LIMIT 1`
	account := model.Account{}
	if err := q.db.QueryRowContext(ctx, query, accountId).Scan(&account.ID, &account.Owner, &account.Balance, &account.Currency, &account.CreatedAt); err != nil {
		log.Printf("createAccount storage err: %v", err)
		return model.Account{}, err
	}
	return account, nil
}

func (q *Queries) GetAllAccounts(ctx context.Context) ([]model.Account, error) {
	query := `SELECT * FROM accounts 
	ORDER BY id
	LIMIT 5
	OFFSET 5`
	accounts := []model.Account{}
	rows, err := q.db.QueryContext(ctx, query)
	if err != nil {
		log.Printf("GetAllAccounts err1: %v", err)
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		account := model.Account{}
		if err := rows.Scan(&account.ID, &account.Owner, &account.Balance, &account.Currency, &account.CreatedAt); err != nil {
			log.Printf("getAllAccounts scan.err: %v", err)
			return nil, err
		}
		accounts = append(accounts, account)
	}
	if err := rows.Err(); err != nil {
		log.Printf("getAllAccounts err: %v", err)
		return nil, err
	}
	return accounts, nil
}

type UpdateAccountParams struct {
	ID      int `json:"id"`
	Balance int `json:"balance"`
}

func (q *Queries) UpdateAccount(ctx context.Context, arg UpdateAccountParams) (model.Account, error) {
	query := `Update accounts
	SET balance = $1
	WHERE id =$2
	RETURNING id,owner,balance,currency,created_at`
	account := model.Account{}
	if err := q.db.QueryRowContext(ctx, query, arg.Balance, arg.ID).Scan(&account.ID, &account.Owner, &account.Balance, &account.Currency, &account.CreatedAt); err != nil {
		log.Printf("updateAccount storage err: %v", err)
		return model.Account{}, err
	}
	return account, nil
}

func (q *Queries) DeleteAccount(ctx context.Context, accountId int) error {
	query := `DELETE FROM accounts 
	WHERE id =$1`
	_, err := q.db.ExecContext(ctx, query, accountId)
	if err != nil {
		log.Printf("DeleteAccount err: %v", err)
		return err
	}
	return nil
}
