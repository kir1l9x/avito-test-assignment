package models

import (
	"fmt"

	domainErrors "github.com/kir1l9x/avito-test-assignment/internal/domain/errors"
	"github.com/kir1l9x/avito-test-assignment/pkg/validators"
)

type Team struct {
	Name    string        `validate:"required"`
	Members []*TeamMember `validate:"required,min=1,dive,required"`
}

func NewTeam(name string, members []*TeamMember) (*Team, error) {
	team := &Team{
		Name:    name,
		Members: members,
	}

	if err := team.validate(); err != nil {
		return nil, err
	}

	return team, nil
}

func (t *Team) validate() error {
	if err := validators.Validator.Struct(t); err != nil {
		return domainErrors.New(domainErrors.CodeInternal, fmt.Errorf("validation error on domain team model: %w", err))
	}

	return nil
}
