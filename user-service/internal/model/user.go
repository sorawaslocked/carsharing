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
	ProfileImage *Image

	Roles              []sharedmodel.Role
	IsDocumentVerified bool
	IsEmailVerified    bool
	IsSuspended        bool

	CreatedAt time.Time
	UpdatedAt time.Time
}

// UserUpdate is the repo-layer struct with pre-processed fields (hashed
// password, resolved image URL). Constructed by the service before calling the repo.
type UserUpdate struct {
	Email           *string
	PhoneNumber     *string
	FirstName       *string
	LastName        *string
	BirthDate       *time.Time
	PasswordHash    *[]byte
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

	Pagination *Pagination
}

type Pagination struct {
	Limit  int64
	Offset int64
}

type Credentials struct {
	Email       *string `validate:"required_without=PhoneNumber,omitempty,email"`
	PhoneNumber *string `validate:"required_without=Email,omitempty,e164"`
	Password    string  `validate:"required"`
}
