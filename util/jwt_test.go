package util

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestJWT(t *testing.T) {
	var username = "Test"
	accessToken, _, err := CreateToken(username, time.Minute)
	require.NoError(t, err)
	require.NotEmpty(t, accessToken)
	claims, ok := ParseToken(accessToken)
	require.True(t, ok)
	require.Equal(t, username, claims.Username)
}
