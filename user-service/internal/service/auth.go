package service

import (
	"car-rental-user-service/internal/model"
	"car-rental-user-service/internal/pkg/jwt"
	"car-rental-user-service/internal/pkg/logger"
	"car-rental-user-service/internal/pkg/security"
	"context"
	"errors"
	"github.com/go-playground/validator/v10"
	"log/slog"
	"strings"
	"time"
)

type AuthService struct {
	log         *slog.Logger
	validate    *validator.Validate
	jwtProvider *jwt.JwtProvider
	userRepo    UserRepository
}

func NewAuthService(
	log *slog.Logger,
	validate *validator.Validate,
	jwtProvider *jwt.JwtProvider,
	userRepo UserRepository,
) *AuthService {
	return &AuthService{
		log:         log,
		validate:    validate,
		jwtProvider: jwtProvider,
		userRepo:    userRepo,
	}
}

func (s *AuthService) Register(ctx context.Context, cred model.Credentials) (uint64, map[string]error) {
	errs := make(map[string]error)

	if cred.Email == nil {
		errs["email"] = model.ErrRequiredField
	}
	if cred.PhoneNumber == nil {
		errs["phoneNumber"] = model.ErrRequiredField
	}
	if cred.Password == "" {
		errs["password"] = model.ErrRequiredField
	}
	if cred.PasswordConfirmation == nil {
		errs["passwordConfirmation"] = model.ErrRequiredField
	}
	if cred.FirstName == nil {
		errs["firstName"] = model.ErrRequiredField
	}
	if cred.LastName == nil {
		errs["lastName"] = model.ErrRequiredField
	}
	if cred.BirthDate == nil {
		errs["birthDate"] = model.ErrRequiredField
	}
	if len(errs) != 0 {
		return 0, errs
	}

	err := s.validate.Struct(cred)
	if err != nil {
		var validationErrors validator.ValidationErrors
		errors.As(err, &validationErrors)

		for _, fieldErr := range validationErrors {
			var field string

			switch fieldErr.Field() {
			case "PhoneNumber":
				field = "phoneNumber"
			case "PasswordConfirmation":
				field = "passwordConfirmation"
			case "FirstName":
				field = "firstName"
			case "LastName":
				field = "lastName"
			case "BirthDate":
				field = "birthDate"
			default:
				field = strings.ToLower(fieldErr.Field())
			}

			if _, ok := errs[field]; ok {
				continue
			}
			errs[field] = validationError(fieldErr)
		}
	}
	if len(errs) != 0 {
		return 0, errs
	}

	s.log.Info("registering user", slog.String("email", *cred.Email))

	passwordHash, err := security.HashPassword(cred.Password)
	if err != nil {
		s.log.Error("hashing password", logger.Err(err))
		errs["bcrypt"] = model.ErrBcrypt

		return 0, errs
	}

	user := model.User{
		Email:        *cred.Email,
		PhoneNumber:  *cred.PhoneNumber,
		FirstName:    *cred.FirstName,
		LastName:     *cred.LastName,
		BirthDate:    *cred.BirthDate,
		PasswordHash: passwordHash,
		Role:         model.RoleUser,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
		IsActive:     false,
		IsConfirmed:  false,
	}

	createdID, err := s.userRepo.Insert(ctx, user)
	if err != nil {
		s.log.Error("inserting user", logger.Err(err))
		errs["repository"] = err
	}

	return createdID, errs
}

func (s *AuthService) Login(ctx context.Context, cred model.Credentials) (model.Token, map[string]error) {
	errs := make(map[string]error)

	if cred.Email == nil && cred.PhoneNumber == nil {
		errs["phoneNumber"] = model.ErrRequiredField
		errs["email"] = model.ErrRequiredField
	}
	if cred.Password == "" {
		errs["password"] = model.ErrRequiredField
	}
	if len(errs) != 0 {
		return model.Token{}, errs
	}

	err := s.validate.Struct(cred)
	if err != nil {
		var validationErrors validator.ValidationErrors
		errors.As(err, &validationErrors)

		for _, fieldErr := range validationErrors {
			var field string

			switch fieldErr.Field() {
			case "PhoneNumber":
				field = "phoneNumber"
			default:
				field = strings.ToLower(fieldErr.Field())
			}

			if _, ok := errs[field]; ok {
				continue
			}
			errs[field] = validationError(fieldErr)
		}
	}
	if len(errs) != 0 {
		return model.Token{}, errs
	}

	filter := model.UserFilter{}
	switch {
	case cred.Email != nil:
		filter.Email = cred.Email
		s.log.Info("logging in user", slog.String("email", *cred.Email))
	case cred.PhoneNumber != nil:
		filter.PhoneNumber = cred.PhoneNumber
		s.log.Info("logging in user", slog.String("phoneNumber", *cred.PhoneNumber))
	}

	// TODO: add not found error handling
	user, err := s.userRepo.FindOne(ctx, filter)
	if err != nil {
		s.log.Error("finding user", logger.Err(err))
		errs["repository"] = model.ErrNotFound

		return model.Token{}, errs
	}

	err = security.CheckPassword(cred.Password, user.PasswordHash)
	if err != nil {
		errs["password"] = model.ErrPasswordsDoNotMatch

		return model.Token{}, errs
	}

	accessToken, err := s.jwtProvider.GenerateAccessToken(user.ID, user.Role.String())
	if err != nil {
		s.log.Error(
			"generating access token",
			logger.Err(err),
			slog.Uint64("userId", user.ID),
		)
		errs["jwt"] = model.ErrJwt

		return model.Token{}, errs
	}
	refreshToken, err := s.jwtProvider.GenerateRefreshToken(user.ID, user.Role.String())
	if err != nil {
		s.log.Error(
			"generating refresh token",
			logger.Err(err),
			slog.Uint64("userId", user.ID),
		)
		errs["jwt"] = model.ErrJwt

		return model.Token{}, errs
	}

	return model.Token{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *AuthService) RefreshToken(ctx context.Context, refreshToken string) (model.Token, error) {
	if refreshToken == "" {
		return model.Token{}, model.ErrRequiredField
	}

	err := s.validate.Var(refreshToken, "jwt")
	if err != nil {
		return model.Token{}, model.ErrInvalidToken
	}

	claims, err := s.jwtProvider.VerifyAndParseClaims(refreshToken)
	if err != nil {
		s.log.Error(
			"verifying refresh token",
			logger.Err(err),
			slog.String("refreshToken", refreshToken),
		)

		return model.Token{}, model.ErrJwt
	}

	newAccessToken, err := s.jwtProvider.GenerateAccessToken(claims.UserID, claims.Role)
	if err != nil {
		s.log.Error(
			"generating access token",
			logger.Err(err),
			slog.Uint64("userId", claims.UserID),
		)

		return model.Token{}, model.ErrJwt
	}
	newRefreshToken, err := s.jwtProvider.GenerateRefreshToken(claims.UserID, claims.Role)
	if err != nil {
		s.log.Error(
			"generating refresh token",
			logger.Err(err),
			slog.Uint64("userId", claims.UserID),
		)

		return model.Token{}, model.ErrJwt
	}

	return model.Token{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
	}, nil
}
