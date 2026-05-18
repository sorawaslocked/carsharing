package interceptor

import (
	"context"

	"carsharing/booking-service/internal/pkg/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type BaseInterceptor struct{}

func NewBaseInterceptor() *BaseInterceptor {
	return &BaseInterceptor{}
}

func (i *BaseInterceptor) Unary() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		ctx = i.extractMetadata(ctx)
		return handler(ctx, req)
	}
}

func (i *BaseInterceptor) extractMetadata(ctx context.Context) context.Context {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ctx
	}

	return utils.SetMetadata(
		ctx,
		firstVal(md.Get("x-request-id")),
		firstVal(md.Get("x-client-ip")),
		firstVal(md.Get("x-user-id")),
		firstVal(md.Get("x-user-roles")),
	)
}

func firstVal(vals []string) string {
	if len(vals) > 0 {
		return vals[0]
	}
	return ""
}
