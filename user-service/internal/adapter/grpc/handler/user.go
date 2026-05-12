package handler

import (
	"context"
	"log/slog"

	"google.golang.org/protobuf/types/known/emptypb"

	baseuserpb "github.com/sorawaslocked/car-rental-protos/gen/base/user"
	usersvc "github.com/sorawaslocked/car-rental-protos/gen/service/user"
	"github.com/sorawaslocked/car-rental-user-service/internal/adapter/grpc/dto"
	"github.com/sorawaslocked/car-rental-user-service/internal/model"
	pkglog "github.com/sorawaslocked/car-rental-user-service/internal/pkg/log"
	"github.com/sorawaslocked/car-rental-user-service/internal/pkg/utils"
)

type UserHandler struct {
	log         *slog.Logger
	userService UserService
	usersvc.UnimplementedUserServiceServer
}

func NewUserHandler(log *slog.Logger, userService UserService) *UserHandler {
	return &UserHandler{
		log:         pkglog.WithComponent(log, "grpc.UserHandler"),
		userService: userService,
	}
}

func (h *UserHandler) CreateUser(ctx context.Context, req *usersvc.CreateUserRequest) (*usersvc.CreateUserResponse, error) {
	logger := pkglog.WithMethod(h.log, "CreateUser")
	logger = pkglog.WithMetadata(logger, utils.MetadataFromCtx(ctx))

	data, err := dto.FromCreateUserRequest(req)
	if err != nil {
		return nil, dto.ToStatusError(err)
	}

	id, err := h.userService.Create(ctx, data)
	if err != nil {
		return nil, dto.ToStatusError(err)
	}

	_ = logger
	return &usersvc.CreateUserResponse{Id: &id}, nil
}

func (h *UserHandler) GetUser(ctx context.Context, req *usersvc.GetUserRequest) (*usersvc.GetUserResponse, error) {
	user, err := h.userService.Get(ctx, req.GetId())
	if err != nil {
		return nil, dto.ToStatusError(err)
	}

	return &usersvc.GetUserResponse{User: dto.UserToProto(user)}, nil
}

func (h *UserHandler) ListUsers(ctx context.Context, req *usersvc.ListUsersRequest) (*usersvc.ListUsersResponse, error) {
	filter := dto.FromListUsersRequest(req)

	users, err := h.userService.List(ctx, filter)
	if err != nil {
		return nil, dto.ToStatusError(err)
	}

	protoUsers := make([]*baseuserpb.User, len(users))
	for i, u := range users {
		protoUsers[i] = dto.UserToProto(u)
	}

	return &usersvc.ListUsersResponse{Users: protoUsers}, nil
}

func (h *UserHandler) UpdateUser(ctx context.Context, req *usersvc.UpdateUserRequest) (*emptypb.Empty, error) {
	data, err := dto.FromUpdateUserRequest(req)
	if err != nil {
		return nil, dto.ToStatusError(err)
	}

	if err := h.userService.Update(ctx, req.GetId(), data); err != nil {
		return nil, dto.ToStatusError(err)
	}

	return &emptypb.Empty{}, nil
}

func (h *UserHandler) DeleteUser(ctx context.Context, req *usersvc.DeleteUserRequest) (*emptypb.Empty, error) {
	if err := h.userService.Delete(ctx, req.GetId()); err != nil {
		return nil, dto.ToStatusError(err)
	}

	return &emptypb.Empty{}, nil
}

func (h *UserHandler) Register(ctx context.Context, req *usersvc.RegisterRequest) (*usersvc.RegisterResponse, error) {
	data, err := dto.FromRegisterRequest(req)
	if err != nil {
		return nil, dto.ToStatusError(err)
	}

	id, err := h.userService.Register(ctx, data)
	if err != nil {
		return nil, dto.ToStatusError(err)
	}

	return &usersvc.RegisterResponse{Id: id}, nil
}

func (h *UserHandler) SignIn(ctx context.Context, req *usersvc.SignInRequest) (*usersvc.SignInResponse, error) {
	creds := dto.FromSignInRequest(req)

	id, err := h.userService.SignIn(ctx, creds)
	if err != nil {
		return nil, dto.ToStatusError(err)
	}

	return &usersvc.SignInResponse{Id: id}, nil
}

func (h *UserHandler) SendActivationCode(ctx context.Context, _ *emptypb.Empty) (*emptypb.Empty, error) {
	md := utils.MetadataFromCtx(ctx)
	if md.UserID == nil {
		return nil, dto.ToStatusError(model.ErrMissingMetadata)
	}

	if err := h.userService.SendActivationCode(ctx, *md.UserID); err != nil {
		return nil, dto.ToStatusError(err)
	}

	return &emptypb.Empty{}, nil
}

func (h *UserHandler) CheckActivationCode(ctx context.Context, req *usersvc.CheckActivationCodeRequest) (*emptypb.Empty, error) {
	md := utils.MetadataFromCtx(ctx)
	if md.UserID == nil {
		return nil, dto.ToStatusError(model.ErrMissingMetadata)
	}

	if err := h.userService.CheckActivationCode(ctx, *md.UserID, req.GetCode()); err != nil {
		return nil, dto.ToStatusError(err)
	}

	return &emptypb.Empty{}, nil
}
