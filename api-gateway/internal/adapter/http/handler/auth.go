package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/sorawaslocked/car-rental-api-gateway/internal/adapter/http/dto"
	"github.com/sorawaslocked/car-rental-api-gateway/internal/config"
	"log"
	"net/http"
)

type Auth struct {
	svc    AuthService
	cookie config.Cookie
}

func NewAuth(svc AuthService, cookie config.Cookie) *Auth {
	return &Auth{
		svc:    svc,
		cookie: cookie,
	}
}

func (h *Auth) Register(ctx *gin.Context) {
	data, err := dto.FromRegisterRequest(ctx)
	if err != nil {
		dto.MalformedJson(ctx)

		return
	}

	id, err := h.svc.Register(ctx, data)
	if err != nil {
		dto.FromError(ctx, err)

		return
	}

	res := dto.RegisterResponse{ID: &id}
	dto.Created(ctx, res)
}

func (h *Auth) Login(ctx *gin.Context) {
	cred, err := dto.FromLoginRequest(ctx)
	if err != nil {
		dto.MalformedJson(ctx)

		return
	}

	token, err := h.svc.Login(ctx, cred)
	if err != nil {
		dto.FromError(ctx, err)

		return
	}

	h.setRefreshCookies(ctx, token.RefreshToken, token.RefreshTokenExpiresIn)
	res := dto.LoginResponse{
		AccessToken: &token.AccessToken,
		ExpiresIn:   &token.AccessTokenExpiresIn,
	}
	dto.Ok(ctx, res)
}

func (h *Auth) RefreshToken(ctx *gin.Context) {
	refreshToken := h.getRefreshTokenFromRequest(ctx)

	token, err := h.svc.RefreshToken(ctx, refreshToken)
	if err != nil {
		dto.FromError(ctx, err)
		h.clearRefreshCookies(ctx)

		return
	}

	res := dto.RefreshTokenResponse{
		AccessToken: &token.AccessToken,
		ExpiresIn:   &token.AccessTokenExpiresIn,
	}
	h.setRefreshCookies(ctx, token.RefreshToken, token.RefreshTokenExpiresIn)
	dto.Ok(ctx, res)
}

func (h *Auth) Logout(ctx *gin.Context) {
	refreshToken := h.getRefreshTokenFromRequest(ctx)

	h.clearRefreshCookies(ctx)
	err := h.svc.Logout(ctx, refreshToken)
	if err != nil {
		dto.FromError(ctx, err)

		return
	}

	dto.NoContent(ctx)
}

func (h *Auth) setRefreshCookies(ctx *gin.Context, refreshToken string, expiresIn int64) {
	path := "/api/v1/auth"
	maxAge := int(expiresIn)
	httpOnly := true

	ctx.SetSameSite(http.SameSiteLaxMode)
	ctx.SetCookie("refresh_token", refreshToken, maxAge, path, h.cookie.Domain, h.cookie.Secure, httpOnly)
}

func (h *Auth) clearRefreshCookies(ctx *gin.Context) {
	path := "/api/v1/auth"
	httpOnly := true

	ctx.SetSameSite(http.SameSiteLaxMode)
	ctx.SetCookie("refresh_token", "", -1, path, h.cookie.Domain, h.cookie.Secure, httpOnly)
}

func (h *Auth) getRefreshTokenFromRequest(ctx *gin.Context) string {
	refreshToken, err := ctx.Cookie("refresh_token")
	log.Print(err)

	return refreshToken
}
