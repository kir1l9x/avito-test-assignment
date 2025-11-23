package services

import (
	"context"

	"github.com/avito-tech/go-transaction-manager/trm/v2"
	"go.uber.org/zap"

	domainErrors "github.com/kir1l9x/avito-test-assignment/internal/domain/errors"
	"github.com/kir1l9x/avito-test-assignment/internal/domain/models"
	"github.com/kir1l9x/avito-test-assignment/internal/domain/repositories"
	"github.com/kir1l9x/avito-test-assignment/internal/infrastructure/mappers"
	"github.com/kir1l9x/avito-test-assignment/pkg/logger"
)

type TeamService struct {
	teamRepo repositories.TeamsRepository
	userRepo repositories.UserRepository
	trM      trm.Manager
}

func NewTeamService(teamRepo repositories.TeamsRepository, userRepo repositories.UserRepository, trM trm.Manager) *TeamService {
	return &TeamService{teamRepo: teamRepo, userRepo: userRepo, trM: trM}
}

func (ts *TeamService) CreateTeamWithMembers(ctx context.Context, team *models.Team) (*models.Team, error) {
	logg := logger.FromContext(ctx)

	var createdTeam *models.Team

	err := ts.trM.Do(ctx, func(ctx context.Context) error {
		t, err := ts.teamRepo.Create(ctx, team)
		if err != nil {
			logg.Error("error while trying to create team", zap.Error(err), zap.String("team_name", team.Name))
			return err
		}

		createdTeam = t

		users, err := mappers.FromMembersToUsers(team.Members, team.Name)
		if err != nil {
			logg.Error("error while trying to map members into users", zap.Error(err))
			return domainErrors.New(domainErrors.CodeInternal, err)
		}

		createdUsers := make([]models.User, 0, len(users))

		for _, user := range users {
			createdUser, err := ts.userRepo.CreateOrUpdate(ctx, &user)
			if err != nil {
				logg.Error("error while trying to create user", zap.Error(err))
				return err
			}

			createdUsers = append(createdUsers, *createdUser)
		}

		createdTeam.Members, err = mappers.FromUsersToMembers(createdUsers)
		if err != nil {
			logg.Error("error while trying to map users into members", zap.Error(err))
			return domainErrors.New(domainErrors.CodeInternal, err)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return createdTeam, nil
}

func (ts *TeamService) GetTeamByName(ctx context.Context, name string) (*models.Team, error) {
	logg := logger.FromContext(ctx)

	team, err := ts.teamRepo.GetByName(ctx, name)
	if err != nil {
		logg.Error("error while trying to get team", zap.Error(err), zap.String("team_name", name))
		return nil, err
	}

	teamUsers, err := ts.userRepo.GetUsersByTeam(ctx, name)
	if err != nil {
		logg.Error("error while trying to get users from team", zap.Error(err), zap.String("team_name", name))
		return nil, err
	}

	teamMembers, err := mappers.FromUsersToMembers(teamUsers)
	if err != nil {
		logg.Error("error while trying to map users into members", zap.Error(err))
		return nil, domainErrors.New(domainErrors.CodeInternal, err)
	}

	team.Members = teamMembers

	return team, nil
}
