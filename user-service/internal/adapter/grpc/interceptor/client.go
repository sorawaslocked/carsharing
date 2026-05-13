package interceptor

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type ClientBaseInterceptor struct{}

func NewClientBaseInterceptor() *ClientBaseInterceptor {
	return &ClientBaseInterceptor{}
}

// Unary forwards request-scoped metadata (request ID, client IP) from the
// local context into outgoing gRPC metadata so downstream services can trace calls.
func (i *ClientBaseInterceptor) Unary(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	md := metadata.New(nil)

	if v, ok := ctx.Value(CtxRequestIDKey).(string); ok && v != "" {
		md.Set("x-request-id", v)
	}
	if v, ok := ctx.Value(CtxClientIPKey).(string); ok && v != "" {
		md.Set("x-client-ip", v)
	}

	ctx = metadata.NewOutgoingContext(ctx, md)
	return invoker(ctx, method, req, reply, cc, opts...)
}
