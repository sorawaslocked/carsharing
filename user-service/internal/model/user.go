package model

import "time"

type User struct {
	ID           uint64
	Email        string
	PhoneNumber  string
	FirstName    string
	LastName     string
	BirthDate    time.Time
	PasswordHash []byte
	Roles        []Role
	CreatedAt    time.Time
	UpdatedAt    time.Time

	IsActive    bool
	IsConfirmed bool
}

type UserFilter struct {
	ID          *uint64
	Email       *string
	PhoneNumber *string
	FirstName   *string
	LastName    *string
	Roles       []Role

	IsActive    *bool
	IsConfirmed *bool
}

type UserUpdateData struct {
	Email                *string    `validate:"omitempty,email"`
	PhoneNumber          *string    `validate:"omitempty,e164"`
	FirstName            *string    `validate:"omitempty,min=1,max=100,alphaunicode"`
	LastName             *string    `validate:"omitempty,min=1,max=100,alphaunicode"`
	BirthDate            *time.Time `validate:"omitempty,min_age=18"`
	Password             *string    `validate:"omitempty,min=8,max=20,complex_password"`
	PasswordConfirmation *string    `validate:"required_with=Password,min=8,max=20,complex_password"`
	Roles                []Role
	UpdatedAt            time.Time

	IsActive    *bool
	IsConfirmed *bool
}
