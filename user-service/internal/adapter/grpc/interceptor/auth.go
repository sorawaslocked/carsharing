package interceptor

import (
	"context"
	"log/slog"

	pkglog "carsharing/shared/pkg/log"
	"carsharing/shared/pkg/utils"
	"carsharing/user-service/internal/adapter/grpc/dto"
	"carsharing/user-service/internal/model"
	usersvc "github.com/sorawaslocked/car-rental-protos/gen/service/user"
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
type idCarrier interface{ GetId() string }
type userIDCarrier interface{ GetUserId() string }

func extractByID(req any) (string, bool) {
	c, ok := req.(idCarrier)
	if !ok {
		return "", false
	}
	id := c.GetId()
	return id, id != ""
}

func extractByUserID(req any) (string, bool) {
	c, ok := req.(userIDCarrier)
	if !ok {
		return "", false
	}
	id := c.GetUserId()
	return id, id != ""
}

var privilegedRoles = []model.Role{model.RoleAdmin, model.RoleUserManager}

func buildPolicies() map[string]methodPolicy {
	return map[string]methodPolicy{
		// Public — no authentication required.
		usersvc.HealthService_Health_FullMethodName: {public: true},
		usersvc.UserService_Register_FullMethodName: {public: true},
		usersvc.UserService_SignIn_FullMethodName:   {public: true},

		// Privileged roles only.
		usersvc.UserService_CreateUser_FullMethodName:    {allowedRoles: privilegedRoles},
		usersvc.UserService_ListUsers_FullMethodName:     {allowedRoles: privilegedRoles},
		usersvc.UserService_CheckDocument_FullMethodName: {allowedRoles: privilegedRoles},

		// Privileged roles OR the resource owner.
		usersvc.UserService_GetUser_FullMethodName:    {allowedRoles: privilegedRoles, ownerExtract: extractByID},
		usersvc.UserService_UpdateUser_FullMethodName: {allowedRoles: privilegedRoles, ownerExtract: extractByID},
		usersvc.UserService_DeleteUser_FullMethodName: {allowedRoles: privilegedRoles, ownerExtract: extractByID},
		usersvc.UserService_GetProcessedDocumentsForUser_FullMethodName: {
			allowedRoles: privilegedRoles,
			ownerExtract: extractByUserID,
		},

		// Any authenticated user.
		usersvc.UserService_SendActivationCode_FullMethodName:        {},
		usersvc.UserService_CheckActivationCode_FullMethodName:       {},
		usersvc.UserService_GetProfileImageUploadData_FullMethodName: {},
		usersvc.UserService_CreateDocument_FullMethodName:            {},
		usersvc.UserService_GetUploadDocumentData_FullMethodName:     {},
	}
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
			if model.Role(callerRole) == allowed {
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
