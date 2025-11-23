package models

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/kir1l9x/avito-test-assignment/internal/domain/errors"
)

func TestNewTeamMember_Success(t *testing.T) {
	m, err := NewTeamMember("u1", "Alice", true)
	require.NoError(t, err)
	require.NotNil(t, m)
	require.Equal(t, "u1", m.ID)
	require.Equal(t, "Alice", m.Username)
	require.True(t, m.IsActive)
}

func TestNewTeamMember_ValidationError(t *testing.T) {
	m, err := NewTeamMember("", "Alice", true)
	require.Nil(t, m)
	require.Error(t, err)

	var derr *errors.DomainError
	require.ErrorAs(t, err, &derr)
}
