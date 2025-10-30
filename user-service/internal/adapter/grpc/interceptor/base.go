package interceptor

import (
	"context"
	"github.com/sorawaslocked/car-rental-user-service/internal/adapter/grpc/dto"
	"github.com/sorawaslocked/car-rental-user-service/internal/model"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

const (
	CtxClientIPKey  = "client-ip"
	CtxRequestIDKey = "request-id"
)

type BaseInterceptor struct{}

func NewBaseInterceptor() *BaseInterceptor {
	return &BaseInterceptor{}
}

func (i *BaseInterceptor) Unary(ctx context.Context, req any, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, dto.ToStatusCodeError(model.ErrMissingMetadata)
	}

	ctx = context.WithValue(ctx, CtxRequestIDKey, requestIDFromMetadata(md))
	ctx = context.WithValue(ctx, CtxClientIPKey, clientIPFromMetadata(md))

	return handler(ctx, req)
}
