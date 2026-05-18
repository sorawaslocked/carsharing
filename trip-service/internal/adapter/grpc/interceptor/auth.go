package interceptor

import (
	"context"
	"log/slog"

	"google.golang.org/grpc"

	"carsharing/trip-service/internal/adapter/grpc/dto"
	"carsharing/trip-service/internal/model"
	pkglog "carsharing/trip-service/internal/pkg/log"
	"carsharing/trip-service/internal/pkg/utils"
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
		ctx = extractMetadata(ctx)
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
}

func (i *AuthInterceptor) Stream() grpc.StreamServerInterceptor {
	return func(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		ctx := extractMetadata(ss.Context())
		md := utils.MetadataFromCtx(ctx)
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
}
