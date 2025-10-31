package service

import (
	"context"
	"github.com/sorawaslocked/car-rental-api-gateway/internal/model"
)

type UserService struct {
	presenter UserPresenter
}

func NewUserService(presenter UserPresenter) UserService {
	return UserService{presenter: presenter}
}

func (s *UserService) Insert(ctx context.Context, data model.UserCreateData) (uint64, error) {
	return s.presenter.Create(ctx, data)
}

func (s *UserService) FindOne(ctx context.Context, filter model.UserFilter) (model.User, error) {
	return s.presenter.Get(ctx, filter)
}

func (s *UserService) Find(ctx context.Context, filter model.UserFilter) ([]model.User, error) {
	return s.presenter.GetAll(ctx, filter)
}

func (s *UserService) Update(ctx context.Context, filter model.UserFilter, data model.UserUpdateData) error {
	return s.presenter.Update(ctx, filter, data)
}

func (s *UserService) Delete(ctx context.Context, filter model.UserFilter) error {
	return s.presenter.Delete(ctx, filter)
}
