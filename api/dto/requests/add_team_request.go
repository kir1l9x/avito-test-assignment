package requests

import (
	"fmt"

	domainErrors "github.com/kir1l9x/avito-test-assignment/internal/domain/errors"
	"github.com/kir1l9x/avito-test-assignment/internal/domain/models"
)

type TeamMemberRequest struct {
	UserID   string `json:"user_id" binding:"required,user_id"`
	Username string `json:"username" binding:"required,ascii"`
	IsActive bool   `json:"is_active"`
}

type AddTeamRequest struct {
	TeamName string              `json:"team_name" binding:"required,ascii"`
	Members  []TeamMemberRequest `json:"members" binding:"required,min=1,dive"`
}

func RequestMembersToDomain(reqMembers []TeamMemberRequest) ([]*models.TeamMember, error) {
	domainMembers := make([]*models.TeamMember, 0, len(reqMembers))
	for _, reqMember := range reqMembers {
		domainMember, err := models.NewTeamMember(reqMember.UserID, reqMember.Username, reqMember.IsActive)
		if err != nil {
			return nil, domainErrors.New(
				domainErrors.CodeInternal,
				fmt.Errorf("failed to map request members to domain: %w", err),
			)
		}

		domainMembers = append(domainMembers, domainMember)
	}

	return domainMembers, nil
}
