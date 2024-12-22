package db

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/jxgzzztang/simplebank/util"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func RandomAccount(t *testing.T) Account {
	ctx := context.Background()

	accountsParams := CreateAccountParams{
		Owner:    util.RandomOwner(),
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}

	result, err := testQuery.CreateAccount(ctx, accountsParams)

	require.NoError(t, err)
	require.NotEmpty(t, result)

	require.NotZero(t, result.ID)
	require.NotZero(t, result.CreatedAt)

	require.Equal(t, accountsParams.Owner, result.Owner)
	require.Equal(t, accountsParams.Balance, result.Balance)
	require.Equal(t, accountsParams.Currency, result.Currency)

	return result

}

func TestCreateAccount(t *testing.T) {
	RandomAccount(t)
}

func TestGetAccount(t *testing.T) {
	account1 := RandomAccount(t)
	account2, err := testQuery.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, account2)

	require.Equal(t, account1.ID, account2.ID)
	require.Equal(t, account1.Owner, account2.Owner)
	require.Equal(t, account1.Balance, account2.Balance)
	require.Equal(t, account1.Currency, account2.Currency)
	require.WithinDuration(t, account1.CreatedAt.Time, account2.CreatedAt.Time, time.Second)
}

func TestUpdateAccount(t *testing.T) {
	account1 := RandomAccount(t)
	money := util.RandomMoney()
	account2, err := testQuery.UpdateAccount(context.Background(), UpdateAccountParams{
		ID:      account1.ID,
		Balance: money,
	})

	require.NoError(t, err)
	require.NotEmpty(t, account2)
	require.Equal(t, account1.ID, account2.ID)
	require.Equal(t, money, account2.Balance)
	require.Equal(t, account1.Currency, account2.Currency)
	require.WithinDuration(t, account1.CreatedAt.Time, account2.CreatedAt.Time, time.Second)
	require.NotEmpty(t, account2.Balance, account1.Balance)
}

func TestDeleteAccount(t *testing.T) {
	account1 := RandomAccount(t)
	err := testQuery.DeleteAccount(context.Background(), account1.ID)

	require.NoError(t, err)

	account2, err := testQuery.GetAccount(context.Background(), account1.ID)
	require.Error(t, err)
	require.EqualError(t, err, pgx.ErrNoRows.Error())
	require.Empty(t, account2)
}

func TestListAccounts(t *testing.T) {
	for i := 0; i < 10; i++ {
		RandomAccount(t)
	}

	listAccounts, err := testQuery.ListAccount(context.Background(), ListAccountParams{
		Limit:  5,
		Offset: 5,
	})
	require.NoError(t, err)
	require.Len(t, listAccounts, 5)

	for _, account := range listAccounts {
		require.NotEmpty(t, account)
	}

}
