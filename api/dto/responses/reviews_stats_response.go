package responses

import "github.com/kir1l9x/avito-test-assignment/api/dto/responses/responseObjs"

type ReviewsStatsResponse struct {
	Stats []responseObjs.UserReviewStat `json:"stats"`
}
