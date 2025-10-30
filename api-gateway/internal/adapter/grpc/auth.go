package grpc

import (
	"context"
	"github.com/sorawaslocked/car-rental-api-gateway/internal/model"
	authsvc "github.com/sorawaslocked/car-rental-protos/gen/service/auth"
)

type AuthHandler struct {
	client authsvc.AuthServiceClient
}

func NewAuthHandler(client authsvc.AuthServiceClient) *AuthHandler {
	return &AuthHandler{client: client}
}

func (h *AuthHandler) Register(ctx context.Context, cred model.Credentials) (uint64, error) {
	res, err := h.client.Register(ctx, &authsvc.RegisterRequest{
		Email:                cred.Email,
		PhoneNumber:          cred.PhoneNumber,
		Password:             cred.Password,
		PasswordConfirmation: cred.PasswordConfirmation,
		FirstName:            cred.FirstName,
		LastName:             cred.LastName,
		BirthDate:            cred.BirthDate,
	})
	if err != nil {
		return 0, fromGrpcErr(err)
	}

	return *res.Id, nil
}

func (h *AuthHandler) Login(ctx context.Context, cred model.Credentials) (model.Token, error) {
	req := &authsvc.LoginRequest{
		Password: cred.Password,
	}
	if cred.Email != "" {
		req.Email = &cred.Email
	}
	if cred.PhoneNumber != "" {
		req.PhoneNumber = &cred.PhoneNumber
	}

	res, err := h.client.Login(ctx, req)
	if err != nil {
		return model.Token{}, fromGrpcErr(err)
	}

	return model.Token{
		AccessToken:  *res.AccessToken,
		RefreshToken: *res.RefreshToken,
	}, nil
}

func (h *AuthHandler) RefreshToken(ctx context.Context, refreshToken string) (model.Token, error) {
	res, err := h.client.RefreshToken(ctx, &authsvc.RefreshTokenRequest{
		RefreshToken: refreshToken,
	})
	if err != nil {
		return model.Token{}, fromGrpcErr(err)
	}

	return model.Token{
		AccessToken:  *res.AccessToken,
		RefreshToken: *res.RefreshToken,
	}, nil
}
