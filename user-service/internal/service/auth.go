package service

import (
	"car-rental-user-service/internal/model"
	"context"
	"errors"
	"github.com/go-playground/validator/v10"
	"strings"
)

type AuthService struct {
	validate *validator.Validate
	userRepo UserRepository
}

func NewAuthService(
	validate *validator.Validate,
	userRepo UserRepository,
) *AuthService {
	return &AuthService{
		validate: validate,
		userRepo: userRepo,
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

	return 0, nil
}

func (s *AuthService) Login(ctx context.Context, cred model.Credentials) (model.Token, map[string]error) {
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

	return model.Token{}, nil
}

func (s *AuthService) RefreshToken(ctx context.Context, refreshToken string) (model.Token, error) {
	if refreshToken == "" {
		return model.Token{}, ErrRequiredField
	}

	err := s.validate.Var(refreshToken, "jwt")
	if err != nil {
		return model.Token{}, ErrInvalidToken
	}

	return model.Token{}, nil
}
