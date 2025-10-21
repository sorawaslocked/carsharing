package service

import (
	"car-rental-user-service/internal/model"
	"car-rental-user-service/internal/pkg/jwt"
	"car-rental-user-service/internal/pkg/security"
	"context"
	"errors"
	"github.com/go-playground/validator/v10"
	"strings"
	"time"
)

type AuthService struct {
	validate    *validator.Validate
	jwtProvider *jwt.JwtProvider
	userRepo    UserRepository
}

func NewAuthService(
	validate *validator.Validate,
	jwtProvider *jwt.JwtProvider,
	userRepo UserRepository,
) *AuthService {
	return &AuthService{
		validate:    validate,
		jwtProvider: jwtProvider,
		userRepo:    userRepo,
	}
}

func (s *AuthService) Register(ctx context.Context, cred model.Credentials) (uint64, map[string]error) {
	errs := make(map[string]error)

	if cred.Email == nil {
		errs["email"] = ErrRequiredField
	}
	if cred.PhoneNumber == nil {
		errs["phoneNumber"] = ErrRequiredField
	}
	if cred.Password == "" {
		errs["password"] = ErrRequiredField
	}
	if cred.PasswordConfirmation == nil {
		errs["passwordConfirmation"] = ErrRequiredField
	}
	if cred.FirstName == nil {
		errs["firstName"] = ErrRequiredField
	}
	if cred.LastName == nil {
		errs["lastName"] = ErrRequiredField
	}
	if cred.BirthDate == nil {
		errs["birthDate"] = ErrRequiredField
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

	passwordHash, err := security.HashPassword(cred.Password)
	if err != nil {
		errs["bcrypt"] = ErrBcrypt

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
		errs["repository"] = err
	}

	return createdID, errs
}

func (s *AuthService) Login(ctx context.Context, cred model.Credentials) (model.Token, map[string]error) {
	errs := make(map[string]error)

	if cred.Email == nil && cred.PhoneNumber == nil {
		errs["phoneNumber"] = ErrRequiredField
		errs["email"] = ErrRequiredField
	}
	if cred.Password == "" {
		errs["password"] = ErrRequiredField
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
	case cred.PhoneNumber != nil:
		filter.PhoneNumber = cred.PhoneNumber
	}

	// TODO: add not found error handling
	user, err := s.userRepo.Find(ctx, filter)
	if err != nil {
		errs["repository"] = ErrNotFound

		return model.Token{}, errs
	}

	err = security.CheckPassword(cred.Password, user.PasswordHash)
	if err != nil {
		errs["password"] = ErrPasswordsDoNotMatch

		return model.Token{}, errs
	}

	accessToken, err := s.jwtProvider.GenerateAccessToken(user.ID, user.Role.String())
	if err != nil {
		errs["jwt"] = ErrJwt

		return model.Token{}, errs
	}
	refreshToken, err := s.jwtProvider.GenerateRefreshToken(user.ID, user.Role.String())
	if err != nil {
		errs["jwt"] = ErrJwt

		return model.Token{}, errs
	}

	return model.Token{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *AuthService) RefreshToken(ctx context.Context, refreshToken string) (model.Token, error) {
	if refreshToken == "" {
		return model.Token{}, ErrRequiredField
	}

	err := s.validate.Var(refreshToken, "jwt")
	if err != nil {
		return model.Token{}, ErrInvalidToken
	}

	claims, err := s.jwtProvider.VerifyAndParseClaims(refreshToken)
	if err != nil {
		return model.Token{}, ErrJwt
	}

	newAccessToken, err := s.jwtProvider.GenerateAccessToken(claims.UserID, claims.Role)
	if err != nil {
		return model.Token{}, ErrJwt
	}
	newRefreshToken, err := s.jwtProvider.GenerateRefreshToken(claims.UserID, claims.Role)
	if err != nil {
		return model.Token{}, ErrJwt
	}

	return model.Token{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
	}, nil
}
