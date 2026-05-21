package validation

import (
	"time"

	sharedmodel "carsharing/shared/model"
)

type UserCreate struct {
	Email                string    `validate:"required,email"`
	PhoneNumber          *string   `validate:"omitempty,e164"`
	FirstName            string    `validate:"required,min=1,max=100"`
	LastName             string    `validate:"required,min=1,max=100"`
	BirthDate            time.Time `validate:"required,min_age=18"`
	Password             string    `validate:"required,min=8,max=20,complex_password"`
	PasswordConfirmation string    `validate:"required,min=8,max=20,complex_password,eqfield=Password"`
}

type UserUpdate struct {
	Email                *string    `validate:"omitempty,email"`
	PhoneNumber          *string    `validate:"omitempty,e164"`
	FirstName            *string    `validate:"omitempty,min=1,max=100"`
	LastName             *string    `validate:"omitempty,min=1,max=100"`
	BirthDate            *time.Time `validate:"omitempty,min_age=18"`
	Password             *string    `validate:"omitempty,min=8,max=20,complex_password"`
	PasswordConfirmation *string    `validate:"omitempty,min=8,max=20,complex_password"`
	ProfileImageKey      *string

	Roles              []sharedmodel.Role
	IsDocumentVerified *bool
	IsEmailVerified    *bool
	IsSuspended        *bool
}
