package middleware

import (
	"github.com/gin-gonic/gin"
	"log/slog"
)

func Logger(log *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		logger := log.With(
			slog.String("requestId", c.GetString("x-request-id")),
			slog.String("clientIP", c.GetString("x-client-ip")),
			slog.String("method", c.Request.Method),
			slog.String("path", c.Request.URL.Path),
			slog.String("query", c.Request.URL.RawQuery),
		)
		logger.Info("api request")

		c.Next()

		logger.Info(
			"api response",
			slog.Int("status", c.Writer.Status()),
		)
	}
}
