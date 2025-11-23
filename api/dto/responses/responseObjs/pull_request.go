package responseObjs

import (
	"time"

	"github.com/kir1l9x/avito-test-assignment/internal/domain/models"
)

type PullRequest struct {
	PullRequestID     string   `json:"pull_request_id"`
	PullRequestName   string   `json:"pull_request_name"`
	AuthorID          string   `json:"author_id"`
	Status            string   `json:"status"`
	AssignedReviewers []string `json:"assigned_reviewers"`
	CreatedAt         string   `json:"createdAt,omitempty"`
	MergedAt          *string  `json:"mergedAt,omitempty"`
}

type PullRequestShort struct {
	PullRequestID   string `json:"pull_request_id"`
	PullRequestName string `json:"pull_request_name"`
	AuthorID        string `json:"author_id"`
	Status          string `json:"status"`
}

func FromDomainPullRequest(pr *models.PullRequest) PullRequest {
	created := pr.CreatedAt.UTC().Format(time.RFC3339)

	var merged *string
	if pr.MergedAt != nil {
		m := pr.MergedAt.UTC().Format(time.RFC3339)
		merged = &m
	}

	return PullRequest{
		PullRequestID:     pr.ID,
		PullRequestName:   pr.Name,
		AuthorID:          pr.AuthorID,
		Status:            pr.Status.String(),
		AssignedReviewers: pr.Reviewers,
		CreatedAt:         created,
		MergedAt:          merged,
	}
}

func FromDomainPullRequestsShort(prs []models.PullRequest) []PullRequestShort {
	prsShort := make([]PullRequestShort, 0, len(prs))
	for _, pr := range prs {
		prsShort = append(prsShort, PullRequestShort{
			PullRequestID:   pr.ID,
			PullRequestName: pr.Name,
			AuthorID:        pr.AuthorID,
			Status:          pr.Status.String(),
		})
	}

	return prsShort
}
