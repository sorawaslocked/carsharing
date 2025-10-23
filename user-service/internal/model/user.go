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
	Roles       []*Role

	IsActive    *bool
	IsConfirmed *bool
}

type UserUpdateData struct {
	Email        *string
	PhoneNumber  *string
	FirstName    *string
	LastName     *string
	BirthDate    *time.Time
	PasswordHash *string
	Roles        []*Role
	UpdatedAt    time.Time

	IsActive    *bool
	IsConfirmed *bool
}
