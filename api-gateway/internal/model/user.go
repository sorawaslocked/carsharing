package model

import "time"

type User struct {
	ID              string
	Email           string
	PhoneNumber     *string
	FirstName       string
	LastName        string
	BirthDate       string
	Password        Password
	ProfileImageURL *string

	Roles              []string
	IsDocumentVerified bool
	IsEmailVerified    bool
	IsSuspended        bool

	CreatedAt time.Time
	UpdatedAt time.Time
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

type UserCreate struct {
	Email       string
	PhoneNumber *string
	FirstName   string
	LastName    string
	BirthDate   string
	Password    Password
}

type UserUpdate struct {
	Email           *string
	PhoneNumber     *string
	FirstName       *string
	LastName        *string
	BirthDate       *string
	Password        Password
	ProfileImageKey *string

	Roles              []string
	IsDocumentVerified *bool
	IsEmailVerified    *bool
	IsSuspended        *bool
}

type UserProfileUpdate struct {
	PhoneNumber     *string
	FirstName       *string
	LastName        *string
	BirthDate       *string
	Password        Password
	ProfileImageKey *string
}

type Password struct {
	Hash             []byte
	Text             *string
	TextConfirmation *string
}

type Credentials struct {
	Email       *string
	PhoneNumber *string
	Password    Password
}

type AccessToken struct {
	Token     string
	ExpiresIn int64
}

type RefreshToken struct {
	Token     string
	ExpiresIn int64
}

type Document struct {
	ID        string
	UserID    string
	ImageType string
	Status    string
	Reason    *string
	ImageURL  string

	CreatedAt time.Time
	UpdatedAt time.Time
}
