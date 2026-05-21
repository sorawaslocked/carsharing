package model

import (
	"time"

	sharedmodel "carsharing/shared/model"
)

type User struct {
	ID           string
	Email        string
	PhoneNumber  *string
	FirstName    string
	LastName     string
	BirthDate    time.Time
	PasswordHash []byte
	ProfileImage sharedmodel.Image

	Roles              []sharedmodel.Role
	IsDocumentVerified bool
	IsEmailVerified    bool
	IsSuspended        bool

	CreatedAt time.Time
	UpdatedAt time.Time
}

type UserUpdate struct {
	Email           *string
	PhoneNumber     *string
	FirstName       *string
	LastName        *string
	BirthDate       *time.Time
	PasswordHash    []byte
	ProfileImageKey *string

	Roles              []sharedmodel.Role
	IsDocumentVerified *bool
	IsEmailVerified    *bool
	IsSuspended        *bool
	UpdatedAt          time.Time
}

type UserFilter struct {
	Email       *string
	PhoneNumber *string
	FirstName   *string
	LastName    *string

	IsDocumentVerified *bool
	IsEmailVerified    *bool
	IsSuspended        *bool

	Pagination *sharedmodel.Pagination
}
