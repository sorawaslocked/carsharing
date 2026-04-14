package middleware

import (
	"context"
	"time"
)

type UserPermissionsCache interface {
	GetRoles(ctx context.Context, userID string) ([]string, error)
	GetVerified(ctx context.Context, userID string) (bool, error)
	GetSuspended(ctx context.Context, userID string) (bool, error)
}

type TokenManager interface {
	GenerateAccessToken(userID string) (token string, exp time.Time, err error)
	GenerateRefreshToken(userID string) (token string, exp time.Time, err error)
	ParseToken(token string) (string, error)
}
