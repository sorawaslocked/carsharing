package http

import (
	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	"log/slog"
)

func LoggerMiddleware(log *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		requestLog := log.With(
			slog.String("requestId", requestid.Get(c)),
			slog.String("clientIP", c.ClientIP()),
			slog.String("method", c.Request.Method),
			slog.String("path", c.Request.URL.Path),
			slog.String("query", c.Request.URL.RawQuery),
		)
		requestLog.Info("api request")

		c.Next()

		requestLog.Info(
			"api response",
			slog.Int("status", c.Writer.Status()),
		)
	}
}
