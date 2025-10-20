package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/sorawaslocked/car-rental-api-gateway/internal/adapters/http/handlers/dto"
	"github.com/sorawaslocked/car-rental-api-gateway/internal/model"
)

type Auth struct {
	svc AuthService
}

func NewAuth(svc AuthService) *Auth {
	return &Auth{svc: svc}
}

func (handler *Auth) Register(ctx *gin.Context) {
	var req dto.RegisterRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		malformedJson(ctx)

		return
	}

	id, errors := handler.svc.Register(ctx, model.Credentials{
		FirstName:            &req.FirstName,
		LastName:             &req.LastName,
		DateOfBirth:          &req.DateOfBirth,
		Email:                &req.Email,
		PhoneNumber:          &req.PhoneNumber,
		Password:             req.Password,
		PasswordConfirmation: &req.PasswordConfirmation,
	})
	if len(errors) > 0 {
		badRequest(ctx, errors)

		return
	}

	body := make(map[string]any)
	body["id"] = id
	ok(ctx, body)
}

func (handler *Auth) Login(ctx *gin.Context) {
	var req dto.LoginRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		malformedJson(ctx)

		return
	}

	token, errors := handler.svc.Login(ctx, model.Credentials{
		Email:       req.Email,
		PhoneNumber: req.PhoneNumber,
		Password:    req.Password,
	})
	if len(errors) > 0 {
		badRequest(ctx, errors)

		return
	}

	body := make(map[string]any)
	body["accessToken"] = token.AccessToken
	body["refreshToken"] = token.RefreshToken

	ok(ctx, body)
}

func (handler *Auth) RefreshToken(ctx *gin.Context) {
	var req dto.RefreshTokenRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		malformedJson(ctx)

		return
	}

	token, err := handler.svc.RefreshToken(ctx, req.RefreshToken)
	if err != nil {
		errors := make(map[string]string)
		errors["refreshToken"] = err.Error()

		badRequest(ctx, errors)

		return
	}

	body := make(map[string]any)
	body["accessToken"] = token.AccessToken
	body["refreshToken"] = token.RefreshToken

	ok(ctx, body)
}
