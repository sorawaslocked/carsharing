package service

import (
	"context"
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/sorawaslocked/car-rental-user-service/internal/model"
	"github.com/sorawaslocked/car-rental-user-service/internal/pkg/logger"
	"github.com/sorawaslocked/car-rental-user-service/internal/pkg/security"
	"log/slog"
)

type AuthService struct {
	log         *slog.Logger
	validate    *validator.Validate
	jwtProvider JwtProvider
	userService UserService
	userRepo    UserRepository
}

func NewAuthService(
	log *slog.Logger,
	validate *validator.Validate,
	jwtProvider JwtProvider,
	userService UserService,
	userRepo UserRepository,
) *AuthService {
	return &AuthService{
		log:         log,
		validate:    validate,
		jwtProvider: jwtProvider,
		userService: userService,
		userRepo:    userRepo,
	}
}

func (s *AuthService) Register(ctx context.Context, cred model.Credentials) (uint64, error) {
	input := registerValidation{
		Email:                cred.Email,
		PhoneNumber:          cred.PhoneNumber,
		Password:             cred.Password,
		PasswordConfirmation: cred.PasswordConfirmation,
		FirstName:            cred.FirstName,
		LastName:             cred.LastName,
		BirthDate:            cred.BirthDate,
	}
	err := validateInput(s.validate, input)
	if err != nil {
		return 0, err
	}

	s.log.Info("registering user", slog.String("email", cred.Email))
	passwordHash, err := security.HashPassword(cred.Password)
	if err != nil {
		s.log.Error("bcrypt: hashing password", logger.Err(err))

		return 0, model.ErrBcrypt
	}

	user := model.User{
		Email:        cred.Email,
		PhoneNumber:  cred.PhoneNumber,
		FirstName:    cred.FirstName,
		LastName:     cred.LastName,
		BirthDate:    cred.BirthDate,
		PasswordHash: passwordHash,
	}

	createdID, err := s.userService.Insert(ctx, user)
	if err != nil {
		if errors.Is(err, model.ErrDuplicateEmail) {
			return 0, model.ValidationErrors{
				"email": model.ErrDuplicateEmail,
			}
		}

		s.log.Error("sql: inserting user", logger.Err(err))

		return 0, model.ErrSql
	}

	return createdID, nil
}

func (s *AuthService) Login(ctx context.Context, cred model.Credentials) (model.Token, error) {
	input := loginValidation{
		Email:       cred.Email,
		PhoneNumber: cred.PhoneNumber,
		Password:    cred.Password,
	}
	err := validateInput(s.validate, input)
	if err != nil {
		return model.Token{}, err
	}

	filter := model.UserFilter{}
	if input.Email != "" {
		filter.Email = &input.Email
		s.log.Info("logging in user", slog.String("email", input.Email))
	}
	if cred.PhoneNumber != "" {
		filter.PhoneNumber = &input.PhoneNumber
		s.log.Info("logging in user", slog.String("phoneNumber", input.PhoneNumber))
	}

	// TODO: add not found error handling
	user, err := s.userRepo.FindOne(ctx, filter)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return model.Token{}, model.ErrNotFound
		}
		s.log.Error("sql: finding user", logger.Err(err))

		return model.Token{}, model.ErrSql
	}

	err = security.CheckPassword(cred.Password, user.PasswordHash)
	if err != nil {
		return model.Token{}, model.ValidationErrors{
			"password": model.ErrPasswordsDoNotMatch,
		}
	}

	userRoles := toRoleStrings(user.Roles)
	accessToken, err := s.jwtProvider.GenerateAccessToken(user.ID, userRoles)
	if err != nil {
		s.log.Error(
			"jwt: generating access token",
			logger.Err(err),
			slog.Uint64("userId", user.ID),
		)

		return model.Token{}, model.ErrJwt
	}
	refreshToken, err := s.jwtProvider.GenerateRefreshToken(user.ID, userRoles)
	if err != nil {
		s.log.Error(
			"jwt: generating refresh token",
			logger.Err(err),
			slog.Uint64("userId", user.ID),
		)

		return model.Token{}, model.ErrJwt
	}

	return model.Token{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *AuthService) RefreshToken(_ context.Context, refreshToken string) (model.Token, error) {
	input := refreshTokenValidation{
		RefreshToken: refreshToken,
	}
	err := validateInput(s.validate, input)
	if err != nil {
		return model.Token{}, err
	}

	id, roles, err := s.jwtProvider.VerifyAndParseClaims(refreshToken)
	if err != nil {
		s.log.Error(
			"jwt: verifying refresh token",
			logger.Err(err),
			slog.String("refreshToken", refreshToken),
		)

		return model.Token{}, model.ErrJwt
	}

	newAccessToken, err := s.jwtProvider.GenerateAccessToken(id, roles)
	if err != nil {
		s.log.Error(
			"jwt: generating access token",
			logger.Err(err),
			slog.Uint64("userId", id),
		)

		return model.Token{}, model.ErrJwt
	}
	newRefreshToken, err := s.jwtProvider.GenerateRefreshToken(id, roles)
	if err != nil {
		s.log.Error(
			"jwt: generating refresh token",
			logger.Err(err),
			slog.Uint64("userId", id),
		)

		return model.Token{}, model.ErrJwt
	}

	return model.Token{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
	}, nil
}
