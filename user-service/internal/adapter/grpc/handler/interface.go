package handler

import (
	"context"
	"github.com/sorawaslocked/car-rental-user-service/internal/model"
)

type AuthService interface {
	Register(ctx context.Context, cred model.Credentials) (uint64, error)
	Login(ctx context.Context, cred model.Credentials) (model.Token, error)
	RefreshToken(ctx context.Context, refreshToken string) (model.Token, error)
}
