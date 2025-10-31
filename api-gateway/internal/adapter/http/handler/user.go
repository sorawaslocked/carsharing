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

func (h *User) GetAll(ctx *gin.Context) {
	users, err := h.svc.Find(ctx, model.UserFilter{})
	if err != nil {
		dto.FromError(ctx, err)

		return
	}

	dto.Ok(ctx, gin.H{"users": users})
}

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
