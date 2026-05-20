package auth

import (
	"context"
	"log/slog"

	"carsharing/car-service/internal/adapter/grpc/dto"
	"carsharing/car-service/internal/model"
	sharedmodel "carsharing/shared/model"
	pkglog "carsharing/shared/pkg/log"
	"carsharing/shared/pkg/utils"

	"google.golang.org/grpc"
)

type methodPolicy struct {
	public       bool
	allowedRoles []sharedmodel.Role
}

type AuthInterceptor struct {
	log      *slog.Logger
	policies map[string]methodPolicy
}

func NewAuthInterceptor(log *slog.Logger) *AuthInterceptor {
	return &AuthInterceptor{
		log:      pkglog.WithComponent(log, "grpc.AuthInterceptor"),
		policies: buildPolicies(),
	}
}

func (i *AuthInterceptor) Unary(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	md := utils.MetadataFromCtx(ctx)
	log := pkglog.WithMetadata(pkglog.WithMethod(i.log, info.FullMethod), md)

	policy, known := i.policies[info.FullMethod]
	if !known {
		return nil, dto.FromErrorToStatusCode(model.ErrInsufficientPermissions)
	}

	if policy.public {
		return handler(ctx, req)
	}

	if md.UserID == nil {
		return nil, dto.FromErrorToStatusCode(model.ErrUnauthenticated)
	}

	// No role restrictions — any authenticated caller may proceed.
	if len(policy.allowedRoles) == 0 {
		return handler(ctx, req)
	}

	// Role check: any matching role grants access.
	for _, allowed := range policy.allowedRoles {
		for _, callerRole := range md.UserRoles {
			if callerRole == allowed {
				return handler(ctx, req)
			}
		}
	}

	log.Warn("permission denied", slog.Any("userRoles", md.UserRoles))
	return nil, dto.FromErrorToStatusCode(model.ErrInsufficientPermissions)
}
