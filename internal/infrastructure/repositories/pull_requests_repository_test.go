package repositories

import (
	"context"
	"testing"
	"time"

	pgxTm "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/stretchr/testify/require"

	domainErrors "github.com/kir1l9x/avito-test-assignment/internal/domain/errors"
	"github.com/kir1l9x/avito-test-assignment/internal/domain/models"
	"github.com/kir1l9x/avito-test-assignment/internal/domain/valueObjects"
)

func TestPullRequestsRepository(t *testing.T) {
	pool := SetupTestDB(t)
	getter := pgxTm.DefaultCtxGetter
	teamsRepo := NewTeamsRepository(pool, getter)
	usersRepo := NewUsersRepository(pool, getter)

	t.Run("Create/Get/Update lifecycle rewrites reviewers list", func(t *testing.T) {
		TruncateAll(t, pool)
		ctx := context.Background()

		team := &models.Team{Name: "alpha"}
		createdTeam, err := teamsRepo.Create(ctx, team)
		require.NoError(t, err)

		author := seedUser(t, usersRepo, "u1", "author", createdTeam.Name, true)
		reviewerA := seedUser(t, usersRepo, "u2", "rev-a", createdTeam.Name, true)
		reviewerB := seedUser(t, usersRepo, "u3", "rev-b", createdTeam.Name, true)

		prRepo := NewPullRequestsRepository(pool, getter)

		prModel := &models.PullRequest{
			ID:        "pr-1",
			AuthorID:  author.ID,
			Name:      "Initial PR",
			Status:    valueObjects.OpenPullRequest,
			Reviewers: []string{reviewerA.ID},
			CreatedAt: time.Now().UTC(),
		}

		createdPR, err := prRepo.Create(ctx, prModel)
		require.NoError(t, err)
		require.Equal(t, prModel.ID, createdPR.ID)
		require.ElementsMatch(t, []string{reviewerA.ID}, createdPR.Reviewers)

		gotByID, err := prRepo.GetByID(ctx, prModel.ID)
		require.NoError(t, err)
		require.Equal(t, createdPR.ID, gotByID.ID)
		require.ElementsMatch(t, createdPR.Reviewers, gotByID.Reviewers)

		listByReviewer, err := prRepo.GetListByUser(ctx, reviewerA.ID)
		require.NoError(t, err)
		require.Len(t, listByReviewer, 1)
		require.Equal(t, prModel.ID, listByReviewer[0].ID)

		noAssignments, err := prRepo.GetListByUser(ctx, reviewerB.ID)
		require.NoError(t, err)
		require.Len(t, noAssignments, 0)

		updatedInput := *createdPR
		updatedInput.Status = valueObjects.MergedPullRequest
		mergedAt := time.Now().UTC()
		updatedInput.MergedAt = &mergedAt
		updatedInput.Reviewers = []string{reviewerB.ID}

		updatedPR, err := prRepo.Update(ctx, &updatedInput)
		require.NoError(t, err)
		require.Equal(t, valueObjects.MergedPullRequest, updatedPR.Status)
		require.NotNil(t, updatedPR.MergedAt)
		require.ElementsMatch(t, []string{reviewerB.ID}, updatedPR.Reviewers)

		listOldReviewer, err := prRepo.GetListByUser(ctx, reviewerA.ID)
		require.NoError(t, err)
		require.Empty(t, listOldReviewer)

		listNewReviewer, err := prRepo.GetListByUser(ctx, reviewerB.ID)
		require.NoError(t, err)
		require.Len(t, listNewReviewer, 1)
		require.Equal(t, prModel.ID, listNewReviewer[0].ID)
		require.ElementsMatch(t, []string{reviewerB.ID}, listNewReviewer[0].Reviewers)
	})

	t.Run("GetByID returns PR_NOT_FOUND", func(t *testing.T) {
		TruncateAll(t, pool)
		ctx := context.Background()
		prRepo := NewPullRequestsRepository(pool, getter)

		_, err := prRepo.GetByID(ctx, "unknown")
		require.Error(t, err)

		var de *domainErrors.DomainError
		require.ErrorAs(t, err, &de)
		require.Equal(t, domainErrors.CodePRNotFound, de.Code)
	})

	t.Run("Update returns PR_NOT_FOUND for missing entity", func(t *testing.T) {
		TruncateAll(t, pool)
		ctx := context.Background()
		prRepo := NewPullRequestsRepository(pool, getter)

		pr := newPullRequestModel("pr-missing", "u-x", "Ghost PR", nil)

		_, err := prRepo.Update(ctx, pr)
		require.Error(t, err)

		var de *domainErrors.DomainError
		require.ErrorAs(t, err, &de)
		require.Equal(t, domainErrors.CodePRNotFound, de.Code)
	})

	t.Run("Create returns PR_EXISTS on duplicate ID", func(t *testing.T) {
		TruncateAll(t, pool)
		ctx := context.Background()

		team := &models.Team{Name: "duplicates"}
		_, err := teamsRepo.Create(ctx, team)
		require.NoError(t, err)

		author := seedUser(t, usersRepo, "u10", "author", team.Name, true)
		reviewer := seedUser(t, usersRepo, "u11", "reviewer", team.Name, true)

		prRepo := NewPullRequestsRepository(pool, getter)

		pr := &models.PullRequest{
			ID:        "pr-dup",
			AuthorID:  author.ID,
			Name:      "Dup PR",
			Status:    valueObjects.OpenPullRequest,
			Reviewers: []string{reviewer.ID},
			CreatedAt: time.Now().UTC(),
		}

		_, err = prRepo.Create(ctx, pr)
		require.NoError(t, err)

		_, err = prRepo.Create(ctx, pr)
		require.Error(t, err)

		var de *domainErrors.DomainError
		require.ErrorAs(t, err, &de)
		require.Equal(t, domainErrors.CodePRExists, de.Code)
	})

	t.Run("Stats queries return aggregated data", func(t *testing.T) {
		TruncateAll(t, pool)
		ctx := context.Background()
		prRepo := NewPullRequestsRepository(pool, getter)

		allExpected, openExpected := seedReviewStatsData(t, ctx, teamsRepo, usersRepo, prRepo)

		statsAll, err := prRepo.GetReviewsAmountByUser(ctx)
		require.NoError(t, err)
		assertReviewStatsEqual(t, statsAll, allExpected)

		statsTop, err := prRepo.GetTop5ReviewsAmountByUser(ctx)
		require.NoError(t, err)
		assertReviewStatsEqual(t, statsTop, allExpected)

		statsOpen, err := prRepo.GetOpenReviewsAmountByUser(ctx)
		require.NoError(t, err)
		assertReviewStatsEqual(t, statsOpen, openExpected)

		statsTopOpen, err := prRepo.GetTop5OpenReviewsAmountByUser(ctx)
		require.NoError(t, err)
		assertReviewStatsEqual(t, statsTopOpen, openExpected)
	})
}

func seedReviewStatsData(
	t *testing.T,
	ctx context.Context,
	teamsRepo *TeamsRepository,
	usersRepo *UsersRepository,
	prRepo *PullRequestsRepository,
) (map[string]int64, map[string]int64) {
	t.Helper()

	team := &models.Team{Name: "stats"}
	_, err := teamsRepo.Create(ctx, team)
	require.NoError(t, err)

	author := seedUser(t, usersRepo, "us-author", "Author", team.Name, true)
	reviewerA := seedUser(t, usersRepo, "us-a", "Rev-A", team.Name, true)
	reviewerB := seedUser(t, usersRepo, "us-b", "Rev-B", team.Name, true)

	createTestPR(t, ctx, prRepo, "pr-open-1", author.ID, "Open A", valueObjects.OpenPullRequest, []string{reviewerA.ID, reviewerB.ID})
	createTestPR(t, ctx, prRepo, "pr-open-2", author.ID, "Open B", valueObjects.OpenPullRequest, []string{reviewerA.ID})
	createTestPR(t, ctx, prRepo, "pr-merged", author.ID, "Merged", valueObjects.MergedPullRequest, []string{reviewerB.ID})

	return map[string]int64{
			reviewerA.ID: 2,
			reviewerB.ID: 2,
		}, map[string]int64{
			reviewerA.ID: 2,
			reviewerB.ID: 1,
		}
}

func createTestPR(
	t *testing.T,
	ctx context.Context,
	prRepo *PullRequestsRepository,
	id string,
	authorID string,
	name string,
	status valueObjects.PullRequestStatus,
	reviewers []string,
) {
	t.Helper()

	pr := &models.PullRequest{
		ID:        id,
		AuthorID:  authorID,
		Name:      name,
		Status:    status,
		Reviewers: reviewers,
		CreatedAt: time.Now().UTC(),
	}

	if status == valueObjects.MergedPullRequest {
		now := time.Now().UTC()
		pr.MergedAt = &now
	}

	_, err := prRepo.Create(ctx, pr)
	require.NoError(t, err)
}

func assertReviewStatsEqual(t *testing.T, stats []models.UserReviewStats, expected map[string]int64) {
	t.Helper()

	actual := make(map[string]int64, len(stats))
	for _, stat := range stats {
		actual[stat.UserID] = stat.ReviewCount
	}

	require.Equal(t, expected, actual)
}
