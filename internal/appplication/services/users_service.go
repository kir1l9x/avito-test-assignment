package services

import (
	"context"

	"go.uber.org/zap"

	"github.com/kir1l9x/avito-test-assignment/internal/domain/models"
	"github.com/kir1l9x/avito-test-assignment/internal/domain/repositories"
	"github.com/kir1l9x/avito-test-assignment/pkg/logger"
)

type UsersService struct {
	usersRepo repositories.UserRepository
}

func NewUsersService(usersRepo repositories.UserRepository) *UsersService {
	return &UsersService{usersRepo: usersRepo}
}

func (us *UsersService) SetUserIsActive(ctx context.Context, userID string, isActive bool) (*models.User, error) {
	logg := logger.FromContext(ctx)
	user, err := us.usersRepo.SetIsActive(ctx, userID, isActive)
	if err != nil {
		logg.Error("failed setting user isActive", zap.Error(err), zap.String("user_id", userID))
		return nil, err
	}

	return user, nil
}
