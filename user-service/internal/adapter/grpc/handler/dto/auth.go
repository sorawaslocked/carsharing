package dto

import (
	authsvc "github.com/sorawaslocked/car-rental-protos/gen/service/auth"
	"github.com/sorawaslocked/car-rental-user-service/internal/model"
	"time"
)

func FromRegisterRequest(req *authsvc.RegisterRequest) (model.UserCreateData, model.ValidationErrors) {
	birthDate, err := time.Parse("2006-01-02", req.BirthDate)
	if err != nil {
		return model.UserCreateData{}, model.ValidationErrors{
			"email": model.ErrInvalidDateFormat,
		}
	}

	return model.UserCreateData{
		Email:                req.Email,
		PhoneNumber:          req.PhoneNumber,
		Password:             req.Password,
		PasswordConfirmation: req.PasswordConfirmation,
		FirstName:            req.FirstName,
		LastName:             req.LastName,
		BirthDate:            birthDate,
	}, nil
}

func FromLoginRequest(req *authsvc.LoginRequest) model.Credentials {
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
