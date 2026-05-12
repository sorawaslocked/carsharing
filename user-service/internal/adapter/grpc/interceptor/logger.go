package interceptor

import (
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
	requestID, _ := ctx.Value(CtxRequestIDKey).(string)
	clientIP, _ := ctx.Value(CtxClientIPKey).(string)

	logger := i.log.With(
		slog.String("requestId", requestID),
		slog.String("clientIP", clientIP),
		slog.String("method", info.FullMethod),
	)
	logger.Info("grpc request")

	m, err := handler(ctx, req)

	var statusCode string
	if err != nil {
		if st, ok := status.FromError(err); ok {
			statusCode = st.Code().String()
		}
	} else {
		statusCode = codes.OK.String()
	}

	logger.Info("grpc response", slog.String("status", statusCode))

	return m, err
}
