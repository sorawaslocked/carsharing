package log

import (
	"carsharing/shared/pkg/utils"
	"log/slog"
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

	userNotNil := md.UserID != nil
	rolesNotEmpty := len(md.UserRoles) > 0

	if userNotNil || rolesNotEmpty {
		userArgs := make([]any, 0, 2)

		if userNotNil {
			userArgs = append(userArgs, slog.String("id", *md.UserID))
		}
		if rolesNotEmpty {
			userArgs = append(userArgs, slog.Any("roles", md.UserRoles))
		}

		args = append(args, slog.Group("user", userArgs...))
	}

	return log.With(args)
}
