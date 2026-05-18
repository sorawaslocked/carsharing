package interceptor

import (
	"context"

	"google.golang.org/grpc/metadata"

	"carsharing/trip-service/internal/pkg/utils"
)

func extractMetadata(ctx context.Context) context.Context {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ctx
	}
	return utils.SetMetadata(ctx,
		firstMD(md, "x-request-id"),
		firstMD(md, "x-client-ip"),
		firstMD(md, "x-user-id"),
		firstMD(md, "x-user-roles"),
	)
}

func firstMD(md metadata.MD, key string) string {
	if vals := md.Get(key); len(vals) > 0 {
		return vals[0]
	}
	return ""
}
