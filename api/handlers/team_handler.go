package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/kir1l9x/avito-test-assignment/api/dto/requests"
	"github.com/kir1l9x/avito-test-assignment/api/dto/responses"
	"github.com/kir1l9x/avito-test-assignment/api/dto/responses/responseObjs"
	"github.com/kir1l9x/avito-test-assignment/api/middleware"
	"github.com/kir1l9x/avito-test-assignment/internal/appplication/services"
	"github.com/kir1l9x/avito-test-assignment/internal/domain/models"
)

type TeamHandler struct {
	teamService *services.TeamService
}

func NewTeamHandler(teamService *services.TeamService) *TeamHandler {
	return &TeamHandler{teamService: teamService}
}

func (th *TeamHandler) AddTeam(c *gin.Context) {
	var req requests.AddTeamRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.WriteValidationError(c, err)
		return
	}

	domainMembers, err := requests.RequestMembersToDomain(req.Members)
	if err != nil {
		middleware.WriteError(c, err)
		return
	}

	team, err := models.NewTeam(req.TeamName, domainMembers)
	if err != nil {
		middleware.WriteError(c, err)
		return
	}

	createdTeam, err := th.teamService.CreateTeamWithMembers(c.Request.Context(), team)
	if err != nil {
		middleware.WriteError(c, err)
		return
	}

	resp := responses.AddTeamResponse{
		Team: responseObjs.FromDomainTeam(createdTeam),
	}

	c.JSON(http.StatusCreated, resp)
}

func (th *TeamHandler) GetTeam(c *gin.Context) {
	var req requests.GetTeamRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		middleware.WriteValidationError(c, err)
		return
	}

	team, err := th.teamService.GetTeamByName(c.Request.Context(), req.TeamName)
	if err != nil {
		middleware.WriteError(c, err)
		return
	}

	resp := responseObjs.FromDomainTeam(team)
	c.JSON(http.StatusOK, resp)
}
