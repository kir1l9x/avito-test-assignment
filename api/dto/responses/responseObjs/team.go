package responseObjs

import "github.com/kir1l9x/avito-test-assignment/internal/domain/models"

type TeamMember struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	IsActive bool   `json:"is_active"`
}

type Team struct {
	TeamName string       `json:"team_name"`
	Members  []TeamMember `json:"members"`
}

func FromDomainTeam(t *models.Team) Team {
	members := make([]TeamMember, 0, len(t.Members))
	for _, m := range t.Members {
		members = append(members, TeamMember{
			UserID:   m.ID,
			Username: m.Username,
			IsActive: m.IsActive,
		})
	}

	return Team{
		TeamName: t.Name,
		Members:  members,
	}
}
