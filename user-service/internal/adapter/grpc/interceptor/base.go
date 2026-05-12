package interceptor

import (
	"context"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// Context keys — must match the constants in internal/pkg/utils/metadata.go.
const (
	CtxRequestIDKey = "x-request-id"
	CtxClientIPKey  = "x-client-ip"
	CtxUserIDKey    = "x-user-id"
	CtxUserRolesKey = "x-user-roles"
)

type BaseInterceptor struct{}

func NewBaseInterceptor() *BaseInterceptor {
	return &BaseInterceptor{}
}

func (i *BaseInterceptor) Unary(ctx context.Context, req any, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	md, _ := metadata.FromIncomingContext(ctx)

	ctx = context.WithValue(ctx, CtxRequestIDKey, stringFromMD(md, "x-request-id"))
	ctx = context.WithValue(ctx, CtxClientIPKey, stringFromMD(md, "x-client-ip"))

	if userID := stringFromMD(md, "x-user-id"); userID != "" {
		ctx = context.WithValue(ctx, CtxUserIDKey, userID)
	}
	if roles := stringsFromMD(md, "x-user-roles"); len(roles) > 0 {
		ctx = context.WithValue(ctx, CtxUserRolesKey, roles)
	}

	return handler(ctx, req)
}

func stringFromMD(md metadata.MD, key string) string {
	if values := md.Get(key); len(values) > 0 {
		return values[0]
	}
	return ""
}

func stringsFromMD(md metadata.MD, key string) []string {
	values := md.Get(key)
	if len(values) == 0 {
		return nil
	}
	// x-user-roles is comma-separated: "admin,user"
	if len(values) == 1 {
		return strings.Split(values[0], ",")
	}
	return values
}
