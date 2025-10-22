package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/sorawaslocked/car-rental-api-gateway/internal/adapter/http/handler/dto"
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
	res := dto.RegisterResponse{
		Errors: errors,
	}

	if id != 0 {
		res.ID = &id
	}
	if res.Errors == nil {
		created(ctx, res)

		return
	}
	if _, exists := res.Errors["grpc"]; exists {
		internalServerError(ctx)

		return
	}

	badRequest(ctx, res)
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
	res := dto.LoginResponse{
		Errors: errors,
	}
	if token.AccessToken != "" && token.RefreshToken != "" {
		res.AccessToken = &token.AccessToken
		res.RefreshToken = &token.RefreshToken
	}
	if res.Errors == nil {
		ok(ctx, res)

		return
	}
	if _, exists := res.Errors["grpc"]; exists {
		internalServerError(ctx)

		return
	}

	badRequest(ctx, res)
}

func (handler *Auth) RefreshToken(ctx *gin.Context) {
	var req dto.RefreshTokenRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		malformedJson(ctx)

		return
	}

	token, errors := handler.svc.RefreshToken(ctx, req.RefreshToken)

	res := dto.RefreshTokenResponse{
		Errors: errors,
	}
	if token.AccessToken != "" && token.RefreshToken != "" {
		res.AccessToken = &token.AccessToken
		res.RefreshToken = &token.RefreshToken
	}
	if res.Errors == nil {
		ok(ctx, res)

		return
	}
	if _, exists := res.Errors["grpc"]; exists {
		internalServerError(ctx)

		return
	}

	badRequest(ctx, res)
}
