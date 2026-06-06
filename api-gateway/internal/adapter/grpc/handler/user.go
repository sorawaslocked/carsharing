package handler

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"time"

	"carsharing/api-gateway/internal/adapter/grpc/dto"
	"carsharing/api-gateway/internal/model"
	basepb "carsharing/protos/gen/base"
	usersvc "carsharing/protos/gen/service/user"
	pkglog "carsharing/shared/pkg/log"
	"carsharing/shared/pkg/utils"
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
	log := pkglog.WithMetadata(pkglog.WithMethod(h.log, "Create"), utils.MetadataFromCtx(ctx))
	log.Debug("calling user service")

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
		log.Warn("creating user", pkglog.Err(err))

		return "", dto.FromGrpcErr(err)
	}

	log.Debug("user created", slog.String("id", res.GetId()))

	return res.GetId(), nil
}

func (h *UserHandler) Get(ctx context.Context, id string) (model.User, error) {
	log := pkglog.WithMetadata(pkglog.WithMethod(h.log, "Get"), utils.MetadataFromCtx(ctx))
	log.Debug("calling user service")

	res, err := h.client.GetUser(ctx, &usersvc.GetUserRequest{Id: id})
	if err != nil {
		log.Warn("getting user", pkglog.Err(err))

		return model.User{}, dto.FromGrpcErr(err)
	}

	return dto.UserFromProto(res.GetUser()), nil
}

func (h *UserHandler) List(ctx context.Context, filter model.UserFilter) ([]model.User, error) {
	log := pkglog.WithMetadata(pkglog.WithMethod(h.log, "List"), utils.MetadataFromCtx(ctx))
	log.Debug("calling user service")

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
		log.Warn("listing users", pkglog.Err(err))

		return nil, dto.FromGrpcErr(err)
	}

	users := make([]model.User, len(res.GetUsers()))
	for i, u := range res.GetUsers() {
		users[i] = dto.UserFromProto(u)
	}

	return users, nil
}

func (h *UserHandler) Update(ctx context.Context, id string, data model.UserUpdate) error {
	log := pkglog.WithMetadata(pkglog.WithMethod(h.log, "Update"), utils.MetadataFromCtx(ctx))
	log.Debug("calling user service")

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
		log.Warn("updating user", pkglog.Err(err))

		return dto.FromGrpcErr(err)
	}

	return nil
}

func (h *UserHandler) Delete(ctx context.Context, id string) error {
	log := pkglog.WithMetadata(pkglog.WithMethod(h.log, "Delete"), utils.MetadataFromCtx(ctx))
	log.Debug("calling user service")

	_, err := h.client.DeleteUser(ctx, &usersvc.DeleteUserRequest{Id: id})
	if err != nil {
		log.Warn("deleting user", pkglog.Err(err))

		return dto.FromGrpcErr(err)
	}

	return nil
}

func (h *UserHandler) Register(ctx context.Context, data model.UserCreate) (string, error) {
	log := pkglog.WithMetadata(pkglog.WithMethod(h.log, "Register"), utils.MetadataFromCtx(ctx))
	log.Debug("calling user service")

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
		log.Warn("registering user", pkglog.Err(err))

		return "", dto.FromGrpcErr(err)
	}

	log.Debug("user registered", slog.String("id", res.GetId()))

	return res.GetId(), nil
}

func (h *UserHandler) SignIn(ctx context.Context, creds model.Credentials) (string, error) {
	log := pkglog.WithMetadata(pkglog.WithMethod(h.log, "SignIn"), utils.MetadataFromCtx(ctx))

	req := &usersvc.SignInRequest{
		Email:       creds.Email,
		PhoneNumber: creds.PhoneNumber,
	}
	if creds.Password.Text != nil {
		req.Password = *creds.Password.Text
	}

	res, err := h.client.SignIn(ctx, req)
	if err != nil {
		log.Warn("signing in", pkglog.Err(err))

		return "", dto.FromGrpcErr(err)
	}

	return res.GetId(), nil
}

func (h *UserHandler) SendActivationCode(ctx context.Context) error {
	log := pkglog.WithMetadata(pkglog.WithMethod(h.log, "SendActivationCode"), utils.MetadataFromCtx(ctx))

	_, err := h.client.SendActivationCode(ctx, &emptypb.Empty{})
	if err != nil {
		log.Warn("sending activation code", pkglog.Err(err))

		return dto.FromGrpcErr(err)
	}

	return nil
}

func (h *UserHandler) CheckActivationCode(ctx context.Context, code string) error {
	log := pkglog.WithMetadata(pkglog.WithMethod(h.log, "CheckActivationCode"), utils.MetadataFromCtx(ctx))

	_, err := h.client.CheckActivationCode(ctx, &usersvc.CheckActivationCodeRequest{Code: code})
	if err != nil {
		log.Warn("checking activation code", pkglog.Err(err))

		return dto.FromGrpcErr(err)
	}

	return nil
}

func (h *UserHandler) StreamDocumentAnalyzed(ctx context.Context, userID *string, passed *bool, send func(model.DocumentAnalyzedEvent) error) error {
	log := pkglog.WithMetadata(pkglog.WithMethod(h.log, "StreamDocumentAnalyzed"), utils.MetadataFromCtx(ctx))

	for {
		if ctx.Err() != nil {
			return nil
		}

		streamCtx, cancelStream := context.WithCancel(ctx)
		stream, err := h.client.StreamDocumentAnalyzed(streamCtx, &usersvc.StreamDocumentAnalyzedRequest{
			UserId: userID,
			Passed: passed,
		})
		if err != nil {
			cancelStream()
			if ctx.Err() != nil {
				return nil
			}
			if isUnavailable(err) {
				log.Warn("transient error opening document stream, reconnecting", pkglog.Err(err))
				select {
				case <-time.After(streamReconnectDelay):
				case <-ctx.Done():
					return nil
				}
				continue
			}
			log.Warn("streaming document analyzed", pkglog.Err(err))
			return dto.FromGrpcErr(err)
		}

		for {
			msg, err := stream.Recv()
			if errors.Is(err, io.EOF) {
				cancelStream()
				return nil
			}
			if err != nil {
				cancelStream()
				if ctx.Err() != nil {
					return nil
				}
				if isUnavailable(err) {
					log.Warn("document stream interrupted, reconnecting", pkglog.Err(err))
					select {
					case <-time.After(streamReconnectDelay):
					case <-ctx.Done():
						return nil
					}
					break
				}
				log.Warn("receiving document analyzed stream", pkglog.Err(err))
				return dto.FromGrpcErr(err)
			}

			defects := make([]model.DocumentDefect, len(msg.GetDefects()))
			for i, d := range msg.GetDefects() {
				defects[i] = model.DocumentDefect{
					Type:        d.GetType(),
					Description: d.GetDescription(),
				}
			}

			if err = send(model.DocumentAnalyzedEvent{
				DocumentID: msg.GetDocumentId(),
				UserID:     msg.GetUserId(),
				Passed:     msg.GetPassed(),
				Defects:    defects,
			}); err != nil {
				cancelStream()
				return err
			}
		}
	}
}
