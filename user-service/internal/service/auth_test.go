package service

import (
	"context"
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/sorawaslocked/car-rental-user-service/internal/model"
	validatecfg "github.com/sorawaslocked/car-rental-user-service/internal/pkg/validate"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"log/slog"
	"os"
	"testing"
	"time"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Insert(ctx context.Context, user model.User) (uint64, error) {
	args := m.Called(ctx, user)

	return args.Get(0).(uint64), args.Error(1)
}

func (m *MockUserRepository) FindOne(ctx context.Context, filter model.UserFilter) (model.User, error) {
	args := m.Called(ctx, filter)

	return args.Get(0).(model.User), args.Error(1)
}

func (m *MockUserRepository) Find(ctx context.Context, filter model.UserFilter) ([]model.User, error) {
	args := m.Called(ctx, filter)

	return args.Get(0).([]model.User), args.Error(1)
}

func (m *MockUserRepository) Update(ctx context.Context, filter model.UserFilter, update model.UserUpdateData) error {
	args := m.Called(ctx, filter, update)

	return args.Error(0)
}

func (m *MockUserRepository) Delete(ctx context.Context, filter model.UserFilter) error {
	args := m.Called(ctx, filter)

	return args.Error(0)
}

type MockJWTProvider struct {
	mock.Mock
}

func (m *MockJWTProvider) GenerateAccessToken(userID uint64, roles []string) (string, error) {
	args := m.Called(userID, roles)
	return args.String(0), args.Error(1)
}

func (m *MockJWTProvider) GenerateRefreshToken(userID uint64, roles []string) (string, error) {
	args := m.Called(userID, roles)
	return args.String(0), args.Error(1)
}

func (m *MockJWTProvider) VerifyAndParseClaims(token string) (uint64, []string, error) {
	args := m.Called(token)
	return args.Get(0).(uint64), args.Get(1).([]string), args.Error(2)
}

func setupAuthService() (*AuthService, *MockUserRepository, *MockJWTProvider) {
	log := slog.New(slog.NewTextHandler(os.Stdout, nil))
	validate := validator.New()
	validate.RegisterValidation("min_age", validatecfg.MinAge)
	validate.RegisterValidation("complex_password", validatecfg.ComplexPassword)
	mockRepo := new(MockUserRepository)
	mockJWT := new(MockJWTProvider)

	service := &AuthService{
		log:         log,
		validate:    validate,
		jwtProvider: mockJWT,
		userRepo:    mockRepo,
	}

	return service, mockRepo, mockJWT
}

func TestAuthService_Register_Success(t *testing.T) {
	service, mockRepo, _ := setupAuthService()
	ctx := context.Background()

	cred := model.Credentials{
		Email:                "test@example.com",
		PhoneNumber:          "+1234567890",
		Password:             "StrongPass123!",
		PasswordConfirmation: "StrongPass123!",
		FirstName:            "John",
		LastName:             "Doe",
		BirthDate:            time.Now().AddDate(-25, 0, 0),
	}

	expectedID := uint64(123)
	mockRepo.On("Insert", ctx, mock.AnythingOfType("model.User")).Return(expectedID, nil)

	userID, err := service.Register(ctx, cred)

	assert.NoError(t, err)
	assert.Equal(t, expectedID, userID)
	mockRepo.AssertExpectations(t)
}

func TestAuthService_Register_ValidationError_PasswordMismatch(t *testing.T) {
	service, _, _ := setupAuthService()
	ctx := context.Background()

	cred := model.Credentials{
		Email:                "test@example.com",
		PhoneNumber:          "+1234567890",
		Password:             "StrongPass123!",
		PasswordConfirmation: "DifferentPass123!",
		FirstName:            "John",
		LastName:             "Doe",
		BirthDate:            time.Now().AddDate(-25, 0, 0),
	}

	userID, err := service.Register(ctx, cred)

	assert.Error(t, err)
	assert.Equal(t, uint64(0), userID)
}

func TestAuthService_Register_ValidationError_InvalidEmail(t *testing.T) {
	service, _, _ := setupAuthService()
	ctx := context.Background()

	cred := model.Credentials{
		Email:                "invalid-email",
		PhoneNumber:          "+1234567890",
		Password:             "StrongPass123!",
		PasswordConfirmation: "StrongPass123!",
		FirstName:            "John",
		LastName:             "Doe",
		BirthDate:            time.Now().AddDate(-25, 0, 0),
	}

	userID, err := service.Register(ctx, cred)

	assert.Error(t, err)
	assert.Equal(t, uint64(0), userID)
}

func TestAuthService_Register_RepositoryError(t *testing.T) {
	service, mockRepo, _ := setupAuthService()
	ctx := context.Background()

	cred := model.Credentials{
		Email:                "test@example.com",
		PhoneNumber:          "+1234567890",
		Password:             "StrongPass123!",
		PasswordConfirmation: "StrongPass123!",
		FirstName:            "John",
		LastName:             "Doe",
		BirthDate:            time.Now().AddDate(-25, 0, 0),
	}

	mockRepo.On("Insert", ctx, mock.AnythingOfType("model.User")).Return(uint64(0), errors.New("database error"))

	userID, err := service.Register(ctx, cred)

	assert.Error(t, err)
	assert.Equal(t, model.ErrSql, err)
	assert.Equal(t, uint64(0), userID)
	mockRepo.AssertExpectations(t)
}

func TestAuthService_Login_Success_WithEmail(t *testing.T) {
	service, mockRepo, mockJWT := setupAuthService()
	ctx := context.Background()

	cred := model.Credentials{
		Email:    "test@example.com",
		Password: "StrongPass123!",
	}

	expectedUser := model.User{
		ID:           123,
		Email:        "test@example.com",
		PasswordHash: []byte("$2a$10$SCON6tCMIKvGcItvCo5zsOjeUjlzSjEfdQs2yonQsVWDO3mzmfuzG"), // hash of "StrongPass123!"
		Roles:        []model.Role{model.RoleUser},
	}

	email := cred.Email
	mockRepo.On("FindOne", ctx, mock.MatchedBy(func(f model.UserFilter) bool {
		return f.Email != nil && *f.Email == email
	})).Return(expectedUser, nil)

	mockJWT.On("GenerateAccessToken", uint64(123), []string{"user"}).Return("access_token_123", nil)
	mockJWT.On("GenerateRefreshToken", uint64(123), []string{"user"}).Return("refresh_token_123", nil)

	token, err := service.Login(ctx, cred)

	assert.NoError(t, err)
	assert.Equal(t, "access_token_123", token.AccessToken)
	assert.Equal(t, "refresh_token_123", token.RefreshToken)
	mockRepo.AssertExpectations(t)
	mockJWT.AssertExpectations(t)
}

func TestAuthService_Login_Success_WithPhoneNumber(t *testing.T) {
	service, mockRepo, mockJWT := setupAuthService()
	ctx := context.Background()

	cred := model.Credentials{
		PhoneNumber: "+1234567890",
		Password:    "StrongPass123!",
	}

	expectedUser := model.User{
		ID:           123,
		PhoneNumber:  "+1234567890",
		PasswordHash: []byte("$2a$10$SCON6tCMIKvGcItvCo5zsOjeUjlzSjEfdQs2yonQsVWDO3mzmfuzG"),
		Roles:        []model.Role{model.RoleUser},
	}

	phoneNumber := cred.PhoneNumber
	mockRepo.On("FindOne", ctx, mock.MatchedBy(func(f model.UserFilter) bool {
		return f.PhoneNumber != nil && *f.PhoneNumber == phoneNumber
	})).Return(expectedUser, nil)

	mockJWT.On("GenerateAccessToken", uint64(123), []string{"user"}).Return("access_token_123", nil)
	mockJWT.On("GenerateRefreshToken", uint64(123), []string{"user"}).Return("refresh_token_123", nil)

	token, err := service.Login(ctx, cred)

	assert.NoError(t, err)
	assert.Equal(t, "access_token_123", token.AccessToken)
	assert.Equal(t, "refresh_token_123", token.RefreshToken)
	mockRepo.AssertExpectations(t)
	mockJWT.AssertExpectations(t)
}

func TestAuthService_Login_ValidationError(t *testing.T) {
	service, _, _ := setupAuthService()
	ctx := context.Background()

	cred := model.Credentials{
		Email:    "", // Missing both email and phone
		Password: "StrongPass123!",
	}

	token, err := service.Login(ctx, cred)

	assert.Error(t, err)
	assert.Empty(t, token.AccessToken)
	assert.Empty(t, token.RefreshToken)
}

func TestAuthService_Login_UserNotFound(t *testing.T) {
	service, mockRepo, _ := setupAuthService()
	ctx := context.Background()

	cred := model.Credentials{
		Email:    "notfound@example.com",
		Password: "StrongPass123!",
	}

	email := cred.Email
	mockRepo.On("FindOne", ctx, mock.MatchedBy(func(f model.UserFilter) bool {
		return f.Email != nil && *f.Email == email
	})).Return(model.User{}, errors.New("not found"))

	token, err := service.Login(ctx, cred)

	assert.Error(t, err)
	assert.Equal(t, model.ErrNotFound, err)
	assert.Empty(t, token.AccessToken)
	mockRepo.AssertExpectations(t)
}

func TestAuthService_Login_WrongPassword(t *testing.T) {
	service, mockRepo, _ := setupAuthService()
	ctx := context.Background()

	cred := model.Credentials{
		Email:    "test@example.com",
		Password: "WrongPassword123!",
	}

	expectedUser := model.User{
		ID:           123,
		Email:        "test@example.com",
		PasswordHash: []byte("$2a$10$SCON6tCMIKvGcItvCo5zsOjeUjlzSjEfdQs2yonQsVWDO3mzmfuzG"), // hash of "StrongPass123!"
		Roles:        []model.Role{model.RoleUser},
	}

	email := cred.Email
	mockRepo.On("FindOne", ctx, mock.MatchedBy(func(f model.UserFilter) bool {
		return f.Email != nil && *f.Email == email
	})).Return(expectedUser, nil)

	token, err := service.Login(ctx, cred)

	assert.Error(t, err)
	assert.Equal(t, model.ErrPasswordsDoNotMatch, err)
	assert.Empty(t, token.AccessToken)
	mockRepo.AssertExpectations(t)
}

func TestAuthService_Login_AccessTokenGenerationError(t *testing.T) {
	service, mockRepo, mockJWT := setupAuthService()
	ctx := context.Background()

	cred := model.Credentials{
		Email:    "test@example.com",
		Password: "StrongPass123!",
	}

	expectedUser := model.User{
		ID:           123,
		Email:        "test@example.com",
		PasswordHash: []byte("$2a$10$SCON6tCMIKvGcItvCo5zsOjeUjlzSjEfdQs2yonQsVWDO3mzmfuzG"),
		Roles:        []model.Role{model.RoleUser},
	}

	email := cred.Email
	mockRepo.On("FindOne", ctx, mock.MatchedBy(func(f model.UserFilter) bool {
		return f.Email != nil && *f.Email == email
	})).Return(expectedUser, nil)

	mockJWT.On("GenerateAccessToken", uint64(123), []string{"user"}).Return("", errors.New("jwt error"))

	token, err := service.Login(ctx, cred)

	assert.Error(t, err)
	assert.Equal(t, model.ErrJwt, err)
	assert.Empty(t, token.AccessToken)
	mockRepo.AssertExpectations(t)
	mockJWT.AssertExpectations(t)
}

func TestAuthService_Login_RefreshTokenGenerationError(t *testing.T) {
	service, mockRepo, mockJWT := setupAuthService()
	ctx := context.Background()

	cred := model.Credentials{
		Email:    "test@example.com",
		Password: "StrongPass123!",
	}

	expectedUser := model.User{
		ID:           123,
		Email:        "test@example.com",
		PasswordHash: []byte("$2a$10$SCON6tCMIKvGcItvCo5zsOjeUjlzSjEfdQs2yonQsVWDO3mzmfuzG"),
		Roles:        []model.Role{model.RoleUser},
	}

	email := cred.Email
	mockRepo.On("FindOne", ctx, mock.MatchedBy(func(f model.UserFilter) bool {
		return f.Email != nil && *f.Email == email
	})).Return(expectedUser, nil)

	mockJWT.On("GenerateAccessToken", uint64(123), []string{"user"}).Return("access_token_123", nil)
	mockJWT.On("GenerateRefreshToken", uint64(123), []string{"user"}).Return("", errors.New("jwt error"))

	token, err := service.Login(ctx, cred)

	assert.Error(t, err)
	assert.Equal(t, model.ErrJwt, err)
	assert.Empty(t, token.RefreshToken)
	mockRepo.AssertExpectations(t)
	mockJWT.AssertExpectations(t)
}

func TestAuthService_RefreshToken_Success(t *testing.T) {
	service, _, mockJWT := setupAuthService()
	ctx := context.Background()

	refreshToken := "valid.refresh.token"
	userID := uint64(123)
	roles := []string{"user", "admin"}

	mockJWT.On("VerifyAndParseClaims", refreshToken).Return(userID, roles, nil)
	mockJWT.On("GenerateAccessToken", userID, roles).Return("new_access_token", nil)
	mockJWT.On("GenerateRefreshToken", userID, roles).Return("new_refresh_token", nil)

	token, err := service.RefreshToken(ctx, refreshToken)

	assert.NoError(t, err)
	assert.Equal(t, "new_access_token", token.AccessToken)
	assert.Equal(t, "new_refresh_token", token.RefreshToken)
	mockJWT.AssertExpectations(t)
}

func TestAuthService_RefreshToken_EmptyToken(t *testing.T) {
	service, _, _ := setupAuthService()
	ctx := context.Background()

	token, err := service.RefreshToken(ctx, "")

	assert.Error(t, err)
	assert.Equal(t, model.ErrRequiredField, err)
	assert.Empty(t, token.AccessToken)
}

func TestAuthService_RefreshToken_InvalidToken(t *testing.T) {
	service, _, _ := setupAuthService()
	ctx := context.Background()

	refreshToken := "invalid_token_format"

	token, err := service.RefreshToken(ctx, refreshToken)

	assert.Error(t, err)
	assert.Equal(t, model.ErrInvalidToken, err)
	assert.Empty(t, token.AccessToken)
}

func TestAuthService_RefreshToken_VerificationError(t *testing.T) {
	service, _, mockJWT := setupAuthService()
	ctx := context.Background()

	refreshToken := "valid.format.token"

	mockJWT.On("VerifyAndParseClaims", refreshToken).Return(uint64(0), []string(nil), errors.New("verification failed"))

	token, err := service.RefreshToken(ctx, refreshToken)

	assert.Error(t, err)
	assert.Equal(t, model.ErrJwt, err)
	assert.Empty(t, token.AccessToken)
	mockJWT.AssertExpectations(t)
}

func TestAuthService_RefreshToken_NewAccessTokenError(t *testing.T) {
	service, _, mockJWT := setupAuthService()
	ctx := context.Background()

	refreshToken := "valid.refresh.token"
	userID := uint64(123)
	roles := []string{"user"}

	mockJWT.On("VerifyAndParseClaims", refreshToken).Return(userID, roles, nil)
	mockJWT.On("GenerateAccessToken", userID, roles).Return("", errors.New("generation failed"))

	token, err := service.RefreshToken(ctx, refreshToken)

	assert.Error(t, err)
	assert.Equal(t, model.ErrJwt, err)
	assert.Empty(t, token.AccessToken)
	mockJWT.AssertExpectations(t)
}

func TestAuthService_RefreshToken_NewRefreshTokenError(t *testing.T) {
	service, _, mockJWT := setupAuthService()
	ctx := context.Background()

	refreshToken := "valid.refresh.token"
	userID := uint64(123)
	roles := []string{"user"}

	mockJWT.On("VerifyAndParseClaims", refreshToken).Return(userID, roles, nil)
	mockJWT.On("GenerateAccessToken", userID, roles).Return("new_access_token", nil)
	mockJWT.On("GenerateRefreshToken", userID, roles).Return("", errors.New("generation failed"))

	token, err := service.RefreshToken(ctx, refreshToken)

	assert.Error(t, err)
	assert.Equal(t, model.ErrJwt, err)
	assert.Empty(t, token.RefreshToken)
	mockJWT.AssertExpectations(t)
}
