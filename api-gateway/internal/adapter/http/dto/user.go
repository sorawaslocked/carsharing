package dto

import (
	"github.com/gin-gonic/gin"
	"github.com/sorawaslocked/car-rental-api-gateway/internal/model"
	"strconv"
	"time"
)

type User struct {
	ID           uint64    `json:"id"`
	Email        string    `json:"email"`
	PhoneNumber  string    `json:"phoneNumber,omitempty"`
	FirstName    string    `json:"firstName"`
	LastName     string    `json:"lastName"`
	BirthDate    string    `json:"birthDate"`
	PasswordHash []byte    `json:"passwordHash"`
	Roles        []string  `json:"roles"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
	IsActive     bool      `json:"isActive"`
	IsConfirmed  bool      `json:"isConfirmed"`
}

type UserCreateRequest struct {
	Email                string   `json:"email"`
	PhoneNumber          *string  `json:"phoneNumber"`
	Password             string   `json:"password"`
	PasswordConfirmation string   `json:"passwordConfirmation"`
	FirstName            string   `json:"firstName"`
	LastName             string   `json:"lastName"`
	BirthDate            string   `json:"birthDate"`
	Roles                []string `json:"roles"`
	IsActive             bool     `json:"isActive"`
	IsConfirmed          bool     `json:"isConfirmed"`
}

type UserCreateResponse struct {
	ID *uint64 `json:"id,omitempty"`
}

type UserGetResponse struct {
	User *User `json:"user,omitempty"`
}

type UserGetAllResponse struct {
	Users *[]User `json:"users,omitempty"`
}

type UserUpdateRequest struct {
	Email                *string   `json:"email"`
	PhoneNumber          *string   `json:"phoneNumber"`
	Password             *string   `json:"password"`
	PasswordConfirmation *string   `json:"passwordConfirmation"`
	FirstName            *string   `json:"firstName"`
	LastName             *string   `json:"lastName"`
	BirthDate            *string   `json:"birthDate"`
	Roles                *[]string `json:"roles"`
	IsActive             *bool     `json:"isActive"`
	IsConfirmed          *bool     `json:"isConfirmed"`
}

type CheckActivationCodeRequest struct {
	Code string `json:"code"`
}

func FromCreateUserRequest(ctx *gin.Context) (model.UserCreateData, error) {
	var req UserCreateRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		return model.UserCreateData{}, err
	}

	data := model.UserCreateData{
		Email:                req.Email,
		FirstName:            req.FirstName,
		LastName:             req.LastName,
		BirthDate:            req.BirthDate,
		Password:             req.Password,
		PasswordConfirmation: req.PasswordConfirmation,
		IsActive:             &req.IsActive,
		IsConfirmed:          &req.IsConfirmed,
	}
	if req.PhoneNumber != nil {
		data.PhoneNumber = *req.PhoneNumber
	}
	if req.Roles != nil {
		data.Roles = &req.Roles
	}

	return data, nil
}

func FromUpdateRequest(ctx *gin.Context) (model.UserUpdateData, error) {
	var req UserUpdateRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		return model.UserUpdateData{}, err
	}

	return model.UserUpdateData{
		Email:                req.Email,
		PhoneNumber:          req.PhoneNumber,
		Password:             req.Password,
		PasswordConfirmation: req.PasswordConfirmation,
		FirstName:            req.FirstName,
		LastName:             req.LastName,
		BirthDate:            req.BirthDate,
		Roles:                req.Roles,
		IsActive:             req.IsActive,
		IsConfirmed:          req.IsConfirmed,
	}, nil
}

func FromCheckActivationCodeRequest(ctx *gin.Context) (string, error) {
	var req CheckActivationCodeRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		return "", err
	}

	return req.Code, nil
}

func FilterFromCtx(ctx *gin.Context) (model.UserFilter, error) {
	idStr := ctx.Query("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil && idStr != "" {
		return model.UserFilter{}, model.ErrInvalidQueryParam
	}

	email := ctx.Query("email")
	filter := model.UserFilter{}

	if idStr != "" {
		filter.ID = &id
	}
	if email != "" {
		filter.Email = &email
	}

	return filter, nil
}
