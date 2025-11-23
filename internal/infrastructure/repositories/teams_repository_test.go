package repositories

import (
	"context"
	"testing"

	pgxTm "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/stretchr/testify/require"

	domainErrors "github.com/kir1l9x/avito-test-assignment/internal/domain/errors"
	"github.com/kir1l9x/avito-test-assignment/internal/domain/models"
)

func TestTeamsRepository(t *testing.T) {
	pool := SetupTestDB(t)
	getter := pgxTm.DefaultCtxGetter

	t.Run("Create / GetByName — success", func(t *testing.T) {
		TruncateAll(t, pool)
		ctx := context.Background()

		repo := NewTeamsRepository(pool, getter)

		team := &models.Team{Name: "backend"}

		created, err := repo.Create(ctx, team)
		require.NoError(t, err)
		require.Equal(t, team.Name, created.Name)

		got, err := repo.GetByName(ctx, "backend")
		require.NoError(t, err)
		require.Equal(t, "backend", got.Name)
		require.Empty(t, got.Members)
	})

	t.Run("GetByName — TEAM_NOT_FOUND", func(t *testing.T) {
		TruncateAll(t, pool)
		ctx := context.Background()

		repo := NewTeamsRepository(pool, getter)

		_, err := repo.GetByName(ctx, "unknown-team")
		require.Error(t, err)

		var de *domainErrors.DomainError
		require.ErrorAs(t, err, &de)
		require.Equal(t, domainErrors.CodeTeamNotFound, de.Code)
	})

	t.Run("Create — duplicate team → TEAM_EXISTS", func(t *testing.T) {
		TruncateAll(t, pool)
		ctx := context.Background()

		repo := NewTeamsRepository(pool, getter)

		team := &models.Team{Name: "devops"}

		_, err := repo.Create(ctx, team)
		require.NoError(t, err)

		_, err = repo.Create(ctx, team)
		require.Error(t, err)

		var de *domainErrors.DomainError
		require.ErrorAs(t, err, &de)
		require.Equal(t, domainErrors.CodeTeamExists, de.Code)
	})
}
