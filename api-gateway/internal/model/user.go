package model

import "time"

type User struct {
	ID           uint64
	Email        string
	PhoneNumber  string
	FirstName    string
	LastName     string
	BirthDate    string
	PasswordHash []byte
	Roles        []string
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
	Roles       []string

	IsActive    *bool
	IsConfirmed *bool
}

type UserUpdate struct {
	Email        *string
	PhoneNumber  *string
	FirstName    *string
	LastName     *string
	BirthDate    *string
	PasswordHash *[]byte
	Roles        *[]string
	UpdatedAt    time.Time

	IsActive    *bool
	IsConfirmed *bool
}

type UserCreateData struct {
	Email                string
	PhoneNumber          string
	Password             string
	PasswordConfirmation string
	FirstName            string
	LastName             string
	BirthDate            string
	Roles                *[]string

	IsActive    *bool
	IsConfirmed *bool
}

type UserUpdateData struct {
	Email                *string
	PhoneNumber          *string
	FirstName            *string
	LastName             *string
	BirthDate            *string
	Password             *string
	PasswordConfirmation *string
	Roles                *[]string

	IsActive    *bool
	IsConfirmed *bool
}
