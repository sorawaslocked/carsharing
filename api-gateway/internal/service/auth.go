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

func (s *AuthService) Register(ctx context.Context, data model.UserCreateData) (uint64, error) {
	return s.presenter.Register(ctx, data)
}

func (s *AuthService) Login(ctx context.Context, cred model.Credentials) (model.Token, error) {
	return s.presenter.Login(ctx, cred)
}

func (s *AuthService) RefreshToken(ctx context.Context, refreshToken string) (model.Token, error) {
	return s.presenter.RefreshToken(ctx, refreshToken)
}

func (s *AuthService) Logout(ctx context.Context, refreshToken string) error {
	return s.presenter.Logout(ctx, refreshToken)
}
