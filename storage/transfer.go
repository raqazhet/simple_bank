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
		// move money out of account1
		fmt.Println(txName, "get account 1")
		account1, err := q.GetAccountForUpdate(ctx, arg.FromAccountID)
		if err != nil {
			return err
		}
		fmt.Println(txName, "update account 1")
		result.FromAccount, err = q.UpdateAccount(ctx, UpdateAccountParams{
			ID:      arg.FromAccountID,
			Balance: account1.Balance - arg.Amount,
		})
		if err != nil {
			return err
		}
		// Move money out of account2
		fmt.Println(txName, "get account 2")
		account2, err := q.GetAccountForUpdate(ctx, arg.ToAccountID)
		if err != nil {
			return err
		}
		fmt.Println(txName, "update account 2")
		result.ToAccount, err = q.UpdateAccount(ctx, UpdateAccountParams{
			ID:      arg.ToAccountID,
			Balance: arg.Amount + account2.Balance,
		})
		if err != nil {
			return err
		}
		// TODO: update accoun'ts balance
		return nil
	})
	// fmt.Println("result.transfer", result.Transfer)
	// fmt.Println("res.FromEnt", result.FromEntry)
	// fmt.Println("res.ToEnt", result.ToEntry)
	// fmt.Println(result)
	if err != nil {
		return model.TransferTxResult{}, err
	}
	return result, nil
}
