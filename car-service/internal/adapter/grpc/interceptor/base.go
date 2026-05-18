package interceptor

import (
	"context"

	"github.com/sorawaslocked/car-rental-car-service/internal/adapter/grpc/dto"
	"github.com/sorawaslocked/car-rental-car-service/internal/model"
	"github.com/sorawaslocked/car-rental-car-service/internal/pkg/utils"

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

	ctx = utils.SetMetadata(
		ctx,
		firstVal(md.Get(mdKeyRequestID)),
		firstVal(md.Get(mdKeyClientIP)),
		firstVal(md.Get(mdKeyUserID)),
		firstVal(md.Get(mdKeyUserRoles)),
	)

	return handler(ctx, req)
}

func firstVal(vals []string) string {
	if len(vals) > 0 {
		return vals[0]
	}
	return ""
}
