package repositories

import (
	"context"
	"testing"

	pgxTm "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/stretchr/testify/require"

	domainErrors "github.com/kir1l9x/avito-test-assignment/internal/domain/errors"
	"github.com/kir1l9x/avito-test-assignment/internal/domain/models"
)

func TestUsersRepository(t *testing.T) {
	pool := SetupTestDB(t)
	getter := pgxTm.DefaultCtxGetter
	usersRepo := NewUsersRepository(pool, getter)

	t.Run("CreateOrUpdate inserts and updates user", func(t *testing.T) {
		TruncateAll(t, pool)
		ctx := context.Background()
		ensureTeamExists(t, pool, "core")

		input := &models.User{
			ID:       "u1",
			Username: "John",
			TeamName: "core",
			IsActive: true,
		}

		created, err := usersRepo.CreateOrUpdate(ctx, input)
		require.NoError(t, err)
		require.Equal(t, input, created)

		input.Username = "Johnny"
		input.IsActive = false

		updated, err := usersRepo.CreateOrUpdate(ctx, input)
		require.NoError(t, err)
		require.Equal(t, input.Username, updated.Username)
		require.Equal(t, input.IsActive, updated.IsActive)
	})

	t.Run("CreateOrUpdate returns TEAM_NOT_FOUND when FK fails", func(t *testing.T) {
		TruncateAll(t, pool)
		ctx := context.Background()

		_, err := usersRepo.CreateOrUpdate(ctx, &models.User{
			ID:       "u2",
			Username: "Ann",
			TeamName: "ghost",
			IsActive: true,
		})
		require.Error(t, err)

		var de *domainErrors.DomainError
		require.ErrorAs(t, err, &de)
		require.Equal(t, domainErrors.CodeTeamNotFound, de.Code)
	})

	t.Run("GetByID returns stored user", func(t *testing.T) {
		TruncateAll(t, pool)
		ctx := context.Background()
		ensureTeamExists(t, pool, "analytics")

		expected := seedUser(t, usersRepo, "u3", "Kate", "analytics", true)

		got, err := usersRepo.GetByID(ctx, expected.ID)
		require.NoError(t, err)
		require.Equal(t, expected, got)
	})

	t.Run("GetByID returns USER_NOT_FOUND", func(t *testing.T) {
		TruncateAll(t, pool)
		ctx := context.Background()

		_, err := usersRepo.GetByID(ctx, "missing")
		require.Error(t, err)

		var de *domainErrors.DomainError
		require.ErrorAs(t, err, &de)
		require.Equal(t, domainErrors.CodeUserNotFound, de.Code)
	})

	t.Run("GetUsersByTeam returns all members", func(t *testing.T) {
		TruncateAll(t, pool)
		ctx := context.Background()
		ensureTeamExists(t, pool, "dev")

		first := seedUser(t, usersRepo, "u10", "Alice", "dev", true)
		second := seedUser(t, usersRepo, "u11", "Bob", "dev", false)

		users, err := usersRepo.GetUsersByTeam(ctx, "dev")
		require.NoError(t, err)
		require.Len(t, users, 2)

		got := map[string]models.User{}
		for _, u := range users {
			got[u.ID] = u
		}

		require.Equal(t, *first, got[first.ID])
		require.Equal(t, *second, got[second.ID])
	})

	t.Run("GetUsersByTeam returns empty slice for unknown team", func(t *testing.T) {
		TruncateAll(t, pool)
		ctx := context.Background()

		users, err := usersRepo.GetUsersByTeam(ctx, "no-team")
		require.NoError(t, err)
		require.Empty(t, users)
	})

	t.Run("GetActiveUsersInTeam returns only active members", func(t *testing.T) {
		TruncateAll(t, pool)
		ctx := context.Background()
		ensureTeamExists(t, pool, "ml")

		seedUser(t, usersRepo, "u20", "Active", "ml", true)
		seedUser(t, usersRepo, "u21", "Inactive", "ml", false)

		active, err := usersRepo.GetActiveUsersInTeam(ctx, "ml")
		require.NoError(t, err)
		require.Len(t, active, 1)
		require.Equal(t, "u20", active[0].ID)
	})

	t.Run("SetIsActive toggles flag and returns updated user", func(t *testing.T) {
		TruncateAll(t, pool)
		ctx := context.Background()
		ensureTeamExists(t, pool, "design")

		seedUser(t, usersRepo, "u30", "Mark", "design", false)

		updated, err := usersRepo.SetIsActive(ctx, "u30", true)
		require.NoError(t, err)
		require.True(t, updated.IsActive)
	})

	t.Run("SetIsActive returns USER_NOT_FOUND", func(t *testing.T) {
		TruncateAll(t, pool)
		ctx := context.Background()

		_, err := usersRepo.SetIsActive(ctx, "unknown", true)
		require.Error(t, err)

		var de *domainErrors.DomainError
		require.ErrorAs(t, err, &de)
		require.Equal(t, domainErrors.CodeUserNotFound, de.Code)
	})
}
