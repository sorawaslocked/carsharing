package log

import (
	"log/slog"

	"carsharing/api-gateway/internal/pkg/utils"
)

func Err(err error) slog.Attr {
	return slog.String("error", err.Error())
}

func WithComponent(log *slog.Logger, component string) *slog.Logger {
	return log.With(
		slog.Group("src",
			slog.String("component", component),
		),
	)
}

func WithMethod(log *slog.Logger, method string) *slog.Logger {
	return log.With(
		slog.Group("src",
			slog.String("method", method),
		),
	)
}

func WithMetadata(log *slog.Logger, md utils.Metadata) *slog.Logger {
	args := make([]any, 0, 3)
	args = append(args, slog.String("clientIP", md.ClientIP))
	args = append(args, slog.String("requestID", md.RequestID))
	if md.UserID != nil {
		args = append(args, slog.String("userID", *md.UserID))
	}
	return log.With(slog.Group("metadata", args...))
}
