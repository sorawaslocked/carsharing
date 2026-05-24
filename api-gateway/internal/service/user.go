package service

import (
	"context"
	"log/slog"

	"carsharing/api-gateway/internal/model"
	pkglog "carsharing/shared/pkg/log"
	"carsharing/shared/pkg/utils"
)

const (
	ctxDeviceIDKey = "x-device-id"
	ctxUserIDKey   = "x-user-id"
)

type UserService struct {
	log          *slog.Logger
	presenter    UserPresenter
	tokenManager TokenManager
	sessionCache UserSessionCache
}

func NewUserService(presenter UserPresenter, tokenManager TokenManager, sessionCache UserSessionCache, log *slog.Logger) *UserService {
	return &UserService{
		log:          pkglog.WithComponent(log, "service.UserService"),
		presenter:    presenter,
		tokenManager: tokenManager,
		sessionCache: sessionCache,
	}
}

func (s *UserService) Create(ctx context.Context, data model.UserCreate) (string, error) {
	log := pkglog.WithMetadata(pkglog.WithMethod(s.log, "Create"), utils.MetadataFromCtx(ctx))

	id, err := s.presenter.Create(ctx, data)
	if err != nil {
		log.Warn("creating user", pkglog.Err(err))

		return "", err
	}

	return id, nil
}

func (s *UserService) Get(ctx context.Context, id string) (model.User, error) {
	log := pkglog.WithMetadata(pkglog.WithMethod(s.log, "Get"), utils.MetadataFromCtx(ctx))

	user, err := s.presenter.Get(ctx, id)
	if err != nil {
		log.Warn("getting user", pkglog.Err(err))

		return model.User{}, err
	}

	return user, nil
}

func (s *UserService) List(ctx context.Context, filter model.UserFilter) ([]model.User, error) {
	log := pkglog.WithMetadata(pkglog.WithMethod(s.log, "List"), utils.MetadataFromCtx(ctx))

	users, err := s.presenter.List(ctx, filter)
	if err != nil {
		log.Warn("listing users", pkglog.Err(err))

		return nil, err
	}

	return users, nil
}

func (s *UserService) Update(ctx context.Context, id string, data model.UserUpdate) error {
	log := pkglog.WithMetadata(pkglog.WithMethod(s.log, "Update"), utils.MetadataFromCtx(ctx))

	if err := s.presenter.Update(ctx, id, data); err != nil {
		log.Warn("updating user", pkglog.Err(err))

		return err
	}

	return nil
}

func (s *UserService) Delete(ctx context.Context, id string) error {
	log := pkglog.WithMetadata(pkglog.WithMethod(s.log, "Delete"), utils.MetadataFromCtx(ctx))

	if err := s.presenter.Delete(ctx, id); err != nil {
		log.Warn("deleting user", pkglog.Err(err))

		return err
	}

	return nil
}

func (s *UserService) Register(ctx context.Context, data model.UserCreate) (string, error) {
	log := pkglog.WithMetadata(pkglog.WithMethod(s.log, "Register"), utils.MetadataFromCtx(ctx))

	id, err := s.presenter.Register(ctx, data)
	if err != nil {
		log.Warn("registering user", pkglog.Err(err))

		return "", err
	}

	return id, nil
}

func (s *UserService) SignIn(ctx context.Context, creds model.Credentials) (model.AccessToken, model.RefreshToken, error) {
	log := pkglog.WithMetadata(pkglog.WithMethod(s.log, "SignIn"), utils.MetadataFromCtx(ctx))

	deviceID := ctx.Value(ctxDeviceIDKey).(string)

	id, err := s.presenter.SignIn(ctx, creds)
	if err != nil {
		return model.AccessToken{}, model.RefreshToken{}, err
	}

	if err = s.sessionCache.SetSignedIn(ctx, id, deviceID, true); err != nil {
		log.Error("setting session", pkglog.Err(err))

		return model.AccessToken{}, model.RefreshToken{}, err
	}

	accessToken, accessTokenExp, err := s.tokenManager.GenerateAccessToken(ctx, id)
	if err != nil {
		log.Error("generating access token", pkglog.Err(err))

		return model.AccessToken{}, model.RefreshToken{}, err
	}

	refreshToken, refreshTokenExp, err := s.tokenManager.GenerateRefreshToken(ctx, id)
	if err != nil {
		log.Error("generating refresh token", pkglog.Err(err))

		return model.AccessToken{}, model.RefreshToken{}, err
	}

	return model.AccessToken{
			Token:     accessToken,
			ExpiresIn: accessTokenExp.Unix(),
		}, model.RefreshToken{
			Token:     refreshToken,
			ExpiresIn: refreshTokenExp.Unix(),
		}, nil
}

func (s *UserService) RefreshToken(ctx context.Context, refreshToken string) (model.AccessToken, model.RefreshToken, error) {
	log := pkglog.WithMetadata(pkglog.WithMethod(s.log, "RefreshToken"), utils.MetadataFromCtx(ctx))

	deviceID := ctx.Value(ctxDeviceIDKey).(string)

	id, _, err := s.tokenManager.ParseToken(ctx, refreshToken)
	if err != nil {
		return model.AccessToken{}, model.RefreshToken{}, err
	}

	isLoggedIn, err := s.sessionCache.IsSignedIn(ctx, id, deviceID)
	if err != nil {
		log.Error("checking session", pkglog.Err(err))

		return model.AccessToken{}, model.RefreshToken{}, err
	}
	if !isLoggedIn {
		return model.AccessToken{}, model.RefreshToken{}, model.ErrUnauthorized
	}

	if err = s.sessionCache.SetSignedIn(ctx, id, deviceID, true); err != nil {
		log.Error("setting session", pkglog.Err(err))

		return model.AccessToken{}, model.RefreshToken{}, err
	}

	newAccessToken, newAccessTokenExp, err := s.tokenManager.GenerateAccessToken(ctx, id)
	if err != nil {
		log.Error("generating access token", pkglog.Err(err))

		return model.AccessToken{}, model.RefreshToken{}, err
	}

	newRefreshToken, newRefreshTokenExp, err := s.tokenManager.GenerateRefreshToken(ctx, id)
	if err != nil {
		log.Error("generating refresh token", pkglog.Err(err))

		return model.AccessToken{}, model.RefreshToken{}, err
	}

	return model.AccessToken{
			Token:     newAccessToken,
			ExpiresIn: newAccessTokenExp.Unix(),
		}, model.RefreshToken{
			Token:     newRefreshToken,
			ExpiresIn: newRefreshTokenExp.Unix(),
		}, nil
}

func (s *UserService) SignOut(ctx context.Context) error {
	log := pkglog.WithMetadata(pkglog.WithMethod(s.log, "SignOut"), utils.MetadataFromCtx(ctx))

	deviceID := ctx.Value(ctxDeviceIDKey).(string)
	userID := ctx.Value(ctxUserIDKey).(string)

	if err := s.sessionCache.SetSignedIn(ctx, userID, deviceID, false); err != nil {
		log.Error("clearing session", pkglog.Err(err))

		return err
	}

	return nil
}

func (s *UserService) GetProfile(ctx context.Context) (model.User, error) {
	userID := ctx.Value(ctxUserIDKey).(string)

	return s.Get(ctx, userID)
}

func (s *UserService) UpdateProfile(ctx context.Context, data model.UserProfileUpdate) error {
	userID := ctx.Value(ctxUserIDKey).(string)

	return s.Update(ctx, userID, model.UserUpdate{
		PhoneNumber:     data.PhoneNumber,
		FirstName:       data.FirstName,
		LastName:        data.LastName,
		BirthDate:       data.BirthDate,
		Password:        data.Password,
		ProfileImageKey: data.ProfileImageKey,
	})
}

func (s *UserService) SendActivationCode(ctx context.Context) error {
	log := pkglog.WithMetadata(pkglog.WithMethod(s.log, "SendActivationCode"), utils.MetadataFromCtx(ctx))

	if err := s.presenter.SendActivationCode(ctx); err != nil {
		log.Warn("sending activation code", pkglog.Err(err))

		return err
	}

	return nil
}

func (s *UserService) CheckActivationCode(ctx context.Context, code string) error {
	log := pkglog.WithMetadata(pkglog.WithMethod(s.log, "CheckActivationCode"), utils.MetadataFromCtx(ctx))

	if err := s.presenter.CheckActivationCode(ctx, code); err != nil {
		log.Warn("checking activation code", pkglog.Err(err))

		return err
	}

	return nil
}
