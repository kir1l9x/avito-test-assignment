package services

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/kir1l9x/avito-test-assignment/internal/domain/models"
	"github.com/kir1l9x/avito-test-assignment/internal/domain/valueObjects"
)

func TestPullRequestsService_CreatePullRequest_Success(t *testing.T) {
	ctx := context.Background()

	prRepo := new(mockPullRequestsRepository)
	userRepo := new(mockUserRepository)
	trm := &fakeTrManager{}

	svc := NewPullRequestService(prRepo, userRepo, trm)

	author := &models.User{
		ID:       "u1",
		Username: "Alice",
		TeamName: "backend",
		IsActive: true,
	}

	activeUsers := []models.User{
		*author,
		{ID: "u2", Username: "Bob", TeamName: "backend", IsActive: true},
		{ID: "u3", Username: "Charlie", TeamName: "backend", IsActive: true},
	}

	userRepo.On("GetByID", mock.Anything, "u1").
		Return(author, nil)

	userRepo.On("GetActiveUsersInTeam", mock.Anything, "backend").
		Return(activeUsers, nil)

	createdPR := &models.PullRequest{
		ID:        "pr-1",
		Name:      "Initial setup",
		AuthorID:  "u1",
		Status:    valueObjects.OpenPullRequest,
		Reviewers: []string{"u2", "u3"},
		CreatedAt: time.Now().UTC(),
	}

	prRepo.On("Create", mock.Anything, mock.AnythingOfType("*models.PullRequest")).
		Return(createdPR, nil)

	pr, err := svc.CreatePullRequest(ctx, "pr-1", "Initial setup", "u1")
	require.NoError(t, err)
	require.NotNil(t, pr)
	require.Equal(t, createdPR, pr)

	require.True(t, trm.called)

	prRepo.AssertExpectations(t)
	userRepo.AssertExpectations(t)
}

func TestPullRequestsService_CreatePullRequest_GetAuthorError(t *testing.T) {
	ctx := context.Background()

	prRepo := new(mockPullRequestsRepository)
	userRepo := new(mockUserRepository)
	trm := &fakeTrManager{}

	svc := NewPullRequestService(prRepo, userRepo, trm)

	userRepo.On("GetByID", mock.Anything, "u1").
		Return((*models.User)(nil), errors.New("db error"))

	pr, err := svc.CreatePullRequest(ctx, "pr-1", "Initial setup", "u1")
	require.Error(t, err)
	require.Nil(t, pr)
	require.False(t, trm.called)
}

func TestPullRequestsService_CreatePullRequest_GetActiveUsersError(t *testing.T) {
	ctx := context.Background()

	prRepo := new(mockPullRequestsRepository)
	userRepo := new(mockUserRepository)
	trm := &fakeTrManager{}
	svc := NewPullRequestService(prRepo, userRepo, trm)

	author := &models.User{ID: "u1", Username: "Alice", TeamName: "backend"}

	userRepo.On("GetByID", mock.Anything, "u1").
		Return(author, nil)

	userRepo.On("GetActiveUsersInTeam", mock.Anything, "backend").
		Return(nil, errors.New("team error"))

	pr, err := svc.CreatePullRequest(ctx, "pr-1", "Initial setup", "u1")
	require.Error(t, err)
	require.Nil(t, pr)
	require.False(t, trm.called)
}

func TestPullRequestsService_CreatePullRequest_RepoCreateError(t *testing.T) {
	ctx := context.Background()

	prRepo := new(mockPullRequestsRepository)
	userRepo := new(mockUserRepository)
	trm := &fakeTrManager{}
	svc := NewPullRequestService(prRepo, userRepo, trm)

	author := &models.User{ID: "u1", Username: "Alice", TeamName: "backend"}
	active := []models.User{
		*author, {ID: "u2", Username: "Bob", TeamName: "backend", IsActive: true},
	}

	userRepo.On("GetByID", mock.Anything, "u1").
		Return(author, nil)

	userRepo.On("GetActiveUsersInTeam", mock.Anything, "backend").
		Return(active, nil)

	prRepo.On("Create", mock.Anything, mock.Anything).
		Return((*models.PullRequest)(nil), errors.New("insert error"))

	pr, err := svc.CreatePullRequest(ctx, "pr-1", "Initial setup", "u1")
	require.Error(t, err)
	require.Nil(t, pr)
	require.True(t, trm.called)
}

func TestPullRequestsService_GetPRsWhereUserIsReviewer_Success(t *testing.T) {
	ctx := context.Background()
	prRepo := new(mockPullRequestsRepository)

	svc := NewPullRequestService(prRepo, nil, &fakeTrManager{})

	list := []models.PullRequest{
		{ID: "pr-1"},
		{ID: "pr-2"},
	}

	prRepo.On("GetListByUser", mock.Anything, "u1").
		Return(list, nil)

	got, err := svc.GetPRsWhereUserIsReviewer(ctx, "u1")
	require.NoError(t, err)
	require.Equal(t, list, got)
}

func TestPullRequestsService_Merge_AlreadyMerged(t *testing.T) {
	ctx := context.Background()

	prRepo := new(mockPullRequestsRepository)
	svc := NewPullRequestService(prRepo, nil, &fakeTrManager{})

	existing := &models.PullRequest{
		ID:     "pr-9",
		Status: valueObjects.MergedPullRequest,
	}

	prRepo.On("GetByID", mock.Anything, "pr-9").
		Return(existing, nil)

	res, err := svc.Merge(ctx, "pr-9")
	require.NoError(t, err)
	require.Equal(t, existing, res)
}

func TestPullRequestsService_Merge_Success(t *testing.T) {
	ctx := context.Background()

	prRepo := new(mockPullRequestsRepository)
	trm := &fakeTrManager{}
	svc := NewPullRequestService(prRepo, nil, trm)

	existing := &models.PullRequest{
		ID:        "pr-1",
		Status:    valueObjects.OpenPullRequest,
		CreatedAt: time.Now().Add(-time.Hour),
	}

	prRepo.On("GetByID", mock.Anything, "pr-1").
		Return(existing, nil)

	updated := &models.PullRequest{
		ID:        "pr-1",
		AuthorID:  existing.AuthorID,
		Name:      existing.Name,
		Status:    valueObjects.MergedPullRequest,
		CreatedAt: existing.CreatedAt,
		MergedAt:  ptrTime(time.Now().UTC()),
	}

	prRepo.On("Update", mock.Anything, mock.Anything).
		Return(updated, nil)

	res, err := svc.Merge(ctx, "pr-1")
	require.NoError(t, err)
	require.Equal(t, updated, res)
	require.True(t, trm.called)
}

func TestPullRequestsService_ReassignReviewer_Success(t *testing.T) {
	ctx := context.Background()

	prRepo := new(mockPullRequestsRepository)
	userRepo := new(mockUserRepository)
	trm := &fakeTrManager{}
	svc := NewPullRequestService(prRepo, userRepo, trm)

	pr := &models.PullRequest{
		ID:        "pr-1",
		Status:    valueObjects.OpenPullRequest,
		AuthorID:  "u1",
		Reviewers: []string{"u2"},
	}

	author := &models.User{ID: "u1", TeamName: "backend"}

	active := []models.User{
		{ID: "u1", TeamName: "backend"},
		{ID: "u2", TeamName: "backend"},
		{ID: "u3", TeamName: "backend"},
	}

	prRepo.On("GetByID", mock.Anything, "pr-1").
		Return(pr, nil)

	userRepo.On("GetByID", mock.Anything, "u1").
		Return(author, nil)

	userRepo.On("GetActiveUsersInTeam", mock.Anything, "backend").
		Return(active, nil)

	updated := &models.PullRequest{
		ID:        "pr-1",
		Status:    valueObjects.OpenPullRequest,
		AuthorID:  "u1",
		Reviewers: []string{"u3"},
		CreatedAt: time.Now().UTC(),
	}

	prRepo.On("Update", mock.Anything, mock.Anything).
		Return(updated, nil)

	res, replacedBy, err := svc.ReassignReviewer(ctx, "pr-1", "u2")
	require.NoError(t, err)
	require.Equal(t, "u3", replacedBy)
	require.Equal(t, updated, res)
	require.True(t, trm.called)
}

func TestPullRequestsService_GetReviewsStats(t *testing.T) {
	ctx := context.Background()
	prRepo := new(mockPullRequestsRepository)
	svc := NewPullRequestService(prRepo, nil, &fakeTrManager{})

	expected := []models.UserReviewStats{
		{UserID: "u1", ReviewCount: 3},
		{UserID: "u2", ReviewCount: 1},
	}

	prRepo.On("GetReviewsAmountByUser", mock.Anything).Return(expected, nil)

	stats, err := svc.GetReviewsAmountPerUser(ctx)
	require.NoError(t, err)
	require.Equal(t, expected, stats)
}

func TestPullRequestsService_GetTop5ReviewersAllTime(t *testing.T) {
	ctx := context.Background()
	prRepo := new(mockPullRequestsRepository)
	svc := NewPullRequestService(prRepo, nil, &fakeTrManager{})

	expected := []models.UserReviewStats{{UserID: "u1", ReviewCount: 5}}
	prRepo.On("GetTop5ReviewsAmountByUser", mock.Anything).Return(expected, nil)

	stats, err := svc.GetTop5ReviewersAllTime(ctx)
	require.NoError(t, err)
	require.Equal(t, expected, stats)
}

func TestPullRequestsService_GetOpenReviewsAmount(t *testing.T) {
	ctx := context.Background()
	prRepo := new(mockPullRequestsRepository)
	svc := NewPullRequestService(prRepo, nil, &fakeTrManager{})

	expected := []models.UserReviewStats{{UserID: "u1", ReviewCount: 2}}
	prRepo.On("GetOpenReviewsAmountByUser", mock.Anything).Return(expected, nil)

	stats, err := svc.GetOpenReviewsAmountPerUser(ctx)
	require.NoError(t, err)
	require.Equal(t, expected, stats)
}

func TestPullRequestsService_GetTop5ReviewersCurrent(t *testing.T) {
	ctx := context.Background()
	prRepo := new(mockPullRequestsRepository)
	svc := NewPullRequestService(prRepo, nil, &fakeTrManager{})

	expected := []models.UserReviewStats{{UserID: "u1", ReviewCount: 2}}
	prRepo.On("GetTop5OpenReviewsAmountByUser", mock.Anything).Return(expected, nil)

	stats, err := svc.GetTop5ReviewersCurrent(ctx)
	require.NoError(t, err)
	require.Equal(t, expected, stats)
}

func ptrTime(t time.Time) *time.Time {
	return &t
}
