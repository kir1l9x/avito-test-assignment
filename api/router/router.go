package router

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/kir1l9x/avito-test-assignment/api/handlers"
	"github.com/kir1l9x/avito-test-assignment/pkg/logger"
)

func NewRouter(
	log *zap.Logger,
	teamsHandler *handlers.TeamHandler,
	usersHandler *handlers.UsersHandler,
	prHandler *handlers.PullRequestHandler,
) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)

	r := gin.New()
	r.Use(logger.Middleware(log))

	r.GET("/health", func(c *gin.Context) {
		c.Status(200)
	})

	team := r.Group("/team")
	{
		team.POST("/add", teamsHandler.AddTeam)
		team.GET("/get", teamsHandler.GetTeam)
	}

	users := r.Group("/users")
	{
		users.POST("/setIsActive", usersHandler.SetIsActive)
		users.GET("/getReview", prHandler.GetUsersReviews)
	}

	pr := r.Group("/pullRequest")
	{
		pr.POST("/create", prHandler.Create)
		pr.POST("/merge", prHandler.Merge)
		pr.POST("/reassign", prHandler.Reassign)
		pr.GET("/stats/all", prHandler.GetReviewsAmountPerUser)
		pr.GET("/stats/open", prHandler.GetOpenReviewsAmountPerUser)
		pr.GET("/stats/top", prHandler.GetTop5ReviewersAllTime)
		pr.GET("/stats/top-open", prHandler.GetTop5ReviewersCurrent)
	}

	return r
}
