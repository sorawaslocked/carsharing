package handler

import (
	"context"

	"github.com/sorawaslocked/car-rental-user-service/internal/model"
)

type UserService interface {
	Create(ctx context.Context, data model.UserCreate) (string, error)
	Get(ctx context.Context, id string) (model.User, error)
	List(ctx context.Context, filter model.UserFilter) ([]model.User, error)
	Update(ctx context.Context, id string, data model.UserUpdate) error
	Delete(ctx context.Context, id string) error

	Register(ctx context.Context, data model.UserCreate) (string, error)
	SignIn(ctx context.Context, creds model.Credentials) (string, error)

	SendActivationCode(ctx context.Context, userID string) error
	CheckActivationCode(ctx context.Context, userID, code string) error
}
