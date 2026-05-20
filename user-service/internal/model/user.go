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

// UserCreate is the input for both Create (admin) and Register (self-service).
type UserCreate struct {
	Email                string    `validate:"required,email"`
	PhoneNumber          *string   `validate:"omitempty,e164"`
	FirstName            string    `validate:"required,min=1,max=100"`
	LastName             string    `validate:"required,min=1,max=100"`
	BirthDate            time.Time `validate:"required,min_age=18"`
	Password             string    `validate:"required,min=8,max=20,complex_password"`
	PasswordConfirmation string    `validate:"required,min=8,max=20,complex_password,eqfield=Password"`
}

// UserUpdate is the service-layer input for partial user updates.
// Password fields are plaintext; the service hashes them before persisting.
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

// UserRepoUpdate is the repo-layer struct with pre-processed fields (hashed
// password, resolved image URL). Constructed by the service before calling the repo.
type UserRepoUpdate struct {
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
