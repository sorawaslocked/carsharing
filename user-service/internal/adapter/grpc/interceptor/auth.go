package interceptor

import (
	"context"
	usersvc "github.com/sorawaslocked/car-rental-protos/gen/service/user"
	"github.com/sorawaslocked/car-rental-user-service/internal/adapter/grpc/dto"
	"github.com/sorawaslocked/car-rental-user-service/internal/model"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"strings"
)

type _claims struct {
	id    uint64
	roles []model.Role
}

// AuthInterceptor is a middleware struct to handle authorization and authentication
type AuthInterceptor struct {
	jwtProvider    JwtProvider
	permittedRoles map[string]map[model.Role]bool // permittedRoles maps endpoints to the roles which can access it
}

func NewAuthInterceptor(jwtProvider JwtProvider) *AuthInterceptor {
	return &AuthInterceptor{
		jwtProvider:    jwtProvider,
		permittedRoles: createPermittedRoles(),
	}
}

func (i *AuthInterceptor) Unary(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, dto.ToStatusCodeError(model.ErrMissingMetadata)
	}

	claims, err := i.authenticateAndGetClaims(md["authorization"], info.FullMethod)
	if err != nil {
		return nil, dto.ToStatusCodeError(err)
	}
	// No-auth method
	if claims.id < 1 {
		return handler(ctx, req)
	}

	err = i.authorize(claims, info.FullMethod)
	if err != nil {
		return nil, dto.ToStatusCodeError(err)
	}

	ctx = context.WithValue(ctx, "userID", claims.id)

	m, err := handler(ctx, req)
	if err != nil {
		return nil, err
	}

	if info.FullMethod == UserServiceGet {
		res := m.(*usersvc.GetResponse)

		if res.User.ID != claims.id {
			return nil, dto.ToStatusCodeError(model.ErrInsufficientPermissions)
		}
	}

	return m, err
}

func (i *AuthInterceptor) authenticateAndGetClaims(authorization []string, method string) (_claims, error) {
	if i.permittedRoles[method] == nil {
		return _claims{}, nil
	}

	if len(authorization) < 1 {
		return _claims{}, model.ErrInvalidToken
	}

	token := strings.TrimPrefix(authorization[0], "Bearer ")

	id, roleStrings, err := i.jwtProvider.VerifyAndParseClaims(token)
	if err != nil {
		return _claims{}, model.ErrInvalidToken
	}

	roles := make([]model.Role, len(roleStrings))
	for idx, roleString := range roleStrings {
		role, err := model.FromStringToRole(roleString)
		if err != nil {
			return _claims{}, model.ErrInvalidToken
		}

		roles[idx] = role
	}

	return _claims{
		id:    id,
		roles: roles,
	}, nil
}

func (i *AuthInterceptor) authorize(claims _claims, method string) error {
	permittedRolesForMethod := i.permittedRoles[method]

	for _, role := range claims.roles {
		if _, ok := permittedRolesForMethod[role]; ok {
			return nil
		}
	}

	return model.ErrInsufficientPermissions
}
