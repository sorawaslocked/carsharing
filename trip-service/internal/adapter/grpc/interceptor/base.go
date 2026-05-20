package interceptor

import (
	"context"
	"strings"

	sharedmodel "carsharing/shared/model"
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
		var roles []sharedmodel.Role
		for _, r := range strings.Split(rolesStr, ",") {
			if r = strings.TrimSpace(r); r != "" {
				roles = append(roles, sharedmodel.Role(r))
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
