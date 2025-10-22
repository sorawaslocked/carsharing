package handler

import (
	"car-rental-user-service/internal/model"
	"context"
	svc "github.com/sorawaslocked/car-rental-protos/gen/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log/slog"
	"time"
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
	birthDate, err := time.Parse("2006-01-04", req.DateOfBirth)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	cred := model.Credentials{
		Email:                &req.Email,
		PhoneNumber:          &req.PhoneNumber,
		Password:             req.Password,
		PasswordConfirmation: &req.PasswordConfirmation,
		FirstName:            &req.FirstName,
		LastName:             &req.LastName,
		BirthDate:            &birthDate,
	}
	newID, errs := h.authService.Register(ctx, cred)

	res := &svc.RegisterResponse{
		Errors: fromErrToStringMap(errs),
	}
	if newID != 0 {
		res.Id = &newID
	}

	return res, nil
}

func (h *AuthHandler) Login(ctx context.Context, req *svc.LoginRequest) (*svc.LoginResponse, error) {
	cred := model.Credentials{
		Email:       req.Email,
		PhoneNumber: req.PhoneNumber,
		Password:    req.Password,
	}
	token, errs := h.authService.Login(ctx, cred)

	res := &svc.LoginResponse{
		Errors: fromErrToStringMap(errs),
	}
	if token.AccessToken != "" && token.RefreshToken != "" {
		res.AccessToken = &token.AccessToken
		res.RefreshToken = &token.RefreshToken
	}

	return res, nil
}

func (h *AuthHandler) RefreshToken(ctx context.Context, req *svc.RefreshTokenRequest) (*svc.RefreshTokenResponse, error) {
	token, err := h.authService.RefreshToken(ctx, req.RefreshToken)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}

	res := &svc.RefreshTokenResponse{
		AccessToken:  &token.AccessToken,
		RefreshToken: &token.RefreshToken,
	}

	return res, nil
}
