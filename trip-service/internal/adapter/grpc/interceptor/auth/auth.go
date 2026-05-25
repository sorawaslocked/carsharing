package auth

import (
	"context"
	"log/slog"

	"google.golang.org/grpc"

	pkglog "carsharing/shared/pkg/log"
	pkgutils "carsharing/shared/pkg/utils"
	"carsharing/trip-service/internal/adapter/grpc/dto"
	"carsharing/trip-service/internal/model"
)

type Interceptor struct {
	log      *slog.Logger
	policies map[string]methodPolicy
}

func NewInterceptor(log *slog.Logger) *Interceptor {
	return &Interceptor{
		log:      pkglog.WithComponent(log, "grpc.interceptor.auth.Interceptor"),
		policies: buildPolicies(),
	}
}

func (i *Interceptor) Unary(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	md := pkgutils.MetadataFromCtx(ctx)
	log := pkglog.WithMetadata(pkglog.WithMethod(i.log, info.FullMethod), md)

	policy, known := i.policies[info.FullMethod]
	if !known {
		return nil, dto.ToStatusError(model.ErrInsufficientPermissions)
	}
	if policy.public {
		return handler(ctx, req)
	}
	if md.UserID == nil {
		return nil, dto.ToStatusError(model.ErrUnauthenticated)
	}
	if len(policy.allowedRoles) == 0 {
		return handler(ctx, req)
	}
	for _, allowed := range policy.allowedRoles {
		for _, callerRole := range md.UserRoles {
			if callerRole == allowed {
				return handler(ctx, req)
			}
		}
	}
	log.Warn("permission denied", slog.Any("userRoles", md.UserRoles))
	return nil, dto.ToStatusError(model.ErrInsufficientPermissions)
}

func (i *Interceptor) Stream(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	md := pkgutils.MetadataFromCtx(ss.Context())
	log := pkglog.WithMetadata(pkglog.WithMethod(i.log, info.FullMethod), md)

	policy, known := i.policies[info.FullMethod]
	if !known {
		return dto.ToStatusError(model.ErrInsufficientPermissions)
	}
	if policy.public {
		return handler(srv, ss)
	}
	if md.UserID == nil {
		return dto.ToStatusError(model.ErrUnauthenticated)
	}
	if len(policy.allowedRoles) == 0 {
		return handler(srv, ss)
	}
	for _, allowed := range policy.allowedRoles {
		for _, callerRole := range md.UserRoles {
			if callerRole == allowed {
				return handler(srv, ss)
			}
		}
	}
	log.Warn("permission denied", slog.Any("userRoles", md.UserRoles))
	return dto.ToStatusError(model.ErrInsufficientPermissions)
}
