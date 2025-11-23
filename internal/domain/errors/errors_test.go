package errors

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDomainError_Error(t *testing.T) {
	e := New(CodePRNotFound, ErrPRNotFound)
	require.Contains(t, e.Error(), string(CodePRNotFound))
	require.Contains(t, e.Error(), ErrPRNotFound.Error())
}

func TestDomainError_Unwrap(t *testing.T) {
	e := New(CodeUserNotFound, ErrUserNotFound)

	require.Equal(t, ErrUserNotFound, errors.Unwrap(e))
}

func TestNewDomainError(t *testing.T) {
	e := New(CodeInternal, ErrTeamExists)

	require.Equal(t, CodeInternal, e.Code)
	require.Equal(t, ErrTeamExists, e.Err)
}
