package storage

import (
	"context"

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
		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: arg.FromAccountID, ToAccountID: arg.ToAccountID, Amount: arg.Amount,
		})
		if err != nil {
			return err
		}
		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount:    -arg.Amount,
		})
		if err != nil {
			return err
		}
		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount:    arg.Amount,
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
