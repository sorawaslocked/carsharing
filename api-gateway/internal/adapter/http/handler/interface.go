package handler

import (
	"context"
	"github.com/sorawaslocked/car-rental-api-gateway/internal/model"
)

type AuthService interface {
	Register(ctx context.Context, cred model.Credentials) (uint64, map[string]string)
	Login(ctx context.Context, cred model.Credentials) (model.Token, map[string]string)
	RefreshToken(ctx context.Context, refreshToken string) (model.Token, map[string]string)
}
