package interceptor

import (
	"context"
	"strings"

	"carsharing/car-service/internal/adapter/grpc/dto"
	"carsharing/car-service/internal/model"
	sharedmodel "carsharing/shared/model"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

const (
	mdKeyRequestID = "x-request-id"
	mdKeyClientIP  = "x-client-ip"
	mdKeyUserID    = "x-user-id"
	mdKeyUserRoles = "x-user-roles"
)

type BaseInterceptor struct{}

func NewBaseInterceptor() *BaseInterceptor {
	return &BaseInterceptor{}
}

func (i *BaseInterceptor) Unary(ctx context.Context, req any, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, dto.FromErrorToStatusCode(model.ErrMissingMetadata)
	}

	ctx = context.WithValue(ctx, mdKeyRequestID, firstVal(md.Get(mdKeyRequestID)))
	ctx = context.WithValue(ctx, mdKeyClientIP, firstVal(md.Get(mdKeyClientIP)))
	ctx = context.WithValue(ctx, mdKeyUserID, firstVal(md.Get(mdKeyUserID)))
	ctx = context.WithValue(ctx, mdKeyUserRoles, parseRoles(firstVal(md.Get(mdKeyUserRoles))))

	return handler(ctx, req)
}

func firstVal(vals []string) string {
	if len(vals) > 0 {
		return vals[0]
	}
	return ""
}

func parseRoles(s string) []sharedmodel.Role {
	if s == "" {
		return nil
	}
	var roles []sharedmodel.Role
	for _, r := range strings.Split(s, ",") {
		if r = strings.TrimSpace(r); r != "" {
			roles = append(roles, sharedmodel.Role(r))
		}
	}
	return roles
}
