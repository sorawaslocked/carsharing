package interceptor

import (
	"context"
	"log/slog"
	"time"

	pkglog "carsharing/shared/pkg/log"
	"carsharing/shared/pkg/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

type LoggerInterceptor struct {
	log *slog.Logger
}

func NewLoggerInterceptor(log *slog.Logger) *LoggerInterceptor {
	return &LoggerInterceptor{
		log: pkglog.WithComponent(log, "grpc.LoggerInterceptor"),
	}
}

func (i *LoggerInterceptor) Unary() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		start := time.Now()
		log := pkglog.WithMethod(i.log, info.FullMethod)
		log = pkglog.WithMetadata(log, utils.MetadataFromCtx(ctx))

		resp, err := handler(ctx, req)

		log.Info("grpc call",
			slog.String("code", status.Code(err).String()),
			slog.Duration("duration", time.Since(start)),
		)

		return resp, err
	}
}
