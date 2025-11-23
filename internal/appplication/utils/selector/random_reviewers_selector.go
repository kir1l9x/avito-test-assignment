package selector

import (
	"math/rand"
	"time"

	domainErrors "github.com/kir1l9x/avito-test-assignment/internal/domain/errors"
	"github.com/kir1l9x/avito-test-assignment/internal/domain/models"
)

func PickReviewersFromUsers(users []models.User, authorID string, required int) ([]string, error) {
	candidates := make([]models.User, 0, len(users)-1)
	for _, u := range users {
		if u.ID != authorID {
			candidates = append(candidates, u)
		}
	}

	if len(candidates) == 0 {
		return []string{}, nil
	}

	if len(candidates) <= required {
		ids := make([]string, 0, len(candidates))
		for _, c := range candidates {
			ids = append(ids, c.ID)
		}
		return ids, nil
	}

	rnd := rand.New(rand.NewSource(time.Now().UnixNano())) //nolint:gosec // pseudo-random is enough for reviewer selection
	for i := len(candidates) - 1; i > 0; i-- {
		j := rnd.Intn(i + 1)
		candidates[i], candidates[j] = candidates[j], candidates[i]
	}

	result := make([]string, required)
	for i := 0; i < required; i++ {
		result[i] = candidates[i].ID
	}

	return result, nil
}

func PickNewReviewer(
	activeUsers []models.User,
	existingReviewers []string,
	authorID string,
	oldReviewerID string,
) (string, error) {
	existing := make(map[string]struct{}, len(existingReviewers))
	for _, r := range existingReviewers {
		existing[r] = struct{}{}
	}

	if _, ok := existing[oldReviewerID]; !ok {
		return "", domainErrors.New(domainErrors.CodeNotAssigned, domainErrors.ErrReviewerNotAssigned)
	}

	candidates := make([]string, 0, len(activeUsers))
	for _, u := range activeUsers {
		if u.ID == authorID {
			continue
		}
		if u.ID == oldReviewerID {
			continue
		}
		if _, ok := existing[u.ID]; ok {
			continue
		}

		candidates = append(candidates, u.ID)
	}

	if len(candidates) == 0 {
		return "", domainErrors.New(domainErrors.CodeNoCandidate, domainErrors.ErrNoReplacement)
	}

	rnd := rand.New(rand.NewSource(time.Now().UnixNano())) //nolint:gosec // pseudo-random is enough for reviewer selection
	idx := rnd.Intn(len(candidates))

	return candidates[idx], nil
}
