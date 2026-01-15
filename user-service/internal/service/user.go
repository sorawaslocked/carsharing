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
	log                   *slog.Logger
	validate              *validator.Validate
	jwtProvider           JwtProvider
	userRepo              UserRepository
	activationCodeStorage ActivationCodeStorage
	mailer                Mailer
}

func NewUserService(
	log *slog.Logger,
	validate *validator.Validate,
	jwtProvider JwtProvider,
	userRepo UserRepository,
	activationCodeStorage ActivationCodeStorage,
	mailer Mailer,
) *UserService {
	return &UserService{
		log:                   log,
		validate:              validate,
		jwtProvider:           jwtProvider,
		userRepo:              userRepo,
		activationCodeStorage: activationCodeStorage,
		mailer:                mailer,
	}
}

func (s *UserService) Insert(ctx context.Context, data model.UserCreateData) (uint64, error) {
	err := validateInput(s.validate, data)
	if err != nil {
		return 0, err
	}

	_, err = s.userRepo.FindOne(ctx, model.UserFilter{Email: &data.Email})
	if err == nil {
		return 0, model.ErrDuplicateEmail
	}

	passwordHash, err := security.HashString(data.Password)
	if err != nil {
		s.log.Error("bcrypt: hashing password", logger.Err(err))

		return 0, model.ErrBcrypt
	}

	user := model.User{
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
	formatFilter(&filter)

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
	formatFilter(&filter)

	_, err = s.FindOne(ctx, filter)
	if err != nil {
		return err
	}

	err = validateInput(s.validate, data)
	if err != nil {
		return err
	}

	update := model.UserUpdate{
		Email:       data.Email,
		PhoneNumber: data.PhoneNumber,
		FirstName:   data.FirstName,
		LastName:    data.LastName,
		BirthDate:   data.BirthDate,
		Roles:       data.Roles,
		UpdatedAt:   time.Now(),
		IsActive:    data.IsActive,
		IsConfirmed: data.IsConfirmed,
	}

	if data.Password != nil {
		passwordHash, err := security.HashString(*data.Password)
		if err != nil {
			s.log.Error("bcrypt: hashing password", logger.Err(err))

			return model.ErrBcrypt
		}
		update.PasswordHash = &passwordHash
	}

	err = s.userRepo.Update(ctx, filter, update)

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
	formatFilter(&filter)

	return s.userRepo.Delete(ctx, filter)
}

func (s *UserService) Me(ctx context.Context) (model.User, error) {
	id, ok := ctx.Value("userID").(uint64)
	if !ok {
		return model.User{}, model.ErrInvalidToken
	}

	filter := model.UserFilter{
		ID: &id,
	}

	return s.userRepo.FindOne(ctx, filter)
}

func (s *UserService) SendActivationCode(ctx context.Context) error {
	id, ok := ctx.Value("userID").(uint64)
	if !ok {
		return model.ErrInvalidToken
	}

	filter := model.UserFilter{
		ID: &id,
	}

	user, err := s.userRepo.FindOne(ctx, filter)
	if err != nil {
		return err
	}

	code, err := s.activationCodeStorage.Create(ctx, user.ID)
	if err != nil {
		return err
	}

	err = s.mailer.SendActivationCode(ctx, user.Email, code)
	if err != nil {
		s.log.Error("Mailer", logger.Err(err))

		return err
	}

	return nil
}

func (s *UserService) CheckActivationCode(ctx context.Context, code string) error {
	id, ok := ctx.Value("userID").(uint64)
	if !ok {
		return model.ErrInvalidToken
	}

	filter := model.UserFilter{
		ID: &id,
	}

	user, err := s.userRepo.FindOne(ctx, filter)
	if err != nil {
		return err
	}

	if user.IsActive {
		return model.ErrActivatedUser
	}

	codeValidation := &activationCodeValidation{
		Code: code,
	}

	err = validateInput(s.validate, codeValidation)
	if err != nil {
		return err
	}

	codeHash, err := s.activationCodeStorage.Get(ctx, id)
	if err != nil {
		return err
	}

	err = security.CheckStringHash(code, codeHash)
	if err != nil {
		return model.ValidationErrors{
			"code": model.ErrInvalidActivationCode,
		}
	}

	isActive := true
	update := model.UserUpdate{
		UpdatedAt: time.Now(),
		IsActive:  &isActive,
	}

	err = s.userRepo.Update(ctx, filter, update)
	if err != nil {
		return err
	}

	return nil
}
