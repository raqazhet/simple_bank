package util

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPassword(t *testing.T) {
	password := RandomString(6)
	hashPassword, err := HashedPassword(password)
	require.NoError(t, err)
	require.NotEmpty(t, hashPassword)
	err = CHeckPassword(password, hashPassword)
	require.NoError(t, err)
}
