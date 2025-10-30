package service

import (
	"context"
	"github.com/sorawaslocked/car-rental-api-gateway/internal/model"
)

type AuthPresenter interface {
	Register(ctx context.Context, data model.UserCreateData) (uint64, error)
	Login(ctx context.Context, cred model.Credentials) (model.Token, error)
	RefreshToken(ctx context.Context, refreshToken string) (model.Token, error)
}
