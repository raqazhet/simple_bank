package storage

import (
	"context"
	"fmt"
	"testing"

	"bank/model"

	"github.com/stretchr/testify/require"
)

type keystr string

const txkey = keystr("txkey")

func TestTransferTx(t *testing.T) {
	store := NewStorage(testDB)
	// Next we create 2 random accounts using the createRandomAccount
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)
	n := 2
	amount := 10
	// run n concurent transfer transaction
	errs := make(chan error)
	resullts := make(chan model.TransferTxResult)
	for i := 0; i < n; i++ {
		txNamee := fmt.Sprintf("tx %d", i+1)
		go func() {
			ctx := context.WithValue(context.Background(), txkey, txNamee)
			res, err := store.TransferTx(ctx, TransferTxParams{
				FromAccountID: account1.ID,
				ToAccountID:   account2.ID,
				Amount:        amount,
			})
			errs <- err
			resullts <- res
		}()
	}
	existed := make(map[int]bool)
	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)
		result := <-resullts
		require.NotEmpty(t, result)
		// Check transfer
		transfer := result.Transfer
		require.NotEmpty(t, transfer)
		require.Equal(t, account1.ID, transfer.FromAccountID)
		require.Equal(t, account2.ID, transfer.ToAccountID)
		require.Equal(t, amount, transfer.Amount)
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)
		_, err = store.GetTransfer(context.Background(), transfer.ID)
		require.NoError(t, err)
		// Check entries
		fromEntry := result.FromEntry
		require.NotEmpty(t, fromEntry)
		require.Equal(t, account1.ID, fromEntry.AccountID)
		require.Equal(t, -amount, fromEntry.Amount)
		require.NotZero(t, fromEntry.ID)
		require.NotZero(t, fromEntry.CreatedAt)
		_, err = store.GetEntry(context.Background(), fromEntry.ID)
		require.NoError(t, err)
		// check to entries
		toEntry := result.ToEntry
		require.NotEmpty(t, toEntry)
		require.Equal(t, account2.ID, toEntry.AccountID)
		require.Equal(t, amount, toEntry.Amount)
		require.NotZero(t, toEntry.ID)
		require.NotZero(t, toEntry.CreatedAt)
		_, err = store.GetEntry(context.Background(), toEntry.ID)
		require.NoError(t, err)
		// Check accounts
		fromAccount := result.FromAccount
		require.NotEmpty(t, fromAccount)
		require.Equal(t, account1.ID, fromAccount.ID)

		toaccount := result.ToAccount
		require.NotEmpty(t, toaccount)
		require.Equal(t, account2.ID, toaccount.ID)

		fmt.Println(">>tx: ", fromAccount.Balance, toaccount.Balance)
		// check account's balance
		diff1 := account1.Balance - fromAccount.Balance
		diff2 := toaccount.Balance - account2.Balance
		require.Equal(t, diff1, diff2)
		require.True(t, diff1 > 0)
		require.True(t, diff1%amount == 0) // 1*amount
		// Check results

		// Check accounts balance
		k := diff1 / amount
		require.True(t, k >= 1 && k <= n)
		require.NotContains(t, existed, k)
		existed[k] = true
	}
	// Check the final update balance
	updateAccount1, err := store.GetAccountById(context.Background(), account1.ID)
	require.NoError(t, err)
	updateAccount2, err := store.GetAccountById(context.Background(), account2.ID)
	require.NoError(t, err)
	fmt.Println(">>after:", updateAccount1.Balance, updateAccount2.Balance)
	require.Equal(t, account1.Balance-n*amount, updateAccount1.Balance)
	require.Equal(t, account2.Balance+n*amount, updateAccount2.Balance)
}
