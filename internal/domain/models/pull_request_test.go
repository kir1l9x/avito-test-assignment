package models

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/kir1l9x/avito-test-assignment/internal/domain/errors"
	"github.com/kir1l9x/avito-test-assignment/internal/domain/valueObjects"
)

func TestNewPullRequest_Success(t *testing.T) {
	pr, err := NewPullRequest("pr1", "u1", "Test PR", valueObjects.OpenPullRequest, []string{"u2"})
	require.NoError(t, err)
	require.NotNil(t, pr)

	require.Equal(t, "pr1", pr.ID)
	require.Equal(t, "u1", pr.AuthorID)
	require.Equal(t, "Test PR", pr.Name)
	require.Equal(t, valueObjects.OpenPullRequest, pr.Status)
	require.Len(t, pr.Reviewers, 1)
	require.NotZero(t, pr.CreatedAt)
}

func TestNewPullRequest_ValidationError(t *testing.T) {
	pr, err := NewPullRequest("", "u1", "name", valueObjects.OpenPullRequest, nil)
	require.Nil(t, pr)
	require.Error(t, err)

	var derr *errors.DomainError
	require.ErrorAs(t, err, &derr)
	require.Equal(t, errors.CodeInternal, derr.Code)
}

func TestPullRequest_TooManyReviewers(t *testing.T) {
	_, err := NewPullRequest("pr1", "u1", "Test", valueObjects.OpenPullRequest, []string{"u2", "u3", "u4"})
	require.Error(t, err)

	var derr *errors.DomainError
	require.ErrorAs(t, err, &derr)
}

func TestPullRequest_MergedAt_Ok(t *testing.T) {
	now := time.Now()
	pr := &PullRequest{
		ID:        "pr1",
		AuthorID:  "u1",
		Name:      "test",
		Status:    valueObjects.OpenPullRequest,
		Reviewers: nil,
		CreatedAt: now,
		MergedAt:  nil,
	}
	err := pr.validate()
	require.NoError(t, err)
}
