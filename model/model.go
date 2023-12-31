package model

import (
	"time"

	"github.com/google/uuid"
)

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
	Currency string `json:"currency" binding:"required,currency"`
}

// User struct

type User struct {
	Username          string    `json:"username" binding:"required,alphanum"`
	HashedPassword    string    `json:"hashed_password" binding:"required,min=8"`
	Fullname          string    `json:"full_name" binding:"required"`
	Email             string    `json:"email" binding:"required,email"`
	PasswordChangedAt time.Time `json:"password_changed_at"`
	CreatedAt         time.Time `json:"created_at"`
}

// Session params
type CreateSessionParams struct {
	ID           uuid.UUID `json:"id"`
	Username     string    `json:"username"`
	RefreshToken string    `json:"refresh_token"`
	UserAgent    string    `json:"user_agent"`
	ClientIp     string    `json:"client_ip"`
	IsBlocked    bool      `json:"is_blocked"`
	ExpiresAt    time.Time `json:"expires_at"`
	CreatedAt    time.Time `json:"created_at"`
}
