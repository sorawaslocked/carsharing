package service

import (
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/sorawaslocked/car-rental-user-service/internal/model"
	"time"
)

type registerValidation struct {
	Email                string    `validate:"required,email"`
	PhoneNumber          string    `validate:"omitempty,e164"`
	Password             string    `validate:"required,min=8,max=20,complex_password"`
	PasswordConfirmation string    `validate:"required,min=8,max=20,complex_password,eqfield=Password"`
	FirstName            string    `validate:"required,min=1,max=100,alphaunicode"`
	LastName             string    `validate:"required,min=1,max=100,alphaunicode"`
	BirthDate            time.Time `validate:"required,min_age=18"`
}

type loginValidation struct {
	Email       string `validate:"required_without=PhoneNumber,omitempty,email"`
	PhoneNumber string `validate:"required_without=Email,omitempty,e164"`
	Password    string `validate:"required"`
}

type refreshTokenValidation struct {
	RefreshToken string `validate:"required,jwt"`
}

func validateInput(v *validator.Validate, input any) error {
	err := v.Struct(input)
	if err == nil {
		return nil
	}

	var errs model.ValidationErrors
	errs = make(map[string]error)
	var validationErrors validator.ValidationErrors
	errors.As(err, &validationErrors)

	for _, fieldErr := range validationErrors {
		fieldErrField := fieldErr.Field()
		field := uncapitalize(fieldErrField)

		if _, ok := errs[field]; ok {
			continue
		}
		errs[field] = validationError(fieldErr)
	}

	return errs
}
