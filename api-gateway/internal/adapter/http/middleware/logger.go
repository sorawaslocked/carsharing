package middleware

import (
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
)

func Logger(log *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		args := []any{
			slog.String("method", c.Request.Method),
			slog.String("path", c.Request.URL.Path),
			slog.String("query", c.Request.URL.RawQuery),
			slog.Int("status", c.Writer.Status()),
			slog.Duration("duration", time.Since(start)),
			slog.String("requestID", c.GetString("x-request-id")),
			slog.String("clientIP", c.GetString("x-client-ip")),
		}
		if userID := c.GetString("x-user-id"); userID != "" {
			args = append(args, slog.String("userID", userID))
		}

		log.Info("request", args...)
	}
}
