package responseObjs

import "github.com/kir1l9x/avito-test-assignment/internal/domain/models"

type UserReviewStat struct {
	UserID      string `json:"user_id"`
	ReviewCount int64  `json:"review_count"`
}

func FromDomainUserReviewStats(stats []models.UserReviewStats) []UserReviewStat {
	result := make([]UserReviewStat, 0, len(stats))
	for _, stat := range stats {
		result = append(result, UserReviewStat{
			UserID:      stat.UserID,
			ReviewCount: stat.ReviewCount,
		})
	}

	return result
}
