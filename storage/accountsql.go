package storage

import (
	"context"
	"log"

	"bank/model"
)

func (q *Queries) CreateAccount(ctx context.Context, arg model.CreateAccountParams) (model.Account, error) {
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
				LIMIT 1
				FOR NO KEY UPDATE`
	account := model.Account{}
	if err := q.db.QueryRowContext(ctx, query, accountId).Scan(&account.ID, &account.Owner, &account.Balance, &account.Currency, &account.CreatedAt); err != nil {
		log.Printf("getAccountByID storage err: %v", err)
		return model.Account{}, err
	}
	return account, nil
}

func (q Queries) GetAccountForUpdate(ctx context.Context, accountId int) (model.Account, error) {
	query := `SELECT id,owner,balance,currency,created_at from accounts
	WHERE id = $1 LIMIT 1
	FOR NO KEY UPDATE`
	account := model.Account{}
	if err := q.db.QueryRowContext(ctx, query, accountId).Scan(&account.ID, &account.Owner, &account.Balance, &account.Currency, &account.CreatedAt); err != nil {
		log.Printf("getAccountForUpdate storage err: %v", err)
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

type AddAccountBalance struct {
	Amount int `json:"amount"`
	Id     int `json:"id"`
}

func (q *Queries) AddAccountBalance(ctx context.Context, arg AddAccountBalance) (model.Account, error) {
	query := `Update accounts SET balance = balance + $1
	WHERE id = $2
	RETURNING id,owner,balance,currency,created_at`
	account := model.Account{}
	args := []any{arg.Amount, arg.Id}
	if err := q.db.QueryRowContext(ctx, query, args...).Scan(
		&account.ID, &account.Owner, &account.Balance,
		&account.Currency, &account.CreatedAt); err != nil {
		log.Printf("AddAccountBalance err: %v", err)
		return model.Account{}, err
	}
	return account, nil
}
