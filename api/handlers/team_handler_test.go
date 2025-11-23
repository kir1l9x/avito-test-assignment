package handlers_test

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/kir1l9x/avito-test-assignment/api/dto/requests"
	"github.com/kir1l9x/avito-test-assignment/api/dto/responses"
	"github.com/kir1l9x/avito-test-assignment/api/dto/responses/responseObjs"
	"github.com/kir1l9x/avito-test-assignment/api/handlers/handlersTestsHelp"
)

func TestTeamHandler_AddTeam_Success(t *testing.T) {
	srv, _ := handlersTestsHelp.SetupTestServer(t)
	defer srv.Close()

	req := requests.AddTeamRequest{
		TeamName: "team-add-success",
		Members: []requests.TeamMemberRequest{
			{UserID: "u1001", Username: "alice", IsActive: true},
			{UserID: "u1002", Username: "bob", IsActive: false},
		},
	}

	resp := postJSON(t, srv.URL+"/team/add", req)
	require.Equal(t, http.StatusCreated, resp.StatusCode)

	var parsed responses.AddTeamResponse
	decodeResponse(t, resp, &parsed)

	require.Equal(t, req.TeamName, parsed.Team.TeamName)
	require.Len(t, parsed.Team.Members, len(req.Members))
	require.ElementsMatch(t, []string{"u1001", "u1002"}, []string{
		parsed.Team.Members[0].UserID,
		parsed.Team.Members[1].UserID,
	})
}

func TestTeamHandler_AddTeam_Duplicate(t *testing.T) {
	srv, _ := handlersTestsHelp.SetupTestServer(t)
	defer srv.Close()

	req := requests.AddTeamRequest{
		TeamName: "team-add-duplicate",
		Members: []requests.TeamMemberRequest{
			{UserID: "u1003", Username: "alice", IsActive: true},
		},
	}

	createTeam(t, srv.URL, req.TeamName, req.Members)

	resp := postJSON(t, srv.URL+"/team/add", req)
	require.Equal(t, http.StatusBadRequest, resp.StatusCode)

	var parsed responses.ErrorResponse
	decodeResponse(t, resp, &parsed)
	require.Equal(t, "TEAM_EXISTS", parsed.Error.Code)
}

func TestTeamHandler_GetTeam_Success(t *testing.T) {
	srv, _ := handlersTestsHelp.SetupTestServer(t)
	defer srv.Close()

	createTeam(t, srv.URL, "team-get-success", []requests.TeamMemberRequest{
		{UserID: "u1004", Username: "alice", IsActive: true},
		{UserID: "u1005", Username: "bob", IsActive: true},
	})

	resp, err := http.Get(srv.URL + "/team/get?team_name=" + url.QueryEscape("team-get-success"))
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	var parsed responseObjs.Team
	decodeResponse(t, resp, &parsed)

	require.Equal(t, "team-get-success", parsed.TeamName)
	require.Len(t, parsed.Members, 2)
	require.ElementsMatch(t, []string{"u1004", "u1005"}, []string{
		parsed.Members[0].UserID,
		parsed.Members[1].UserID,
	})
}

func TestTeamHandler_GetTeam_NotFound(t *testing.T) {
	srv, _ := handlersTestsHelp.SetupTestServer(t)
	defer srv.Close()

	resp, err := http.Get(srv.URL + "/team/get?team_name=ghost")
	require.NoError(t, err)
	require.Equal(t, http.StatusNotFound, resp.StatusCode)

	var parsed responses.ErrorResponse
	decodeResponse(t, resp, &parsed)
	require.Equal(t, "NOT_FOUND", parsed.Error.Code)
}
