package middleware

import (
	"log/slog"
	"strings"
	"time"

	"carsharing/api-gateway/internal/adapter/http/dto"
	"carsharing/api-gateway/internal/model"
	sharedmodel "carsharing/shared/model"
	pkglog "carsharing/shared/pkg/log"
	"carsharing/shared/pkg/utils"

	"github.com/gin-gonic/gin"
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
	log                  *slog.Logger
}

func NewAuthentication(
	tokenParser TokenParser,
	userPermissionsCache UserPermissionsCache,
	userSessionCache UserSessionCache,
	logger *slog.Logger,
) *Authentication {
	return &Authentication{
		tokenParser:          tokenParser,
		userPermissionsCache: userPermissionsCache,
		userSessionCache:     userSessionCache,
		log:                  pkglog.WithComponent(logger, "middleware.Authentication"),
	}
}

func (a *Authentication) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		log := pkglog.WithMetadata(pkglog.WithMethod(a.log, "Middleware"), utils.MetadataFromCtx(c))

		log.Debug("authenticating request")

		header, err := authHeader(c)
		if err != nil {
			log.Debug("missing authorization header")
			dto.FromError(c, err)
			c.Abort()

			return
		}

		userID, exp, err := a.parseClaims(c, header)
		if err != nil {
			log.Warn("parsing token claims", pkglog.Err(err))
			dto.FromError(c, err)
			c.Abort()

			return
		}

		log.Debug("token parsed", slog.String("userID", userID))

		deviceID := c.GetString(ctxDeviceIDKey)

		isSignedIn, err := a.userSessionCache.IsSignedIn(c, userID, deviceID)
		if err != nil {
			log.Error("checking session", pkglog.Err(err))
			dto.FromError(c, err)
			c.Abort()

			return
		}
		if !isSignedIn {
			log.Debug("session not found", slog.String("userID", userID), slog.String("deviceID", deviceID))
			dto.FromError(c, model.ErrUnauthorized)
			c.Abort()

			return
		}

		rawRoles, err := a.userPermissionsCache.GetRoles(c, userID)
		if err != nil {
			log.Error("getting roles", pkglog.Err(err))
			dto.FromError(c, model.ErrInternalServerError)
			c.Abort()

			return
		}
		roles := make([]sharedmodel.Role, len(rawRoles))
		for i, r := range rawRoles {
			roles[i] = sharedmodel.Role(r)
		}

		isDocumentVerified, err := a.userPermissionsCache.IsDocumentVerified(c, userID)
		if err != nil {
			log.Error("getting document verified", pkglog.Err(err))
			dto.FromError(c, model.ErrInternalServerError)
			c.Abort()

			return
		}

		isEmailVerified, err := a.userPermissionsCache.IsEmailVerified(c, userID)
		if err != nil {
			log.Error("getting email verified", pkglog.Err(err))
			dto.FromError(c, model.ErrInternalServerError)
			c.Abort()

			return
		}

		isSuspended, err := a.userPermissionsCache.IsSuspended(c, userID)
		if err != nil {
			log.Error("getting suspended", pkglog.Err(err))
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

		log.Debug("authentication complete", slog.String("userID", userID))

		c.Next()
	}
}

func authHeader(c *gin.Context) (string, error) {
	header := c.GetHeader("Authorization")
	if header != "" {
		return header, nil
	}

	// WebSocket clients cannot set custom headers, so accept the token as a query param.
	if token := c.Query("token"); token != "" {
		return "Bearer " + token, nil
	}

	return "", model.ErrUnauthorized
}

func (a *Authentication) parseClaims(c *gin.Context, authHeader string) (userID string, exp time.Time, err error) {
	token := strings.TrimPrefix(authHeader, "Bearer ")

	userID, exp, err = a.tokenParser.ParseToken(c.Request.Context(), token)
	if err != nil {
		return "", time.Time{}, model.ErrUnauthorized
	}

	return userID, exp, nil
}
