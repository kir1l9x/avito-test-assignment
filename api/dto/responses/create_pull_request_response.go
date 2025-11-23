package responses

import "github.com/kir1l9x/avito-test-assignment/api/dto/responses/responseObjs"

type CreatePullRequestResponse struct {
	PR responseObjs.PullRequest `json:"pr"`
}
