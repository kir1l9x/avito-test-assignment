package mappers

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/kir1l9x/avito-test-assignment/internal/domain/valueObjects"
	"github.com/kir1l9x/avito-test-assignment/internal/infrastructure/dbModels"
)

func TestPullRequestStatusesMapper_ToDomain_Success(t *testing.T) {
	mapper := NewPullRequestMapper(
		map[int16]string{1: "OPEN", 2: "MERGED"},
		map[string]int16{},
	)

	s, err := mapper.ToDomain(1)
	require.NoError(t, err)
	require.Equal(t, "OPEN", s)
}

func TestPullRequestStatusesMapper_ToDomain_Error(t *testing.T) {
	mapper := NewPullRequestMapper(
		map[int16]string{1: "OPEN"},
		map[string]int16{},
	)

	s, err := mapper.ToDomain(999)
	require.Error(t, err)
	require.Empty(t, s)
}

func TestPullRequestStatusesMapper_FromDomain_Success(t *testing.T) {
	mapper := NewPullRequestMapper(
		map[int16]string{},
		map[string]int16{"OPEN": 1, "MERGED": 2},
	)

	id, err := mapper.FromDomain("MERGED")
	require.NoError(t, err)
	require.Equal(t, int16(2), id)
}

func TestPullRequestStatusesMapper_FromDomain_Error(t *testing.T) {
	mapper := NewPullRequestMapper(
		map[int16]string{},
		map[string]int16{"OPEN": 1},
	)

	id, err := mapper.FromDomain("UNKNOWN")
	require.Error(t, err)
	require.Equal(t, int16(0), id)
}

func TestPullRequestStatusesMapper_ToDomainModel_Success(t *testing.T) {
	mapper := NewPullRequestMapper(
		map[int16]string{1: "OPEN"},
		map[string]int16{},
	)

	now := time.Now()
	dbModel := &dbModels.PullRequest{
		ID:        "pr-1",
		Name:      "Init",
		AuthorID:  "u1",
		StatusID:  1,
		Reviewers: []string{"u2", "u3"},
		CreatedAt: now,
		MergedAt:  nil,
	}

	domainPR, err := mapper.ToDomainModel(dbModel)
	require.NoError(t, err)

	require.Equal(t, "pr-1", domainPR.ID)
	require.Equal(t, "Init", domainPR.Name)
	require.Equal(t, "u1", domainPR.AuthorID)
	require.Equal(t, valueObjects.OpenPullRequest, domainPR.Status)
	require.Equal(t, dbModel.Reviewers, domainPR.Reviewers)
	require.Equal(t, now, domainPR.CreatedAt)
	require.Nil(t, domainPR.MergedAt)
}

func TestPullRequestStatusesMapper_ToDomainModel_Error(t *testing.T) {
	mapper := NewPullRequestMapper(
		map[int16]string{},
		map[string]int16{},
	)

	dbModel := &dbModels.PullRequest{
		ID:       "pr-1",
		StatusID: 99,
	}

	_, err := mapper.ToDomainModel(dbModel)
	require.Error(t, err)
}
