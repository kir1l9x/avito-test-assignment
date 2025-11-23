package models

type UserReviewStats struct {
	UserID      string `validate:"required"`
	ReviewCount int64
}
