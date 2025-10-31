package grpc

import (
	"context"
	"github.com/sorawaslocked/car-rental-api-gateway/internal/adapter/grpc/dto"
	"github.com/sorawaslocked/car-rental-api-gateway/internal/model"
	usersvc "github.com/sorawaslocked/car-rental-protos/gen/service/user"
)

type UserHandler struct {
	client usersvc.UserServiceClient
}

func NewUserHandler(client usersvc.UserServiceClient) *UserHandler {
	return &UserHandler{client: client}
}

func (h *UserHandler) Create(ctx context.Context, data model.UserCreateData) (uint64, error) {
	req := &usersvc.CreateRequest{
		Email:                data.Email,
		Password:             data.Password,
		PasswordConfirmation: data.PasswordConfirmation,
		FirstName:            data.FirstName,
		LastName:             data.LastName,
		BirthDate:            data.BirthDate,
	}
	if data.PhoneNumber != "" {
		req.PhoneNumber = &data.PhoneNumber
	}
	if data.Roles != nil {
		req.Roles = *data.Roles
	}
	if data.IsActive != nil {
		req.IsActive = *data.IsActive
	}
	if data.IsConfirmed != nil {
		req.IsConfirmed = *data.IsConfirmed
	}

	res, err := h.client.Create(ctx, req)
	if err != nil {
		return 0, fromGrpcErr(err)
	}

	return *res.ID, nil
}

func (h *UserHandler) Get(ctx context.Context, filter model.UserFilter) (model.User, error) {
	req := &usersvc.GetRequest{
		ID:    filter.ID,
		Email: filter.Email,
	}

	res, err := h.client.Get(ctx, req)
	if err != nil {
		return model.User{}, fromGrpcErr(err)
	}

	return dto.FromProto(res.User), nil
}

func (h *UserHandler) GetAll(ctx context.Context, _ model.UserFilter) ([]model.User, error) {
	req := &usersvc.GetAllRequest{}

	res, err := h.client.GetAll(ctx, req)
	if err != nil {
		return nil, fromGrpcErr(err)
	}

	users := make([]model.User, len(res.Users))
	for i, user := range res.Users {
		users[i] = dto.FromProto(user)
	}

	return users, nil
}

func (h *UserHandler) Update(ctx context.Context, filter model.UserFilter, data model.UserUpdateData) error {
	req := &usersvc.UpdateRequest{
		ID:                   filter.ID,
		Email:                filter.Email,
		NewEmail:             data.Email,
		PhoneNumber:          data.PhoneNumber,
		FirstName:            data.FirstName,
		LastName:             data.LastName,
		BirthDate:            data.BirthDate,
		Password:             data.Password,
		PasswordConfirmation: data.PasswordConfirmation,
		IsActive:             data.IsActive,
		IsConfirmed:          data.IsConfirmed,
	}
	if data.Roles != nil {
		req.Roles = *data.Roles
	}

	_, err := h.client.Update(ctx, req)
	if err != nil {
		return fromGrpcErr(err)
	}

	return nil
}

func (h *UserHandler) Delete(ctx context.Context, filter model.UserFilter) error {
	req := &usersvc.DeleteRequest{
		ID:    filter.ID,
		Email: filter.Email,
	}

	_, err := h.client.Delete(ctx, req)
	if err != nil {
		return fromGrpcErr(err)
	}

	return nil
}
