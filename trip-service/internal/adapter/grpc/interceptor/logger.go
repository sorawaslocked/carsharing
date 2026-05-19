package interceptor

import (
	"context"
	"log/slog"
	"time"

	"google.golang.org/grpc"

	pkglog "carsharing/shared/pkg/log"
	pkgutils "carsharing/shared/pkg/utils"
)

func LoggerUnaryInterceptor(log *slog.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		log := pkglog.WithMetadata(log, pkgutils.MetadataFromCtx(ctx))
		start := time.Now()

		resp, err := handler(ctx, req)

		attrs := []any{
			slog.String("method", info.FullMethod),
			slog.Duration("duration", time.Since(start)),
		}
		if err != nil {
			log.Error("grpc request failed", append(attrs, pkglog.Err(err))...)
		} else {
			log.Info("grpc request", attrs...)
		}

		return resp, err
	}
}

func LoggerStreamInterceptor(log *slog.Logger) grpc.StreamServerInterceptor {
	return func(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		log := pkglog.WithMetadata(log, pkgutils.MetadataFromCtx(ss.Context()))
		start := time.Now()

		err := handler(srv, ss)

		attrs := []any{
			slog.String("method", info.FullMethod),
			slog.Duration("duration", time.Since(start)),
		}
		if err != nil {
			log.Error("grpc stream failed", append(attrs, pkglog.Err(err))...)
		} else {
			log.Info("grpc stream", attrs...)
		}

		return err
	}
}
