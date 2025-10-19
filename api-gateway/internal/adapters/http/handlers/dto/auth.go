package dto

import (
	"github.com/gin-gonic/gin"
)

type RegisterRequest struct {
	Email                string `json:"email"`
	PhoneNumber          string `json:"phoneNumber"`
	Password             string `json:"password"`
	PasswordConfirmation string `json:"passwordConfirmation"`
	FirstName            string `json:"firstName"`
	LastName             string `json:"lastName"`
	DateOfBirth          string `json:"dateOfBirth"`
}

type RegisterResponse struct {
	ID     uint64            `json:"id"`
	Errors map[string]string `json:"errors,omitempty"`
}

type LoginRequest struct {
	Email       *string `json:"email"`
	PhoneNumber *string `json:"phoneNumber"`
	Password    string  `json:"password"`
}

type LoginResponse struct {
	AccessToken  string            `json:"accessToken"`
	RefreshToken string            `json:"refreshToken"`
	Errors       map[string]string `json:"errors,omitempty"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refreshToken"`
}

type RefreshTokenResponse struct {
	AccessToken  string            `json:"accessToken"`
	RefreshToken string            `json:"refreshToken"`
	Errors       map[string]string `json:"errors,omitempty"`
}

func FromRefreshTokenRequest(ctx *gin.Context) (string, error) {
	var req RefreshTokenRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		return "", err
	}

	return req.RefreshToken, nil
}
