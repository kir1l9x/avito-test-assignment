package models

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/kir1l9x/avito-test-assignment/internal/domain/errors"
)

func TestNewUser_Success(t *testing.T) {
	u, err := NewUser("u1", "Alice", "backend", true)
	require.NoError(t, err)
	require.NotNil(t, u)
	require.Equal(t, "u1", u.ID)
	require.Equal(t, "Alice", u.Username)
	require.Equal(t, "backend", u.TeamName)
	require.True(t, u.IsActive)
}

func TestNewUser_ValidationError(t *testing.T) {
	u, err := NewUser("", "Alice", "backend", true)
	require.Nil(t, u)
	require.Error(t, err)

	var derr *errors.DomainError
	require.ErrorAs(t, err, &derr)
}
