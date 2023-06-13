package storage

import (
	"context"
	"fmt"

	"bank/model"
)

type TransferTxParams struct {
	// is id of the account where money will be sent from
	FromAccountID int `json:"from_account_id"`
	// is the id of the account where money will be sent to
	ToAccountID int `json:"to_account_id"`
	// And the last field is the Amount of money to be sent
	Amount int `json:"amount"`
}

func (r *Storage) TransferTx(ctx context.Context, arg TransferTxParams) (model.TransferTxResult, error) {
	var result model.TransferTxResult
	err := r.execTx(ctx, func(q *Queries) error {
		var err error
		txName := ctx.Value("txkey")
		fmt.Println("create transfer: ", txName)
		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: arg.FromAccountID, ToAccountID: arg.ToAccountID, Amount: arg.Amount,
		})
		if err != nil {
			return err
		}
		fmt.Println(txName, "create entry 1")
		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount:    -arg.Amount,
		})
		if err != nil {
			return err
		}
		fmt.Println(txName, "create entry 2")
		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount:    arg.Amount,
		})
		if err != nil {
			return err
		}
		if arg.FromAccountID < arg.ToAccountID {
			fmt.Println(txName, "update account 1")
			result.FromAccount, result.ToAccount, err = addMoney(ctx, q, arg.FromAccountID, arg.ToAccountID, -arg.Amount, arg.Amount)
			if err != nil {
				return err
			}

		} else {
			fmt.Println(txName, "update account 2")
			result.ToAccount, result.FromAccount, err = addMoney(ctx, q, arg.FromAccountID, arg.ToAccountID, -arg.Amount, -arg.Amount)
			if err != nil {
				return err
			}
		}

		// TODO: update accoun'ts balance
		return err
	})
	if err != nil {
		return model.TransferTxResult{}, err
	}
	return result, nil
}

func addMoney(ctx context.Context,
	q *Queries, accountid1, accountid2,
	amount1, amount2 int,
) (account1, account2 model.Account, err error) {
	account1, err = q.AddAccountBalance(ctx, AddAccountBalance{
		Id:     accountid1,
		Amount: amount1,
	})
	if err != nil {
		return
	}
	account2, err = q.AddAccountBalance(ctx, AddAccountBalance{
		Id:     accountid2,
		Amount: amount2,
	})
	if err != nil {
		return
	}
	return
}
