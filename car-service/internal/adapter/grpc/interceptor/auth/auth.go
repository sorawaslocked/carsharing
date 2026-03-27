package auth

import (
	"context"
	"github.com/sorawaslocked/car-rental-car-service/internal/adapter/grpc/dto"
	"github.com/sorawaslocked/car-rental-car-service/internal/model"
	"github.com/sorawaslocked/car-rental-car-service/internal/pkg/utils"

	"google.golang.org/grpc"
)

type Interceptor struct{}

func NewInterceptor() *Interceptor {
	return &Interceptor{}
}

func (i *Interceptor) Unary(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	md, ok := utils.MetadataFromCtx(ctx)
	if !ok {
		return nil, dto.FromErrorToStatusCode(model.ErrMissingMetadata)
	}

	if !isAllowed(info.FullMethod, md.UserRoles) {
		return nil, dto.FromErrorToStatusCode(model.ErrUnauthorized)
	}

	return handler(ctx, req)
}
