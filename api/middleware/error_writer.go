package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	apiResp "github.com/kir1l9x/avito-test-assignment/api/dto/responses"
	domainErrors "github.com/kir1l9x/avito-test-assignment/internal/domain/errors"
)

func WriteError(c *gin.Context, err error) {
	if err == nil {
		return
	}

	var de *domainErrors.DomainError
	if errors.As(err, &de) {
		code, status := mapDomainCodeToAPI(de.Code)
		c.JSON(status, apiResp.ErrorResponse{
			Error: apiResp.ErrorBody{
				Code:    code,
				Message: de.Err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusInternalServerError, apiResp.ErrorResponse{
		Error: apiResp.ErrorBody{
			Code:    "INTERNAL",
			Message: "internal error",
		},
	})
}

func WriteValidationError(c *gin.Context, err error) {
	if err == nil {
		return
	}

	var gerr *gin.Error
	if errors.As(err, &gerr) {
		c.JSON(http.StatusBadRequest, apiResp.ErrorResponse{
			Error: apiResp.ErrorBody{
				Code:    "BAD_REQUEST",
				Message: err.Error(),
			},
		})
		return
	}

	var verrs validator.ValidationErrors
	if ok := errors.As(err, &verrs); ok {
		var messages []string
		for _, fe := range verrs {
			field := fe.Field()
			tag := fe.Tag()
			param := fe.Param()
			msg := fmt.Sprintf("%s failed on '%s' %s", field, tag, param)
			messages = append(messages, msg)
		}

		c.JSON(http.StatusBadRequest, apiResp.ErrorResponse{
			Error: apiResp.ErrorBody{
				Code:    "BAD_REQUEST",
				Message: strings.Join(messages, "; "),
			},
		})
		return
	}

	c.JSON(http.StatusBadRequest, apiResp.ErrorResponse{
		Error: apiResp.ErrorBody{
			Code:    "BAD_REQUEST",
			Message: err.Error(),
		},
	})
}

func mapDomainCodeToAPI(code domainErrors.Code) (string, int) {
	switch code {

	case domainErrors.CodeTeamExists:
		return "TEAM_EXISTS", http.StatusBadRequest

	case domainErrors.CodePRExists:
		return "PR_EXISTS", http.StatusConflict

	case domainErrors.CodePRMerged:
		return "PR_MERGED", http.StatusConflict

	case domainErrors.CodeNotAssigned:
		return "NOT_ASSIGNED", http.StatusConflict

	case domainErrors.CodeNoCandidate:
		return "NO_CANDIDATE", http.StatusConflict

	case domainErrors.CodeTeamNotFound,
		domainErrors.CodeUserNotFound,
		domainErrors.CodePRNotFound:
		return "NOT_FOUND", http.StatusNotFound

	case domainErrors.CodeInternal:
		return "INTERNAL", http.StatusInternalServerError

	default:
		return "INTERNAL", http.StatusInternalServerError
	}
}
