package db

import (
	"context"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestTransferTx(t *testing.T) {
	store := NewStore(testDB)

	account1 := RandomAccount(t)
	account2 := RandomAccount(t)

	n := 5

	amount := int64(10)

	errors := make(chan error)
	results := make(chan TransferTxResult)

	for i := 0; i < n; i++ {
		go func() {
			err, result := store.TransferTx(context.Background(), TransferTxParams{
				FromAccountID: account1.ID,
				ToAccountID:   account2.ID,
				Amount:        amount,
			})
			errors <- err
			results <- result
		}()
	}

	exsits := make(map[int]bool)
	for i := 0; i < n; i++ {
		err := <-errors
		require.NoError(t, err)

		result := <-results
		require.NotEmpty(t, result)
		require.NotZero(t, result.Transfer.ID)
		require.NotZero(t, result.Transfer.CreatedAt)
		require.Equal(t, result.Transfer.FromAccountID, account1.ID)
		require.Equal(t, result.Transfer.ToAccountID, account2.ID)

		fromEntry := result.FromEntry
		require.NotEmpty(t, fromEntry)
		require.NotZero(t, fromEntry.ID)
		require.NotZero(t, fromEntry.CreatedAt)
		require.Equal(t, fromEntry.AccountID, account1.ID)
		require.Equal(t, fromEntry.Amount, -amount)

		toEntry := result.ToEntry
		require.NotEmpty(t, toEntry)
		require.NotZero(t, toEntry.ID)
		require.NotZero(t, toEntry.CreatedAt)
		require.Equal(t, toEntry.AccountID, account2.ID)
		require.Equal(t, toEntry.Amount, amount)

		fromAccount := result.FromAccount

		require.NotEmpty(t, fromAccount)
		require.Equal(t, account1.ID, fromAccount.ID)

		toAccount := result.ToAccount
		require.NotEmpty(t, toAccount)
		require.Equal(t, account2.ID, toAccount.ID)

		diff1 := account1.Balance - fromAccount.Balance
		diff2 := toAccount.Balance - account2.Balance

		require.Equal(t, diff1, diff2)
		require.True(t, diff1 > 0)
		require.True(t, diff2%amount == 0)

		k := int(diff1 / amount)
		require.True(t, k >= 1 && k <= n)
		require.NotContains(t, exsits, k)
		exsits[k] = true
	}

	updateAccount1, err := testQuery.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)
	require.NotZero(t, updateAccount1.ID)

	updateAccount2, err := testQuery.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)
	require.NotZero(t, updateAccount2.ID)

	require.Equal(t, account1.Balance - int64(n) * amount, updateAccount1.Balance)
	require.Equal(t, account2.Balance + int64(n) * amount, updateAccount2.Balance)
}
