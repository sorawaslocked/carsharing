package handler

import (
	"car-rental-user-service/internal/model"
	"context"
)

type AuthService interface {
	Register(ctx context.Context, cred model.Credentials) (uint64, map[string]error)
	Login(ctx context.Context, cred model.Credentials) (model.Token, map[string]error)
	RefreshToken(ctx context.Context, refreshToken string) (model.Token, error)
}
