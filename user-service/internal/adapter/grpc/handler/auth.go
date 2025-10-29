package handler

import (
	"context"
	authsvc "github.com/sorawaslocked/car-rental-protos/gen/service/auth"
	"github.com/sorawaslocked/car-rental-user-service/internal/adapter/grpc/handler/dto"
	"log/slog"
)

type AuthHandler struct {
	log         *slog.Logger
	authService AuthService
	authsvc.UnimplementedAuthServiceServer
}

func NewAuthHandler(log *slog.Logger, authService AuthService) *AuthHandler {
	return &AuthHandler{
		log:         log,
		authService: authService,
	}
}

func (h *AuthHandler) Register(ctx context.Context, req *authsvc.RegisterRequest) (*authsvc.RegisterResponse, error) {
	data, validationErrors := dto.FromRegisterRequest(req)
	if validationErrors != nil {
		return &authsvc.RegisterResponse{}, dto.ToStatusCodeError(validationErrors)
	}

	id, err := h.authService.Register(ctx, data)
	if err != nil {
		return &authsvc.RegisterResponse{}, dto.ToStatusCodeError(err)
	}

	return &authsvc.RegisterResponse{Id: &id}, nil
}

func (h *AuthHandler) Login(ctx context.Context, req *authsvc.LoginRequest) (*authsvc.LoginResponse, error) {
	cred := dto.FromLoginRequest(req)

	token, err := h.authService.Login(ctx, cred)
	if err != nil {
		return &authsvc.LoginResponse{}, dto.ToStatusCodeError(err)
	}

	return &authsvc.LoginResponse{
		AccessToken:  &token.AccessToken,
		RefreshToken: &token.RefreshToken,
	}, nil
}

func (h *AuthHandler) RefreshToken(ctx context.Context, req *authsvc.RefreshTokenRequest) (*authsvc.RefreshTokenResponse, error) {
	token, err := h.authService.RefreshToken(ctx, req.RefreshToken)
	if err != nil {
		return &authsvc.RefreshTokenResponse{}, dto.ToStatusCodeError(err)
	}

	return &authsvc.RefreshTokenResponse{
		AccessToken:  &token.AccessToken,
		RefreshToken: &token.RefreshToken,
	}, nil
}
