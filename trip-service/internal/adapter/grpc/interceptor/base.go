package interceptor

import (
	"context"
	"strings"

	"google.golang.org/grpc/metadata"
)

func extractMetadata(ctx context.Context) context.Context {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ctx
	}
	if requestID := firstMD(md, "x-request-id"); requestID != "" {
		ctx = context.WithValue(ctx, "x-request-id", requestID)
	}
	if clientIP := firstMD(md, "x-client-ip"); clientIP != "" {
		ctx = context.WithValue(ctx, "x-client-ip", clientIP)
	}
	if userID := firstMD(md, "x-user-id"); userID != "" {
		ctx = context.WithValue(ctx, "x-user-id", userID)
	}
	if rolesStr := firstMD(md, "x-user-roles"); rolesStr != "" {
		var roles []string
		for _, r := range strings.Split(rolesStr, ",") {
			if r = strings.TrimSpace(r); r != "" {
				roles = append(roles, r)
			}
		}
		if len(roles) > 0 {
			ctx = context.WithValue(ctx, "x-user-roles", roles)
		}
	}
	return ctx
}

func firstMD(md metadata.MD, key string) string {
	if vals := md.Get(key); len(vals) > 0 {
		return vals[0]
	}
	return ""
}
