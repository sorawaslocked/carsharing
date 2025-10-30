package dto

import (
	"github.com/gin-gonic/gin"
	"github.com/sorawaslocked/car-rental-api-gateway/internal/model"
)

type RegisterRequest struct {
	Email                string `json:"email"`
	PhoneNumber          string `json:"phoneNumber"`
	Password             string `json:"password"`
	PasswordConfirmation string `json:"passwordConfirmation"`
	FirstName            string `json:"firstName"`
	LastName             string `json:"lastName"`
	BirthDate            string `json:"birthDate"`
}

type RegisterResponse struct {
	ID *uint64 `json:"id,omitempty"`
}

type LoginRequest struct {
	Email       *string `json:"email"`
	PhoneNumber *string `json:"phoneNumber"`
	Password    string  `json:"password"`
}

type LoginResponse struct {
	AccessToken  *string `json:"accessToken,omitempty"`
	RefreshToken *string `json:"refreshToken,omitempty"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refreshToken"`
}

type RefreshTokenResponse struct {
	AccessToken  *string `json:"accessToken,omitempty"`
	RefreshToken *string `json:"refreshToken,omitempty"`
}

func FromRegisterRequest(ctx *gin.Context) (model.UserCreateData, error) {
	var req RegisterRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		return model.UserCreateData{}, err
	}

	return model.UserCreateData{
		Email:                req.Email,
		PhoneNumber:          req.PhoneNumber,
		Password:             req.Password,
		PasswordConfirmation: req.PasswordConfirmation,
		FirstName:            req.FirstName,
		LastName:             req.LastName,
		BirthDate:            req.BirthDate,
	}, nil
}

func FromLoginRequest(ctx *gin.Context) (model.Credentials, error) {
	var req LoginRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		return model.Credentials{}, err
	}

	cred := model.Credentials{
		Password: req.Password,
	}
	if req.Email != nil {
		cred.Email = *req.Email
	}
	if req.PhoneNumber != nil {
		cred.PhoneNumber = *req.PhoneNumber
	}

	return cred, nil
}

func FromRefreshTokenRequest(ctx *gin.Context) (string, error) {
	var req RefreshTokenRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		return "", err
	}

	return req.RefreshToken, nil
}
