package service

import (
	"context"
	"github.com/go-playground/validator/v10"
	"github.com/sorawaslocked/car-rental-user-service/internal/model"
	"github.com/sorawaslocked/car-rental-user-service/internal/pkg/logger"
	"github.com/sorawaslocked/car-rental-user-service/internal/pkg/security"
	"log/slog"
)

type AuthService struct {
	log         *slog.Logger
	validate    *validator.Validate
	jwtProvider JwtProvider
	userService *UserService
	userRepo    UserRepository
}

func NewAuthService(
	log *slog.Logger,
	validate *validator.Validate,
	jwtProvider JwtProvider,
	userService *UserService,
	userRepo UserRepository,
) *AuthService {
	return &AuthService{
		log:         log,
		validate:    validate,
		jwtProvider: jwtProvider,
		userService: userService,
		userRepo:    userRepo,
	}
}

func (s *AuthService) Register(ctx context.Context, data model.UserCreateData) (uint64, error) {
	err := validateInput(s.validate, data)
	if err != nil {
		return 0, err
	}

	createdID, err := s.userService.Insert(ctx, data)
	if err != nil {
		return 0, err
	}

	return createdID, nil
}

func (s *AuthService) Login(ctx context.Context, cred model.Credentials) (model.Token, error) {
	err := validateInput(s.validate, cred)
	if err != nil {
		return model.Token{}, err
	}
	filter := model.UserFilter{
		Email: &cred.Email,
	}

	user, err := s.userService.FindOne(ctx, filter)
	if err != nil {
		return model.Token{}, err
	}

	err = security.CheckPassword(cred.Password, user.PasswordHash)
	if err != nil {
		return model.Token{}, model.ValidationErrors{
			"password": model.ErrPasswordsDoNotMatch,
		}
	}

	userRoles := toRoleStrings(user.Roles)
	accessToken, err := s.jwtProvider.GenerateAccessToken(user.ID, userRoles)
	if err != nil {
		s.log.Error(
			"jwt: generating access token",
			logger.Err(err),
			slog.Uint64("userId", user.ID),
		)

		return model.Token{}, model.ErrJwt
	}

	refreshToken, err := s.jwtProvider.GenerateRefreshToken(user.ID, userRoles)
	if err != nil {
		s.log.Error(
			"jwt: generating refresh token",
			logger.Err(err),
			slog.Uint64("userId", user.ID),
		)

		return model.Token{}, model.ErrJwt
	}

	return model.Token{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *AuthService) RefreshToken(_ context.Context, refreshToken string) (model.Token, error) {
	input := refreshTokenValidation{
		RefreshToken: refreshToken,
	}
	err := validateInput(s.validate, input)
	if err != nil {
		return model.Token{}, err
	}

	id, roles, err := s.jwtProvider.VerifyAndParseClaims(refreshToken)
	if err != nil {
		s.log.Error(
			"jwt: verifying refresh token",
			logger.Err(err),
			slog.String("refreshToken", refreshToken),
		)

		return model.Token{}, model.ErrJwt
	}

	newAccessToken, err := s.jwtProvider.GenerateAccessToken(id, roles)
	if err != nil {
		s.log.Error(
			"jwt: generating access token",
			logger.Err(err),
			slog.Uint64("userId", id),
		)

		return model.Token{}, model.ErrJwt
	}
	newRefreshToken, err := s.jwtProvider.GenerateRefreshToken(id, roles)
	if err != nil {
		s.log.Error(
			"jwt: generating refresh token",
			logger.Err(err),
			slog.Uint64("userId", id),
		)

		return model.Token{}, model.ErrJwt
	}

	return model.Token{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
	}, nil
}
