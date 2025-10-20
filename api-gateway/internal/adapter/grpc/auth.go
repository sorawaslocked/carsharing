package grpc

import (
	"context"
	"github.com/sorawaslocked/car-rental-api-gateway/internal/model"
	svc "github.com/sorawaslocked/car-rental-protos/gen/service"
)

type AuthHandler struct {
	client svc.AuthServiceClient
}

func NewAuthHandler(client svc.AuthServiceClient) *AuthHandler {
	return &AuthHandler{client: client}
}

func (h *AuthHandler) Register(ctx context.Context, cred model.Credentials) (uint64, map[string]string) {
	res, err := h.client.Register(ctx, &svc.RegisterRequest{
		FirstName:            *cred.FirstName,
		LastName:             *cred.LastName,
		Email:                *cred.Email,
		PhoneNumber:          *cred.PhoneNumber,
		DateOfBirth:          *cred.DateOfBirth,
		Password:             cred.Password,
		PasswordConfirmation: *cred.PasswordConfirmation,
	})
	if err != nil {
		return 0, grpcError(err)
	}
	if res.Id == nil {
		return 0, res.Errors
	}

	return *res.Id, nil
}

func (h *AuthHandler) Login(ctx context.Context, cred model.Credentials) (model.Token, map[string]string) {
	req := &svc.LoginRequest{
		Password: cred.Password,
	}
	if cred.Email != nil {
		req.Credentials = &svc.LoginRequest_Email{Email: *cred.Email}
	}
	if cred.PhoneNumber != nil {
		req.Credentials = &svc.LoginRequest_PhoneNumber{PhoneNumber: *cred.PhoneNumber}
	}

	res, err := h.client.Login(ctx, req)
	if err != nil {
		return model.Token{}, grpcError(err)
	}
	if res.AccessToken == nil {
		return model.Token{}, res.Errors
	}

	return model.Token{
		AccessToken:  *res.AccessToken,
		RefreshToken: *res.RefreshToken,
	}, nil
}

func (h *AuthHandler) RefreshToken(ctx context.Context, refreshToken string) (model.Token, map[string]string) {
	res, err := h.client.RefreshToken(ctx, &svc.RefreshTokenRequest{
		RefreshToken: refreshToken,
	})
	if err != nil {
		return model.Token{}, grpcError(err)
	}
	if res.AccessToken == nil {
		return model.Token{}, res.Errors
	}

	return model.Token{
		AccessToken:  *res.AccessToken,
		RefreshToken: *res.RefreshToken,
	}, nil
}
