package repositories

import (
	"context"

	"github.com/kir1l9x/avito-test-assignment/internal/domain/models"
)

type UserRepository interface {
	CreateOrUpdate(ctx context.Context, user *models.User) (*models.User, error)
	GetByID(ctx context.Context, id string) (*models.User, error)
	GetUsersByTeam(ctx context.Context, teamName string) ([]models.User, error)
	GetActiveUsersInTeam(ctx context.Context, teamName string) ([]models.User, error)
	SetIsActive(ctx context.Context, id string, isActive bool) (*models.User, error)
}
