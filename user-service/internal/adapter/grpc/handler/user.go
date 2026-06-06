package handler

import (
	"context"
	"log/slog"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"

	sharedmodel "carsharing/shared/model"
	pkglog "carsharing/shared/pkg/log"
	"carsharing/shared/pkg/utils"
	"carsharing/user-service/internal/adapter/grpc/dto"
	"carsharing/user-service/internal/model"

	baseuserpb "carsharing/protos/gen/base/user"
	usersvc "carsharing/protos/gen/service/user"
)

type UserHandler struct {
	log                *slog.Logger
	userService        UserService
	documentSubscriber DocumentAnalyzedSubscriber
	usersvc.UnimplementedUserServiceServer
}

func NewUserHandler(log *slog.Logger, userService UserService, documentSubscriber DocumentAnalyzedSubscriber) *UserHandler {
	return &UserHandler{
		log:                pkglog.WithComponent(log, "adapter.grpc.handler.UserHandler"),
		userService:        userService,
		documentSubscriber: documentSubscriber,
	}
}

func (h *UserHandler) CreateUser(ctx context.Context, req *usersvc.CreateUserRequest) (*usersvc.CreateUserResponse, error) {
	log := pkglog.WithMetadata(pkglog.WithMethod(h.log, "CreateUser"), utils.MetadataFromCtx(ctx))

	data, err := dto.FromCreateUserRequest(req)
	if err != nil {
		return nil, dto.ToStatusError(err)
	}

	id, err := h.userService.Create(ctx, data)
	if err != nil {
		log.Warn("creating user", pkglog.Err(err))

		return nil, dto.ToStatusError(err)
	}

	return &usersvc.CreateUserResponse{Id: &id}, nil
}

func (h *UserHandler) GetUser(ctx context.Context, req *usersvc.GetUserRequest) (*usersvc.GetUserResponse, error) {
	log := pkglog.WithMetadata(pkglog.WithMethod(h.log, "GetUser"), utils.MetadataFromCtx(ctx))

	user, err := h.userService.Get(ctx, req.GetId())
	if err != nil {
		log.Warn("getting user", pkglog.Err(err))

		return nil, dto.ToStatusError(err)
	}

	return &usersvc.GetUserResponse{User: dto.UserToProto(user)}, nil
}

func (h *UserHandler) ListUsers(ctx context.Context, req *usersvc.ListUsersRequest) (*usersvc.ListUsersResponse, error) {
	log := pkglog.WithMetadata(pkglog.WithMethod(h.log, "ListUsers"), utils.MetadataFromCtx(ctx))

	filter := dto.FromListUsersRequest(req)

	users, err := h.userService.List(ctx, filter)
	if err != nil {
		log.Warn("listing users", pkglog.Err(err))

		return nil, dto.ToStatusError(err)
	}

	protoUsers := make([]*baseuserpb.User, len(users))
	for i, u := range users {
		protoUsers[i] = dto.UserToProto(u)
	}

	return &usersvc.ListUsersResponse{Users: protoUsers}, nil
}

func (h *UserHandler) UpdateUser(ctx context.Context, req *usersvc.UpdateUserRequest) (*emptypb.Empty, error) {
	log := pkglog.WithMetadata(pkglog.WithMethod(h.log, "UpdateUser"), utils.MetadataFromCtx(ctx))

	data, err := dto.FromUpdateUserRequest(req)
	if err != nil {
		return nil, dto.ToStatusError(err)
	}

	if err := h.userService.Update(ctx, req.GetId(), data); err != nil {
		log.Warn("updating user", pkglog.Err(err))

		return nil, dto.ToStatusError(err)
	}

	return &emptypb.Empty{}, nil
}

func (h *UserHandler) DeleteUser(ctx context.Context, req *usersvc.DeleteUserRequest) (*emptypb.Empty, error) {
	log := pkglog.WithMetadata(pkglog.WithMethod(h.log, "DeleteUser"), utils.MetadataFromCtx(ctx))

	if err := h.userService.Delete(ctx, req.GetId()); err != nil {
		log.Warn("deleting user", pkglog.Err(err))

		return nil, dto.ToStatusError(err)
	}

	return &emptypb.Empty{}, nil
}

func (h *UserHandler) Register(ctx context.Context, req *usersvc.RegisterRequest) (*usersvc.RegisterResponse, error) {
	log := pkglog.WithMetadata(pkglog.WithMethod(h.log, "Register"), utils.MetadataFromCtx(ctx))

	data, err := dto.FromRegisterRequest(req)
	if err != nil {
		return nil, dto.ToStatusError(err)
	}

	id, err := h.userService.Register(ctx, data)
	if err != nil {
		log.Warn("registering user", pkglog.Err(err))

		return nil, dto.ToStatusError(err)
	}

	return &usersvc.RegisterResponse{Id: id}, nil
}

func (h *UserHandler) SignIn(ctx context.Context, req *usersvc.SignInRequest) (*usersvc.SignInResponse, error) {
	log := pkglog.WithMetadata(pkglog.WithMethod(h.log, "SignIn"), utils.MetadataFromCtx(ctx))

	creds := dto.FromSignInRequest(req)

	id, err := h.userService.SignIn(ctx, creds)
	if err != nil {
		log.Warn("signing in", pkglog.Err(err))

		return nil, dto.ToStatusError(err)
	}

	return &usersvc.SignInResponse{Id: id}, nil
}

func (h *UserHandler) SendActivationCode(ctx context.Context, _ *emptypb.Empty) (*emptypb.Empty, error) {
	log := pkglog.WithMetadata(pkglog.WithMethod(h.log, "SendActivationCode"), utils.MetadataFromCtx(ctx))

	if err := h.userService.SendActivationCode(ctx); err != nil {
		log.Warn("sending activation code", pkglog.Err(err))

		return nil, dto.ToStatusError(err)
	}

	return &emptypb.Empty{}, nil
}

func (h *UserHandler) CheckActivationCode(ctx context.Context, req *usersvc.CheckActivationCodeRequest) (*emptypb.Empty, error) {
	log := pkglog.WithMetadata(pkglog.WithMethod(h.log, "CheckActivationCode"), utils.MetadataFromCtx(ctx))

	if err := h.userService.CheckActivationCode(ctx, req.GetCode()); err != nil {
		log.Warn("checking activation code", pkglog.Err(err))

		return nil, dto.ToStatusError(err)
	}

	return &emptypb.Empty{}, nil
}

func (h *UserHandler) GetProfileImageUploadData(ctx context.Context, _ *emptypb.Empty) (*usersvc.GetProfileImageUploadDataResponse, error) {
	log := pkglog.WithMetadata(pkglog.WithMethod(h.log, "GetProfileImageUploadData"), utils.MetadataFromCtx(ctx))

	data, err := h.userService.GetUserProfileImageUploadData(ctx)
	if err != nil {
		log.Warn("getting profile image upload data", pkglog.Err(err))

		return nil, dto.ToStatusError(err)
	}

	return &usersvc.GetProfileImageUploadDataResponse{
		UploadData: dto.ImageUploadDataToProto(data),
	}, nil
}

func (h *UserHandler) CreateDocument(ctx context.Context, req *usersvc.CreateDocumentRequest) (*usersvc.CreateDocumentResponse, error) {
	log := pkglog.WithMetadata(pkglog.WithMethod(h.log, "CreateDocument"), utils.MetadataFromCtx(ctx))

	id, err := h.userService.CreateDocument(ctx, dto.FromCreateDocumentRequest(req))
	if err != nil {
		log.Warn("creating document", pkglog.Err(err))

		return nil, dto.ToStatusError(err)
	}

	return &usersvc.CreateDocumentResponse{Id: id}, nil
}

func (h *UserHandler) GetUploadDocumentData(ctx context.Context, req *usersvc.GetUploadDocumentDataRequest) (*usersvc.GetUploadDocumentDataResponse, error) {
	log := pkglog.WithMetadata(pkglog.WithMethod(h.log, "GetUploadDocumentData"), utils.MetadataFromCtx(ctx))

	data, err := h.userService.GetDocumentImageUploadData(ctx, dto.FromGetUploadDocumentDataRequest(req))
	if err != nil {
		log.Warn("getting document upload data", pkglog.Err(err))

		return nil, dto.ToStatusError(err)
	}

	return &usersvc.GetUploadDocumentDataResponse{
		UploadData: dto.ImageUploadDataToProto(data),
	}, nil
}

func (h *UserHandler) ListDocuments(ctx context.Context, req *usersvc.ListDocumentsRequest) (*usersvc.ListDocumentsResponse, error) {
	log := pkglog.WithMetadata(pkglog.WithMethod(h.log, "ListDocuments"), utils.MetadataFromCtx(ctx))

	docs, err := h.userService.ListDocuments(ctx, dto.FromListDocumentsRequest(req))
	if err != nil {
		log.Warn("listing documents", pkglog.Err(err))

		return nil, dto.ToStatusError(err)
	}

	protoDocs := make([]*baseuserpb.Document, len(docs))
	for i, d := range docs {
		protoDocs[i] = dto.DocumentToProto(d)
	}

	return &usersvc.ListDocumentsResponse{Documents: protoDocs}, nil
}

func (h *UserHandler) CheckDocument(ctx context.Context, req *usersvc.CheckDocumentRequest) (*emptypb.Empty, error) {
	log := pkglog.WithMetadata(pkglog.WithMethod(h.log, "CheckDocument"), utils.MetadataFromCtx(ctx))

	docID, data := dto.FromCheckDocumentRequest(req)

	if err := h.userService.CheckDocument(ctx, docID, data); err != nil {
		log.Warn("checking document", pkglog.Err(err))

		return nil, dto.ToStatusError(err)
	}

	return &emptypb.Empty{}, nil
}

func (h *UserHandler) StreamDocumentAnalyzed(req *usersvc.StreamDocumentAnalyzedRequest, stream grpc.ServerStreamingServer[usersvc.StreamDocumentAnalyzedResponse]) error {
	ctx := stream.Context()
	md := utils.MetadataFromCtx(ctx)
	log := pkglog.WithMetadata(pkglog.WithMethod(h.log, "StreamDocumentAnalyzed"), md)

	// Non-privileged callers may only stream events for their own user ID.
	// If they omit userID entirely, default to their own so they don't receive
	// every user's events.
	if !hasPrivilegedRole(md.UserRoles) {
		if req.UserId == nil {
			req.UserId = md.UserID
		} else if *req.UserId != *md.UserID {
			return dto.ToStatusError(model.ErrInsufficientPermissions)
		}
	}

	ch, unsub := h.documentSubscriber.SubscribeStream(req.UserId, req.Passed)
	defer unsub()

	log.Info("document analyzed stream opened")
	defer log.Info("document analyzed stream closed")

	for {
		select {
		case event, ok := <-ch:
			if !ok {
				return nil
			}
			defects := make([]*baseuserpb.Defect, len(event.Defects))
			for i, d := range event.Defects {
				defects[i] = &baseuserpb.Defect{
					Type:        d.Type,
					Description: d.Description,
				}
			}
			if err := stream.Send(&usersvc.StreamDocumentAnalyzedResponse{
				DocumentId: event.DocumentID,
				UserId:     event.UserID,
				Passed:     event.Passed,
				Defects:    defects,
			}); err != nil {
				return err
			}
		case <-ctx.Done():
			return nil
		}
	}
}

func hasPrivilegedRole(roles []sharedmodel.Role) bool {
	for _, r := range roles {
		if r == sharedmodel.RoleAdmin || r == sharedmodel.RoleUserManager {
			return true
		}
	}
	return false
}
