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

type UserPresenter interface {
	Create(ctx context.Context, data model.UserCreateData) (uint64, error)
	Get(ctx context.Context, filter model.UserFilter) (model.User, error)
	GetAll(ctx context.Context, filter model.UserFilter) ([]model.User, error)
	Update(ctx context.Context, filter model.UserFilter, data model.UserUpdateData) error
	Delete(ctx context.Context, filter model.UserFilter) error
	Me(ctx context.Context) (model.User, error)
}
