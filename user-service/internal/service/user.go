package service

import (
	"context"
	"github.com/go-playground/validator/v10"
	"github.com/sorawaslocked/car-rental-user-service/internal/model"
	"log/slog"
	"time"
)

type UserService struct {
	log         *slog.Logger
	validate    *validator.Validate
	jwtProvider JwtProvider
	userRepo    UserRepository
}

func NewUserService(
	log *slog.Logger,
	validate *validator.Validate,
	jwtProvider JwtProvider,
	userRepo UserRepository,
) *UserService {
	return &UserService{
		log:         log,
		validate:    validate,
		jwtProvider: jwtProvider,
		userRepo:    userRepo,
	}
}

func (s *UserService) Insert(ctx context.Context, user model.User) (uint64, error) {
	user.Roles = []model.Role{model.RoleUser}
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	user.IsActive = false
	user.IsConfirmed = false

	createdID, err := s.userRepo.Insert(ctx, user)

	return createdID, err
}

func (s *UserService) FindOne(ctx context.Context, filter model.UserFilter, jwtToken string) (model.User, error) {
	return model.User{}, nil
}

func (s *UserService) Find(ctx context.Context, filter model.UserFilter, jwtToken string) ([]model.User, error) {
	return nil, nil
}

func (s *UserService) Update(ctx context.Context, filter model.UserFilter, update model.UserUpdateData, jwtToken string) error {
	return nil
}

func (s *UserService) Delete(ctx context.Context, filter model.UserFilter, jwtToken string) error {
	return nil
}
