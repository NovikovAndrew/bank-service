package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	store := NewStore(testDB)

	account1 := generateRandomAccount(t)
	account2 := generateRandomAccount(t)

	// run n concurrent transactions
	n := 5
	amount := 10
	resChan := make(chan TransferTxResult)
	errChan := make(chan error)

	fmt.Println(resChan)

	for i := 0; i < n; i++ {
		go func() {
			params := TransferTxParams{
				FromAccountID: account1.ID,
				ToAccountID:   account2.ID,
				Amount:        int64(amount),
			}
			res, err := store.TransferTx(context.Background(), params)

			resChan <- res
			errChan <- err
		}()
	}

	for i := 0; i < n; i++ {
		err := <-errChan
		require.NoError(t, err)

		result := <-resChan
		require.NotEmpty(t, result)

		// check transfer
		transfer := result.Transfer
		require.NotEmpty(t, transfer)
		require.Equal(t, transfer.FromAccountID, account1.ID)
		require.Equal(t, transfer.ToAccountID, account2.ID)
		require.Equal(t, transfer.Amount, amount)
		require.NotEqual(t, transfer.FromAccountID, transfer.ToAccountID)
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)

		_, err = store.GetTransfer(context.Background(), transfer.ID)
		require.NoError(t, err)

		// check entries
		fromEntry := result.FromEntry
		require.NotEmpty(t, fromEntry)
		require.Equal(t, fromEntry.AccountID, account1)
		require.Equal(t, fromEntry.Amount, -amount)
		require.NotEmpty(t, fromEntry.ID)
		require.NotEmpty(t, fromEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), fromEntry.ID)
		require.NoError(t, err)

		toEntry := result.ToEntry
		require.NotEmpty(t, toEntry)
		require.Equal(t, toEntry.AccountID, account2)
		require.Equal(t, toEntry.Amount, amount)
		require.NotEmpty(t, toEntry.ID)
		require.NotEmpty(t, toEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), toEntry.ID)
		require.NoError(t, err)
	}
}
