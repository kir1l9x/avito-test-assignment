package requests

type GetUsersReviewsRequest struct {
	UserID string `form:"user_id" binding:"required,user_id"`
}
