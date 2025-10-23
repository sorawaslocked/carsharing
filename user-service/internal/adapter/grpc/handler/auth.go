package handler

import (
	"context"
	"github.com/sorawaslocked/car-rental-user-service/internal/model"
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
	cred := model.Credentials{}
	if req.Email != "" {
		cred.Email = &req.Email
	}
	if req.PhoneNumber != "" {
		cred.PhoneNumber = &req.PhoneNumber
	}
	if req.PasswordConfirmation != "" {
		cred.PasswordConfirmation = &req.PasswordConfirmation
	}
	if req.FirstName != "" {
		cred.FirstName = &req.FirstName
	}
	if req.LastName != "" {
		cred.LastName = &req.LastName
	}
	cred.Password = req.Password

	birthDate, err := time.Parse("2006-01-04", req.DateOfBirth)
	if err == nil {
		cred.BirthDate = &birthDate
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
