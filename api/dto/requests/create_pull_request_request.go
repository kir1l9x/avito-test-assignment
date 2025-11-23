package requests

type CreatePullRequestRequest struct {
	PullRequestID   string `json:"pull_request_id"   binding:"required,pr_id"`
	PullRequestName string `json:"pull_request_name" binding:"required,ascii"`
	AuthorID        string `json:"author_id"         binding:"required,user_id"`
}
