package responses

import "github.com/kir1l9x/avito-test-assignment/api/dto/responses/responseObjs"

type ReassignPullRequestResponse struct {
	PR         responseObjs.PullRequest `json:"pr"`
	ReplacedBy string                   `json:"replaced_by"`
}
