package models

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/kir1l9x/avito-test-assignment/internal/domain/errors"
)

func TestNewTeam_Success(t *testing.T) {
	members := []*TeamMember{
		{ID: "u1", Username: "Alice", IsActive: true},
	}

	team, err := NewTeam("backend", members)
	require.NoError(t, err)
	require.NotNil(t, team)
	require.Equal(t, "backend", team.Name)
	require.Len(t, team.Members, 1)
}

func TestNewTeam_MissingName(t *testing.T) {
	team, err := NewTeam("", []*TeamMember{
		{ID: "u1", Username: "Alice", IsActive: true},
	})

	require.Nil(t, team)
	require.Error(t, err)

	var derr *errors.DomainError
	require.ErrorAs(t, err, &derr)
}

func TestNewTeam_NoMembers(t *testing.T) {
	team, err := NewTeam("backend", []*TeamMember{})
	require.Nil(t, team)
	require.Error(t, err)

	var derr *errors.DomainError
	require.ErrorAs(t, err, &derr)
}
