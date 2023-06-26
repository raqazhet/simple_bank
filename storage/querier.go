package storage

import (
	"context"

	"bank/model"
)

type Querier interface {
	AddAccountBalance(ctx context.Context, arg AddAccountBalance) (model.Account, error)
	CreateUser(ctx context.Context, arg model.User) (model.User, error)
	GetUser(ctx context.Context, username string) (model.User, error)
	CreateAccount(ctx context.Context, arg model.CreateAccountParams) (model.Account, error)
	CreateEntry(ctx context.Context, arg CreateEntryParams) (model.Entry, error)
	CreateTransfer(ctx context.Context, arg CreateTransferParams) (model.Transfer, error)
	DeleteAccount(ctx context.Context, id int) error
	GetAccountById(ctx context.Context, id int) (model.Account, error)
	GetAccountForUpdate(ctx context.Context, id int) (model.Account, error)
	GetEntry(ctx context.Context, id int) (model.Entry, error)
	GetTransfer(ctx context.Context, id int) (model.Transfer, error)
	GetAllAccounts(ctx context.Context) ([]model.Account, error)
	SaveNewRefreshToken(ctx context.Context, arg model.CreateSessionParams) (model.CreateSessionParams, error)
	// ListEntries(ctx context.Context, arg model.ListEntriesParams) ([]Entry, error)
	// ListTransfers(ctx context.Context, arg ListTransfersParams) ([]Transfer, error)
	UpdateAccount(ctx context.Context, arg UpdateAccountParams) (model.Account, error)
}

var _ Querier = (*Queries)(nil)
