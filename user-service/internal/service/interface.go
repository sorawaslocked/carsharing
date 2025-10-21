package service

import (
	"car-rental-user-service/internal/model"
	"context"
)

type UserRepository interface {
	Insert(ctx context.Context, user model.User) (uint64, error)
	Find(ctx context.Context, filter model.UserFilter) (model.User, error)
	Update(ctx context.Context, update model.UserUpdateData) (uint64, error)
	Delete(ctx context.Context, filter model.UserFilter) (uint64, error)
}
