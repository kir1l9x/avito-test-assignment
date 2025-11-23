package services

import (
	"context"
	"fmt"
	"time"

	"github.com/avito-tech/go-transaction-manager/trm/v2"
	"go.uber.org/zap"

	"github.com/kir1l9x/avito-test-assignment/internal/appplication/utils/selector"
	domainErrors "github.com/kir1l9x/avito-test-assignment/internal/domain/errors"
	"github.com/kir1l9x/avito-test-assignment/internal/domain/models"
	"github.com/kir1l9x/avito-test-assignment/internal/domain/repositories"
	"github.com/kir1l9x/avito-test-assignment/internal/domain/valueObjects"
	"github.com/kir1l9x/avito-test-assignment/pkg/logger"
)

type PullRequestsService struct {
	prRepo   repositories.PullRequestsRepository
	userRepo repositories.UserRepository
	trM      trm.Manager
}

func NewPullRequestService(
	prRepo repositories.PullRequestsRepository,
	userRepo repositories.UserRepository,
	trM trm.Manager,
) *PullRequestsService {
	return &PullRequestsService{prRepo: prRepo, userRepo: userRepo, trM: trM}
}

func (prs *PullRequestsService) CreatePullRequest(
	ctx context.Context,
	prID, prName, authorID string,
) (*models.PullRequest, error) {
	logg := logger.FromContext(ctx)

	author, err := prs.userRepo.GetByID(ctx, authorID)
	if err != nil {
		logg.Error("failed to get author from database", zap.Error(err))
		return nil, err
	}

	authorTeamActives, err := prs.userRepo.GetActiveUsersInTeam(ctx, author.TeamName)
	if err != nil {
		logg.Error("failed to get author team actives from database", zap.Error(err))
		return nil, err
	}

	membersToReview, err := selector.PickReviewersFromUsers(authorTeamActives, author.ID, 2)
	if err != nil {
		logg.Error("failed to pick users to review", zap.Error(err))
		return nil, err
	}

	pullRequest, err := models.NewPullRequest(prID, author.ID, prName, valueObjects.OpenPullRequest, membersToReview)
	if err != nil {
		logg.Error("failed to create new pull request", zap.Error(err))
		return nil, domainErrors.New(domainErrors.CodeInternal, fmt.Errorf("failed to create new pull request: %w", err))
	}

	var createdPullRequest *models.PullRequest

	err = prs.trM.Do(ctx, func(ctx context.Context) error {
		cpr, err := prs.prRepo.Create(ctx, pullRequest)
		if err != nil {
			logg.Error("failed to create new pull request", zap.Error(err))
			return err
		}

		createdPullRequest = cpr

		return nil
	})
	if err != nil {
		logg.Error("failed to create new pull request", zap.Error(err))
		return nil, err
	}

	return createdPullRequest, nil
}

func (prs *PullRequestsService) GetPRsWhereUserIsReviewer(ctx context.Context, userID string) ([]models.PullRequest, error) {
	logg := logger.FromContext(ctx)

	pullRequests, err := prs.prRepo.GetListByUser(ctx, userID)
	if err != nil {
		logg.Error("failed to get pull requests from database", zap.Error(err))
		return nil, err
	}

	return pullRequests, nil
}

func (prs *PullRequestsService) Merge(ctx context.Context, prID string) (*models.PullRequest, error) {
	logg := logger.FromContext(ctx)
	prToMerge, err := prs.prRepo.GetByID(ctx, prID)
	if err != nil {
		logg.Error("failed to get pull request from database", zap.Error(err))
		return nil, err
	}

	if prToMerge.Status == valueObjects.MergedPullRequest {
		return prToMerge, nil
	}

	prToMerge.Status = valueObjects.MergedPullRequest
	now := time.Now().UTC()
	prToMerge.MergedAt = &now

	var updatedPr *models.PullRequest

	err = prs.trM.Do(ctx, func(ctx context.Context) error {
		upr, err := prs.prRepo.Update(ctx, prToMerge)
		if err != nil {
			logg.Error("failed to update pull request", zap.Error(err))
			return err
		}

		updatedPr = upr

		return nil
	})
	if err != nil {
		logg.Error("failed to update pull request", zap.Error(err))
		return nil, err
	}

	return updatedPr, nil
}

func (prs *PullRequestsService) ReassignReviewer(ctx context.Context, prID string, reviewerID string) (*models.PullRequest, string, error) {
	logg := logger.FromContext(ctx)

	pr, err := prs.prRepo.GetByID(ctx, prID)
	if err != nil {
		logg.Error("failed to get pull request from database", zap.Error(err))
		return nil, "", err
	}

	if pr.Status == valueObjects.MergedPullRequest {
		return nil, "", domainErrors.New(domainErrors.CodePRMerged, domainErrors.ErrPRMerged)
	}

	author, err := prs.userRepo.GetByID(ctx, pr.AuthorID)
	if err != nil {
		logg.Error("failed to get author from database", zap.Error(err))
		return nil, "", err
	}

	activeUsers, err := prs.userRepo.GetActiveUsersInTeam(ctx, author.TeamName)
	if err != nil {
		logg.Error("failed to get author team actives from database", zap.Error(err))
		return nil, "", err
	}

	newUserToReview, err := selector.PickNewReviewer(activeUsers, pr.Reviewers, author.ID, reviewerID)
	if err != nil {
		logg.Error("failed to pick users to review", zap.Error(err))
		return nil, "", err
	}

	for i, reviewer := range pr.Reviewers {
		if reviewer == reviewerID {
			pr.Reviewers[i] = newUserToReview
		}
	}

	var updatedPr *models.PullRequest

	err = prs.trM.Do(ctx, func(ctx context.Context) error {
		upr, err := prs.prRepo.Update(ctx, pr)
		if err != nil {
			logg.Error("failed to update pull request", zap.Error(err))
			return err
		}

		updatedPr = upr

		return nil
	})
	if err != nil {
		logg.Error("failed to update pull request", zap.Error(err))
		return nil, "", err
	}

	return updatedPr, newUserToReview, nil
}

func (prs *PullRequestsService) GetReviewsAmountPerUser(ctx context.Context) ([]models.UserReviewStats, error) {
	logg := logger.FromContext(ctx)

	stats, err := prs.prRepo.GetReviewsAmountByUser(ctx)
	if err != nil {
		logg.Error("failed to get reviews amount per user", zap.Error(err))
		return nil, err
	}

	return stats, nil
}

func (prs *PullRequestsService) GetTop5ReviewersAllTime(ctx context.Context) ([]models.UserReviewStats, error) {
	logg := logger.FromContext(ctx)

	stats, err := prs.prRepo.GetTop5ReviewsAmountByUser(ctx)
	if err != nil {
		logg.Error("failed to get top 5 reviewers all time", zap.Error(err))
		return nil, err
	}

	return stats, nil
}

func (prs *PullRequestsService) GetOpenReviewsAmountPerUser(ctx context.Context) ([]models.UserReviewStats, error) {
	logg := logger.FromContext(ctx)

	stats, err := prs.prRepo.GetOpenReviewsAmountByUser(ctx)
	if err != nil {
		logg.Error("failed to get open reviews amount per user", zap.Error(err))
		return nil, err
	}

	return stats, nil
}

func (prs *PullRequestsService) GetTop5ReviewersCurrent(ctx context.Context) ([]models.UserReviewStats, error) {
	logg := logger.FromContext(ctx)

	stats, err := prs.prRepo.GetTop5OpenReviewsAmountByUser(ctx)
	if err != nil {
		logg.Error("failed to get top 5 reviewers for open prs", zap.Error(err))
		return nil, err
	}

	return stats, nil
}
