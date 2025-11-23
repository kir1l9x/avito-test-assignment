package responses

import "github.com/kir1l9x/avito-test-assignment/api/dto/responses/responseObjs"

type SetIsActiveResponse struct {
	User responseObjs.User `json:"user"`
}
