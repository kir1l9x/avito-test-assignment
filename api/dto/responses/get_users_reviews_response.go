package responses

import "github.com/kir1l9x/avito-test-assignment/api/dto/responses/responseObjs"

type GetUsersReviewsResponse struct {
	UserID       string                          `json:"user_id"`
	PullRequests []responseObjs.PullRequestShort `json:"pull_requests"`
}
