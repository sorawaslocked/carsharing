package model

import "time"

type User struct {
	ID           uint64
	Email        string
	PhoneNumber  string
	FirstName    string
	LastName     string
	BirthDate    time.Time
	PasswordHash string
	Role         Role
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
	Role        *Role

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
	Role         *Role

	IsActive    *bool
	IsConfirmed *bool
}
