package handler_test

import (
	"context"
	"io"
	"log/slog"
	"testing"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	usersvc "carsharing/protos/gen/service/user"
	sharedmodel "carsharing/shared/model"
	"carsharing/user-service/internal/adapter/grpc/handler"
	"carsharing/user-service/internal/adapter/grpc/handler/mocks"
	"carsharing/user-service/internal/model"
	"carsharing/user-service/internal/validation"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Context keys mirror internal/adapter/grpc/interceptor/base.go exported constants.
const (
	ctxKeyRequestID = "x-request-id"
	ctxKeyClientIP  = "x-client-ip"
	ctxKeyUserID    = "x-user-id"
	ctxKeyUserRoles = "x-user-roles"
)

const (
	testUserID  = "11111111-1111-1111-1111-111111111111"
	testDocID   = "22222222-2222-2222-2222-222222222222"
	testEmail   = "john@example.com"
	testFName   = "John"
	testLName   = "Doe"
	testPasswd  = "Password1!"
	testObjKey  = "users/1234_abcdef"
	testImgType = "id_front"
	testCode    = "ABC123"
)

var testBirthDate = time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC)

func ctxWithUser(userID string) context.Context {
	ctx := context.Background()
	ctx = context.WithValue(ctx, ctxKeyRequestID, "test-req-id")
	ctx = context.WithValue(ctx, ctxKeyClientIP, "127.0.0.1")
	ctx = context.WithValue(ctx, ctxKeyUserID, userID)
	ctx = context.WithValue(ctx, ctxKeyUserRoles, []sharedmodel.Role{sharedmodel.RoleUser})
	return ctx
}

func ctxAnon() context.Context {
	ctx := context.Background()
	ctx = context.WithValue(ctx, ctxKeyRequestID, "test-req-id")
	ctx = context.WithValue(ctx, ctxKeyClientIP, "127.0.0.1")
	return ctx
}

func ptr[T any](v T) *T { return &v }

func newHandler(t *testing.T) (*handler.UserHandler, *mocks.MockUserService) {
	t.Helper()
	svc := mocks.NewMockUserService(t)
	h := handler.NewUserHandler(slog.New(slog.NewTextHandler(io.Discard, nil)), svc, &noopSubscriber{})
	return h, svc
}

// noopSubscriber satisfies handler.DocumentAnalyzedSubscriber for tests that
// don't exercise streaming; returns a pre-closed channel so nothing blocks.
type noopSubscriber struct{}

func (n *noopSubscriber) SubscribeStream(_ *string, _ *bool) (<-chan model.DocumentAnalyzedEvent, func()) {
	ch := make(chan model.DocumentAnalyzedEvent)
	close(ch)
	return ch, func() {}
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

func grpcCode(err error) codes.Code {
	if s, ok := status.FromError(err); ok {
		return s.Code()
	}
	return codes.Unknown
}

// --- CreateUser ---

func TestCreateUser_Success(t *testing.T) {
	h, svc := newHandler(t)
	ctx := ctxWithUser(testUserID)

	svc.EXPECT().Create(ctx, validation.UserCreate{
		Email:                testEmail,
		FirstName:            testFName,
		LastName:             testLName,
		BirthDate:            testBirthDate,
		Password:             testPasswd,
		PasswordConfirmation: testPasswd,
	}).Return(testUserID, nil)

	resp, err := h.CreateUser(ctx, &usersvc.CreateUserRequest{
		Email:                testEmail,
		FirstName:            testFName,
		LastName:             testLName,
		BirthDate:            "1990-01-01",
		Password:             testPasswd,
		PasswordConfirmation: testPasswd,
	})

	require.NoError(t, err)
	require.NotNil(t, resp.Id)
	assert.Equal(t, testUserID, *resp.Id)
}

func TestCreateUser_InvalidBirthDate(t *testing.T) {
	h, _ := newHandler(t)
	ctx := ctxWithUser(testUserID)

	_, err := h.CreateUser(ctx, &usersvc.CreateUserRequest{
		Email:     testEmail,
		BirthDate: "not-a-date",
	})

	assert.Equal(t, codes.InvalidArgument, grpcCode(err))
}

func TestCreateUser_ServiceError_DuplicateEmail(t *testing.T) {
	h, svc := newHandler(t)
	ctx := ctxWithUser(testUserID)

	svc.EXPECT().Create(ctx, validation.UserCreate{
		Email:                testEmail,
		FirstName:            testFName,
		LastName:             testLName,
		BirthDate:            testBirthDate,
		Password:             testPasswd,
		PasswordConfirmation: testPasswd,
	}).Return("", model.ErrDuplicateEmail)

	_, err := h.CreateUser(ctx, &usersvc.CreateUserRequest{
		Email:                testEmail,
		FirstName:            testFName,
		LastName:             testLName,
		BirthDate:            "1990-01-01",
		Password:             testPasswd,
		PasswordConfirmation: testPasswd,
	})

	assert.Equal(t, codes.AlreadyExists, grpcCode(err))
}

// --- GetUser ---

func TestGetUser_Success(t *testing.T) {
	h, svc := newHandler(t)
	ctx := ctxWithUser(testUserID)

	svc.EXPECT().Get(ctx, testUserID).Return(baseUser(), nil)

	resp, err := h.GetUser(ctx, &usersvc.GetUserRequest{Id: testUserID})

	require.NoError(t, err)
	require.NotNil(t, resp.User)
	assert.Equal(t, testUserID, resp.User.Id)
	assert.Equal(t, testEmail, resp.User.Email)
}

func TestGetUser_NotFound(t *testing.T) {
	h, svc := newHandler(t)
	ctx := ctxWithUser(testUserID)

	svc.EXPECT().Get(ctx, testUserID).Return(model.User{}, model.ErrUserNotFound)

	_, err := h.GetUser(ctx, &usersvc.GetUserRequest{Id: testUserID})

	assert.Equal(t, codes.NotFound, grpcCode(err))
}

func TestGetUser_WithProfileImageURL(t *testing.T) {
	h, svc := newHandler(t)
	ctx := ctxWithUser(testUserID)
	user := baseUser()
	user.ProfileImage = sharedmodel.Image{Key: testObjKey, URL: "https://minio/signed"}

	svc.EXPECT().Get(ctx, testUserID).Return(user, nil)

	resp, err := h.GetUser(ctx, &usersvc.GetUserRequest{Id: testUserID})

	require.NoError(t, err)
	require.NotNil(t, resp.User.ProfileImage)
	assert.Equal(t, testObjKey, resp.User.ProfileImage.Key)
	assert.Equal(t, "https://minio/signed", resp.User.ProfileImage.Url)
}

// --- ListUsers ---

func TestListUsers_Success(t *testing.T) {
	h, svc := newHandler(t)
	ctx := ctxWithUser(testUserID)

	svc.EXPECT().List(ctx, validation.UserFilter{}).Return([]model.User{baseUser()}, nil)

	resp, err := h.ListUsers(ctx, &usersvc.ListUsersRequest{})

	require.NoError(t, err)
	require.Len(t, resp.Users, 1)
	assert.Equal(t, testUserID, resp.Users[0].Id)
}

func TestListUsers_Empty(t *testing.T) {
	h, svc := newHandler(t)
	ctx := ctxWithUser(testUserID)

	svc.EXPECT().List(ctx, validation.UserFilter{}).Return(nil, nil)

	resp, err := h.ListUsers(ctx, &usersvc.ListUsersRequest{})

	require.NoError(t, err)
	assert.Empty(t, resp.Users)
}

func TestListUsers_WithFilter(t *testing.T) {
	h, svc := newHandler(t)
	ctx := ctxWithUser(testUserID)

	svc.EXPECT().List(ctx, validation.UserFilter{Email: ptr(testEmail)}).Return([]model.User{baseUser()}, nil)

	resp, err := h.ListUsers(ctx, &usersvc.ListUsersRequest{Email: ptr(testEmail)})

	require.NoError(t, err)
	require.Len(t, resp.Users, 1)
}

// --- UpdateUser ---

func TestUpdateUser_Success(t *testing.T) {
	h, svc := newHandler(t)
	ctx := ctxWithUser(testUserID)
	newName := "Jane"

	svc.EXPECT().Update(ctx, testUserID, validation.UserUpdate{FirstName: &newName}).Return(nil)

	_, err := h.UpdateUser(ctx, &usersvc.UpdateUserRequest{Id: testUserID, FirstName: &newName})

	require.NoError(t, err)
}

func TestUpdateUser_NotFound(t *testing.T) {
	h, svc := newHandler(t)
	ctx := ctxWithUser(testUserID)
	newName := "Jane"

	svc.EXPECT().Update(ctx, testUserID, validation.UserUpdate{FirstName: &newName}).Return(model.ErrUserNotFound)

	_, err := h.UpdateUser(ctx, &usersvc.UpdateUserRequest{Id: testUserID, FirstName: &newName})

	assert.Equal(t, codes.NotFound, grpcCode(err))
}

func TestUpdateUser_InvalidBirthDate(t *testing.T) {
	h, _ := newHandler(t)
	ctx := ctxWithUser(testUserID)

	_, err := h.UpdateUser(ctx, &usersvc.UpdateUserRequest{Id: testUserID, BirthDate: ptr("bad-date")})

	assert.Equal(t, codes.InvalidArgument, grpcCode(err))
}

// --- DeleteUser ---

func TestDeleteUser_Success(t *testing.T) {
	h, svc := newHandler(t)
	ctx := ctxWithUser(testUserID)

	svc.EXPECT().Delete(ctx, testUserID).Return(nil)

	resp, err := h.DeleteUser(ctx, &usersvc.DeleteUserRequest{Id: testUserID})

	require.NoError(t, err)
	assert.Equal(t, &emptypb.Empty{}, resp)
}

func TestDeleteUser_NotFound(t *testing.T) {
	h, svc := newHandler(t)
	ctx := ctxWithUser(testUserID)

	svc.EXPECT().Delete(ctx, testUserID).Return(model.ErrUserNotFound)

	_, err := h.DeleteUser(ctx, &usersvc.DeleteUserRequest{Id: testUserID})

	assert.Equal(t, codes.NotFound, grpcCode(err))
}

// --- Register ---

func TestRegister_Success(t *testing.T) {
	h, svc := newHandler(t)
	ctx := ctxAnon()

	svc.EXPECT().Register(ctx, validation.UserCreate{
		Email:                testEmail,
		FirstName:            testFName,
		LastName:             testLName,
		BirthDate:            testBirthDate,
		Password:             testPasswd,
		PasswordConfirmation: testPasswd,
	}).Return(testUserID, nil)

	resp, err := h.Register(ctx, &usersvc.RegisterRequest{
		Email:                testEmail,
		FirstName:            testFName,
		LastName:             testLName,
		BirthDate:            "1990-01-01",
		Password:             testPasswd,
		PasswordConfirmation: testPasswd,
	})

	require.NoError(t, err)
	assert.Equal(t, testUserID, resp.Id)
}

// --- SignIn ---

func TestSignIn_Success(t *testing.T) {
	h, svc := newHandler(t)
	ctx := ctxAnon()

	svc.EXPECT().SignIn(ctx, validation.Credentials{Email: ptr(testEmail), Password: testPasswd}).
		Return(testUserID, nil)

	resp, err := h.SignIn(ctx, &usersvc.SignInRequest{Email: ptr(testEmail), Password: testPasswd})

	require.NoError(t, err)
	assert.Equal(t, testUserID, resp.Id)
}

func TestSignIn_Unauthenticated(t *testing.T) {
	h, svc := newHandler(t)
	ctx := ctxAnon()

	svc.EXPECT().SignIn(ctx, validation.Credentials{Email: ptr(testEmail), Password: testPasswd}).
		Return("", model.ErrUnauthenticated)

	_, err := h.SignIn(ctx, &usersvc.SignInRequest{Email: ptr(testEmail), Password: testPasswd})

	assert.Equal(t, codes.Unauthenticated, grpcCode(err))
}

// --- SendActivationCode ---

func TestSendActivationCode_Success(t *testing.T) {
	h, svc := newHandler(t)
	ctx := ctxWithUser(testUserID)

	svc.EXPECT().SendActivationCode(ctx).Return(nil)

	_, err := h.SendActivationCode(ctx, &emptypb.Empty{})

	require.NoError(t, err)
}

func TestSendActivationCode_Unauthenticated(t *testing.T) {
	h, svc := newHandler(t)
	ctx := ctxWithUser(testUserID)

	svc.EXPECT().SendActivationCode(ctx).Return(model.ErrUnauthenticated)

	_, err := h.SendActivationCode(ctx, &emptypb.Empty{})

	assert.Equal(t, codes.Unauthenticated, grpcCode(err))
}

func TestSendActivationCode_Throttled(t *testing.T) {
	h, svc := newHandler(t)
	ctx := ctxWithUser(testUserID)

	svc.EXPECT().SendActivationCode(ctx).Return(model.ErrActivationCodeResendTooSoon)

	_, err := h.SendActivationCode(ctx, &emptypb.Empty{})

	assert.Equal(t, codes.ResourceExhausted, grpcCode(err))
}

// --- CheckActivationCode ---

func TestCheckActivationCode_Success(t *testing.T) {
	h, svc := newHandler(t)
	ctx := ctxWithUser(testUserID)

	svc.EXPECT().CheckActivationCode(ctx, testCode).Return(nil)

	_, err := h.CheckActivationCode(ctx, &usersvc.CheckActivationCodeRequest{Code: testCode})

	require.NoError(t, err)
}

func TestCheckActivationCode_InvalidCode(t *testing.T) {
	h, svc := newHandler(t)
	ctx := ctxWithUser(testUserID)
	ve := validation.Errors{"code": validation.ErrInvalidActivationCode}

	svc.EXPECT().CheckActivationCode(ctx, testCode).Return(ve)

	_, err := h.CheckActivationCode(ctx, &usersvc.CheckActivationCodeRequest{Code: testCode})

	assert.Equal(t, codes.InvalidArgument, grpcCode(err))
}

// --- GetProfileImageUploadData ---

func TestGetProfileImageUploadData_Success(t *testing.T) {
	h, svc := newHandler(t)
	ctx := ctxWithUser(testUserID)
	data := sharedmodel.ImageUploadData{PresignedPutURL: "https://minio/put", ObjectKey: testObjKey}

	svc.EXPECT().GetUserProfileImageUploadData(ctx).Return(data, nil)

	resp, err := h.GetProfileImageUploadData(ctx, &emptypb.Empty{})

	require.NoError(t, err)
	require.NotNil(t, resp.UploadData)
	assert.Equal(t, "https://minio/put", resp.UploadData.PresignedPutUrl)
	assert.Equal(t, testObjKey, resp.UploadData.ObjectKey)
}

// --- CreateDocument ---

func TestCreateDocument_Success(t *testing.T) {
	h, svc := newHandler(t)
	ctx := ctxWithUser(testUserID)

	svc.EXPECT().CreateDocument(ctx, validation.DocumentCreate{ObjectKey: testObjKey, ImageType: testImgType}).Return(testDocID, nil)

	resp, err := h.CreateDocument(ctx, &usersvc.CreateDocumentRequest{
		ObjectKey: testObjKey,
		ImageType: testImgType,
	})

	require.NoError(t, err)
	assert.Equal(t, testDocID, resp.Id)
}

func TestCreateDocument_InvalidImageType(t *testing.T) {
	h, svc := newHandler(t)
	ctx := ctxWithUser(testUserID)

	svc.EXPECT().CreateDocument(ctx, validation.DocumentCreate{ObjectKey: testObjKey, ImageType: "invalid_type"}).
		Return("", validation.Errors{"image_type": validation.ErrInvalidImageType})

	_, err := h.CreateDocument(ctx, &usersvc.CreateDocumentRequest{
		ObjectKey: testObjKey,
		ImageType: "invalid_type",
	})

	assert.Equal(t, codes.InvalidArgument, grpcCode(err))
}

// --- GetUploadDocumentData ---

func TestGetUploadDocumentData_Success(t *testing.T) {
	h, svc := newHandler(t)
	ctx := ctxWithUser(testUserID)
	data := sharedmodel.ImageUploadData{PresignedPutURL: "https://minio/put", ObjectKey: "documents/id_front/key"}

	svc.EXPECT().GetDocumentImageUploadData(ctx, testImgType).Return(data, nil)

	resp, err := h.GetUploadDocumentData(ctx, &usersvc.GetUploadDocumentDataRequest{ImageType: testImgType})

	require.NoError(t, err)
	require.NotNil(t, resp.UploadData)
	assert.Equal(t, "https://minio/put", resp.UploadData.PresignedPutUrl)
}

// --- ListDocuments ---

func TestListDocuments_Success(t *testing.T) {
	h, svc := newHandler(t)
	ctx := ctxWithUser(testUserID)
	docs := []model.Document{
		{
			ID:        testDocID,
			UserID:    testUserID,
			ImageType: model.DocumentImageTypeIDFront,
			Status:    model.DocumentStatusApproved,
			Image:     sharedmodel.Image{Key: testObjKey, URL: "https://minio/signed"},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	svc.EXPECT().ListDocuments(ctx, validation.DocumentFilter{UserID: testUserID}).Return(docs, nil)

	resp, err := h.ListDocuments(ctx, &usersvc.ListDocumentsRequest{UserId: testUserID})

	require.NoError(t, err)
	require.Len(t, resp.Documents, 1)
	assert.Equal(t, testDocID, resp.Documents[0].Id)
	require.NotNil(t, resp.Documents[0].Image)
	assert.Equal(t, "https://minio/signed", resp.Documents[0].Image.Url)
}

func TestListDocuments_UserNotFound(t *testing.T) {
	h, svc := newHandler(t)
	ctx := ctxWithUser(testUserID)

	svc.EXPECT().ListDocuments(ctx, validation.DocumentFilter{UserID: testUserID}).Return(nil, model.ErrUserNotFound)

	_, err := h.ListDocuments(ctx, &usersvc.ListDocumentsRequest{UserId: testUserID})

	assert.Equal(t, codes.NotFound, grpcCode(err))
}

// --- CheckDocument ---

func TestCheckDocument_Success(t *testing.T) {
	h, svc := newHandler(t)
	ctx := ctxWithUser(testUserID)

	svc.EXPECT().CheckDocument(ctx, testDocID, validation.DocumentUpdate{Status: "approved"}).Return(nil)

	_, err := h.CheckDocument(ctx, &usersvc.CheckDocumentRequest{
		DocId:  testDocID,
		Status: "approved",
	})

	require.NoError(t, err)
}

func TestCheckDocument_WithError(t *testing.T) {
	h, svc := newHandler(t)
	ctx := ctxWithUser(testUserID)
	errMsg := ptr("blurry image")

	svc.EXPECT().CheckDocument(ctx, testDocID, validation.DocumentUpdate{Status: "rejected", Error: errMsg}).Return(nil)

	_, err := h.CheckDocument(ctx, &usersvc.CheckDocumentRequest{
		DocId:  testDocID,
		Status: "rejected",
		Error:  errMsg,
	})

	require.NoError(t, err)
}

func TestCheckDocument_InvalidStatus(t *testing.T) {
	h, svc := newHandler(t)
	ctx := ctxWithUser(testUserID)

	svc.EXPECT().CheckDocument(ctx, testDocID, validation.DocumentUpdate{Status: "unknown_status"}).
		Return(validation.Errors{"status": validation.ErrInvalidDocumentStatus})

	_, err := h.CheckDocument(ctx, &usersvc.CheckDocumentRequest{
		DocId:  testDocID,
		Status: "unknown_status",
	})

	assert.Equal(t, codes.InvalidArgument, grpcCode(err))
}

func TestCheckDocument_NonReviewableStatus(t *testing.T) {
	h, svc := newHandler(t)
	ctx := ctxWithUser(testUserID)

	svc.EXPECT().CheckDocument(ctx, testDocID, validation.DocumentUpdate{Status: "pending"}).
		Return(validation.Errors{"status": validation.ErrDocumentStatusNotReviewable})

	_, err := h.CheckDocument(ctx, &usersvc.CheckDocumentRequest{
		DocId:  testDocID,
		Status: "pending",
	})

	assert.Equal(t, codes.InvalidArgument, grpcCode(err))
}
