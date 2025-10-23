package handler

import (
	"context"
	svc "github.com/sorawaslocked/car-rental-protos/gen/service"
	"github.com/sorawaslocked/car-rental-user-service/internal/adapter/grpc/handler/dto"
	"log/slog"
)

type AuthHandler struct {
	log         *slog.Logger
	authService AuthService
	svc.UnimplementedAuthServiceServer
}

func NewAuthHandler(log *slog.Logger, authService AuthService) *AuthHandler {
	return &AuthHandler{
		log:         log,
		authService: authService,
	}
}

func (h *AuthHandler) Register(ctx context.Context, req *svc.RegisterRequest) (*svc.RegisterResponse, error) {
	cred, validationErrors := dto.FromRegisterRequest(req)
	if validationErrors != nil {
		return &svc.RegisterResponse{}, dto.ToStatusCodeError(validationErrors)
	}

	id, err := h.authService.Register(ctx, cred)
	if err != nil {
		return &svc.RegisterResponse{}, dto.ToStatusCodeError(err)
	}

	return &svc.RegisterResponse{Id: &id}, nil
}

func (h *AuthHandler) Login(ctx context.Context, req *svc.LoginRequest) (*svc.LoginResponse, error) {
	cred := dto.FromLoginRequest(req)

	token, err := h.authService.Login(ctx, cred)
	if err != nil {
		return &svc.LoginResponse{}, dto.ToStatusCodeError(err)
	}

	return &svc.LoginResponse{
		AccessToken:  &token.AccessToken,
		RefreshToken: &token.RefreshToken,
	}, nil
}

func (h *AuthHandler) RefreshToken(ctx context.Context, req *svc.RefreshTokenRequest) (*svc.RefreshTokenResponse, error) {
	token, err := h.authService.RefreshToken(ctx, req.RefreshToken)
	if err != nil {
		return &svc.RefreshTokenResponse{}, dto.ToStatusCodeError(err)
	}

	return &svc.RefreshTokenResponse{
		AccessToken:  &token.AccessToken,
		RefreshToken: &token.RefreshToken,
	}, nil
}
