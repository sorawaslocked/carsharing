package handler

import (
	"log/slog"
	"net/http"

	"carsharing/api-gateway/internal/adapter/http/dto"
	"carsharing/api-gateway/internal/config"
	"carsharing/api-gateway/internal/model"
	pkglog "carsharing/shared/pkg/log"
	"carsharing/shared/pkg/utils"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	svc    UserService
	cookie config.Cookie
	log    *slog.Logger
}

func NewUser(svc UserService, cookie config.Cookie, log *slog.Logger) *UserHandler {
	return &UserHandler{
		svc:    svc,
		cookie: cookie,
		log:    pkglog.WithComponent(log, "http.UserHandler"),
	}
}

// Create (UserHandler) godoc
// @Summary      Create user (admin)
// @Description  Admin endpoint to create a user with explicit roles and activation state.
// @Tags         users
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body  body      dto.UserCreateRequest  true  "Create payload"
// @Success      201   {object}  dto.IDResponse
// @Failure      400   {object}  dto.ErrorResponse
// @Failure      401   {object}  dto.ErrorResponse
// @Failure      409   {object}  dto.ErrorResponse
// @Failure      500   {object}  dto.ErrorResponse
// @Router       /users [post]
func (h *UserHandler) Create(c *gin.Context) {
	log := pkglog.WithMetadata(pkglog.WithMethod(h.log, "Create"), utils.MetadataFromCtx(c))

	data, err := dto.FromCreateUserRequest(c)
	if err != nil {
		dto.MalformedJson(c)

		return
	}

	id, err := h.svc.Create(c, data)
	if err != nil {
		log.Warn("creating user", pkglog.Err(err))

		dto.FromError(c, err)

		return
	}

	dto.Created(c, gin.H{"id": id})
}

// Get (UserHandler) godoc
// @Summary      Get user by ID
// @Description  Returns a single user by their ID.
// @Tags         users
// @Produce      json
// @Security     BearerAuth
// @Param        id    path      string  true  "User ID"
// @Success      200   {object}  dto.UserResponse
// @Failure      400   {object}  dto.ErrorResponse
// @Failure      401   {object}  dto.ErrorResponse
// @Failure      404   {object}  dto.ErrorResponse
// @Failure      500   {object}  dto.ErrorResponse
// @Router       /users/{id} [get]
func (h *UserHandler) Get(c *gin.Context) {
	log := pkglog.WithMetadata(pkglog.WithMethod(h.log, "Get"), utils.MetadataFromCtx(c))

	id, err := dto.IDParam(c)
	if err != nil {
		dto.FromError(c, err)

		return
	}

	user, err := h.svc.Get(c, id)
	if err != nil {
		log.Warn("getting user", pkglog.Err(err))

		dto.FromError(c, err)

		return
	}

	dto.Ok(c, gin.H{"user": dto.ToUserResponse(user)})
}

// List (UserHandler) godoc
// @Summary      Get users with filter
// @Description  Returns a list of users, optionally filtered by query params.
// @Tags         users
// @Produce      json
// @Security     BearerAuth
// @Param        email               query     string   false  "Filter by email"
// @Param        phoneNumber         query     string   false  "Filter by phone number"
// @Param        firstName           query     string   false  "Filter by first name"
// @Param        lastName            query     string   false  "Filter by last name"
// @Param        isDocumentVerified  query     boolean  false  "Filter by document verification status"
// @Param        isEmailVerified     query     boolean  false  "Filter by email verification status"
// @Param        isSuspended         query     boolean  false  "Filter by suspension status"
// @Param        limit               query     integer  false  "Pagination limit"
// @Param        offset              query     integer  false  "Pagination offset"
// @Success      200   {object}  dto.UsersResponse
// @Failure      400   {object}  dto.ErrorResponse
// @Failure      401   {object}  dto.ErrorResponse
// @Failure      500   {object}  dto.ErrorResponse
// @Router       /users [get]
func (h *UserHandler) List(c *gin.Context) {
	log := pkglog.WithMetadata(pkglog.WithMethod(h.log, "List"), utils.MetadataFromCtx(c))

	filter, err := dto.UserFilterFromCtx(c)
	if err != nil {
		dto.FromError(c, err)

		return
	}

	users, err := h.svc.List(c, filter)
	if err != nil {
		log.Warn("listing users", pkglog.Err(err))

		dto.FromError(c, err)

		return
	}

	userResponse := make([]dto.User, len(users))
	for i, user := range users {
		userResponse[i] = dto.ToUserResponse(user)
	}

	dto.Ok(c, gin.H{"users": userResponse})
}

// Update (UserHandler) godoc
// @Summary      Update user
// @Description  Partially updates a user by ID.
// @Tags         users
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id    path      string                 true  "User ID"
// @Param        body  body      dto.UserUpdateRequest  true  "Fields to update"
// @Success      204
// @Failure      400   {object}  dto.ErrorResponse
// @Failure      401   {object}  dto.ErrorResponse
// @Failure      404   {object}  dto.ErrorResponse
// @Failure      500   {object}  dto.ErrorResponse
// @Router       /users/{id} [patch]
func (h *UserHandler) Update(c *gin.Context) {
	log := pkglog.WithMetadata(pkglog.WithMethod(h.log, "Update"), utils.MetadataFromCtx(c))

	id, err := dto.IDParam(c)
	if err != nil {
		dto.FromError(c, err)

		return
	}

	data, err := dto.FromUpdateRequest(c)
	if err != nil {
		dto.FromError(c, err)

		return
	}

	if err = h.svc.Update(c, id, data); err != nil {
		log.Warn("updating user", pkglog.Err(err))

		dto.FromError(c, err)

		return
	}

	dto.NoContent(c)
}

// Delete (UserHandler) godoc
// @Summary      Delete user
// @Description  Deletes a user by ID.
// @Tags         users
// @Produce      json
// @Security     BearerAuth
// @Param        id    path      string  true  "User ID"
// @Success      204
// @Failure      400   {object}  dto.ErrorResponse
// @Failure      401   {object}  dto.ErrorResponse
// @Failure      404   {object}  dto.ErrorResponse
// @Failure      500   {object}  dto.ErrorResponse
// @Router       /users/{id} [delete]
func (h *UserHandler) Delete(c *gin.Context) {
	log := pkglog.WithMetadata(pkglog.WithMethod(h.log, "Delete"), utils.MetadataFromCtx(c))

	id, err := dto.IDParam(c)
	if err != nil {
		dto.FromError(c, err)

		return
	}

	if err = h.svc.Delete(c, id); err != nil {
		log.Warn("deleting user", pkglog.Err(err))

		dto.FromError(c, err)

		return
	}

	dto.NoContent(c)
}

// Register godoc
// @Summary      Register a new user
// @Description  Creates a new user account. Returns the new user's ID.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body  body      dto.RegisterRequest  true  "Registration payload"
// @Success      201   {object}  dto.IDResponse
// @Failure      400   {object}  dto.ErrorResponse
// @Failure      409   {object}  dto.ErrorResponse
// @Failure      500   {object}  dto.ErrorResponse
// @Router       /auth/register [post]
func (h *UserHandler) Register(c *gin.Context) {
	log := pkglog.WithMetadata(pkglog.WithMethod(h.log, "Register"), utils.MetadataFromCtx(c))

	data, err := dto.FromRegisterRequest(c)
	if err != nil {
		dto.MalformedJson(c)

		return
	}

	id, err := h.svc.Register(c, data)
	if err != nil {
		log.Warn("registering user", pkglog.Err(err))

		dto.FromError(c, err)

		return
	}

	dto.Created(c, gin.H{"id": id})
}

// SignIn godoc
// @Summary      Sign in
// @Description  Authenticates a user by email or phone + password. Returns an access token and sets an HttpOnly refresh token cookie.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body  body      dto.LoginRequest  true  "Sign-in credentials"
// @Success      200   {object}  dto.AccessTokenResponse
// @Failure      400   {object}  dto.ErrorResponse
// @Failure      401   {object}  dto.ErrorResponse
// @Failure      500   {object}  dto.ErrorResponse
// @Router       /auth/sign-in [post]
func (h *UserHandler) SignIn(c *gin.Context) {
	log := pkglog.WithMetadata(pkglog.WithMethod(h.log, "SignIn"), utils.MetadataFromCtx(c))

	creds, err := dto.FromLoginRequest(c)
	if err != nil {
		dto.MalformedJson(c)

		return
	}

	accessToken, refreshToken, err := h.svc.SignIn(c, creds)
	if err != nil {
		log.Warn("signing in", pkglog.Err(err))

		dto.FromError(c, err)

		return
	}

	h.setRefreshCookies(c, refreshToken.Token, refreshToken.ExpiresIn)

	dto.Ok(c, gin.H{
		"accessToken": gin.H{
			"token":     accessToken.Token,
			"expiresIn": accessToken.ExpiresIn,
		},
	})
}

// RefreshToken godoc
// @Summary      Refresh access token
// @Description  Issues a new access token using the refresh token stored in the HttpOnly cookie.
// @Tags         auth
// @Produce      json
// @Success      200  {object}  dto.AccessTokenResponse
// @Failure      401  {object}  dto.ErrorResponse
// @Failure      500  {object}  dto.ErrorResponse
// @Router       /auth/refresh-token [post]
func (h *UserHandler) RefreshToken(c *gin.Context) {
	log := pkglog.WithMetadata(pkglog.WithMethod(h.log, "RefreshToken"), utils.MetadataFromCtx(c))

	refreshToken, err := h.getRefreshTokenFromRequest(c)
	if err != nil {
		dto.FromError(c, err)

		return
	}

	newAccessToken, newRefreshToken, err := h.svc.RefreshToken(c, refreshToken)
	if err != nil {
		log.Warn("refreshing token", pkglog.Err(err))

		dto.FromError(c, err)

		h.clearRefreshCookies(c)

		return
	}

	h.setRefreshCookies(c, newRefreshToken.Token, newRefreshToken.ExpiresIn)

	dto.Ok(c, gin.H{
		"accessToken": gin.H{
			"token":     newAccessToken.Token,
			"expiresIn": newAccessToken.ExpiresIn,
		},
	})
}

// SignOut godoc
// @Summary      Sign out
// @Description  Invalidates the current session and clears the refresh token cookie.
// @Tags         auth
// @Produce      json
// @Security     BearerAuth
// @Success      204
// @Failure      401  {object}  dto.ErrorResponse
// @Failure      500  {object}  dto.ErrorResponse
// @Router       /auth/sign-out [post]
func (h *UserHandler) SignOut(c *gin.Context) {
	log := pkglog.WithMetadata(pkglog.WithMethod(h.log, "SignOut"), utils.MetadataFromCtx(c))

	h.clearRefreshCookies(c)

	if err := h.svc.SignOut(c); err != nil {
		log.Warn("signing out", pkglog.Err(err))

		dto.FromError(c, err)

		return
	}

	dto.NoContent(c)
}

// GetProfile godoc
// @Summary      Get current user profile
// @Description  Returns the profile of the authenticated user.
// @Tags         users
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  dto.UserResponse
// @Failure      401  {object}  dto.ErrorResponse
// @Failure      500  {object}  dto.ErrorResponse
// @Router       /users/profile [get]
func (h *UserHandler) GetProfile(c *gin.Context) {
	log := pkglog.WithMetadata(pkglog.WithMethod(h.log, "GetProfile"), utils.MetadataFromCtx(c))

	user, err := h.svc.GetProfile(c)
	if err != nil {
		log.Warn("getting profile", pkglog.Err(err))

		dto.FromError(c, err)

		return
	}

	dto.Ok(c, gin.H{"user": dto.ToUserResponse(user)})
}

// UpdateProfile godoc
// @Summary      Update current user profile
// @Description  Partially updates the authenticated user's own profile. Critical fields (roles, verification flags) are not accepted.
// @Tags         users
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body  body      dto.UserProfileUpdateRequest  true  "Profile fields to update"
// @Success      204
// @Failure      400  {object}  dto.ErrorResponse
// @Failure      401  {object}  dto.ErrorResponse
// @Failure      404  {object}  dto.ErrorResponse
// @Failure      500  {object}  dto.ErrorResponse
// @Router       /users/profile [patch]
func (h *UserHandler) UpdateProfile(c *gin.Context) {
	log := pkglog.WithMetadata(pkglog.WithMethod(h.log, "UpdateProfile"), utils.MetadataFromCtx(c))

	data, err := dto.FromProfileUpdateRequest(c)
	if err != nil {
		dto.MalformedJson(c)

		return
	}

	if err = h.svc.UpdateProfile(c, data); err != nil {
		log.Warn("updating profile", pkglog.Err(err))

		dto.FromError(c, err)

		return
	}

	dto.NoContent(c)
}

// SendActivationCode godoc
// @Summary      Send activation code
// @Description  Sends an activation code to the current user's email or phone.
// @Tags         users
// @Produce      json
// @Security     BearerAuth
// @Success      204
// @Failure      401  {object}  dto.ErrorResponse
// @Failure      500  {object}  dto.ErrorResponse
// @Router       /users/activation-code/send [post]
func (h *UserHandler) SendActivationCode(c *gin.Context) {
	log := pkglog.WithMetadata(pkglog.WithMethod(h.log, "SendActivationCode"), utils.MetadataFromCtx(c))

	if err := h.svc.SendActivationCode(c); err != nil {
		log.Warn("sending activation code", pkglog.Err(err))

		dto.FromError(c, err)

		return
	}

	dto.NoContent(c)
}

// CheckActivationCode godoc
// @Summary      Check activation code
// @Description  Verifies the activation code submitted by the user.
// @Tags         users
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body  body      dto.CheckActivationCodeRequest  true  "Activation code"
// @Success      204
// @Failure      400   {object}  dto.ErrorResponse
// @Failure      401   {object}  dto.ErrorResponse
// @Failure      500   {object}  dto.ErrorResponse
// @Router       /users/activation-code/check [post]
func (h *UserHandler) CheckActivationCode(c *gin.Context) {
	log := pkglog.WithMetadata(pkglog.WithMethod(h.log, "CheckActivationCode"), utils.MetadataFromCtx(c))

	code, err := dto.FromCheckActivationCodeRequest(c)
	if err != nil {
		dto.FromError(c, err)

		return
	}

	if err = h.svc.CheckActivationCode(c, code); err != nil {
		log.Warn("checking activation code", pkglog.Err(err))

		dto.FromError(c, err)

		return
	}

	dto.NoContent(c)
}

func (h *UserHandler) setRefreshCookies(c *gin.Context, refreshToken string, expiresIn int64) {
	const path = "/api/v1/auth"
	maxAge := int(expiresIn)
	const httpOnly = true

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("refresh_token", refreshToken, maxAge, path, h.cookie.Domain, h.cookie.Secure, httpOnly)
}

func (h *UserHandler) clearRefreshCookies(c *gin.Context) {
	const path = "/api/v1/auth"
	const httpOnly = true

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("refresh_token", "", -1, path, h.cookie.Domain, h.cookie.Secure, httpOnly)
}

func (h *UserHandler) getRefreshTokenFromRequest(c *gin.Context) (string, error) {
	refresh, err := c.Cookie("refresh_token")
	if err != nil {
		return "", model.ErrUnauthorized
	}

	return refresh, nil
}
