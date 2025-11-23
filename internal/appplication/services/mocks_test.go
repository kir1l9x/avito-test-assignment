package services

import (
	"context"

	"github.com/avito-tech/go-transaction-manager/trm/v2"
	"github.com/stretchr/testify/mock"

	"github.com/kir1l9x/avito-test-assignment/internal/domain/models"
	"github.com/kir1l9x/avito-test-assignment/internal/domain/repositories"
)

type mockPullRequestsRepository struct {
	mock.Mock
}

func (m *mockPullRequestsRepository) Create(ctx context.Context, pr *models.PullRequest) (*models.PullRequest, error) {
	args := m.Called(ctx, pr)
	res, _ := args.Get(0).(*models.PullRequest)
	return res, args.Error(1)
}

func (m *mockPullRequestsRepository) GetByID(ctx context.Context, id string) (*models.PullRequest, error) {
	args := m.Called(ctx, id)
	res, _ := args.Get(0).(*models.PullRequest)
	return res, args.Error(1)
}

func (m *mockPullRequestsRepository) Update(ctx context.Context, pr *models.PullRequest) (*models.PullRequest, error) {
	args := m.Called(ctx, pr)
	res, _ := args.Get(0).(*models.PullRequest)
	return res, args.Error(1)
}

func (m *mockPullRequestsRepository) GetListByUser(ctx context.Context, userID string) ([]models.PullRequest, error) {
	args := m.Called(ctx, userID)
	res, _ := args.Get(0).([]models.PullRequest)
	return res, args.Error(1)
}

func (m *mockPullRequestsRepository) GetReviewsAmountByUser(ctx context.Context) ([]models.UserReviewStats, error) {
	args := m.Called(ctx)
	res, _ := args.Get(0).([]models.UserReviewStats)
	return res, args.Error(1)
}

func (m *mockPullRequestsRepository) GetTop5ReviewsAmountByUser(ctx context.Context) ([]models.UserReviewStats, error) {
	args := m.Called(ctx)
	res, _ := args.Get(0).([]models.UserReviewStats)
	return res, args.Error(1)
}

func (m *mockPullRequestsRepository) GetOpenReviewsAmountByUser(ctx context.Context) ([]models.UserReviewStats, error) {
	args := m.Called(ctx)
	res, _ := args.Get(0).([]models.UserReviewStats)
	return res, args.Error(1)
}

func (m *mockPullRequestsRepository) GetTop5OpenReviewsAmountByUser(ctx context.Context) ([]models.UserReviewStats, error) {
	args := m.Called(ctx)
	res, _ := args.Get(0).([]models.UserReviewStats)
	return res, args.Error(1)
}

type mockUserRepository struct {
	mock.Mock
}

func (m *mockUserRepository) CreateOrUpdate(ctx context.Context, user *models.User) (*models.User, error) {
	args := m.Called(ctx, user)
	res, _ := args.Get(0).(*models.User)
	return res, args.Error(1)
}

func (m *mockUserRepository) GetByID(ctx context.Context, id string) (*models.User, error) {
	args := m.Called(ctx, id)
	res, _ := args.Get(0).(*models.User)
	return res, args.Error(1)
}

func (m *mockUserRepository) GetUsersByTeam(ctx context.Context, teamName string) ([]models.User, error) {
	args := m.Called(ctx, teamName)
	res, _ := args.Get(0).([]models.User)
	return res, args.Error(1)
}

func (m *mockUserRepository) GetActiveUsersInTeam(ctx context.Context, teamName string) ([]models.User, error) {
	args := m.Called(ctx, teamName)
	res, _ := args.Get(0).([]models.User)
	return res, args.Error(1)
}

func (m *mockUserRepository) SetIsActive(ctx context.Context, id string, isActive bool) (*models.User, error) {
	args := m.Called(ctx, id, isActive)
	res, _ := args.Get(0).(*models.User)
	return res, args.Error(1)
}

var (
	_ repositories.UserRepository         = (*mockUserRepository)(nil)
	_ repositories.PullRequestsRepository = (*mockPullRequestsRepository)(nil)
)

type mockTeamsRepository struct {
	mock.Mock
}

func (m *mockTeamsRepository) Create(ctx context.Context, team *models.Team) (*models.Team, error) {
	args := m.Called(ctx, team)
	res, _ := args.Get(0).(*models.Team)
	return res, args.Error(1)
}

func (m *mockTeamsRepository) GetByName(ctx context.Context, name string) (*models.Team, error) {
	args := m.Called(ctx, name)
	res, _ := args.Get(0).(*models.Team)
	return res, args.Error(1)
}

var _ repositories.TeamsRepository = (*mockTeamsRepository)(nil)

type fakeTrManager struct {
	called bool
	err    error
}

func (m *fakeTrManager) Do(ctx context.Context, f func(ctx context.Context) error) error {
	m.called = true
	if m.err != nil {
		return m.err
	}
	return f(ctx)
}

func (m *fakeTrManager) DoWithSettings(ctx context.Context, _ trm.Settings, f func(ctx context.Context) error) error {
	m.called = true
	if m.err != nil {
		return m.err
	}
	return f(ctx)
}
