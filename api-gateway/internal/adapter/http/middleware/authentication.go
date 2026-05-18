package middleware

import (
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sorawaslocked/car-rental-api-gateway/internal/adapter/http/dto"
	"github.com/sorawaslocked/car-rental-api-gateway/internal/model"
)

const (
	ctxUserIDKey               = "x-user-id"
	ctxUserRolesKey            = "x-user-roles"
	ctxUserDocumentVerifiedKey = "x-user-document-verified"
	ctxUserEmailVerifiedKey    = "x-user-email-verified"
	ctxUserSuspendedKey        = "x-user-suspended"
	ctxTokenExpKey             = "x-token-exp"
)

type Authentication struct {
	tokenParser          TokenParser
	userPermissionsCache UserPermissionsCache
	userSessionCache     UserSessionCache
}

func NewAuthentication(
	tokenParser TokenParser,
	userPermissionsCache UserPermissionsCache,
	userSessionCache UserSessionCache,
) *Authentication {
	return &Authentication{
		tokenParser:          tokenParser,
		userPermissionsCache: userPermissionsCache,
		userSessionCache:     userSessionCache,
	}
}

func (a *Authentication) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		header, err := authHeader(c)
		if err != nil {
			dto.FromError(c, err)
			c.Abort()

			return
		}

		userID, exp, err := a.parseClaims(header)
		if err != nil {
			dto.FromError(c, err)
			c.Abort()

			return
		}

		deviceID := c.GetString(ctxDeviceIDKey)

		isSignedIn, err := a.userSessionCache.IsSignedIn(c, userID, deviceID)
		if err != nil {
			dto.FromError(c, err)
			c.Abort()

			return
		}
		if !isSignedIn {
			dto.FromError(c, model.ErrUnauthorized)
			c.Abort()

			return
		}

		roles, err := a.userPermissionsCache.GetRoles(c, userID)
		if err != nil {
			dto.FromError(c, model.ErrInternalServerError)
			c.Abort()

			return
		}

		isDocumentVerified, err := a.userPermissionsCache.IsDocumentVerified(c, userID)
		if err != nil {
			dto.FromError(c, model.ErrInternalServerError)
			c.Abort()

			return
		}

		isEmailVerified, err := a.userPermissionsCache.IsEmailVerified(c, userID)
		if err != nil {
			dto.FromError(c, model.ErrInternalServerError)
			c.Abort()

			return
		}

		isSuspended, err := a.userPermissionsCache.IsSuspended(c, userID)
		if err != nil {
			dto.FromError(c, model.ErrInternalServerError)
			c.Abort()

			return
		}

		c.Set(ctxUserIDKey, userID)
		c.Set(ctxUserRolesKey, roles)
		c.Set(ctxUserDocumentVerifiedKey, isDocumentVerified)
		c.Set(ctxUserEmailVerifiedKey, isEmailVerified)
		c.Set(ctxUserSuspendedKey, isSuspended)
		c.Set(ctxTokenExpKey, exp)

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

func (a *Authentication) parseClaims(authHeader string) (userID string, exp time.Time, err error) {
	token := strings.TrimPrefix(authHeader, "Bearer ")

	userID, exp, err = a.tokenParser.ParseToken(token)
	if err != nil {
		return "", time.Time{}, model.ErrUnauthorized
	}

	return userID, exp, nil
}
