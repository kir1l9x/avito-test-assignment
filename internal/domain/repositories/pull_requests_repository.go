package repositories

import (
	"context"

	"github.com/kir1l9x/avito-test-assignment/internal/domain/models"
)

type PullRequestsRepository interface {
	Create(ctx context.Context, pullRequest *models.PullRequest) (*models.PullRequest, error)
	GetByID(ctx context.Context, pullRequestID string) (*models.PullRequest, error)
	Update(ctx context.Context, pullRequest *models.PullRequest) (*models.PullRequest, error)
	GetListByUser(ctx context.Context, userID string) ([]models.PullRequest, error)
	GetReviewsAmountByUser(ctx context.Context) ([]models.UserReviewStats, error)
	GetTop5ReviewsAmountByUser(ctx context.Context) ([]models.UserReviewStats, error)
	GetOpenReviewsAmountByUser(ctx context.Context) ([]models.UserReviewStats, error)
	GetTop5OpenReviewsAmountByUser(ctx context.Context) ([]models.UserReviewStats, error)
}
