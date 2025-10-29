package service

import (
	"context"
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/sorawaslocked/car-rental-user-service/internal/model"
	"github.com/sorawaslocked/car-rental-user-service/internal/pkg/logger"
	"github.com/sorawaslocked/car-rental-user-service/internal/pkg/security"
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

func (s *UserService) Insert(ctx context.Context, data model.UserCreateData) (uint64, error) {
	err := validateInput(s.validate, data)
	if err != nil {
		return 0, err
	}

	passwordHash, err := security.HashPassword(data.Password)
	if err != nil {
		s.log.Error("bcrypt: hashing password", logger.Err(err))

		return 0, model.ErrBcrypt
	}

	user := model.User{
		ID:           0,
		Email:        data.Email,
		PhoneNumber:  data.PhoneNumber,
		FirstName:    data.FirstName,
		LastName:     data.LastName,
		BirthDate:    data.BirthDate,
		PasswordHash: passwordHash,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	// TODO: add permission checks
	if data.Roles == nil {
		user.Roles = []model.Role{model.RoleUser}
	}
	if data.IsActive != nil {
		user.IsActive = *data.IsActive
	}
	if data.IsConfirmed != nil {
		user.IsConfirmed = *data.IsConfirmed
	}

	return s.userRepo.Insert(ctx, user)
}

func (s *UserService) FindOne(ctx context.Context, filter model.UserFilter) (model.User, error) {
	err := checkQueryParams(s.validate, filter)
	if err != nil {
		return model.User{}, err
	}

	user, err := s.userRepo.FindOne(ctx, filter)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return model.User{}, model.ErrNotFound
		}
		s.log.Error("sql: finding user", logger.Err(err))

		return model.User{}, model.ErrSql
	}

	return user, nil
}

func (s *UserService) Find(ctx context.Context, filter model.UserFilter) ([]model.User, error) {
	users, err := s.userRepo.Find(ctx, filter)
	if err != nil {
		s.log.Error("sql: finding users", logger.Err(err))

		return []model.User{}, model.ErrSql
	}

	return users, nil
}

func (s *UserService) Update(ctx context.Context, filter model.UserFilter, data model.UserUpdateData) error {
	err := checkQueryParams(s.validate, filter)
	if err != nil {
		return err
	}

	err = validateInput(s.validate, data)
	if err != nil {
		return err
	}

	passwordHash, err := security.HashPassword(*data.Password)
	if err != nil {
		s.log.Error("bcrypt: hashing password", logger.Err(err))

		return model.ErrBcrypt
	}

	err = s.userRepo.Update(ctx, filter, model.UserUpdate{
		Email:        data.Email,
		PhoneNumber:  data.PhoneNumber,
		FirstName:    data.FirstName,
		LastName:     data.LastName,
		BirthDate:    data.BirthDate,
		PasswordHash: &passwordHash,
		Roles:        data.Roles,
		UpdatedAt:    time.Now(),
		IsActive:     data.IsActive,
		IsConfirmed:  data.IsConfirmed,
	})

	if err != nil {
		if errors.Is(err, model.ErrDuplicateEmail) {
			return model.ValidationErrors{
				"email": model.ErrDuplicateEmail,
			}
		}

		return err
	}

	return nil
}

func (s *UserService) Delete(ctx context.Context, filter model.UserFilter) error {
	err := checkQueryParams(s.validate, filter)
	if err != nil {
		return err
	}

	return s.userRepo.Delete(ctx, filter)
}
