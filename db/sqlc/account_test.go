package db

import (
	"context"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCreateAccount(t *testing.T) {
	ctx := context.Background()

	accountsParams := CreateAccountParams{
		Owner:    "tom",
		Balance:  100,
		Currency: "USD",
	}

	result, err := testQuery.CreateAccount(ctx, accountsParams)

	require.NoError(t, err)
	require.NotEmpty(t, result)

	require.NotZero(t, result.ID)
	require.NotZero(t, result.CreatedAt)

	require.Equal(t, accountsParams.Owner, result.Owner)
	require.Equal(t, accountsParams.Balance, result.Balance)
	require.Equal(t, accountsParams.Currency, result.Currency)

}
