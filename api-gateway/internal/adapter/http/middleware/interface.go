package middleware

import (
	"context"
	"time"
)

type TokenParser interface {
	ParseToken(token string) (userID string, exp time.Time, err error)
}

type UserPermissionsCache interface {
	GetRoles(ctx context.Context, userID string) ([]string, error)
	IsDocumentVerified(ctx context.Context, userID string) (bool, error)
	IsEmailVerified(ctx context.Context, userID string) (bool, error)
	IsSuspended(ctx context.Context, userID string) (bool, error)
}

type UserSessionCache interface {
	IsSignedIn(ctx context.Context, userID, deviceID string) (bool, error)
}
