package mappers

import (
	"fmt"

	"github.com/kir1l9x/avito-test-assignment/internal/domain/models"
)

func fromUserToMember(user *models.User) (*models.TeamMember, error) {
	member, err := models.NewTeamMember(user.ID, user.Username, user.IsActive)
	if err != nil {
		return nil, fmt.Errorf("failed to map user to member: %w", err)
	}

	return member, nil
}

func FromUsersToMembers(users []models.User) ([]*models.TeamMember, error) {
	members := make([]*models.TeamMember, len(users))
	for i, user := range users {
		member, err := fromUserToMember(&user)
		if err != nil {
			return nil, fmt.Errorf("failed to map user to member: %w", err)
		}

		members[i] = member
	}

	return members, nil
}

func fromMemberToUser(member *models.TeamMember, teamName string) (*models.User, error) {
	user, err := models.NewUser(member.ID, member.Username, teamName, member.IsActive)
	if err != nil {
		return nil, fmt.Errorf("failed to map member to user: %w", err)
	}

	return user, nil
}

func FromMembersToUsers(members []*models.TeamMember, teamName string) ([]models.User, error) {
	users := make([]models.User, len(members))
	for i, member := range members {
		user, err := fromMemberToUser(member, teamName)
		if err != nil {
			return nil, fmt.Errorf("failed to map member to user: %w", err)
		}

		users[i] = *user
	}

	return users, nil
}
