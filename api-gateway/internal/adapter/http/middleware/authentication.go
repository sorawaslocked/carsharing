package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sorawaslocked/car-rental-api-gateway/internal/adapter/http/dto"
	"github.com/sorawaslocked/car-rental-api-gateway/internal/model"
)

const (
	ctxUserIDKey        = "x-user-id"
	ctxUserRolesKey     = "x-user-roles"
	ctxUserVerifiedKey  = "x-user-verified"
	ctxUserSuspendedKey = "x-user-suspended"
)

type Authentication struct {
	tokenManager         TokenManager
	userPermissionsCache UserPermissionsCache
}

func NewAuthentication(tokenManager TokenManager, userPermissionsCache UserPermissionsCache) *Authentication {
	return &Authentication{
		tokenManager:         tokenManager,
		userPermissionsCache: userPermissionsCache,
	}
}

func (a *Authentication) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		header, err := authHeader(c)
		if err != nil {
			dto.FromError(c, err)

			return
		}

		userID, err := a.parseClaims(header)
		if err != nil {
			dto.FromError(c, err)

			return
		}

		ctx := c.Request.Context()

		roles, err := a.userPermissionsCache.GetRoles(ctx, userID)
		if err != nil {
			dto.FromError(c, model.ErrInternalServerError)

			return
		}

		verified, err := a.userPermissionsCache.GetVerified(ctx, userID)
		if err != nil {
			dto.FromError(c, model.ErrInternalServerError)

			return
		}

		suspended, err := a.userPermissionsCache.GetSuspended(ctx, userID)
		if err != nil {
			dto.FromError(c, model.ErrInternalServerError)

			return
		}

		c.Set(ctxUserIDKey, userID)
		c.Set(ctxUserRolesKey, roles)
		c.Set(ctxUserVerifiedKey, verified)
		c.Set(ctxUserSuspendedKey, suspended)

		c.Next()
	}
}

func authHeader(c *gin.Context) (string, error) {
	header := c.GetHeader("Authorization")
	if header == "" {
		return "", model.ErrUnauthorized
	}

	return header, nil
}

func (a *Authentication) parseClaims(authHeader string) (userID string, err error) {
	token := strings.TrimPrefix(authHeader, "Bearer ")

	userID, err = a.tokenManager.ParseToken(token)
	if err != nil {
		return "", model.ErrUnauthorized
	}

	return userID, nil
}
