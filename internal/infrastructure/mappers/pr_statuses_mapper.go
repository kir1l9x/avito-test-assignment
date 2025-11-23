package mappers

import (
	"fmt"

	"github.com/kir1l9x/avito-test-assignment/internal/domain/models"
	"github.com/kir1l9x/avito-test-assignment/internal/domain/valueObjects"
	"github.com/kir1l9x/avito-test-assignment/internal/infrastructure/dbModels"
)

type PullRequestStatusesMapper struct {
	IDToStatus map[int16]string
	StatusToID map[string]int16
}

func NewPullRequestMapper(idToStatus map[int16]string, statusToID map[string]int16) *PullRequestStatusesMapper {
	return &PullRequestStatusesMapper{IDToStatus: idToStatus, StatusToID: statusToID}
}

func (mapper *PullRequestStatusesMapper) ToDomain(id int16) (string, error) {
	s, ok := mapper.IDToStatus[id]
	if !ok {
		return "", fmt.Errorf("failed to map to domain status: %d", id)
	}

	return s, nil
}

func (mapper *PullRequestStatusesMapper) FromDomain(status string) (int16, error) {
	id, ok := mapper.StatusToID[status]
	if !ok {
		return 0, fmt.Errorf("failed to map to db status: %s", status)
	}

	return id, nil
}

func (mapper *PullRequestStatusesMapper) ToDomainModel(dbModelPR *dbModels.PullRequest) (*models.PullRequest, error) {
	statusStr, err := mapper.ToDomain(dbModelPR.StatusID)
	if err != nil {
		return nil, fmt.Errorf("failed to map to domain pr model: %w", err)
	}

	return &models.PullRequest{
		ID:        dbModelPR.ID,
		AuthorID:  dbModelPR.AuthorID,
		Name:      dbModelPR.Name,
		Status:    valueObjects.PullRequestStatus(statusStr),
		Reviewers: dbModelPR.Reviewers,
		CreatedAt: dbModelPR.CreatedAt,
		MergedAt:  dbModelPR.MergedAt,
	}, nil
}
