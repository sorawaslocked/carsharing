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

func (s *UserService) Insert(ctx context.Context, user model.User) (uint64, error) {
	user.Roles = []model.Role{model.RoleUser}
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	user.IsActive = false
	user.IsConfirmed = false

	createdID, err := s.userRepo.Insert(ctx, user)

	return createdID, err
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

func (s *UserService) Update(ctx context.Context, filter model.UserFilter, updateData model.UserUpdateData) error {
	err := checkQueryParams(s.validate, filter)
	if err != nil {
		return err
	}

	err = validateInput(s.validate, updateData)
	if err != nil {
		return err
	}

	passwordHash, err := security.HashPassword(*updateData.Password)
	if err != nil {
		s.log.Error("bcrypt: hashing password", logger.Err(err))

		return model.ErrBcrypt
	}

	err = s.userRepo.Update(ctx, filter, model.UserUpdate{
		Email:        updateData.Email,
		PhoneNumber:  updateData.PhoneNumber,
		FirstName:    updateData.FirstName,
		LastName:     updateData.LastName,
		BirthDate:    updateData.BirthDate,
		PasswordHash: &passwordHash,
		Roles:        updateData.Roles,
		UpdatedAt:    time.Now(),
		IsActive:     updateData.IsActive,
		IsConfirmed:  updateData.IsConfirmed,
	})

	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
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
