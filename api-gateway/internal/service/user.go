package service

import (
	"context"

	"github.com/sorawaslocked/car-rental-api-gateway/internal/model"
)

const (
	ctxDeviceIDKey = "x-device-id"
	ctxUserIDKey   = "x-user-id"
)

type UserService struct {
	presenter    UserPresenter
	tokenManager TokenManager
	sessionCache UserSessionCache
}

func NewUserService(presenter UserPresenter, tokenManager TokenManager, sessionCache UserSessionCache) *UserService {
	return &UserService{
		presenter:    presenter,
		tokenManager: tokenManager,
		sessionCache: sessionCache,
	}
}

func (s *UserService) Create(ctx context.Context, data model.UserCreate) (string, error) {
	return s.presenter.Create(ctx, data)
}

func (s *UserService) Get(ctx context.Context, id string) (model.User, error) {
	return s.presenter.Get(ctx, id)
}

func (s *UserService) List(ctx context.Context, filter model.UserFilter) ([]model.User, error) {
	return s.presenter.List(ctx, filter)
}

func (s *UserService) Update(ctx context.Context, id string, data model.UserUpdate) error {
	return s.presenter.Update(ctx, id, data)
}

func (s *UserService) Delete(ctx context.Context, id string) error {
	return s.presenter.Delete(ctx, id)
}

func (s *UserService) Register(ctx context.Context, data model.UserCreate) (string, error) {
	return s.presenter.Register(ctx, data)
}

func (s *UserService) SignIn(ctx context.Context, creds model.Credentials) (model.AccessToken, model.RefreshToken, error) {
	deviceID := ctx.Value(ctxDeviceIDKey).(string)

	id, err := s.presenter.SignIn(ctx, creds)
	if err != nil {
		return model.AccessToken{}, model.RefreshToken{}, err
	}

	err = s.sessionCache.SetSignedIn(ctx, id, deviceID, true)
	if err != nil {
		return model.AccessToken{}, model.RefreshToken{}, err
	}

	accessToken, accessTokenExp, err := s.tokenManager.GenerateAccessToken(id)
	if err != nil {
		return model.AccessToken{}, model.RefreshToken{}, err
	}
	refreshToken, refreshTokenExp, err := s.tokenManager.GenerateRefreshToken(id)
	if err != nil {
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
	deviceID := ctx.Value(ctxDeviceIDKey).(string)

	id, err := s.tokenManager.ParseToken(refreshToken)
	if err != nil {
		return model.AccessToken{}, model.RefreshToken{}, err
	}

	isLoggedIn, err := s.sessionCache.IsSignedIn(ctx, id, deviceID)
	if err != nil {
		return model.AccessToken{}, model.RefreshToken{}, err
	}
	if !isLoggedIn {
		return model.AccessToken{}, model.RefreshToken{}, model.ErrUnauthorized
	}

	err = s.sessionCache.SetSignedIn(ctx, id, deviceID, true)
	if err != nil {
		return model.AccessToken{}, model.RefreshToken{}, err
	}

	newAccessToken, newAccessTokenExp, err := s.tokenManager.GenerateAccessToken(id)
	if err != nil {
		return model.AccessToken{}, model.RefreshToken{}, err
	}
	newRefreshToken, newRefreshTokenExp, err := s.tokenManager.GenerateRefreshToken(id)
	if err != nil {
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
	deviceID := ctx.Value(ctxDeviceIDKey).(string)
	userID := ctx.Value(ctxUserIDKey).(string)

	err := s.sessionCache.SetSignedIn(ctx, userID, deviceID, false)

	return err
}

func (s *UserService) Me(ctx context.Context) (model.User, error) {
	userID := ctx.Value(ctxUserIDKey).(string)

	return s.Get(ctx, userID)
}

func (s *UserService) SendActivationCode(ctx context.Context) error {
	return s.presenter.SendActivationCode(ctx)
}

func (s *UserService) CheckActivationCode(ctx context.Context, code string) error {
	return s.presenter.CheckActivationCode(ctx, code)
}

func (s *UserService) CreateDocument(ctx context.Context, objectKey, imageType string) (string, error) {
	return s.presenter.CreateDocument(ctx, objectKey, imageType)
}

func (s *UserService) GetUploadDocumentData(ctx context.Context, imageType string) (model.ImageUploadData, error) {
	return s.presenter.GetDocumentImageUploadData(ctx, imageType)
}

func (s *UserService) GetProcessedDocumentsForUser(ctx context.Context, userID string) ([]model.Document, error) {
	return s.presenter.GetProcessedDocumentsForUser(ctx, userID)
}

func (s *UserService) CheckDocument(ctx context.Context, docID string, status string, documentError *string) error {
	return s.presenter.CheckDocument(ctx, docID, status, documentError)
}
