package service_test

import (
	"context"
	"io"
	"log/slog"
	"testing"
	"time"

	sharedmodel "carsharing/shared/model"
	"carsharing/user-service/internal/model"
	"carsharing/user-service/internal/pkg/security"
	"carsharing/user-service/internal/service"
	"carsharing/user-service/internal/service/mocks"
	"carsharing/user-service/internal/validation"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/require"
)

// Context key constants mirror internal/adapter/grpc/interceptor/base.go.
const (
	ctxKeyRequestID = "x-request-id"
	ctxKeyClientIP  = "x-client-ip"
	ctxKeyUserID    = "x-user-id"
	ctxKeyUserRoles = "x-user-roles"
)

// Shared test fixtures.
const (
	testUserID  = "11111111-1111-4111-8111-111111111111"
	testDocID   = "22222222-2222-4222-8222-222222222222"
	testEmail   = "john@example.com"
	testPhone   = "+12345678901"
	testFName   = "John"
	testLName   = "Doe"
	testPasswd  = "Password1!"
	testObjKey  = "users/1234_abcdef"
	testImgType = "id_front"
	testCode    = "ABC123"
)

var testBirthDate = time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC)

// ctxWithUser returns a context populated the same way BaseInterceptor would.
func ctxWithUser(userID string) context.Context {
	ctx := context.Background()
	ctx = context.WithValue(ctx, ctxKeyRequestID, "test-req-id")
	ctx = context.WithValue(ctx, ctxKeyClientIP, "127.0.0.1")
	ctx = context.WithValue(ctx, ctxKeyUserID, userID)
	ctx = context.WithValue(ctx, ctxKeyUserRoles, []sharedmodel.Role{sharedmodel.RoleUser})
	return ctx
}

// ctxAnon returns a context without a user ID (unauthenticated request).
func ctxAnon() context.Context {
	ctx := context.Background()
	ctx = context.WithValue(ctx, ctxKeyRequestID, "test-req-id")
	ctx = context.WithValue(ctx, ctxKeyClientIP, "127.0.0.1")
	return ctx
}

func ptr[T any](v T) *T { return &v }

type deps struct {
	userRepo    *mocks.MockUserRepository
	docRepo     *mocks.MockDocumentRepository
	storage     *mocks.MockObjectStorage
	analyzer    *mocks.MockDocumentAnalyzer
	publisher   *mocks.MockUserPublisher
	codeStorage *mocks.MockActivationCodeStorage
	mailer      *mocks.MockMailer
}

func newDeps(t *testing.T) *deps {
	t.Helper()
	return &deps{
		userRepo:    mocks.NewMockUserRepository(t),
		docRepo:     mocks.NewMockDocumentRepository(t),
		storage:     mocks.NewMockObjectStorage(t),
		analyzer:    mocks.NewMockDocumentAnalyzer(t),
		publisher:   mocks.NewMockUserPublisher(t),
		codeStorage: mocks.NewMockActivationCodeStorage(t),
		mailer:      mocks.NewMockMailer(t),
	}
}

func newService(t *testing.T, d *deps) *service.UserService {
	t.Helper()
	v := validator.New()
	require.NoError(t, validation.RegisterCustomValidators(v, slog.New(slog.NewTextHandler(io.Discard, nil))))
	return service.NewUserService(
		slog.New(slog.NewTextHandler(io.Discard, nil)),
		v,
		d.userRepo, d.docRepo, d.storage, d.analyzer, d.publisher, d.codeStorage, d.mailer,
	)
}

func baseUser() model.User {
	return model.User{
		ID:        testUserID,
		Email:     testEmail,
		FirstName: testFName,
		LastName:  testLName,
		BirthDate: testBirthDate,
		Roles:     []sharedmodel.Role{sharedmodel.RoleUser},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func validUserCreate() validation.UserCreate {
	return validation.UserCreate{
		Email:                testEmail,
		FirstName:            testFName,
		LastName:             testLName,
		BirthDate:            testBirthDate,
		Password:             testPasswd,
		PasswordConfirmation: testPasswd,
	}
}

func mustHashPassword(t *testing.T, s string) []byte {
	t.Helper()
	h, err := security.HashString(s)
	require.NoError(t, err)
	return h
}
