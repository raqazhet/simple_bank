package storage

import (
	"context"
	"testing"
	"time"

	"bank/model"

	"github.com/stretchr/testify/require"
)

func TestCreateAccount(t *testing.T) {
	createRandomAccount(t)
}

func createRandomAccount(t *testing.T) model.Account {
	arg := CreateAccountParams{
		Owner:    "razaq",
		Balance:  200,
		Currency: "USD",
	}
	account, err := testQueries.CreateAccount(context.Background(), arg)
	// Basically,this command will check that the error must be nil and will automatically fail the test if
	// it's not
	require.NoError(t, err)
	// Next,we require that the returned account should not be empty object using
	// require.NotEmpty() function
	require.NotEmpty(t, account)
	require.Equal(t, arg.Owner, account.Owner)
	require.Equal(t, arg.Balance, account.Balance)
	require.Equal(t, arg.Currency, account.Currency)
	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)
	return account
}

func TestGetAccount(t *testing.T) {
	account1 := createRandomAccount(t)
	account2, err := testQueries.GetAccountById(context.Background(), account1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, account2)
	require.Equal(t, account1.ID, account2.ID)
	require.Equal(t, account1.Owner, account2.Owner)
	require.Equal(t, account1.Balance, account2.Balance)
	require.Equal(t, account1.Currency, account2.Currency)
	require.WithinDuration(t, account1.CreatedAt, account2.CreatedAt, time.Second)
}

func TestUpdateAccount(t *testing.T) {
	account1 := createRandomAccount(t)
	arg := UpdateAccountParams{
		ID:      account1.ID,
		Balance: 400,
	}
	account2, err := testQueries.UpdateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, account2)
	require.Equal(t, account1.ID, account2.ID)
	require.Equal(t, account1.Owner, account2.Owner)
	require.Equal(t, account1.Currency, account2.Currency)
	require.WithinDuration(t, account1.CreatedAt, account2.CreatedAt, time.Second)
}

func TestDeleteAccount(t *testing.T) {
	account1 := createRandomAccount(t)
	err := testQueries.DeleteAccount(context.Background(), account1.ID)
	require.NoError(t, err)
}

func TestAllListsAccounts(t *testing.T) {
	for i := 0; i < 5; i++ {
		createRandomAccount(t)
	}
	arg, err := testQueries.GetAllAccounts(context.Background())
	require.NoError(t, err)
	require.Len(t, arg, 5)
	for _, v := range arg {
		require.NotEmpty(t, v)
	}
}
