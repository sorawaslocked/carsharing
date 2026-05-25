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
	log := pkglog.WithMetadata(i.log, utils.MetadataFromCtx(ctx))
	log = log.With(slog.String("method", info.FullMethod))

	start := time.Now()
	resp, err := handler(ctx, req)
	duration := time.Since(start)

	code, errMsg := grpcCodeAndMsg(err)

	attrs := []any{
		slog.String("status", code.String()),
		slog.Int64("durationMs", duration.Milliseconds()),
	}
	if errMsg != "" {
		attrs = append(attrs, slog.String("error", errMsg))
	}
	log.Log(ctx, grpcLogLevel(code), "grpc call", attrs...)

	return resp, err
}

func (i *LoggerInterceptor) Stream(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	log := pkglog.WithMetadata(i.log, utils.MetadataFromCtx(ss.Context()))
	log = log.With(slog.String("method", info.FullMethod))

	start := time.Now()
	err := handler(srv, ss)
	duration := time.Since(start)

	code, errMsg := grpcCodeAndMsg(err)

	attrs := []any{
		slog.String("status", code.String()),
		slog.Int64("durationMs", duration.Milliseconds()),
	}
	if errMsg != "" {
		attrs = append(attrs, slog.String("error", errMsg))
	}
	log.Log(ss.Context(), grpcLogLevel(code), "grpc stream", attrs...)

	return err
}

func grpcCodeAndMsg(err error) (codes.Code, string) {
	if err == nil {
		return codes.OK, ""
	}
	if st, ok := status.FromError(err); ok {
		return st.Code(), st.Message()
	}
	return codes.Unknown, err.Error()
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
