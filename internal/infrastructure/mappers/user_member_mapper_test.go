package mappers

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/kir1l9x/avito-test-assignment/internal/domain/models"
)

func TestFromUsersToMembers_Success(t *testing.T) {
	users := []models.User{
		{ID: "u1", Username: "Alice", IsActive: true},
		{ID: "u2", Username: "Bob", IsActive: false},
	}

	members, err := FromUsersToMembers(users)
	require.NoError(t, err)
	require.Len(t, members, 2)

	require.Equal(t, "u1", members[0].ID)
	require.Equal(t, "Alice", members[0].Username)
	require.True(t, members[0].IsActive)

	require.Equal(t, "u2", members[1].ID)
	require.Equal(t, "Bob", members[1].Username)
	require.False(t, members[1].IsActive)
}

func TestFromUsersToMembers_Error(t *testing.T) {
	users := []models.User{
		{ID: "", Username: "Alice", IsActive: true},
	}

	members, err := FromUsersToMembers(users)
	require.Error(t, err)
	require.Nil(t, members)
}

func TestFromMembersToUsers_Success(t *testing.T) {
	members := []*models.TeamMember{
		{ID: "u1", Username: "Alice", IsActive: true},
		{ID: "u2", Username: "Bob", IsActive: false},
	}

	users, err := FromMembersToUsers(members, "backend")
	require.NoError(t, err)
	require.Len(t, users, 2)

	require.Equal(t, "u1", users[0].ID)
	require.Equal(t, "Alice", users[0].Username)
	require.Equal(t, "backend", users[0].TeamName)
	require.True(t, users[0].IsActive)

	require.Equal(t, "u2", users[1].ID)
	require.Equal(t, "Bob", users[1].Username)
	require.Equal(t, "backend", users[1].TeamName)
	require.False(t, users[1].IsActive)
}

func TestFromMembersToUsers_Error(t *testing.T) {
	members := []*models.TeamMember{
		{ID: "", Username: "Alice", IsActive: true},
	}

	users, err := FromMembersToUsers(members, "backend")
	require.Error(t, err)
	require.Nil(t, users)
}
