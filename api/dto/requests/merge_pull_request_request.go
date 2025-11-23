package requests

type MergePullRequestRequest struct {
	PullRequestID string `json:"pull_request_id" binding:"required,pr_id"`
}
