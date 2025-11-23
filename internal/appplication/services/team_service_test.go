package services

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	domainErrors "github.com/kir1l9x/avito-test-assignment/internal/domain/errors"
	"github.com/kir1l9x/avito-test-assignment/internal/domain/models"
)

func TestTeamService_CreateTeamWithMembers_Success(t *testing.T) {
	ctx := context.Background()

	teamRepo := new(mockTeamsRepository)
	userRepo := new(mockUserRepository)
	trm := &fakeTrManager{}
	svc := NewTeamService(teamRepo, userRepo, trm)

	team := &models.Team{
		Name: "backend",
		Members: []*models.TeamMember{
			{ID: "u1", Username: "Alice", IsActive: true},
			{ID: "u2", Username: "Bob", IsActive: false},
		},
	}

	teamRepo.On("Create", mock.Anything, team).
		Return(&models.Team{Name: "backend"}, nil)

	userRepo.On("CreateOrUpdate", mock.Anything, mock.MatchedBy(func(u *models.User) bool {
		return u.ID == "u1"
	})).
		Return(&models.User{
			ID:       "u1",
			Username: "Alice",
			TeamName: "backend",
			IsActive: true,
		}, nil)

	userRepo.On("CreateOrUpdate", mock.Anything, mock.MatchedBy(func(u *models.User) bool {
		return u.ID == "u2"
	})).
		Return(&models.User{
			ID:       "u2",
			Username: "Bob",
			TeamName: "backend",
			IsActive: false,
		}, nil)

	res, err := svc.CreateTeamWithMembers(ctx, team)
	require.NoError(t, err)
	require.NotNil(t, res)
	require.True(t, trm.called)
}

func TestTeamService_CreateTeamWithMembers_TeamRepoError(t *testing.T) {
	ctx := context.Background()
	teamRepo := new(mockTeamsRepository)
	userRepo := new(mockUserRepository)
	trm := &fakeTrManager{}
	svc := NewTeamService(teamRepo, userRepo, trm)

	team := &models.Team{
		Name:    "backend",
		Members: []*models.TeamMember{{ID: "u1", Username: "Alice", IsActive: true}},
	}

	teamRepo.On("Create", mock.Anything, team).
		Return((*models.Team)(nil), errors.New("insert error"))

	res, err := svc.CreateTeamWithMembers(ctx, team)
	require.Error(t, err)
	require.Nil(t, res)
	require.True(t, trm.called)
}

func TestTeamService_CreateTeamWithMembers_UserRepoError(t *testing.T) {
	ctx := context.Background()
	teamRepo := new(mockTeamsRepository)
	userRepo := new(mockUserRepository)
	trm := &fakeTrManager{}
	svc := NewTeamService(teamRepo, userRepo, trm)

	teamRepo.On("Create", mock.Anything, mock.Anything).
		Return(&models.Team{Name: "backend"}, nil)

	userRepo.On("CreateOrUpdate", mock.Anything, mock.Anything).
		Return((*models.User)(nil), errors.New("db error"))

	res, err := svc.CreateTeamWithMembers(ctx, &models.Team{
		Name:    "backend",
		Members: []*models.TeamMember{{ID: "u1", Username: "Alice", IsActive: true}},
	})

	require.Error(t, err)
	require.Nil(t, res)
	require.True(t, trm.called)
}

func TestTeamService_GetTeamByName_Success(t *testing.T) {
	ctx := context.Background()
	teamRepo := new(mockTeamsRepository)
	userRepo := new(mockUserRepository)
	svc := NewTeamService(teamRepo, userRepo, &fakeTrManager{})

	teamRepo.On("GetByName", mock.Anything, "backend").
		Return(&models.Team{Name: "backend"}, nil)

	userRepo.On("GetUsersByTeam", mock.Anything, "backend").
		Return([]models.User{
			{ID: "u1", Username: "Alice", TeamName: "backend", IsActive: true},
			{ID: "u2", Username: "Bob", TeamName: "backend"},
		}, nil)

	res, err := svc.GetTeamByName(ctx, "backend")
	require.NoError(t, err)
	require.NotNil(t, res)
	require.Len(t, res.Members, 2)
}

func TestTeamService_GetTeamByName_TeamNotFound(t *testing.T) {
	ctx := context.Background()
	teamRepo := new(mockTeamsRepository)
	userRepo := new(mockUserRepository)

	svc := NewTeamService(teamRepo, userRepo, &fakeTrManager{})

	teamRepo.On("GetByName", mock.Anything, "backend").
		Return(nil, domainErrors.New(domainErrors.CodeTeamNotFound, domainErrors.ErrTeamNotFound))

	res, err := svc.GetTeamByName(ctx, "backend")
	require.Error(t, err)
	require.Nil(t, res)
}

func TestTeamService_GetTeamByName_UsersRepoError(t *testing.T) {
	ctx := context.Background()
	teamRepo := new(mockTeamsRepository)
	userRepo := new(mockUserRepository)

	svc := NewTeamService(teamRepo, userRepo, &fakeTrManager{})

	teamRepo.On("GetByName", mock.Anything, "backend").
		Return(&models.Team{Name: "backend"}, nil)

	userRepo.On("GetUsersByTeam", mock.Anything, "backend").
		Return(nil, errors.New("db error"))

	res, err := svc.GetTeamByName(ctx, "backend")
	require.Error(t, err)
	require.Nil(t, res)
}
