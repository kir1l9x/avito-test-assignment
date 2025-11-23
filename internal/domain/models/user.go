package models

import (
	"fmt"

	domainErrors "github.com/kir1l9x/avito-test-assignment/internal/domain/errors"
	"github.com/kir1l9x/avito-test-assignment/pkg/validators"
)

type User struct {
	ID       string `validate:"required"`
	Username string `validate:"required"`
	TeamName string `validate:"required"`
	IsActive bool
}

func NewUser(userID, username, teamName string, isActive bool) (*User, error) {
	user := &User{
		ID:       userID,
		Username: username,
		TeamName: teamName,
		IsActive: isActive,
	}

	if err := user.validate(); err != nil {
		return nil, err
	}

	return user, nil
}

func (u *User) validate() error {
	if err := validators.Validator.Struct(u); err != nil {
		return domainErrors.New(domainErrors.CodeInternal, fmt.Errorf("validation error on domain user model: %w", err))
	}

	return nil
}
