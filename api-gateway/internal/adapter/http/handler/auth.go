package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/sorawaslocked/car-rental-api-gateway/internal/adapter/http/handler/dto"
)

type Auth struct {
	svc AuthService
}

func NewAuth(svc AuthService) *Auth {
	return &Auth{svc: svc}
}

func (handler *Auth) Register(ctx *gin.Context) {
	cred, err := dto.FromRegisterRequest(ctx)
	if err != nil {
		dto.MalformedJson(ctx)

		return
	}

	id, err := handler.svc.Register(ctx, cred)
	if err != nil {
		dto.FromError(ctx, err)

		return
	}

	res := dto.RegisterResponse{ID: &id}
	dto.Ok(ctx, res)
}

func (handler *Auth) Login(ctx *gin.Context) {
	cred, err := dto.FromLoginRequest(ctx)
	if err != nil {
		dto.MalformedJson(ctx)

		return
	}

	token, err := handler.svc.Login(ctx, cred)
	if err != nil {
		dto.FromError(ctx, err)

		return
	}

	res := dto.LoginResponse{
		AccessToken:  &token.AccessToken,
		RefreshToken: &token.RefreshToken,
	}
	dto.Ok(ctx, res)
}

func (handler *Auth) RefreshToken(ctx *gin.Context) {
	cred, err := dto.FromRefreshTokenRequest(ctx)
	if err != nil {
		dto.MalformedJson(ctx)

		return
	}

	token, err := handler.svc.RefreshToken(ctx, cred)
	if err != nil {
		dto.FromError(ctx, err)

		return
	}

	res := dto.RefreshTokenResponse{
		AccessToken:  &token.AccessToken,
		RefreshToken: &token.RefreshToken,
	}
	dto.Ok(ctx, res)
}
