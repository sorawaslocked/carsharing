package dto

import (
	"errors"
	svc "github.com/sorawaslocked/car-rental-protos/gen/service"
	"github.com/sorawaslocked/car-rental-user-service/internal/model"
	"time"
)

func FromRegisterRequest(req *svc.RegisterRequest) (model.Credentials, model.ValidationErrors) {
	birthDate, err := time.Parse("2006-01-02", req.BirthDate)
	if err != nil {
		return model.Credentials{}, model.ValidationErrors{
			"email": model.ErrInvalidDateFormat,
		}
	}

	return model.Credentials{
		Email:                req.Email,
		PhoneNumber:          req.PhoneNumber,
		Password:             req.Password,
		PasswordConfirmation: req.PasswordConfirmation,
		FirstName:            req.FirstName,
		LastName:             req.LastName,
		BirthDate:            birthDate,
	}, nil
}

func ToRegisterResponse(id uint64, err error) *svc.RegisterResponse {
	if err == nil {
		return &svc.RegisterResponse{
			Id: &id,
		}
	}

	var validationErrors model.ValidationErrors
	if errors.As(err, &validationErrors) {
		return &svc.RegisterResponse{
			ValidationErrors: ToGrpcValidationError(validationErrors),
		}
	}

	errStr := err.Error()

	return &svc.RegisterResponse{
		Error: &errStr,
	}
}

func FromLoginRequest(req *svc.LoginRequest) model.Credentials {
	cred := model.Credentials{
		Password: req.Password,
	}
	if req.Email != nil {
		cred.Email = *req.Email
	}
	if req.PhoneNumber != nil {
		cred.PhoneNumber = *req.PhoneNumber
	}

	return cred
}

func ToLoginResponse(token model.Token, err error) *svc.LoginResponse {
	if err == nil {
		return &svc.LoginResponse{
			AccessToken:  &token.AccessToken,
			RefreshToken: &token.RefreshToken,
		}
	}

	var validationErrors model.ValidationErrors
	if errors.As(err, &validationErrors) {
		return &svc.LoginResponse{
			ValidationErrors: ToGrpcValidationError(validationErrors),
		}
	}

	errStr := err.Error()

	return &svc.LoginResponse{
		Error: &errStr,
	}
}

func ToRefreshTokenResponse(token model.Token, err error) *svc.RefreshTokenResponse {
	if err == nil {
		return &svc.RefreshTokenResponse{
			AccessToken:  &token.AccessToken,
			RefreshToken: &token.RefreshToken,
		}
	}

	var validationErrors model.ValidationErrors
	if errors.As(err, &validationErrors) {
		return &svc.RefreshTokenResponse{
			ValidationErrors: ToGrpcValidationError(validationErrors),
		}
	}

	errStr := err.Error()

	return &svc.RefreshTokenResponse{
		Error: &errStr,
	}
}
