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
	"github.com/kir1l9x/avito-test-assignment/internal/domain/valueObjects"
)

func TestPullRequestHandler_Create_Success(t *testing.T) {
	srv, _ := handlersTestsHelp.SetupTestServer(t)
	defer srv.Close()

	createTeam(t, srv.URL, "team-pr-create", []requests.TeamMemberRequest{
		{UserID: "u2001", Username: "author", IsActive: true},
		{UserID: "u2002", Username: "reviewer-1", IsActive: true},
		{UserID: "u2003", Username: "reviewer-2", IsActive: true},
	})

	resp := postJSON(t, srv.URL+"/pullRequest/create", requests.CreatePullRequestRequest{
		PullRequestID:   "pr-1",
		PullRequestName: "Initial PR",
		AuthorID:        "u2001",
	})
	require.Equal(t, http.StatusCreated, resp.StatusCode)

	var parsed responses.CreatePullRequestResponse
	decodeResponse(t, resp, &parsed)

	require.Equal(t, "pr-1", parsed.PR.PullRequestID)
	require.Equal(t, "u2001", parsed.PR.AuthorID)
	require.Equal(t, valueObjects.OpenPullRequest.String(), parsed.PR.Status)
	require.ElementsMatch(t, []string{"u2002", "u2003"}, parsed.PR.AssignedReviewers)
}

func TestPullRequestHandler_Create_AuthorNotFound(t *testing.T) {
	srv, _ := handlersTestsHelp.SetupTestServer(t)
	defer srv.Close()

	resp := postJSON(t, srv.URL+"/pullRequest/create", requests.CreatePullRequestRequest{
		PullRequestID:   "pr-404",
		PullRequestName: "Ghost author",
		AuthorID:        "u9999",
	})
	require.Equal(t, http.StatusNotFound, resp.StatusCode)

	var parsed responses.ErrorResponse
	decodeResponse(t, resp, &parsed)
	require.Equal(t, "NOT_FOUND", parsed.Error.Code)
}

func TestPullRequestHandler_Merge_Idempotent(t *testing.T) {
	srv, _ := handlersTestsHelp.SetupTestServer(t)
	defer srv.Close()

	createTeam(t, srv.URL, "team-pr-merge", []requests.TeamMemberRequest{
		{UserID: "u2011", Username: "author", IsActive: true},
		{UserID: "u2012", Username: "reviewer-1", IsActive: true},
		{UserID: "u2013", Username: "reviewer-2", IsActive: true},
	})

	prID := "pr-2"
	respCreate := postJSON(t, srv.URL+"/pullRequest/create", requests.CreatePullRequestRequest{
		PullRequestID:   prID,
		PullRequestName: "Merge me",
		AuthorID:        "u2011",
	})
	require.Equal(t, http.StatusCreated, respCreate.StatusCode)
	require.NoError(t, respCreate.Body.Close())

	resp := postJSON(t, srv.URL+"/pullRequest/merge", requests.MergePullRequestRequest{
		PullRequestID: prID,
	})
	require.Equal(t, http.StatusOK, resp.StatusCode)

	var first responses.MergePullRequestResponse
	decodeResponse(t, resp, &first)
	require.Equal(t, valueObjects.MergedPullRequest.String(), first.PR.Status)
	require.NotNil(t, first.PR.MergedAt)

	respSecond := postJSON(t, srv.URL+"/pullRequest/merge", requests.MergePullRequestRequest{
		PullRequestID: prID,
	})
	require.Equal(t, http.StatusOK, respSecond.StatusCode)

	var second responses.MergePullRequestResponse
	decodeResponse(t, respSecond, &second)
	require.Equal(t, valueObjects.MergedPullRequest.String(), second.PR.Status)
	require.Equal(t, first.PR.MergedAt, second.PR.MergedAt)
}

func TestPullRequestHandler_ReassignReviewer_Success(t *testing.T) {
	srv, pool := handlersTestsHelp.SetupTestServer(t)
	defer srv.Close()

	createTeam(t, srv.URL, "team-pr-reassign", []requests.TeamMemberRequest{
		{UserID: "u2021", Username: "author", IsActive: true},
		{UserID: "u2022", Username: "reviewer-1", IsActive: true},
		{UserID: "u2023", Username: "reviewer-2", IsActive: true},
		{UserID: "u2024", Username: "reviewer-3", IsActive: true},
	})

	createPullRequestRecord(t, pool, "pr-3", "Needs reassignment", "u2021", []string{"u2022", "u2023"})

	resp := postJSON(t, srv.URL+"/pullRequest/reassign", requests.ReassignPullRequestRequest{
		PullRequestID: "pr-3",
		OldUserID:     "u2022",
	})
	require.Equal(t, http.StatusOK, resp.StatusCode)

	var parsed responses.ReassignPullRequestResponse
	decodeResponse(t, resp, &parsed)

	require.Equal(t, "u2024", parsed.ReplacedBy)
	require.ElementsMatch(t, []string{"u2023", "u2024"}, parsed.PR.AssignedReviewers)
}

func TestPullRequestHandler_ReassignReviewer_MergedPR(t *testing.T) {
	srv, pool := handlersTestsHelp.SetupTestServer(t)
	defer srv.Close()

	createTeam(t, srv.URL, "team-pr-reassign-merged", []requests.TeamMemberRequest{
		{UserID: "u2031", Username: "author", IsActive: true},
		{UserID: "u2032", Username: "reviewer-1", IsActive: true},
		{UserID: "u2033", Username: "reviewer-2", IsActive: true},
	})

	createPullRequestRecord(t, pool, "pr-4", "Merge first", "u2031", []string{"u2032", "u2033"})

	respMerge := postJSON(t, srv.URL+"/pullRequest/merge", requests.MergePullRequestRequest{
		PullRequestID: "pr-4",
	})
	require.Equal(t, http.StatusOK, respMerge.StatusCode)
	require.NoError(t, respMerge.Body.Close())

	resp := postJSON(t, srv.URL+"/pullRequest/reassign", requests.ReassignPullRequestRequest{
		PullRequestID: "pr-4",
		OldUserID:     "u2032",
	})
	require.Equal(t, http.StatusConflict, resp.StatusCode)

	var parsed responses.ErrorResponse
	decodeResponse(t, resp, &parsed)
	require.Equal(t, "PR_MERGED", parsed.Error.Code)
}

func TestPullRequestHandler_GetUsersReviews(t *testing.T) {
	srv, pool := handlersTestsHelp.SetupTestServer(t)
	defer srv.Close()

	createTeam(t, srv.URL, "team-pr-getreviews", []requests.TeamMemberRequest{
		{UserID: "u2041", Username: "author", IsActive: true},
		{UserID: "u2042", Username: "reviewer-1", IsActive: true},
		{UserID: "u2043", Username: "reviewer-2", IsActive: true},
	})

	createPullRequestRecord(t, pool, "pr-5", "First", "u2041", []string{"u2042"})
	createPullRequestRecord(t, pool, "pr-6", "Second", "u2041", []string{"u2042", "u2043"})
	createPullRequestRecord(t, pool, "pr-7", "Skip user", "u2041", []string{"u2043"})

	resp, err := http.Get(srv.URL + "/users/getReview?user_id=" + url.QueryEscape("u2042"))
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	var parsed responses.GetUsersReviewsResponse
	decodeResponse(t, resp, &parsed)

	require.Equal(t, "u2042", parsed.UserID)
	require.Len(t, parsed.PullRequests, 2)
	require.ElementsMatch(t, []string{"pr-5", "pr-6"}, []string{
		parsed.PullRequests[0].PullRequestID,
		parsed.PullRequests[1].PullRequestID,
	})
}

func TestPullRequestHandler_GetUsersReviews_UserWithoutAssignments(t *testing.T) {
	srv, _ := handlersTestsHelp.SetupTestServer(t)
	defer srv.Close()

	resp, err := http.Get(srv.URL + "/users/getReview?user_id=" + url.QueryEscape("u8888"))
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	var parsed responses.GetUsersReviewsResponse
	decodeResponse(t, resp, &parsed)

	require.Equal(t, "u8888", parsed.UserID)
	require.Len(t, parsed.PullRequests, 0)
}

func TestPullRequestHandler_StatsEndpoints(t *testing.T) {
	srv, pool := handlersTestsHelp.SetupTestServer(t)
	defer srv.Close()

	createTeam(t, srv.URL, "team-pr-stats", []requests.TeamMemberRequest{
		{UserID: "u2051", Username: "author", IsActive: true},
		{UserID: "u2052", Username: "reviewer-1", IsActive: true},
		{UserID: "u2053", Username: "reviewer-2", IsActive: true},
	})

	createPullRequestRecord(t, pool, "pr-10", "Open A", "u2051", []string{"u2052", "u2053"})
	createPullRequestRecord(t, pool, "pr-11", "Open B", "u2051", []string{"u2052"})
	createPullRequestRecord(t, pool, "pr-12", "Merged", "u2051", []string{"u2053"})
	markPullRequestMerged(t, pool, "pr-12")

	expectedAll := map[string]int64{"u2052": 2, "u2053": 2}
	expectedOpen := map[string]int64{"u2052": 2, "u2053": 1}

	all := fetchStats(t, srv.URL+"/pullRequest/stats/all")
	assertHandlerStats(t, all.Stats, expectedAll)

	open := fetchStats(t, srv.URL+"/pullRequest/stats/open")
	assertHandlerStats(t, open.Stats, expectedOpen)

	topAll := fetchStats(t, srv.URL+"/pullRequest/stats/top")
	assertHandlerStats(t, topAll.Stats, expectedAll)

	topOpen := fetchStats(t, srv.URL+"/pullRequest/stats/top-open")
	assertHandlerStats(t, topOpen.Stats, expectedOpen)
}

func fetchStats(t *testing.T, url string) responses.ReviewsStatsResponse {
	t.Helper()

	resp, err := http.Get(url)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	var parsed responses.ReviewsStatsResponse
	decodeResponse(t, resp, &parsed)

	return parsed
}

func assertHandlerStats(t *testing.T, stats []responseObjs.UserReviewStat, expected map[string]int64) {
	t.Helper()

	actual := make(map[string]int64, len(stats))
	for _, stat := range stats {
		actual[stat.UserID] = stat.ReviewCount
	}

	require.Equal(t, expected, actual)
}
