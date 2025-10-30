package interceptor

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

const (
	CtxRequestID           = "x-request-id"
	CtxClientIP            = "x-client-ip"
	CtxAuthorizationHeader = "authorization"
)

type BaseClientInterceptor struct{}

func NewBaseInterceptor() *BaseClientInterceptor {
	return &BaseClientInterceptor{}
}

func (i *BaseClientInterceptor) Unary(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	requestID := ctx.Value(CtxRequestID)
	clientIP := ctx.Value(CtxClientIP)
	authorizationHeader := ctx.Value(CtxAuthorizationHeader)

	if requestID != nil {
		ctx = metadata.AppendToOutgoingContext(ctx, CtxRequestID, requestID.(string))
	}
	if clientIP != nil {
		ctx = metadata.AppendToOutgoingContext(ctx, CtxClientIP, clientIP.(string))
	}
	if authorizationHeader != nil {
		ctx = metadata.AppendToOutgoingContext(ctx, CtxAuthorizationHeader, authorizationHeader.(string))
	}

	return invoker(ctx, method, req, reply, cc, opts...)
}
