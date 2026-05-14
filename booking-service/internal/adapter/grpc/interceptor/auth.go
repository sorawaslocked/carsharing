package interceptor

import (
	"context"
	"log/slog"

	"github.com/sorawaslocked/car-rental-booking-service/internal/adapter/grpc/dto"
	"github.com/sorawaslocked/car-rental-booking-service/internal/model"
	pkglog "github.com/sorawaslocked/car-rental-booking-service/internal/pkg/log"
	"github.com/sorawaslocked/car-rental-booking-service/internal/pkg/utils"
	"google.golang.org/grpc"
)

// ownerExtractFn extracts the target user ID from the request so the interceptor
// can check whether the caller is operating on their own resource.
type ownerExtractFn func(req any) (userID string, ok bool)

type methodPolicy struct {
	public       bool
	allowedRoles []model.Role
	ownerExtract ownerExtractFn
}

// Duck-typing interfaces — avoids importing concrete proto request types here.
type userIDCarrier interface{ GetUserId() string }

func extractByUserID(req any) (string, bool) {
	c, ok := req.(userIDCarrier)
	if !ok {
		return "", false
	}
	id := c.GetUserId()
	return id, id != ""
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

		return nil, dto.ToGRPCError(model.ErrInsufficientPermissions)
	}
}
