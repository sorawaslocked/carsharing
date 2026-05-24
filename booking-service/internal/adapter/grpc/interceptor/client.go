package interceptor

import (
	"context"
	"strings"

	sharedmodel "carsharing/shared/model"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type ClientBaseInterceptor struct{}

func NewClientBaseInterceptor() *ClientBaseInterceptor {
	return &ClientBaseInterceptor{}
}

// Unary forwards request-scoped metadata (request ID, client IP, user ID, user roles)
// from the local context into outgoing gRPC metadata so downstream services can trace calls.
func (i *ClientBaseInterceptor) Unary(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	md := metadata.New(nil)

	if v, ok := ctx.Value("x-request-id").(string); ok && v != "" {
		md.Set("x-request-id", v)
	}
	if v, ok := ctx.Value("x-client-ip").(string); ok && v != "" {
		md.Set("x-client-ip", v)
	}
	if v, ok := ctx.Value("x-user-id").(string); ok && v != "" {
		md.Set("x-user-id", v)
	}
	if roles, ok := ctx.Value("x-user-roles").([]sharedmodel.Role); ok && len(roles) > 0 {
		strs := make([]string, len(roles))
		for i, r := range roles {
			strs[i] = string(r)
		}
		md.Set("x-user-roles", strings.Join(strs, ","))
	}

	ctx = metadata.NewOutgoingContext(ctx, md)
	return invoker(ctx, method, req, reply, cc, opts...)
}
