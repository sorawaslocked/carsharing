package handler

import (
	"context"
	"log/slog"

	"google.golang.org/protobuf/types/known/emptypb"

	pkglog "carsharing/shared/pkg/log"
	"carsharing/shared/pkg/utils"
	"carsharing/user-service/internal/adapter/grpc/dto"
	baseuserpb "github.com/sorawaslocked/car-rental-protos/gen/base/user"
	usersvc "github.com/sorawaslocked/car-rental-protos/gen/service/user"
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

func (h *UserHandler) logger(ctx context.Context, method string) *slog.Logger {
	return pkglog.WithMetadata(pkglog.WithMethod(h.log, method), utils.MetadataFromCtx(ctx))
}

func (h *UserHandler) CreateUser(ctx context.Context, req *usersvc.CreateUserRequest) (*usersvc.CreateUserResponse, error) {
	logger := h.logger(ctx, "CreateUser")

	data, err := dto.FromCreateUserRequest(req)
	if err != nil {
		return nil, dto.ToStatusError(err)
	}

	id, err := h.userService.Create(ctx, data)
	if err != nil {
		logger.Error("creating user", pkglog.Err(err))
		return nil, dto.ToStatusError(err)
	}

	return &usersvc.CreateUserResponse{Id: &id}, nil
}

func (h *UserHandler) GetUser(ctx context.Context, req *usersvc.GetUserRequest) (*usersvc.GetUserResponse, error) {
	logger := h.logger(ctx, "GetUser")

	user, err := h.userService.Get(ctx, req.GetId())
	if err != nil {
		logger.Error("getting user", pkglog.Err(err))
		return nil, dto.ToStatusError(err)
	}

	return &usersvc.GetUserResponse{User: dto.UserToProto(user)}, nil
}

func (h *UserHandler) ListUsers(ctx context.Context, req *usersvc.ListUsersRequest) (*usersvc.ListUsersResponse, error) {
	logger := h.logger(ctx, "ListUsers")

	filter := dto.FromListUsersRequest(req)

	users, err := h.userService.List(ctx, filter)
	if err != nil {
		logger.Error("listing users", pkglog.Err(err))
		return nil, dto.ToStatusError(err)
	}

	protoUsers := make([]*baseuserpb.User, len(users))
	for i, u := range users {
		protoUsers[i] = dto.UserToProto(u)
	}

	return &usersvc.ListUsersResponse{Users: protoUsers}, nil
}

func (h *UserHandler) UpdateUser(ctx context.Context, req *usersvc.UpdateUserRequest) (*emptypb.Empty, error) {
	logger := h.logger(ctx, "UpdateUser")

	data, err := dto.FromUpdateUserRequest(req)
	if err != nil {
		return nil, dto.ToStatusError(err)
	}

	if err := h.userService.Update(ctx, req.GetId(), data); err != nil {
		logger.Error("updating user", pkglog.Err(err))
		return nil, dto.ToStatusError(err)
	}

	return &emptypb.Empty{}, nil
}

func (h *UserHandler) DeleteUser(ctx context.Context, req *usersvc.DeleteUserRequest) (*emptypb.Empty, error) {
	logger := h.logger(ctx, "DeleteUser")

	if err := h.userService.Delete(ctx, req.GetId()); err != nil {
		logger.Error("deleting user", pkglog.Err(err))
		return nil, dto.ToStatusError(err)
	}

	return &emptypb.Empty{}, nil
}

func (h *UserHandler) Register(ctx context.Context, req *usersvc.RegisterRequest) (*usersvc.RegisterResponse, error) {
	logger := h.logger(ctx, "Register")

	data, err := dto.FromRegisterRequest(req)
	if err != nil {
		return nil, dto.ToStatusError(err)
	}

	id, err := h.userService.Register(ctx, data)
	if err != nil {
		logger.Error("registering user", pkglog.Err(err))
		return nil, dto.ToStatusError(err)
	}

	return &usersvc.RegisterResponse{Id: id}, nil
}

func (h *UserHandler) SignIn(ctx context.Context, req *usersvc.SignInRequest) (*usersvc.SignInResponse, error) {
	logger := h.logger(ctx, "SignIn")

	creds := dto.FromSignInRequest(req)

	id, err := h.userService.SignIn(ctx, creds)
	if err != nil {
		logger.Error("signing in", pkglog.Err(err))
		return nil, dto.ToStatusError(err)
	}

	return &usersvc.SignInResponse{Id: id}, nil
}

func (h *UserHandler) SendActivationCode(ctx context.Context, _ *emptypb.Empty) (*emptypb.Empty, error) {
	logger := h.logger(ctx, "SendActivationCode")

	if err := h.userService.SendActivationCode(ctx); err != nil {
		logger.Error("sending activation code", pkglog.Err(err))
		return nil, dto.ToStatusError(err)
	}

	return &emptypb.Empty{}, nil
}

func (h *UserHandler) CheckActivationCode(ctx context.Context, req *usersvc.CheckActivationCodeRequest) (*emptypb.Empty, error) {
	logger := h.logger(ctx, "CheckActivationCode")

	if err := h.userService.CheckActivationCode(ctx, req.GetCode()); err != nil {
		logger.Error("checking activation code", pkglog.Err(err))
		return nil, dto.ToStatusError(err)
	}

	return &emptypb.Empty{}, nil
}

func (h *UserHandler) GetProfileImageUploadData(ctx context.Context, _ *emptypb.Empty) (*usersvc.GetProfileImageUploadDataResponse, error) {
	logger := h.logger(ctx, "GetProfileImageUploadData")

	data, err := h.userService.GetUserProfileImageUploadData(ctx)
	if err != nil {
		logger.Error("getting profile image upload data", pkglog.Err(err))
		return nil, dto.ToStatusError(err)
	}

	return &usersvc.GetProfileImageUploadDataResponse{
		UploadData: dto.ImageUploadDataToProto(data),
	}, nil
}

func (h *UserHandler) CreateDocument(ctx context.Context, req *usersvc.CreateDocumentRequest) (*usersvc.CreateDocumentResponse, error) {
	logger := h.logger(ctx, "CreateDocument")

	objectKey, imageType, err := dto.FromCreateDocumentRequest(req)
	if err != nil {
		return nil, dto.ToStatusError(err)
	}

	id, err := h.userService.CreateDocument(ctx, objectKey, imageType)
	if err != nil {
		logger.Error("creating document", pkglog.Err(err))
		return nil, dto.ToStatusError(err)
	}

	return &usersvc.CreateDocumentResponse{Id: id}, nil
}

func (h *UserHandler) GetUploadDocumentData(ctx context.Context, req *usersvc.GetUploadDocumentDataRequest) (*usersvc.GetUploadDocumentDataResponse, error) {
	logger := h.logger(ctx, "GetUploadDocumentData")

	imageType, err := dto.FromGetUploadDocumentDataRequest(req)
	if err != nil {
		return nil, dto.ToStatusError(err)
	}

	data, err := h.userService.GetDocumentImageUploadData(ctx, imageType.String())
	if err != nil {
		logger.Error("getting document upload data", pkglog.Err(err))
		return nil, dto.ToStatusError(err)
	}

	return &usersvc.GetUploadDocumentDataResponse{
		UploadData: dto.ImageUploadDataToProto(data),
	}, nil
}

func (h *UserHandler) GetProcessedDocumentsForUser(ctx context.Context, req *usersvc.GetProcessedDocumentsForUserRequest) (*usersvc.GetProcessedDocumentsForUserResponse, error) {
	logger := h.logger(ctx, "GetProcessedDocumentsForUser")

	docs, err := h.userService.GetProcessedDocumentsForUser(ctx, req.GetUserId())
	if err != nil {
		logger.Error("getting processed documents for user", pkglog.Err(err))
		return nil, dto.ToStatusError(err)
	}

	protoDocs := make([]*baseuserpb.Document, len(docs))
	for i, d := range docs {
		protoDocs[i] = dto.DocumentToProto(d)
	}

	return &usersvc.GetProcessedDocumentsForUserResponse{Documents: protoDocs}, nil
}

func (h *UserHandler) CheckDocument(ctx context.Context, req *usersvc.CheckDocumentRequest) (*emptypb.Empty, error) {
	logger := h.logger(ctx, "CheckDocument")

	docID, status, docError, err := dto.FromCheckDocumentRequest(req)
	if err != nil {
		return nil, dto.ToStatusError(err)
	}

	if err := h.userService.CheckDocument(ctx, docID, status, docError); err != nil {
		logger.Error("checking document", pkglog.Err(err))
		return nil, dto.ToStatusError(err)
	}

	return &emptypb.Empty{}, nil
}
