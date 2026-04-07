package handler

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sorawaslocked/car-rental-api-gateway/internal/adapter/http/dto"
	"github.com/sorawaslocked/car-rental-api-gateway/internal/config"
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

// Register godoc
// @Summary      Register a new user
// @Description  Creates a new user account. Returns the new user's ID.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body  body      dto.RegisterRequest   true  "Registration payload"
// @Success      201   {object}  dto.RegisterResponse
// @Failure      400   {object}  map[string]any  "Malformed JSON or validation error"
// @Failure      409   {object}  map[string]any  "User already exists"
// @Failure      500   {object}  map[string]any
// @Router       /auth/register [post]
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

// Login godoc
// @Summary      Login
// @Description  Authenticates a user by email or phone + password. Returns an access token and sets an HttpOnly refresh token cookie.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body  body      dto.LoginRequest   true  "Login credentials"
// @Success      200   {object}  dto.LoginResponse
// @Failure      400   {object}  map[string]any  "Malformed JSON"
// @Failure      401   {object}  map[string]any  "Invalid credentials"
// @Failure      500   {object}  map[string]any
// @Router       /auth/login [post]
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

// RefreshToken godoc
// @Summary      Refresh access token
// @Description  Issues a new access token using the refresh token stored in the HttpOnly cookie.
// @Tags         auth
// @Produce      json
// @Success      200  {object}  dto.RefreshTokenResponse
// @Failure      401  {object}  map[string]any  "Invalid or expired refresh token"
// @Failure      500  {object}  map[string]any
// @Router       /auth/refresh-token [post]
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

// Logout godoc
// @Summary      Logout
// @Description  Invalidates the refresh token and clears the cookie.
// @Tags         auth
// @Produce      json
// @Success      204  "No Content"
// @Failure      401  {object}  map[string]any
// @Failure      500  {object}  map[string]any
// @Router       /auth/logout [post]
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
