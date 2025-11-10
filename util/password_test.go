package util

import (
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
	"testing"
)

func TestHashPassword(t *testing.T)  {
	password := RandomString(6)

	hashPwd, err := HashPassword(password)
	require.NoError(t, err)
	require.NotEmpty(t, hashPwd)
	err = CheckPassword(password, hashPwd)
	require.NoError(t, err)

	wrongPassword := RandomString(6)
	err = CheckPassword(wrongPassword, hashPwd)
	require.EqualError(t, err, bcrypt.ErrMismatchedHashAndPassword.Error())

}