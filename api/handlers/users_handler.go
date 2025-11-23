package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/kir1l9x/avito-test-assignment/api/dto/requests"
	"github.com/kir1l9x/avito-test-assignment/api/dto/responses"
	"github.com/kir1l9x/avito-test-assignment/api/dto/responses/responseObjs"
	"github.com/kir1l9x/avito-test-assignment/api/middleware"
	"github.com/kir1l9x/avito-test-assignment/internal/appplication/services"
)

type UsersHandler struct {
	usersService *services.UsersService
}

func NewUsersHandler(usersService *services.UsersService) *UsersHandler {
	return &UsersHandler{usersService: usersService}
}

func (uh *UsersHandler) SetIsActive(c *gin.Context) {
	var req requests.SetIsActiveRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.WriteValidationError(c, err)
		return
	}

	user, err := uh.usersService.SetUserIsActive(c.Request.Context(), req.UserID, req.IsActive)
	if err != nil {
		middleware.WriteError(c, err)
		return
	}

	resp := responses.SetIsActiveResponse{
		User: responseObjs.FromDomainUser(user),
	}

	c.JSON(http.StatusOK, resp)
}
