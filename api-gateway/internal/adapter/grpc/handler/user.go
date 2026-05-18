package handler

import (
	"context"
	"log/slog"

	"github.com/sorawaslocked/car-rental-api-gateway/internal/adapter/grpc/dto"
	"github.com/sorawaslocked/car-rental-api-gateway/internal/model"
	pkglog "github.com/sorawaslocked/car-rental-api-gateway/internal/pkg/log"
	"github.com/sorawaslocked/car-rental-api-gateway/internal/pkg/utils"
	basepb "github.com/sorawaslocked/car-rental-protos/gen/base"
	usersvc "github.com/sorawaslocked/car-rental-protos/gen/service/user"
	"google.golang.org/protobuf/types/known/emptypb"
)

type UserHandler struct {
	client usersvc.UserServiceClient
	log    *slog.Logger
}

func NewUserHandler(client usersvc.UserServiceClient, logger *slog.Logger) *UserHandler {
	return &UserHandler{
		client: client,
		log:    pkglog.WithComponent(logger, "grpc.UserHandler"),
	}
}

func (h *UserHandler) Create(ctx context.Context, data model.UserCreate) (string, error) {
	logger := pkglog.WithMethod(h.log, "Create")
	logger = pkglog.WithMetadata(logger, utils.MetadataFromCtx(ctx))

	req := &usersvc.CreateUserRequest{
		Email:       data.Email,
		FirstName:   data.FirstName,
		LastName:    data.LastName,
		BirthDate:   data.BirthDate,
		PhoneNumber: data.PhoneNumber,
	}
	if data.Password.Text != nil {
		req.Password = *data.Password.Text
	}
	if data.Password.TextConfirmation != nil {
		req.PasswordConfirmation = *data.Password.TextConfirmation
	}

	res, err := h.client.CreateUser(ctx, req)
	if err != nil {
		if dto.IsSystemErr(err) {
			logger.Error("grpc call failed", pkglog.Err(err))
		}

		return "", dto.FromGrpcErr(err)
	}

	return res.GetId(), nil
}

func (h *UserHandler) Get(ctx context.Context, id string) (model.User, error) {
	logger := pkglog.WithMethod(h.log, "Get")
	logger = pkglog.WithMetadata(logger, utils.MetadataFromCtx(ctx))

	res, err := h.client.GetUser(ctx, &usersvc.GetUserRequest{Id: id})
	if err != nil {
		if dto.IsSystemErr(err) {
			logger.Error("grpc call failed", pkglog.Err(err))
		}

		return model.User{}, dto.FromGrpcErr(err)
	}

	return dto.UserFromProto(res.GetUser()), nil
}

func (h *UserHandler) List(ctx context.Context, filter model.UserFilter) ([]model.User, error) {
	logger := pkglog.WithMethod(h.log, "List")
	logger = pkglog.WithMetadata(logger, utils.MetadataFromCtx(ctx))

	req := &usersvc.ListUsersRequest{
		Email:              filter.Email,
		PhoneNumber:        filter.PhoneNumber,
		FirstName:          filter.FirstName,
		LastName:           filter.LastName,
		IsDocumentVerified: filter.IsDocumentVerified,
		IsEmailVerified:    filter.IsEmailVerified,
		IsSuspended:        filter.IsSuspended,
	}
	if filter.Pagination != nil {
		req.Pagination = &basepb.Pagination{
			Limit:  filter.Pagination.Limit,
			Offset: filter.Pagination.Offset,
		}
	}

	res, err := h.client.ListUsers(ctx, req)
	if err != nil {
		if dto.IsSystemErr(err) {
			logger.Error("grpc call failed", pkglog.Err(err))
		}

		return nil, dto.FromGrpcErr(err)
	}

	users := make([]model.User, len(res.GetUsers()))
	for i, u := range res.GetUsers() {
		users[i] = dto.UserFromProto(u)
	}

	return users, nil
}

func (h *UserHandler) Update(ctx context.Context, id string, data model.UserUpdate) error {
	logger := pkglog.WithMethod(h.log, "Update")
	logger = pkglog.WithMetadata(logger, utils.MetadataFromCtx(ctx))

	req := &usersvc.UpdateUserRequest{
		Id:                 id,
		Email:              data.Email,
		PhoneNumber:        data.PhoneNumber,
		FirstName:          data.FirstName,
		LastName:           data.LastName,
		BirthDate:          data.BirthDate,
		ProfileImageKey:    data.ProfileImageKey,
		Roles:              data.Roles,
		IsDocumentVerified: data.IsDocumentVerified,
		IsEmailVerified:    data.IsEmailVerified,
		IsSuspended:        data.IsSuspended,
	}
	if data.Password.Text != nil {
		req.Password = data.Password.Text
	}
	if data.Password.TextConfirmation != nil {
		req.PasswordConfirmation = data.Password.TextConfirmation
	}

	_, err := h.client.UpdateUser(ctx, req)
	if err != nil {
		if dto.IsSystemErr(err) {
			logger.Error("grpc call failed", pkglog.Err(err))
		}

		return dto.FromGrpcErr(err)
	}

	return nil
}

func (h *UserHandler) Delete(ctx context.Context, id string) error {
	logger := pkglog.WithMethod(h.log, "Delete")
	logger = pkglog.WithMetadata(logger, utils.MetadataFromCtx(ctx))

	_, err := h.client.DeleteUser(ctx, &usersvc.DeleteUserRequest{Id: id})
	if err != nil {
		if dto.IsSystemErr(err) {
			logger.Error("grpc call failed", pkglog.Err(err))
		}

		return dto.FromGrpcErr(err)
	}

	return nil
}

func (h *UserHandler) Register(ctx context.Context, data model.UserCreate) (string, error) {
	logger := pkglog.WithMethod(h.log, "Register")
	logger = pkglog.WithMetadata(logger, utils.MetadataFromCtx(ctx))

	req := &usersvc.RegisterRequest{
		Email:       data.Email,
		FirstName:   data.FirstName,
		LastName:    data.LastName,
		BirthDate:   data.BirthDate,
		PhoneNumber: data.PhoneNumber,
	}
	if data.Password.Text != nil {
		req.Password = *data.Password.Text
	}
	if data.Password.TextConfirmation != nil {
		req.PasswordConfirmation = *data.Password.TextConfirmation
	}

	res, err := h.client.Register(ctx, req)
	if err != nil {
		if dto.IsSystemErr(err) {
			logger.Error("grpc call failed", pkglog.Err(err))
		}

		return "", dto.FromGrpcErr(err)
	}

	return res.GetId(), nil
}

func (h *UserHandler) SignIn(ctx context.Context, creds model.Credentials) (string, error) {
	logger := pkglog.WithMethod(h.log, "SignIn")
	logger = pkglog.WithMetadata(logger, utils.MetadataFromCtx(ctx))

	req := &usersvc.SignInRequest{
		Email:       creds.Email,
		PhoneNumber: creds.PhoneNumber,
	}
	if creds.Password.Text != nil {
		req.Password = *creds.Password.Text
	}

	res, err := h.client.SignIn(ctx, req)
	if err != nil {
		if dto.IsSystemErr(err) {
			logger.Error("grpc call failed", pkglog.Err(err))
		}

		return "", dto.FromGrpcErr(err)
	}

	return res.GetId(), nil
}

func (h *UserHandler) SendActivationCode(ctx context.Context) error {
	logger := pkglog.WithMethod(h.log, "SendActivationCode")
	logger = pkglog.WithMetadata(logger, utils.MetadataFromCtx(ctx))

	_, err := h.client.SendActivationCode(ctx, &emptypb.Empty{})
	if err != nil {
		if dto.IsSystemErr(err) {
			logger.Error("grpc call failed", pkglog.Err(err))
		}

		return dto.FromGrpcErr(err)
	}

	return nil
}

func (h *UserHandler) CheckActivationCode(ctx context.Context, code string) error {
	logger := pkglog.WithMethod(h.log, "CheckActivationCode")
	logger = pkglog.WithMetadata(logger, utils.MetadataFromCtx(ctx))

	_, err := h.client.CheckActivationCode(ctx, &usersvc.CheckActivationCodeRequest{Code: code})
	if err != nil {
		if dto.IsSystemErr(err) {
			logger.Error("grpc call failed", pkglog.Err(err))
		}

		return dto.FromGrpcErr(err)
	}

	return nil
}
