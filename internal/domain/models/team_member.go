package models

import (
	"fmt"

	domainErrors "github.com/kir1l9x/avito-test-assignment/internal/domain/errors"
	"github.com/kir1l9x/avito-test-assignment/pkg/validators"
)

type TeamMember struct {
	ID       string `validate:"required"`
	Username string `validate:"required"`
	IsActive bool
}

func NewTeamMember(id, username string, isActive bool) (*TeamMember, error) {
	teamMember := &TeamMember{
		ID:       id,
		Username: username,
		IsActive: isActive,
	}

	if err := teamMember.validate(); err != nil {
		return nil, err
	}

	return teamMember, nil
}

func (tm *TeamMember) validate() error {
	if err := validators.Validator.Struct(tm); err != nil {
		return domainErrors.New(
			domainErrors.CodeInternal,
			fmt.Errorf("validation error on domain team member model: %w", err),
		)
	}

	return nil
}
