package interceptor

import (
	pkglog "carsharing/shared/pkg/log"
	"carsharing/shared/pkg/utils"
	"context"
	"log/slog"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type LoggerInterceptor struct {
	log *slog.Logger
}

func NewLoggerInterceptor(log *slog.Logger) *LoggerInterceptor {
	return &LoggerInterceptor{log: log}
}

func (i *LoggerInterceptor) Unary(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	md := utils.MetadataFromCtx(ctx)
	log := pkglog.WithMetadata(i.log, md)
	log = log.With(slog.String("method", info.FullMethod))

	start := time.Now()
	m, err := handler(ctx, req)
	duration := time.Since(start)

	code := codes.OK
	var errMsg string
	if err != nil {
		if st, ok := status.FromError(err); ok {
			code = st.Code()
			errMsg = st.Message()
		} else {
			code = codes.Unknown
			errMsg = err.Error()
		}
	}

	attrs := []any{
		slog.String("status", code.String()),
		slog.Int64("durationMs", duration.Milliseconds()),
	}
	if errMsg != "" {
		attrs = append(attrs, slog.String("error", errMsg))
	}

	log.Log(ctx, grpcLogLevel(code), "grpc call", attrs...)

	return m, err
}

func grpcLogLevel(code codes.Code) slog.Level {
	switch code {
	case codes.OK:
		return slog.LevelInfo
	case codes.Canceled,
		codes.InvalidArgument,
		codes.NotFound,
		codes.AlreadyExists,
		codes.PermissionDenied,
		codes.Unauthenticated,
		codes.ResourceExhausted,
		codes.FailedPrecondition,
		codes.OutOfRange:
		return slog.LevelWarn
	default:
		return slog.LevelError
	}
}
