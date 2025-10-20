package service

import (
	"context"
	"github.com/sorawaslocked/car-rental-api-gateway/internal/model"
)

type AuthService struct {
	presenter AuthPresenter
}

func NewAuthService(presenter AuthPresenter) *AuthService {
	return &AuthService{presenter: presenter}
}

func (svc *AuthService) Register(ctx context.Context, cred model.Credentials) (uint64, map[string]string) {
	return svc.presenter.Register(ctx, cred)
}

func (svc *AuthService) Login(ctx context.Context, cred model.Credentials) (model.Token, map[string]string) {
	return svc.presenter.Login(ctx, cred)
}

func (svc *AuthService) RefreshToken(ctx context.Context, refreshToken string) (model.Token, map[string]string) {
	return svc.presenter.RefreshToken(ctx, refreshToken)
}
