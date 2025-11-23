package selector_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/kir1l9x/avito-test-assignment/internal/appplication/utils/selector"
	"github.com/kir1l9x/avito-test-assignment/internal/domain/errors"
	"github.com/kir1l9x/avito-test-assignment/internal/domain/models"
)

func TestPickReviewersFromUsers_NoCandidates(t *testing.T) {
	users := []models.User{
		{ID: "u1"},
	}

	ids, err := selector.PickReviewersFromUsers(users, "u1", 2)

	require.NoError(t, err)
	require.Len(t, ids, 0)
}

func TestPickReviewersFromUsers_OneCandidate(t *testing.T) {
	users := []models.User{
		{ID: "u1"},
		{ID: "u2"},
	}

	ids, err := selector.PickReviewersFromUsers(users, "u1", 2)

	require.NoError(t, err)
	require.Equal(t, []string{"u2"}, ids)
}

func TestPickReviewersFromUsers_TwoCandidates_RequiredTwo(t *testing.T) {
	users := []models.User{
		{ID: "u1"},
		{ID: "u2"},
		{ID: "u3"},
	}

	ids, err := selector.PickReviewersFromUsers(users, "u1", 2)

	require.NoError(t, err)
	require.Len(t, ids, 2)
	require.Contains(t, ids, "u2")
	require.Contains(t, ids, "u3")
}

func TestPickReviewersFromUsers_MoreCandidatesThanRequired(t *testing.T) {
	users := []models.User{
		{ID: "u1"},
		{ID: "u2"},
		{ID: "u3"},
		{ID: "u4"},
	}

	ids, err := selector.PickReviewersFromUsers(users, "u1", 2)
	require.NoError(t, err)
	require.Len(t, ids, 2)

	require.NotContains(t, ids, "u1")
}

func TestPickNewReviewer_UserNotAssigned(t *testing.T) {
	active := []models.User{
		{ID: "u1"},
		{ID: "u2"},
	}

	_, err := selector.PickNewReviewer(active, []string{"u2"}, "u1", "u3")
	require.Error(t, err)

	var de *errors.DomainError
	require.ErrorAs(t, err, &de)
	require.Equal(t, errors.CodeNotAssigned, de.Code)
}

func TestPickNewReviewer_NoCandidates(t *testing.T) {
	active := []models.User{
		{ID: "u1"},
		{ID: "u2"},
	}

	existing := []string{"u2"}

	_, err := selector.PickNewReviewer(active, existing, "u1", "u2")
	require.Error(t, err)

	var de *errors.DomainError
	require.ErrorAs(t, err, &de)
	require.Equal(t, errors.CodeNoCandidate, de.Code)
}

func TestPickNewReviewer_Success(t *testing.T) {
	active := []models.User{
		{ID: "u1"},
		{ID: "u2"},
		{ID: "u3"},
		{ID: "u4"},
	}

	existing := []string{"u2", "u3"}

	newReviewer, err := selector.PickNewReviewer(active, existing, "u1", "u2")

	require.NoError(t, err)
	require.Equal(t, "u4", newReviewer)
}
