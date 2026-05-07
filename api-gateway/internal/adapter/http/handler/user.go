package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sorawaslocked/car-rental-api-gateway/internal/adapter/http/dto"
	"github.com/sorawaslocked/car-rental-api-gateway/internal/config"
)

type UserHandler struct {
	svc    UserService
	cookie config.Cookie
}

func NewUser(svc UserService, cookie config.Cookie) *UserHandler {
	return &UserHandler{
		svc:    svc,
		cookie: cookie,
	}
}

// Create (UserHandler) godoc
// @Summary      Create user (admin)
// @Description  Admin endpoint to create a user with explicit roles and activation state.
// @Tags         users
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body  body      dto.UserCreateRequest   true  "UserHandler create payload"
// @Success      201   {object}  dto.UserCreateResponse
// @Failure      400   {object}  map[string]any
// @Failure      401   {object}  map[string]any
// @Failure      409   {object}  map[string]any
// @Failure      500   {object}  map[string]any
// @Router       /users [post]
func (h *UserHandler) Create(c *gin.Context) {
	data, err := dto.FromCreateUserRequest(c)
	if err != nil {
		dto.MalformedJson(c)

		return
	}

	id, err := h.svc.Create(c, data)
	if err != nil {
		dto.FromError(c, err)

		return
	}

	dto.Created(c, gin.H{"id": id})
}

// Get (UserHandler) godoc
// @Summary      Get user(s)
// @Description  Returns a single user when id or email query param is provided, otherwise returns all users.
// @Tags         users
// @Produce      json
// @Security     BearerAuth
// @Param        id     query     integer  false  "UserHandler ID"
// @Param        email  query     string   false  "UserHandler email"
// @Success      200    {object}  map[string]any  "user or users"
// @Failure      400    {object}  map[string]any
// @Failure      401    {object}  map[string]any
// @Failure      404    {object}  map[string]any
// @Failure      500    {object}  map[string]any
// @Router       /users [get]
func (h *UserHandler) Get(c *gin.Context) {
	id, err := dto.IDParam(c)

	user, err := h.svc.Get(c, id)
	if err != nil {
		dto.FromError(c, err)

		return
	}

	dto.Ok(c, gin.H{"user": dto.ToUserResponse(user)})
}

// GetAllWithFilter (UserHandler) godoc
// @Summary      Get all users
// @Description  Returns all users.
// @Tags         users
// @Produce      json
// @Security     BearerAuth
// @Success      200    {object}  map[string]any  "users"
// @Failure      401    {object}  map[string]any
// @Failure      500    {object}  map[string]any
// @Router       /users [get]
func (h *UserHandler) GetAllWithFilter(c *gin.Context) {
	filter, err := dto.UserFilterFromCtx(c)
	if err != nil {
		dto.FromError(c, err)

		return
	}

	users, err := h.svc.GetAllWithFilter(c, filter)
	if err != nil {
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
// @Description  Partially updates a user matched by id or email query param.
// @Tags         users
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id     query     integer              false  "UserHandler ID"
// @Param        email  query     string               false  "UserHandler email"
// @Param        body   body      dto.UserUpdateRequest  true   "Fields to update"
// @Success      200    {object}  map[string]any
// @Failure      400    {object}  map[string]any
// @Failure      401    {object}  map[string]any
// @Failure      404    {object}  map[string]any
// @Failure      500    {object}  map[string]any
// @Router       /users [patch]
func (h *UserHandler) Update(c *gin.Context) {
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

	err = h.svc.Update(c, id, data)
	if err != nil {
		dto.FromError(c, err)

		return
	}

	dto.NoContent(c)
}

// Delete (UserHandler) godoc
// @Summary      Delete user
// @Description  Deletes a user matched by id or email query param.
// @Tags         users
// @Produce      json
// @Security     BearerAuth
// @Param        id     query     integer  false  "UserHandler ID"
// @Param        email  query     string   false  "UserHandler email"
// @Success      200    {object}  map[string]any
// @Failure      400    {object}  map[string]any
// @Failure      401    {object}  map[string]any
// @Failure      404    {object}  map[string]any
// @Failure      500    {object}  map[string]any
// @Router       /users [delete]
func (h *UserHandler) Delete(c *gin.Context) {
	id, err := dto.IDParam(c)
	if err != nil {
		dto.FromError(c, err)

		return
	}

	err = h.svc.Delete(c, id)
	if err != nil {
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
// @Param        body  body      dto.RegisterRequest   true  "Registration payload"
// @Success      201   {object}  dto.RegisterResponse
// @Failure      400   {object}  map[string]any  "Malformed JSON or validation error"
// @Failure      409   {object}  map[string]any  "UserHandler already exists"
// @Failure      500   {object}  map[string]any
// @Router       /auth/register [post]
func (h *UserHandler) Register(c *gin.Context) {
	data, err := dto.FromRegisterRequest(c)
	if err != nil {
		dto.MalformedJson(c)

		return
	}

	id, err := h.svc.Register(c, data)
	if err != nil {
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
// @Param        body  body      dto.LoginRequest   true  "Sign-in credentials"
// @Success      200   {object}  dto.LoginResponse
// @Failure      400   {object}  map[string]any  "Malformed JSON"
// @Failure      401   {object}  map[string]any  "Invalid credentials"
// @Failure      500   {object}  map[string]any
// @Router       /auth/sign-in [post]
func (h *UserHandler) SignIn(c *gin.Context) {
	creds, err := dto.FromLoginRequest(c)
	if err != nil {
		dto.MalformedJson(c)

		return
	}

	accessToken, refreshToken, err := h.svc.SignIn(c, creds)
	if err != nil {
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
// @Success      200  {object}  dto.RefreshTokenResponse
// @Failure      401  {object}  map[string]any  "Invalid or expired refresh token"
// @Failure      500  {object}  map[string]any
// @Router       /auth/refresh-token [post]
func (h *UserHandler) RefreshToken(c *gin.Context) {
	refreshToken, err := h.getRefreshTokenFromRequest(c)
	if err != nil {
		dto.FromError(c, err)

		return
	}

	newAccessToken, newRefreshToken, err := h.svc.RefreshToken(c, refreshToken)
	if err != nil {
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
// @Description  Invalidates the refresh token and clears the cookie.
// @Tags         auth
// @Produce      json
// @Success      204  "No Content"
// @Failure      401  {object}  map[string]any
// @Failure      500  {object}  map[string]any
// @Router       /auth/sign-out [post]
func (h *UserHandler) SignOut(c *gin.Context) {
	h.clearRefreshCookies(c)

	err := h.svc.SignOut(c)
	if err != nil {
		dto.FromError(c, err)

		return
	}

	dto.NoContent(c)
}

// Me godoc
// @Summary      Get current user
// @Description  Returns the profile of the authenticated user extracted from the JWT.
// @Tags         users
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  map[string]any
// @Failure      401  {object}  map[string]any
// @Failure      500  {object}  map[string]any
// @Router       /users/me [get]
func (h *UserHandler) Me(c *gin.Context) {
	user, err := h.svc.Me(c)
	if err != nil {
		dto.FromError(c, err)

		return
	}

	dto.Ok(c, gin.H{"user": dto.ToUserResponse(user)})
}

// SendActivationCode godoc
// @Summary      Send activation code
// @Description  Sends an SMS/email activation code to the current user.
// @Tags         users
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  map[string]any
// @Failure      401  {object}  map[string]any
// @Failure      500  {object}  map[string]any
// @Router       /users/activation-code/send [post]
func (h *UserHandler) SendActivationCode(c *gin.Context) {
	err := h.svc.SendActivationCode(c)
	if err != nil {
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
// @Success      200   {object}  map[string]any
// @Failure      400   {object}  map[string]any
// @Failure      401   {object}  map[string]any
// @Failure      500   {object}  map[string]any
// @Router       /users/activation-code/check [post]
func (h *UserHandler) CheckActivationCode(c *gin.Context) {
	code, err := dto.FromCheckActivationCodeRequest(c)
	if err != nil {
		dto.FromError(c, err)

		return
	}

	err = h.svc.CheckActivationCode(c, code)
	if err != nil {
		dto.FromError(c, err)

		return
	}

	dto.NoContent(c)
}

func (h *UserHandler) CreateDocument(c *gin.Context) {
	imageType, objectKey, err := dto.FromCreateDocumentRequest(c)
	if err != nil {
		dto.FromError(c, err)

		return
	}

	id, err := h.svc.CreateDocument(c, imageType, objectKey)
	if err != nil {
		dto.FromError(c, err)

		return
	}

	dto.Created(c, gin.H{"id": id})
}

func (h *UserHandler) GetUploadDocumentData(c *gin.Context) {
	imageType, err := dto.FromGetUploadDocumentDataRequest(c)
	if err != nil {
		dto.FromError(c, err)

		return
	}

	data, err := h.svc.GetUploadDocumentData(c, imageType)
	if err != nil {
		dto.FromError(c, err)

		return
	}

	dto.Ok(c, gin.H{"uploadData": dto.ToImageUploadDataResponse(data)})
}

func (h *UserHandler) GetProcessedDocumentsForUser(c *gin.Context) {
	userID, err := dto.IDParam(c)
	if err != nil {
		dto.FromError(c, err)

		return
	}

	documents, err := h.svc.GetProcessedDocumentsForUser(c, userID)
	if err != nil {
		dto.FromError(c, err)

		return
	}

	documentResponse := make([]dto.Document, len(documents))
	for i, doc := range documents {
		documentResponse[i] = dto.ToDocumentResponse(doc)
	}

	dto.Ok(c, gin.H{"documents": documentResponse})
}

func (h *UserHandler) CheckDocument(c *gin.Context) {
	docID, err := dto.IDParam(c)
	if err != nil {
		dto.FromError(c, err)

		return
	}

	status, documentError, err := dto.FromCheckDocumentRequest(c)
	if err != nil {
		dto.FromError(c, err)

		return
	}

	err = h.svc.CheckDocument(c, docID, status, documentError)
	if err != nil {
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
	return c.Cookie("refresh_token")
}
