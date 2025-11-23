package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/kir1l9x/avito-test-assignment/api/dto/requests"
	"github.com/kir1l9x/avito-test-assignment/api/dto/responses"
	"github.com/kir1l9x/avito-test-assignment/api/dto/responses/responseObjs"
	"github.com/kir1l9x/avito-test-assignment/api/middleware"
	"github.com/kir1l9x/avito-test-assignment/internal/appplication/services"
)

type PullRequestHandler struct {
	prService *services.PullRequestsService
}

func NewPullRequestHandler(prService *services.PullRequestsService) *PullRequestHandler {
	return &PullRequestHandler{prService: prService}
}

func (prh *PullRequestHandler) Create(c *gin.Context) {
	var req requests.CreatePullRequestRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.WriteValidationError(c, err)
		return
	}

	pr, err := prh.prService.CreatePullRequest(c.Request.Context(), req.PullRequestID, req.PullRequestName, req.AuthorID)
	if err != nil {
		middleware.WriteError(c, err)
		return
	}

	resp := responses.CreatePullRequestResponse{PR: responseObjs.FromDomainPullRequest(pr)}

	c.JSON(http.StatusCreated, resp)
}

func (prh *PullRequestHandler) Merge(c *gin.Context) {
	var req requests.MergePullRequestRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.WriteValidationError(c, err)
		return
	}

	pr, err := prh.prService.Merge(c.Request.Context(), req.PullRequestID)
	if err != nil {
		middleware.WriteError(c, err)
		return
	}

	resp := responses.MergePullRequestResponse{PR: responseObjs.FromDomainPullRequest(pr)}

	c.JSON(http.StatusOK, resp)
}

func (prh *PullRequestHandler) Reassign(c *gin.Context) {
	var req requests.ReassignPullRequestRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.WriteValidationError(c, err)
		return
	}

	pr, replacedBy, err := prh.prService.ReassignReviewer(c.Request.Context(), req.PullRequestID, req.OldUserID)
	if err != nil {
		middleware.WriteError(c, err)
		return
	}

	resp := responses.ReassignPullRequestResponse{
		PR:         responseObjs.FromDomainPullRequest(pr),
		ReplacedBy: replacedBy,
	}

	c.JSON(http.StatusOK, resp)
}

func (prh *PullRequestHandler) GetUsersReviews(c *gin.Context) {
	var req requests.GetUsersReviewsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		middleware.WriteValidationError(c, err)
		return
	}

	prs, err := prh.prService.GetPRsWhereUserIsReviewer(c.Request.Context(), req.UserID)
	if err != nil {
		middleware.WriteError(c, err)
		return
	}

	resp := responses.GetUsersReviewsResponse{
		UserID:       req.UserID,
		PullRequests: responseObjs.FromDomainPullRequestsShort(prs),
	}

	c.JSON(http.StatusOK, resp)
}

func (prh *PullRequestHandler) GetReviewsAmountPerUser(c *gin.Context) {
	stats, err := prh.prService.GetReviewsAmountPerUser(c.Request.Context())
	if err != nil {
		middleware.WriteError(c, err)
		return
	}

	resp := responses.ReviewsStatsResponse{
		Stats: responseObjs.FromDomainUserReviewStats(stats),
	}

	c.JSON(http.StatusOK, resp)
}

func (prh *PullRequestHandler) GetOpenReviewsAmountPerUser(c *gin.Context) {
	stats, err := prh.prService.GetOpenReviewsAmountPerUser(c.Request.Context())
	if err != nil {
		middleware.WriteError(c, err)
		return
	}

	resp := responses.ReviewsStatsResponse{
		Stats: responseObjs.FromDomainUserReviewStats(stats),
	}

	c.JSON(http.StatusOK, resp)
}

func (prh *PullRequestHandler) GetTop5ReviewersAllTime(c *gin.Context) {
	stats, err := prh.prService.GetTop5ReviewersAllTime(c.Request.Context())
	if err != nil {
		middleware.WriteError(c, err)
		return
	}

	resp := responses.ReviewsStatsResponse{
		Stats: responseObjs.FromDomainUserReviewStats(stats),
	}

	c.JSON(http.StatusOK, resp)
}

func (prh *PullRequestHandler) GetTop5ReviewersCurrent(c *gin.Context) {
	stats, err := prh.prService.GetTop5ReviewersCurrent(c.Request.Context())
	if err != nil {
		middleware.WriteError(c, err)
		return
	}

	resp := responses.ReviewsStatsResponse{
		Stats: responseObjs.FromDomainUserReviewStats(stats),
	}

	c.JSON(http.StatusOK, resp)
}
