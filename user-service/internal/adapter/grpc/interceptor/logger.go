package interceptor

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"log/slog"
)

type LoggerInterceptor struct {
	log *slog.Logger
}

func NewLoggerInterceptor(log *slog.Logger) *LoggerInterceptor {
	return &LoggerInterceptor{
		log: log,
	}
}

func (i *LoggerInterceptor) Unary(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	requestID := ctx.Value(CtxRequestIDKey).(string)
	clientIP := ctx.Value(CtxClientIPKey).(string)

	logger := i.log.With(
		slog.String("requestId", requestID),
		slog.String("clientIP", clientIP),
		slog.String("method", info.FullMethod),
	)
	logger.Info("grpc request")

	m, err := handler(ctx, req)

	var statusString string

	if err != nil {
		if st, ok := status.FromError(err); ok {
			statusString = st.Code().String()
		}
	} else {
		statusString = codes.OK.String()
	}

	logger.Info(
		"grpc response",
		slog.String("status", statusString),
	)

	return m, err
}

func clientIPFromMetadata(md metadata.MD) string {
	clientIPs := md.Get("x-client-ip")
	if len(clientIPs) > 0 {
		return clientIPs[0]
	}

	return ""
}

func requestIDFromMetadata(md metadata.MD) string {
	requestIDs := md.Get("x-request-id")
	if len(requestIDs) > 0 {
		return requestIDs[0]
	}

	return ""
}
