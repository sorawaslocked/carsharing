package interceptor

import (
	"context"
	"strings"

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

	ctx = context.WithValue(ctx, "x-request-id", firstVal(md.Get("x-request-id")))
	ctx = context.WithValue(ctx, "x-client-ip", firstVal(md.Get("x-client-ip")))

	if userID := firstVal(md.Get("x-user-id")); userID != "" {
		ctx = context.WithValue(ctx, "x-user-id", userID)
	}

	if rolesStr := firstVal(md.Get("x-user-roles")); rolesStr != "" {
		var roles []string
		for _, r := range strings.Split(rolesStr, ",") {
			if s := strings.TrimSpace(r); s != "" {
				roles = append(roles, s)
			}
		}
		if len(roles) > 0 {
			ctx = context.WithValue(ctx, "x-user-roles", roles)
		}
	}

	return ctx
}

func firstVal(vals []string) string {
	if len(vals) > 0 {
		return vals[0]
	}
	return ""
}
