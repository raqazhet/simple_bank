package model

import "time"

type Account struct {
	ID        int       `json:"id"`
	Owner     string    `json:"owner"`
	Balance   int       `json:"balance"`
	Currency  string    `json:"currency"`
	CreatedAt time.Time `json:"created_at"`
}

type Entry struct {
	ID        int `json:"id"`
	AccountID int `json:"account_id"`
	// can be negative or positive
	Amount    int       `json:"amount"`
	CreatedAt time.Time `json:"created_at"`
}

type Transfer struct {
	ID            int `json:"id"`
	FromAccountID int `json:"from_account_id"`
	ToAccountID   int `json:"to_account_id"`
	// must be positive
	Amount    int       `json:"amount"`
	CreatedAt time.Time `json:"created_at"`
}
type TransferTxResult struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
}

// Account go
type CreateAccountParams struct {
	Owner    string `json:"owner" binding:"required"`
	Balance  int    `json:"balance"`
	Currency string `json:"currency" binding:"required"`
}
