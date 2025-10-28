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

type UserService interface {
	Insert(ctx context.Context, user model.User) (uint64, error)
	FindOne(ctx context.Context, filter model.UserFilter, jwtToken string) (model.User, error)
	Find(ctx context.Context, filter model.UserFilter, jwtToken string) ([]model.User, error)
	Update(ctx context.Context, filter model.UserFilter, update model.UserUpdateData, jwtToken string) error
	Delete(ctx context.Context, filter model.UserFilter, jwtToken string) error
}
