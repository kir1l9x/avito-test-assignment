package handlers_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/kir1l9x/avito-test-assignment/api/dto/requests"
	"github.com/kir1l9x/avito-test-assignment/api/dto/responses"
	"github.com/kir1l9x/avito-test-assignment/api/handlers/handlersTestsHelp"
)

func TestUsersHandler_SetIsActive_Success(t *testing.T) {
	srv, pool := handlersTestsHelp.SetupTestServer(t)
	defer srv.Close()

	ctx := context.Background()
	teamName := "team-users-set-active"
	userID := "u3001"

	_, err := pool.Exec(ctx, `DELETE FROM users WHERE id = $1`, userID)
	require.NoError(t, err)
	_, err = pool.Exec(ctx, `DELETE FROM teams WHERE name = $1`, teamName)
	require.NoError(t, err)

	_, err = pool.Exec(ctx, `INSERT INTO teams(name) VALUES ($1)`, teamName)
	require.NoError(t, err)
	_, err = pool.Exec(ctx, `
		INSERT INTO users(id, username, team_name, is_active)
		VALUES ($1, $2, $3, $4)
	`, userID, "john", teamName, false)
	require.NoError(t, err)

	resp := postJSON(t, srv.URL+"/users/setIsActive", requests.SetIsActiveRequest{
		UserID:   userID,
		IsActive: true,
	})
	require.Equal(t, http.StatusOK, resp.StatusCode)

	var parsed responses.SetIsActiveResponse
	decodeResponse(t, resp, &parsed)

	require.Equal(t, userID, parsed.User.UserID)
	require.True(t, parsed.User.IsActive)
}

func TestUsersHandler_SetIsActive_UserNotFound(t *testing.T) {
	srv, pool := handlersTestsHelp.SetupTestServer(t)
	defer srv.Close()

	ctx := context.Background()
	_, err := pool.Exec(ctx, `DELETE FROM users WHERE id = $1`, "u3999")
	require.NoError(t, err)

	resp := postJSON(t, srv.URL+"/users/setIsActive", requests.SetIsActiveRequest{
		UserID:   "u3999",
		IsActive: true,
	})
	require.Equal(t, http.StatusNotFound, resp.StatusCode)

	var parsed responses.ErrorResponse
	decodeResponse(t, resp, &parsed)

	require.Equal(t, "NOT_FOUND", parsed.Error.Code)
}
