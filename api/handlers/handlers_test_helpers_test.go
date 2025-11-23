package handlers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	pgxTm "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"

	"github.com/kir1l9x/avito-test-assignment/api/dto/requests"
	"github.com/kir1l9x/avito-test-assignment/api/dto/responses"
	"github.com/kir1l9x/avito-test-assignment/internal/domain/models"
	"github.com/kir1l9x/avito-test-assignment/internal/domain/valueObjects"
	infraRepo "github.com/kir1l9x/avito-test-assignment/internal/infrastructure/repositories"
)

func postJSON(t *testing.T, url string, payload any) *http.Response {
	t.Helper()

	body, err := json.Marshal(payload)
	require.NoError(t, err)

	resp, err := http.Post(url, "application/json", bytes.NewReader(body)) //nolint:gosec // test helper controls URL
	require.NoError(t, err)

	return resp
}

func decodeResponse(t *testing.T, resp *http.Response, target any) {
	t.Helper()
	t.Cleanup(func() {
		require.NoError(t, resp.Body.Close())
	})

	err := json.NewDecoder(resp.Body).Decode(target)
	require.NoError(t, err)
}

func createTeam(t *testing.T, baseURL string, teamName string, members []requests.TeamMemberRequest) responses.AddTeamResponse {
	t.Helper()

	req := requests.AddTeamRequest{
		TeamName: teamName,
		Members:  members,
	}

	resp := postJSON(t, baseURL+"/team/add", req)
	require.Equal(t, http.StatusCreated, resp.StatusCode)

	var parsed responses.AddTeamResponse
	decodeResponse(t, resp, &parsed)

	return parsed
}

func createPullRequestRecord(
	t *testing.T,
	pool *pgxpool.Pool,
	prID string,
	prName string,
	authorID string,
	reviewers []string,
) *models.PullRequest { //nolint:unparam // authorID kept for future flexibility
	t.Helper()

	_, err := pool.Exec(context.Background(), `DELETE FROM pr_reviewers WHERE pr_id = $1`, prID)
	require.NoError(t, err)
	_, err = pool.Exec(context.Background(), `DELETE FROM pull_requests WHERE id = $1`, prID)
	require.NoError(t, err)

	repo := infraRepo.NewPullRequestsRepository(pool, pgxTm.DefaultCtxGetter)

	pr, err := models.NewPullRequest(prID, authorID, prName, valueObjects.OpenPullRequest, reviewers)
	require.NoError(t, err)

	created, err := repo.Create(context.Background(), pr)
	require.NoError(t, err)

	return created
}

func markPullRequestMerged(t *testing.T, pool *pgxpool.Pool, prID string) {
	t.Helper()

	ctx := context.Background()
	repo := infraRepo.NewPullRequestsRepository(pool, pgxTm.DefaultCtxGetter)

	pr, err := repo.GetByID(ctx, prID)
	require.NoError(t, err)

	pr.Status = valueObjects.MergedPullRequest
	now := time.Now().UTC()
	pr.MergedAt = &now

	_, err = repo.Update(ctx, pr)
	require.NoError(t, err)
}
