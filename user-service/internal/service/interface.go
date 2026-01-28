package service

import (
	"context"
	"github.com/sorawaslocked/car-rental-user-service/internal/model"
	"time"
)

type UserRepository interface {
	Insert(ctx context.Context, user model.User) (uint64, error)
	FindOne(ctx context.Context, filter model.UserFilter) (model.User, error)
	Find(ctx context.Context, filter model.UserFilter) ([]model.User, error)
	Update(ctx context.Context, filter model.UserFilter, update model.UserUpdate) error
	Delete(ctx context.Context, filter model.UserFilter) error
}

type JwtProvider interface {
	GenerateAccessToken(id uint64, roles []string) (string, time.Time, error)
	GenerateRefreshToken(id uint64, roles []string) (string, time.Time, error)
	VerifyAndParseClaims(token string) (uint64, []string, error)
}

type SessionStorage interface {
	Save(ctx context.Context, userID uint64) error
	Exists(ctx context.Context, userID uint64) (bool, error)
	Delete(ctx context.Context, userID uint64) error
}

type ActivationCodeStorage interface {
	Save(ctx context.Context, userID uint64) (string, error)
	Get(ctx context.Context, userID uint64) ([]byte, error)
}

type Mailer interface {
	SendActivationCode(ctx context.Context, receiver, code string) error
}
