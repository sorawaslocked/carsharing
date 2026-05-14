package interceptor

import (
	"context"

	"github.com/sorawaslocked/car-rental-booking-service/internal/pkg/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var publicMethods = map[string]struct{}{
	"/service.booking.HealthService/Health": {},
}

type AuthInterceptor struct{}

func NewAuthInterceptor() *AuthInterceptor {
	return &AuthInterceptor{}
}

func (i *AuthInterceptor) Unary() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		if _, ok := publicMethods[info.FullMethod]; ok {
			return handler(ctx, req)
		}

		md := utils.MetadataFromCtx(ctx)
		if md.UserID == nil || *md.UserID == "" {
			return nil, status.Error(codes.Unauthenticated, "missing user id")
		}

		return handler(ctx, req)
	}
}
