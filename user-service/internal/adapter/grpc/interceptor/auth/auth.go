package auth

import (
	"context"
	"log/slog"

	pkglog "carsharing/shared/pkg/log"
	"carsharing/shared/pkg/utils"
	"carsharing/user-service/internal/adapter/grpc/dto"
	"carsharing/user-service/internal/model"

	"google.golang.org/grpc"
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
	md := utils.MetadataFromCtx(ctx)
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

	// No role or owner restrictions — any authenticated caller may proceed.
	if len(policy.allowedRoles) == 0 && policy.ownerExtract == nil {
		return handler(ctx, req)
	}

	// Role check: any matching privileged role grants access.
	for _, allowed := range policy.allowedRoles {
		for _, callerRole := range md.UserRoles {
			if callerRole == allowed {
				return handler(ctx, req)
			}
		}
	}

	// Owner check: caller operates on their own resource.
	if policy.ownerExtract != nil {
		if targetID, ok := policy.ownerExtract(req); ok && targetID == *md.UserID {
			return handler(ctx, req)
		}
	}

	log.Warn("permission denied", slog.Any("userRoles", md.UserRoles))

	return nil, dto.ToStatusError(model.ErrInsufficientPermissions)
}
