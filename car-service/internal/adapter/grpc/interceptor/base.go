package interceptor

import (
	"car-rental-car-service/internal/adapter/grpc/dto"
	"car-rental-car-service/internal/model"
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

const (
	requestIDKey        = "x-request-id"
	requestClientIPKey  = "x-client-ip"
	requestUserIDKey    = "x-user-id"
	requestUserRoles    = "x-user-roles"
	requestUserVerified = "x-user-verified"
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

	ctx = context.WithValue(ctx, requestIDKey, requestIDFromMetadata(md))
	ctx = context.WithValue(ctx, requestClientIPKey, requestClientIPFromMetadata(md))
	ctx = context.WithValue(ctx, requestUserIDKey, requestUserIDFromMetadata(md))
	ctx = context.WithValue(ctx, requestUserRoles, requestUserRolesFromMetadata(md))
	ctx = context.WithValue(ctx, requestUserVerified, requestUserVerifiedFromMetadata(md))

	return handler(ctx, req)
}

func requestClientIPFromMetadata(md metadata.MD) string {
	clientIP := md.Get(requestClientIPKey)
	if len(clientIP) > 0 {
		return clientIP[0]
	}

	return ""
}

func requestIDFromMetadata(md metadata.MD) string {
	requestID := md.Get(requestIDKey)
	if len(requestID) > 0 {
		return requestID[0]
	}

	return ""
}

func requestUserIDFromMetadata(md metadata.MD) string {
	requestUserID := md.Get(requestUserIDKey)
	if len(requestUserID) > 0 {
		return requestUserID[0]
	}

	return ""
}

func requestUserRolesFromMetadata(md metadata.MD) string {
	requestUserRolesStr := md.Get(requestUserRoles)
	if len(requestUserRolesStr) > 0 {
		return requestUserRolesStr[0]
	}

	return ""
}

func requestUserVerifiedFromMetadata(md metadata.MD) string {
	requestUserVerifiedStr := md.Get(requestUserVerified)
	if len(requestUserVerifiedStr) > 0 {
		return requestUserVerifiedStr[0]
	}

	return ""
}
