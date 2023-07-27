package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"bank-service/util"

	"github.com/stretchr/testify/require"
)

func generateRandomAccount(t *testing.T) Account {
	args := CreateAccountParams{
		Owner:    util.RandomOwner(),
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}

	account, err := testQuery.CreateAccount(context.Background(), args)

	require.NoError(t, err)
	require.NotEmpty(t, args)
	require.Equal(t, args.Owner, account.Owner)
	require.Equal(t, args.Balance, account.Balance)
	require.Equal(t, args.Currency, account.Currency)
	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)

	return account
}

func TestCreateAccount(t *testing.T) {
	generateRandomAccount(t)
}

func TestGetAccount(t *testing.T) {
	account1 := generateRandomAccount(t)
	account2, err := testQuery.GetAccount(context.Background(), account1.ID)

	require.NoError(t, err)
	require.NotEmpty(t, account2)

	require.Equal(t, account1.Owner, account2.Owner)
	require.Equal(t, account1.Balance, account2.Balance)
	require.Equal(t, account1.Currency, account2.Currency)
	require.WithinDuration(t, account1.CreatedAt, account2.CreatedAt, time.Second)
}

func TestUpdateAccount(t *testing.T) {
	account := generateRandomAccount(t)
	args := UpdateAccountParams{
		ID:      account.ID,
		Balance: util.RandomMoney(),
	}

	_, err := testQuery.UpdateAccount(context.Background(), args)
	require.NoError(t, err)
}

func TestDeleteAccount(t *testing.T) {
	account := generateRandomAccount(t)
	err := testQuery.DeleteAccount(context.Background(), account.ID)
	require.NoError(t, err)

	deletedAccount, err := testQuery.GetAccount(context.Background(), account.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, deletedAccount)
}

func TestListAccounts(t *testing.T) {
	for i := 0; i < 10; i++ {
		generateRandomAccount(t)
	}

	args := GetListAccountsParams{
		Limit:  5,
		Offset: 5,
	}

	accounts, err := testQuery.GetListAccounts(context.Background(), args)
	require.NoError(t, err)
	require.Len(t, accounts, 5)

	for _, account := range accounts {
		require.NotEmpty(t, account)
	}
}
