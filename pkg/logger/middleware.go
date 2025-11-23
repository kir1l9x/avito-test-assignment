package logger

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

func Middleware(l *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		reqID := c.GetHeader("X-Request-ID")
		if reqID == "" {
			reqID = uuid.New().String()
		}

		reqLogger := l.With(zap.String("request_id", reqID))

		c.Set("logger", reqLogger)

		ctx := ToContext(c.Request.Context(), reqLogger)
		c.Request = c.Request.WithContext(ctx)

		c.Writer.Header().Set("X-Request-ID", reqID)

		c.Next()
	}
}
