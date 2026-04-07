package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/sorawaslocked/car-rental-api-gateway/internal/adapter/http/dto"
	"github.com/sorawaslocked/car-rental-api-gateway/internal/model"
)

type User struct {
	svc UserService
}

func NewUser(svc UserService) *User {
	return &User{
		svc: svc,
	}
}

// Create (User) godoc
// @Summary      Create user (admin)
// @Description  Admin endpoint to create a user with explicit roles and activation state.
// @Tags         users
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body  body      dto.UserCreateRequest   true  "User create payload"
// @Success      201   {object}  dto.UserCreateResponse
// @Failure      400   {object}  map[string]any
// @Failure      401   {object}  map[string]any
// @Failure      409   {object}  map[string]any
// @Failure      500   {object}  map[string]any
// @Router       /users [post]
func (h *User) Create(ctx *gin.Context) {
	data, err := dto.FromCreateUserRequest(ctx)
	if err != nil {
		dto.MalformedJson(ctx)

		return
	}

	id, err := h.svc.Insert(ctx, data)
	if err != nil {
		dto.FromError(ctx, err)

		return
	}

	res := dto.UserCreateResponse{ID: &id}
	dto.Created(ctx, res)
}

// Get (User) godoc
// @Summary      Get user(s)
// @Description  Returns a single user when id or email query param is provided, otherwise returns all users.
// @Tags         users
// @Produce      json
// @Security     BearerAuth
// @Param        id     query     integer  false  "User ID"
// @Param        email  query     string   false  "User email"
// @Success      200    {object}  map[string]any  "user or users"
// @Failure      400    {object}  map[string]any
// @Failure      401    {object}  map[string]any
// @Failure      404    {object}  map[string]any
// @Failure      500    {object}  map[string]any
// @Router       /users [get]
func (h *User) Get(ctx *gin.Context) {
	filter, err := dto.FilterFromCtx(ctx)
	if err != nil {
		dto.FromError(ctx, err)

		return
	}

	if filter.ID == nil && filter.Email == nil {
		h.GetAll(ctx)

		return
	}

	user, err := h.svc.FindOne(ctx, filter)
	if err != nil {
		dto.FromError(ctx, err)

		return
	}

	dto.Ok(ctx, gin.H{"user": user})
}

// GetAll (User) godoc
// @Summary      Get all users
// @Description  Returns all users.
// @Tags         users
// @Produce      json
// @Security     BearerAuth
// @Success      200    {object}  map[string]any  "users"
// @Failure      401    {object}  map[string]any
// @Failure      500    {object}  map[string]any
// @Router       /users [get]
func (h *User) GetAll(ctx *gin.Context) {
	users, err := h.svc.Find(ctx, model.UserFilter{})
	if err != nil {
		dto.FromError(ctx, err)

		return
	}

	dto.Ok(ctx, gin.H{"users": users})
}

// Update (User) godoc
// @Summary      Update user
// @Description  Partially updates a user matched by id or email query param.
// @Tags         users
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id     query     integer              false  "User ID"
// @Param        email  query     string               false  "User email"
// @Param        body   body      dto.UserUpdateRequest  true   "Fields to update"
// @Success      200    {object}  map[string]any
// @Failure      400    {object}  map[string]any
// @Failure      401    {object}  map[string]any
// @Failure      404    {object}  map[string]any
// @Failure      500    {object}  map[string]any
// @Router       /users [patch]
func (h *User) Update(ctx *gin.Context) {
	filter, err := dto.FilterFromCtx(ctx)
	if err != nil {
		dto.FromError(ctx, err)

		return
	}

	data, err := dto.FromUpdateRequest(ctx)
	if err != nil {
		dto.FromError(ctx, err)

		return
	}

	err = h.svc.Update(ctx, filter, data)
	if err != nil {
		dto.FromError(ctx, err)

		return
	}

	dto.Ok(ctx, nil)
}

// Delete (User) godoc
// @Summary      Delete user
// @Description  Deletes a user matched by id or email query param.
// @Tags         users
// @Produce      json
// @Security     BearerAuth
// @Param        id     query     integer  false  "User ID"
// @Param        email  query     string   false  "User email"
// @Success      200    {object}  map[string]any
// @Failure      400    {object}  map[string]any
// @Failure      401    {object}  map[string]any
// @Failure      404    {object}  map[string]any
// @Failure      500    {object}  map[string]any
// @Router       /users [delete]
func (h *User) Delete(ctx *gin.Context) {
	filter, err := dto.FilterFromCtx(ctx)
	if err != nil {
		dto.FromError(ctx, err)

		return
	}

	err = h.svc.Delete(ctx, filter)
	if err != nil {
		dto.FromError(ctx, err)

		return
	}

	dto.Ok(ctx, nil)
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
func (h *User) Me(ctx *gin.Context) {
	user, err := h.svc.Me(ctx)
	if err != nil {
		dto.FromError(ctx, err)

		return
	}

	dto.Ok(ctx, gin.H{"user": user})
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
func (h *User) SendActivationCode(ctx *gin.Context) {
	err := h.svc.SendActivationCode(ctx)
	if err != nil {
		dto.FromError(ctx, err)

		return
	}

	dto.Ok(ctx, nil)
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
func (h *User) CheckActivationCode(ctx *gin.Context) {
	code, err := dto.FromCheckActivationCodeRequest(ctx)
	if err != nil {
		dto.FromError(ctx, err)

		return
	}

	err = h.svc.CheckActivationCode(ctx, code)
	if err != nil {
		dto.FromError(ctx, err)

		return
	}

	dto.Ok(ctx, nil)
}
