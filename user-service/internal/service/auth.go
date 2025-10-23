package service

import (
	"car-rental-user-service/internal/model"
	"car-rental-user-service/internal/pkg/jwt"
	"car-rental-user-service/internal/pkg/logger"
	"car-rental-user-service/internal/pkg/security"
	"context"
	"github.com/go-playground/validator/v10"
	"log/slog"
	"time"
)

type AuthService struct {
	log         *slog.Logger
	validate    *validator.Validate
	jwtProvider *jwt.Provider
	userRepo    UserRepository
}

func NewAuthService(
	log *slog.Logger,
	validate *validator.Validate,
	jwtProvider *jwt.Provider,
	userRepo UserRepository,
) *AuthService {
	return &AuthService{
		log:         log,
		validate:    validate,
		jwtProvider: jwtProvider,
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
		s.log.Error("hashing password", logger.Err(err))

		return 0, model.ErrBcrypt
	}

	user := model.User{
		Email:        cred.Email,
		PhoneNumber:  cred.PhoneNumber,
		FirstName:    cred.FirstName,
		LastName:     cred.LastName,
		BirthDate:    cred.BirthDate,
		PasswordHash: passwordHash,
		Roles:        []model.Role{model.RoleUser},
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
		IsActive:     false,
		IsConfirmed:  false,
	}

	createdID, err := s.userRepo.Insert(ctx, user)
	if err != nil {
		s.log.Error("sql: inserting user", logger.Err(err))
		err = model.ErrSql
	}

	return createdID, err
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
		s.log.Error("sql: finding user", logger.Err(err))

		return model.Token{}, model.ErrNotFound
	}

	err = security.CheckPassword(cred.Password, user.PasswordHash)
	if err != nil {
		return model.Token{}, model.ErrPasswordsDoNotMatch
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

func (s *AuthService) RefreshToken(ctx context.Context, refreshToken string) (model.Token, error) {
	err := s.validate.Var(refreshToken, "required")
	if err != nil {
		return model.Token{}, model.ErrRequiredField
	}
	err = s.validate.Var(refreshToken, "jwt")
	if err != nil {
		return model.Token{}, model.ErrInvalidToken
	}

	claims, err := s.jwtProvider.VerifyAndParseClaims(refreshToken)
	if err != nil {
		s.log.Error(
			"jwt: verifying refresh token",
			logger.Err(err),
			slog.String("refreshToken", refreshToken),
		)

		return model.Token{}, model.ErrJwt
	}

	newAccessToken, err := s.jwtProvider.GenerateAccessToken(claims.ID, claims.Roles)
	if err != nil {
		s.log.Error(
			"jwt: generating access token",
			logger.Err(err),
			slog.Uint64("userId", claims.ID),
		)

		return model.Token{}, model.ErrJwt
	}
	newRefreshToken, err := s.jwtProvider.GenerateRefreshToken(claims.ID, claims.Roles)
	if err != nil {
		s.log.Error(
			"jwt: generating refresh token",
			logger.Err(err),
			slog.Uint64("userId", claims.ID),
		)

		return model.Token{}, model.ErrJwt
	}

	return model.Token{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
	}, nil
}
