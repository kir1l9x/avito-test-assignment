package repositories

import (
	"context"

	"github.com/kir1l9x/avito-test-assignment/internal/domain/models"
)

type TeamsRepository interface {
	Create(ctx context.Context, team *models.Team) (*models.Team, error)
	GetByName(ctx context.Context, name string) (*models.Team, error)
}
