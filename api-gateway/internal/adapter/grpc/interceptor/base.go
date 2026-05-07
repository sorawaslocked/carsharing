package interceptor

import (
	"context"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

const (
	ctxRequestID = "x-request-id"
	ctxClientIP  = "x-client-ip"
	ctxDeviceID  = "x-device-id"
	ctxUserID    = "x-user-id"
	ctxUserRoles = "x-user-roles"
)

type BaseClientInterceptor struct{}

func NewBaseInterceptor() *BaseClientInterceptor {
	return &BaseClientInterceptor{}
}

func (i *BaseClientInterceptor) Unary(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	if v := ctx.Value(ctxRequestID); v != nil {
		ctx = metadata.AppendToOutgoingContext(ctx, ctxRequestID, v.(string))
	}
	if v := ctx.Value(ctxClientIP); v != nil {
		ctx = metadata.AppendToOutgoingContext(ctx, ctxClientIP, v.(string))
	}
	if v := ctx.Value(ctxDeviceID); v != nil {
		ctx = metadata.AppendToOutgoingContext(ctx, ctxDeviceID, v.(string))
	}
	if v := ctx.Value(ctxUserID); v != nil {
		ctx = metadata.AppendToOutgoingContext(ctx, ctxUserID, v.(string))
	}
	if v := ctx.Value(ctxUserRoles); v != nil {
		ctx = metadata.AppendToOutgoingContext(ctx, ctxUserRoles, strings.Join(v.([]string), ","))
	}

	return invoker(ctx, method, req, reply, cc, opts...)
}
