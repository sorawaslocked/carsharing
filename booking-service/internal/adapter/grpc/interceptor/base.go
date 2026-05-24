package interceptor

import (
	"context"
	"strings"

	"carsharing/booking-service/internal/adapter/grpc/dto"
	"carsharing/booking-service/internal/model"
	sharedmodel "carsharing/shared/model"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type BaseInterceptor struct{}

func NewBaseInterceptor() *BaseInterceptor {
	return &BaseInterceptor{}
}

func (i *BaseInterceptor) Unary(ctx context.Context, req any, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	md, _ := metadata.FromIncomingContext(ctx)

	ctx = context.WithValue(ctx, "x-request-id", stringFromMD(md, "x-request-id"))
	ctx = context.WithValue(ctx, "x-client-ip", stringFromMD(md, "x-client-ip"))

	if userID := stringFromMD(md, "x-user-id"); userID != "" {
		ctx = context.WithValue(ctx, "x-user-id", userID)
	}
	if roleStrs := stringsFromMD(md, "x-user-roles"); len(roleStrs) > 0 {
		roles := make([]sharedmodel.Role, len(roleStrs))
		for j, s := range roleStrs {
			role, ok := sharedmodel.RoleFromString(s)
			if !ok {
				return nil, dto.ToStatusError(model.ErrInvalidMetadata)
			}
			roles[j] = role
		}
		ctx = context.WithValue(ctx, "x-user-roles", roles)
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
	if len(values) == 1 {
		return strings.Split(values[0], ",")
	}
	return values
}
