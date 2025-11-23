package models

import (
	"fmt"
	"time"

	domainErrors "github.com/kir1l9x/avito-test-assignment/internal/domain/errors"
	"github.com/kir1l9x/avito-test-assignment/internal/domain/valueObjects"
	"github.com/kir1l9x/avito-test-assignment/pkg/validators"
)

type PullRequest struct {
	ID        string                         `validate:"required"`
	AuthorID  string                         `validate:"required"`
	Name      string                         `validate:"required"`
	Status    valueObjects.PullRequestStatus `validate:"required"`
	Reviewers []string                       `validate:"max=2"`
	CreatedAt time.Time                      `validate:"required"`
	MergedAt  *time.Time                     `validate:"omitempty"`
}

func NewPullRequest(
	id, authorID, name string,
	status valueObjects.PullRequestStatus,
	reviewers []string,
) (*PullRequest, error) {
	pullRequest := &PullRequest{
		ID:        id,
		AuthorID:  authorID,
		Name:      name,
		Status:    status,
		Reviewers: reviewers,
		CreatedAt: time.Now().UTC(),
	}

	if err := pullRequest.validate(); err != nil {
		return nil, err
	}

	return pullRequest, nil
}

func (pr *PullRequest) validate() error {
	if err := validators.Validator.Struct(pr); err != nil && pr.Status.IsValid() {
		return domainErrors.New(
			domainErrors.CodeInternal,
			fmt.Errorf("validation error on domain pull request model: %w", err),
		)
	}

	return nil
}
