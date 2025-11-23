package requests

type ReassignPullRequestRequest struct {
	PullRequestID string `json:"pull_request_id" binding:"required,pr_id"`
	OldUserID     string `json:"old_user_id"     binding:"required,user_id"`
}
