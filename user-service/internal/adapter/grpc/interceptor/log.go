package interceptor

import (
	"carsharing/shared/pkg/utils"
	"context"
	"log/slog"

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

	log := i.log.With(
		slog.String("requestID", md.RequestID),
		slog.String("clientIP", md.ClientIP),
		slog.String("method", info.FullMethod),
	)
	log.Info("grpc request")

	m, err := handler(ctx, req)

	var statusCode string
	if err != nil {
		if st, ok := status.FromError(err); ok {
			statusCode = st.Code().String()
		}
	} else {
		statusCode = codes.OK.String()
	}

	log.Info("grpc response", slog.String("status", statusCode))

	return m, err
}
