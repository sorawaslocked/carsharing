package interceptor

import (
	"context"
	"log/slog"

	"carsharing/booking-service/internal/adapter/grpc/dto"
	"carsharing/booking-service/internal/model"
	pkglog "carsharing/shared/pkg/log"
	"carsharing/shared/pkg/utils"
	"google.golang.org/grpc"
)

type methodPolicy struct {
	public       bool
	allowedRoles []model.Role
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

func (i *AuthInterceptor) Unary() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		md := utils.MetadataFromCtx(ctx)
		log := pkglog.WithMetadata(pkglog.WithMethod(i.log, info.FullMethod), md)

		policy, known := i.policies[info.FullMethod]
		if !known {
			return nil, dto.ToGRPCError(model.ErrInsufficientPermissions)
		}

		if policy.public {
			return handler(ctx, req)
		}

		if md.UserID == nil {
			return nil, dto.ToGRPCError(model.ErrUnauthenticated)
		}

		if len(policy.allowedRoles) == 0 {
			return handler(ctx, req)
		}

		for _, allowed := range policy.allowedRoles {
			for _, callerRole := range md.UserRoles {
				if model.Role(callerRole) == allowed {
					return handler(ctx, req)
				}
			}
		}

		log.Warn("permission denied", slog.Any("userRoles", md.UserRoles))

		return nil, dto.ToGRPCError(model.ErrInsufficientPermissions)
	}
}
