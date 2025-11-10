package db

import (
	"context"
	"github.com/jxgzzztang/simplebank/util"
	"github.com/stretchr/testify/require"
	"testing"
)

func RandomUser(t *testing.T) User {
	ctx := context.Background()
	hashPassword, err := util.HashPassword(util.RandomString(6))
	require.NoError(t, err)
	accountsParams := CreateUserParams{
		Username: util.RandomOwner(),
		FullName: util.RandomOwner(),
		HashedPassword: hashPassword,
		Email: util.RandomEmail(),
	}

	result, err := testQuery.CreateUser(ctx, accountsParams)

	require.NoError(t, err)
	require.NotEmpty(t, result)

	require.Equal(t, accountsParams.HashedPassword, result.HashedPassword)
	require.Equal(t, accountsParams.FullName, result.FullName)
	require.Equal(t, accountsParams.Username, result.Username)
	require.Equal(t, accountsParams.Email, result.Email)
	require.NotZero(t, result.CreatedAt)
	require.NotEmpty(t, result.HashedPassword)

	return result

}

func TestCreateUser(t *testing.T) {
	RandomUser(t)
}

func TestGetUser(t *testing.T)  {
	user := RandomUser(t)
	ctx := context.Background()

	resultUser, err := testQuery.GetUser(ctx, user.Username)

	require.NoError(t, err)
	require.NotEmpty(t, resultUser)
	require.Equal(t, user.Username, resultUser.Username)
	require.Equal(t, user.FullName, resultUser.FullName)
	require.Equal(t, user.Email, resultUser.Email)
	require.NotZero(t, resultUser.CreatedAt)
}