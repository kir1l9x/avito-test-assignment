package repositories

import (
	"context"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"

	"github.com/kir1l9x/avito-test-assignment/internal/domain/models"
	"github.com/kir1l9x/avito-test-assignment/internal/domain/valueObjects"
)

func ensureTeamExists(t *testing.T, pool *pgxpool.Pool, name string) {
	t.Helper()
	_, err := pool.Exec(context.Background(), `INSERT INTO teams (name) VALUES ($1) ON CONFLICT (name) DO NOTHING`, name)
	require.NoError(t, err)
}

func seedUser(
	t *testing.T,
	repo *UsersRepository,
	userID,
	username,
	teamName string,
	isActive bool,
) *models.User {
	t.Helper()

	user := &models.User{
		ID:       userID,
		Username: username,
		TeamName: teamName,
		IsActive: isActive,
	}

	created, err := repo.CreateOrUpdate(context.Background(), user)
	require.NoError(t, err)

	return created
}

func newPullRequestModel(id, authorID, name string, reviewers []string) *models.PullRequest {
	return &models.PullRequest{
		ID:        id,
		AuthorID:  authorID,
		Name:      name,
		Status:    valueObjects.OpenPullRequest,
		Reviewers: reviewers,
		CreatedAt: time.Now().UTC(),
	}
}
