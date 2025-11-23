package requests

type GetTeamRequest struct {
	TeamName string `form:"team_name" binding:"required,ascii"`
}
