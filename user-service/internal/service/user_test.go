package service_test

import (
	"testing"

	sharedmodel "carsharing/shared/model"
	"carsharing/user-service/internal/model"
	"carsharing/user-service/internal/validation"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// --- Create / Register ---

func TestCreate_Success(t *testing.T) {
	d := newDeps(t)
	svc := newService(t, d)
	ctx := ctxWithUser(testUserID)

	d.userRepo.EXPECT().Insert(ctx, mock.MatchedBy(func(u model.User) bool {
		return u.Email == testEmail && len(u.PasswordHash) > 0
	})).Return(testUserID, nil)
	d.publisher.EXPECT().PublishUserCreated(ctx, testUserID).Return(nil)

	id, err := svc.Create(ctx, validUserCreate())

	require.NoError(t, err)
	assert.Equal(t, testUserID, id)
}

func TestCreate_DuplicateEmail(t *testing.T) {
	d := newDeps(t)
	svc := newService(t, d)
	ctx := ctxWithUser(testUserID)

	d.userRepo.EXPECT().Insert(ctx, mock.Anything).Return("", model.ErrDuplicateEmail)

	_, err := svc.Create(ctx, validUserCreate())

	assert.ErrorIs(t, err, model.ErrDuplicateEmail)
}

func TestCreate_DuplicatePhone(t *testing.T) {
	d := newDeps(t)
	svc := newService(t, d)
	ctx := ctxWithUser(testUserID)
	input := validUserCreate()
	input.PhoneNumber = ptr(testPhone)

	d.userRepo.EXPECT().Insert(ctx, mock.Anything).Return("", model.ErrDuplicatePhone)

	_, err := svc.Create(ctx, input)

	assert.ErrorIs(t, err, model.ErrDuplicatePhone)
}

func TestCreate_ValidationError_InvalidEmail(t *testing.T) {
	d := newDeps(t)
	svc := newService(t, d)
	ctx := ctxWithUser(testUserID)
	input := validUserCreate()
	input.Email = "not-an-email"

	_, err := svc.Create(ctx, input)

	var ve validation.Errors
	require.ErrorAs(t, err, &ve)
	assert.Contains(t, ve, "email")
}

func TestCreate_ValidationError_WeakPassword(t *testing.T) {
	d := newDeps(t)
	svc := newService(t, d)
	ctx := ctxWithUser(testUserID)
	input := validUserCreate()
	input.Password = "weakpassword"
	input.PasswordConfirmation = "weakpassword"

	_, err := svc.Create(ctx, input)

	var ve validation.Errors
	require.ErrorAs(t, err, &ve)
	assert.Contains(t, ve, "password")
}

func TestRegister_Success(t *testing.T) {
	d := newDeps(t)
	svc := newService(t, d)
	ctx := ctxAnon()

	d.userRepo.EXPECT().Insert(ctx, mock.MatchedBy(func(u model.User) bool {
		return u.Email == testEmail
	})).Return(testUserID, nil)
	d.publisher.EXPECT().PublishUserCreated(ctx, testUserID).Return(nil)

	id, err := svc.Register(ctx, validUserCreate())

	require.NoError(t, err)
	assert.Equal(t, testUserID, id)
}

// --- Get ---

func TestGet_Success_NoProfileImage(t *testing.T) {
	d := newDeps(t)
	svc := newService(t, d)
	ctx := ctxWithUser(testUserID)
	user := baseUser()

	d.userRepo.EXPECT().FindByID(ctx, testUserID).Return(user, nil)

	got, err := svc.Get(ctx, testUserID)

	require.NoError(t, err)
	assert.Equal(t, testUserID, got.ID)
	assert.Equal(t, testEmail, got.Email)
	assert.Empty(t, got.ProfileImage.Key)
}

func TestGet_Success_WithProfileImage(t *testing.T) {
	d := newDeps(t)
	svc := newService(t, d)
	ctx := ctxWithUser(testUserID)
	user := baseUser()
	user.ProfileImage = sharedmodel.Image{Key: testObjKey}
	signedURL := "https://minio.example.com/signed-url"

	d.userRepo.EXPECT().FindByID(ctx, testUserID).Return(user, nil)
	d.storage.EXPECT().GetImageURL(ctx, testObjKey).Return(signedURL, nil)

	got, err := svc.Get(ctx, testUserID)

	require.NoError(t, err)
	assert.Equal(t, testObjKey, got.ProfileImage.Key)
	assert.Equal(t, signedURL, got.ProfileImage.URL)
}

func TestGet_NotFound(t *testing.T) {
	d := newDeps(t)
	svc := newService(t, d)
	ctx := ctxWithUser(testUserID)

	d.userRepo.EXPECT().FindByID(ctx, testUserID).Return(model.User{}, model.ErrNotFound)

	_, err := svc.Get(ctx, testUserID)

	assert.ErrorIs(t, err, model.ErrNotFound)
}

// --- List ---

func TestList_Success_Empty(t *testing.T) {
	d := newDeps(t)
	svc := newService(t, d)
	ctx := ctxWithUser(testUserID)

	d.userRepo.EXPECT().Find(ctx, model.UserFilter{
		Pagination: &sharedmodel.Pagination{Limit: sharedmodel.DefaultPaginationLimit, Offset: sharedmodel.DefaultPaginationOffset},
	}).Return(nil, nil)

	users, err := svc.List(ctx, validation.UserFilter{})

	require.NoError(t, err)
	assert.Empty(t, users)
}

func TestList_Success_MultipleUsers(t *testing.T) {
	d := newDeps(t)
	svc := newService(t, d)
	ctx := ctxWithUser(testUserID)
	u1 := baseUser()
	u2 := baseUser()
	u2.ID = "33333333-3333-4333-8333-333333333333"
	u2.Email = "jane@example.com"

	d.userRepo.EXPECT().Find(ctx, model.UserFilter{
		Pagination: &sharedmodel.Pagination{Limit: sharedmodel.DefaultPaginationLimit, Offset: sharedmodel.DefaultPaginationOffset},
	}).Return([]model.User{u1, u2}, nil)

	users, err := svc.List(ctx, validation.UserFilter{})

	require.NoError(t, err)
	require.Len(t, users, 2)
}

func TestList_ResolvesProfileImageURLs(t *testing.T) {
	d := newDeps(t)
	svc := newService(t, d)
	ctx := ctxWithUser(testUserID)
	u := baseUser()
	u.ProfileImage = sharedmodel.Image{Key: testObjKey}
	signedURL := "https://minio.example.com/signed-url"

	d.userRepo.EXPECT().Find(ctx, model.UserFilter{
		Pagination: &sharedmodel.Pagination{Limit: sharedmodel.DefaultPaginationLimit, Offset: sharedmodel.DefaultPaginationOffset},
	}).Return([]model.User{u}, nil)
	d.storage.EXPECT().GetImageURL(ctx, testObjKey).Return(signedURL, nil)

	users, err := svc.List(ctx, validation.UserFilter{})

	require.NoError(t, err)
	require.Len(t, users, 1)
	assert.Equal(t, signedURL, users[0].ProfileImage.URL)
}

// --- Update ---

func TestUpdate_Success_NameChange(t *testing.T) {
	d := newDeps(t)
	svc := newService(t, d)
	ctx := ctxWithUser(testUserID)
	newName := "Jane"

	d.userRepo.EXPECT().Update(ctx, testUserID, mock.MatchedBy(func(u model.UserUpdate) bool {
		return u.FirstName != nil && *u.FirstName == newName
	})).Return(nil)
	d.publisher.EXPECT().PublishUserUpdated(ctx, testUserID, false).Return(nil)

	err := svc.Update(ctx, testUserID, validation.UserUpdate{FirstName: &newName})

	require.NoError(t, err)
}

func TestUpdate_NotFound(t *testing.T) {
	d := newDeps(t)
	svc := newService(t, d)
	ctx := ctxWithUser(testUserID)

	d.userRepo.EXPECT().Update(ctx, testUserID, mock.Anything).Return(model.ErrNotFound)

	err := svc.Update(ctx, testUserID, validation.UserUpdate{FirstName: ptr("Jane")})

	assert.ErrorIs(t, err, model.ErrNotFound)
}

func TestUpdate_PasswordMismatch(t *testing.T) {
	d := newDeps(t)
	svc := newService(t, d)
	ctx := ctxWithUser(testUserID)

	err := svc.Update(ctx, testUserID, validation.UserUpdate{
		Password:             ptr(testPasswd),
		PasswordConfirmation: ptr("DifferentPass1!"),
	})

	var ve validation.Errors
	require.ErrorAs(t, err, &ve)
	assert.Contains(t, ve, "passwordConfirmation")
}

func TestUpdate_IsSecurityUpdate_WhenPasswordChanged(t *testing.T) {
	d := newDeps(t)
	svc := newService(t, d)
	ctx := ctxWithUser(testUserID)

	d.userRepo.EXPECT().Update(ctx, testUserID, mock.MatchedBy(func(u model.UserUpdate) bool {
		return len(u.PasswordHash) > 0
	})).Return(nil)
	d.publisher.EXPECT().PublishUserUpdated(ctx, testUserID, true).Return(nil)

	err := svc.Update(ctx, testUserID, validation.UserUpdate{
		Password:             ptr(testPasswd),
		PasswordConfirmation: ptr(testPasswd),
	})

	require.NoError(t, err)
}

func TestUpdate_IsSecurityUpdate_WhenRolesChanged(t *testing.T) {
	d := newDeps(t)
	svc := newService(t, d)
	ctx := ctxWithUser(testUserID)

	d.userRepo.EXPECT().Update(ctx, testUserID, mock.MatchedBy(func(u model.UserUpdate) bool {
		return len(u.Roles) == 1 && u.Roles[0] == sharedmodel.RoleAdmin
	})).Return(nil)
	d.publisher.EXPECT().PublishUserUpdated(ctx, testUserID, true).Return(nil)

	err := svc.Update(ctx, testUserID, validation.UserUpdate{Roles: []string{"admin"}})

	require.NoError(t, err)
}

// --- Delete ---

func TestDelete_Success(t *testing.T) {
	d := newDeps(t)
	svc := newService(t, d)
	ctx := ctxWithUser(testUserID)

	d.userRepo.EXPECT().Delete(ctx, testUserID).Return(nil)
	d.publisher.EXPECT().PublishUserDeleted(ctx, testUserID).Return(nil)

	err := svc.Delete(ctx, testUserID)

	require.NoError(t, err)
}

func TestDelete_NotFound(t *testing.T) {
	d := newDeps(t)
	svc := newService(t, d)
	ctx := ctxWithUser(testUserID)

	d.userRepo.EXPECT().Delete(ctx, testUserID).Return(model.ErrNotFound)

	err := svc.Delete(ctx, testUserID)

	assert.ErrorIs(t, err, model.ErrNotFound)
}

// --- SignIn ---

func TestSignIn_Success(t *testing.T) {
	d := newDeps(t)
	svc := newService(t, d)
	ctx := ctxAnon()

	user := baseUser()
	user.PasswordHash = mustHashPassword(t, testPasswd)

	d.userRepo.EXPECT().FindOne(ctx, model.UserFilter{Email: ptr(testEmail)}).Return(user, nil)

	id, err := svc.SignIn(ctx, validation.Credentials{Email: ptr(testEmail), Password: testPasswd})

	require.NoError(t, err)
	assert.Equal(t, testUserID, id)
}

func TestSignIn_WrongPassword(t *testing.T) {
	d := newDeps(t)
	svc := newService(t, d)
	ctx := ctxAnon()

	user := baseUser()
	user.PasswordHash = mustHashPassword(t, testPasswd)

	d.userRepo.EXPECT().FindOne(ctx, model.UserFilter{Email: ptr(testEmail)}).Return(user, nil)

	_, err := svc.SignIn(ctx, validation.Credentials{Email: ptr(testEmail), Password: "WrongPass1!"})

	assert.ErrorIs(t, err, model.ErrUnauthenticated)
}

func TestSignIn_UserNotFound(t *testing.T) {
	d := newDeps(t)
	svc := newService(t, d)
	ctx := ctxAnon()

	d.userRepo.EXPECT().FindOne(ctx, model.UserFilter{Email: ptr(testEmail)}).Return(model.User{}, model.ErrNotFound)

	_, err := svc.SignIn(ctx, validation.Credentials{Email: ptr(testEmail), Password: testPasswd})

	assert.ErrorIs(t, err, model.ErrNotFound)
}

func TestSignIn_ByPhone_Success(t *testing.T) {
	d := newDeps(t)
	svc := newService(t, d)
	ctx := ctxAnon()

	user := baseUser()
	user.PhoneNumber = ptr(testPhone)
	user.PasswordHash = mustHashPassword(t, testPasswd)

	d.userRepo.EXPECT().FindOne(ctx, model.UserFilter{PhoneNumber: ptr(testPhone)}).Return(user, nil)

	id, err := svc.SignIn(ctx, validation.Credentials{PhoneNumber: ptr(testPhone), Password: testPasswd})

	require.NoError(t, err)
	assert.Equal(t, testUserID, id)
}

// --- GetUserProfileImageUploadData ---

func TestGetUserProfileImageUploadData_Success(t *testing.T) {
	d := newDeps(t)
	svc := newService(t, d)
	ctx := ctxWithUser(testUserID)
	expected := sharedmodel.ImageUploadData{PresignedPutURL: "https://minio/put", ObjectKey: testObjKey}

	d.storage.EXPECT().GetUserProfileImageUploadData(ctx).Return(expected, nil)

	got, err := svc.GetUserProfileImageUploadData(ctx)

	require.NoError(t, err)
	assert.Equal(t, expected, got)
}

func TestGetUserProfileImageUploadData_StorageError(t *testing.T) {
	d := newDeps(t)
	svc := newService(t, d)
	ctx := ctxWithUser(testUserID)

	d.storage.EXPECT().GetUserProfileImageUploadData(ctx).Return(sharedmodel.ImageUploadData{}, model.ErrObjectStorage)

	_, err := svc.GetUserProfileImageUploadData(ctx)

	assert.ErrorIs(t, err, model.ErrObjectStorage)
}
