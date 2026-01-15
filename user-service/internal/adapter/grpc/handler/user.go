package handler

import (
	"context"
	"github.com/sorawaslocked/car-rental-protos/gen/base"
	usersvc "github.com/sorawaslocked/car-rental-protos/gen/service/user"
	"github.com/sorawaslocked/car-rental-user-service/internal/adapter/grpc/dto"
	"github.com/sorawaslocked/car-rental-user-service/internal/model"
	"log/slog"
)

type UserHandler struct {
	log         *slog.Logger
	userService UserService
	usersvc.UnimplementedUserServiceServer
}

func NewUserHandler(log *slog.Logger, userService UserService) *UserHandler {
	return &UserHandler{
		log:         log,
		userService: userService,
	}
}

func (h *UserHandler) Create(ctx context.Context, req *usersvc.CreateRequest) (*usersvc.CreateResponse, error) {
	data, validationErrs := dto.FromCreateUserRequest(req)
	if validationErrs != nil {
		return nil, dto.ToStatusCodeError(validationErrs)
	}

	id, err := h.userService.Insert(ctx, data)
	if err != nil {
		return nil, dto.ToStatusCodeError(err)
	}

	return &usersvc.CreateResponse{
		ID: &id,
	}, nil
}

func (h *UserHandler) Get(ctx context.Context, req *usersvc.GetRequest) (*usersvc.GetResponse, error) {
	filter := model.UserFilter{
		ID:    req.ID,
		Email: req.Email,
	}

	user, err := h.userService.FindOne(ctx, filter)
	if err != nil {
		return nil, dto.ToStatusCodeError(err)
	}

	return &usersvc.GetResponse{
		User: dto.ToUserProto(user),
	}, nil
}

func (h *UserHandler) GetAll(ctx context.Context, _ *usersvc.GetAllRequest) (*usersvc.GetAllResponse, error) {
	users, err := h.userService.Find(ctx, model.UserFilter{})
	if err != nil {
		return nil, dto.ToStatusCodeError(err)
	}

	usersProto := make([]*base.User, len(users))
	for i, user := range users {
		usersProto[i] = dto.ToUserProto(user)
	}

	return &usersvc.GetAllResponse{
		Users: usersProto,
	}, nil
}

func (h *UserHandler) Update(ctx context.Context, req *usersvc.UpdateRequest) (*usersvc.UpdateResponse, error) {
	filter := model.UserFilter{
		ID:    req.ID,
		Email: req.Email,
	}

	data, err := dto.FromUpdateUserRequest(req)
	if err != nil {
		return nil, dto.ToStatusCodeError(err)
	}

	err = h.userService.Update(ctx, filter, data)
	if err != nil {
		return nil, dto.ToStatusCodeError(err)
	}

	return &usersvc.UpdateResponse{}, nil
}

func (h *UserHandler) Delete(ctx context.Context, req *usersvc.DeleteRequest) (*usersvc.DeleteResponse, error) {
	filter := model.UserFilter{
		ID:    req.ID,
		Email: req.Email,
	}

	err := h.userService.Delete(ctx, filter)
	if err != nil {
		return nil, dto.ToStatusCodeError(err)
	}

	return &usersvc.DeleteResponse{}, nil
}

func (h *UserHandler) Me(ctx context.Context, _ *usersvc.MeRequest) (*usersvc.MeResponse, error) {
	user, err := h.userService.Me(ctx)
	if err != nil {
		return nil, dto.ToStatusCodeError(err)
	}

	return &usersvc.MeResponse{
		User: dto.ToUserProto(user),
	}, nil
}

func (h *UserHandler) SendActivationCode(ctx context.Context, _ *usersvc.SendActivationCodeRequest) (*usersvc.SendActivationCodeResponse, error) {
	err := h.userService.SendActivationCode(ctx)
	if err != nil {
		return nil, dto.ToStatusCodeError(err)
	}

	return &usersvc.SendActivationCodeResponse{}, nil
}

func (h *UserHandler) CheckActivationCode(ctx context.Context, req *usersvc.CheckActivationCodeRequest) (*usersvc.CheckActivationCodeResponse, error) {
	err := h.userService.CheckActivationCode(ctx, req.Code)
	if err != nil {
		return nil, dto.ToStatusCodeError(err)
	}

	return &usersvc.CheckActivationCodeResponse{}, nil
}
