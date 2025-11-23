package services

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/kir1l9x/avito-test-assignment/internal/domain/models"
)

func TestUsersService_SetUserIsActive_Success(t *testing.T) {
	ctx := context.Background()

	userRepo := new(mockUserRepository)
	svc := NewUsersService(userRepo)

	expected := &models.User{
		ID:       "u1",
		Username: "Alice",
		TeamName: "backend",
		IsActive: true,
	}

	userRepo.On("SetIsActive", mock.Anything, "u1", true).
		Return(expected, nil)

	res, err := svc.SetUserIsActive(ctx, "u1", true)
	require.NoError(t, err)
	require.Equal(t, expected, res)
}

func TestUsersService_SetUserIsActive_Error(t *testing.T) {
	ctx := context.Background()

	userRepo := new(mockUserRepository)
	svc := NewUsersService(userRepo)

	userRepo.On("SetIsActive", mock.Anything, "u1", true).
		Return((*models.User)(nil), errors.New("db error"))

	res, err := svc.SetUserIsActive(ctx, "u1", true)
	require.Error(t, err)
	require.Nil(t, res)
}
