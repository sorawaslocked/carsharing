package service_test

import (
	"testing"

	"carsharing/user-service/internal/model"
	"carsharing/user-service/internal/validation"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// --- SendActivationCode ---

func TestSendActivationCode_Success(t *testing.T) {
	d := newDeps(t)
	svc := newService(t, d)
	ctx := ctxWithUser(testUserID)
	user := baseUser()

	d.userRepo.EXPECT().FindByID(ctx, testUserID).Return(user, nil)
	d.codeStorage.EXPECT().Save(ctx, testUserID).Return(testCode, nil)
	d.mailer.EXPECT().SendActivationCode(ctx, user.Email, testCode).Return(nil)

	err := svc.SendActivationCode(ctx)

	require.NoError(t, err)
}

func TestSendActivationCode_NoUserIDInCtx(t *testing.T) {
	d := newDeps(t)
	svc := newService(t, d)

	err := svc.SendActivationCode(ctxAnon())

	assert.ErrorIs(t, err, model.ErrUnauthenticated)
}

func TestSendActivationCode_UserNotFound(t *testing.T) {
	d := newDeps(t)
	svc := newService(t, d)
	ctx := ctxWithUser(testUserID)

	d.userRepo.EXPECT().FindByID(ctx, testUserID).Return(model.User{}, model.ErrNotFound)

	err := svc.SendActivationCode(ctx)

	assert.ErrorIs(t, err, model.ErrNotFound)
}

// --- CheckActivationCode ---

func TestCheckActivationCode_Success(t *testing.T) {
	d := newDeps(t)
	svc := newService(t, d)
	ctx := ctxWithUser(testUserID)
	user := baseUser()
	codeHash := mustHashPassword(t, testCode)

	d.userRepo.EXPECT().FindByID(ctx, testUserID).Return(user, nil)
	d.codeStorage.EXPECT().Get(ctx, testUserID).Return(codeHash, nil)
	d.userRepo.EXPECT().Update(ctx, testUserID, mock.MatchedBy(func(u model.UserUpdate) bool {
		return u.IsEmailVerified != nil && *u.IsEmailVerified
	})).Return(nil)
	d.publisher.EXPECT().PublishUserUpdated(ctx, testUserID, false).Return(nil)

	err := svc.CheckActivationCode(ctx, testCode)

	require.NoError(t, err)
}

func TestCheckActivationCode_NoUserIDInCtx(t *testing.T) {
	d := newDeps(t)
	svc := newService(t, d)

	err := svc.CheckActivationCode(ctxAnon(), testCode)

	assert.ErrorIs(t, err, model.ErrUnauthenticated)
}

func TestCheckActivationCode_AlreadyVerified(t *testing.T) {
	d := newDeps(t)
	svc := newService(t, d)
	ctx := ctxWithUser(testUserID)
	user := baseUser()
	user.IsEmailVerified = true

	d.userRepo.EXPECT().FindByID(ctx, testUserID).Return(user, nil)

	err := svc.CheckActivationCode(ctx, testCode)

	assert.ErrorIs(t, err, model.ErrAlreadyExists)
}

func TestCheckActivationCode_InvalidFormat(t *testing.T) {
	d := newDeps(t)
	svc := newService(t, d)
	ctx := ctxWithUser(testUserID)

	d.userRepo.EXPECT().FindByID(ctx, testUserID).Return(baseUser(), nil)

	// Not 6 chars, not uppercase, not alphanumeric.
	err := svc.CheckActivationCode(ctx, "short")

	var ve validation.Errors
	require.ErrorAs(t, err, &ve)
	assert.Contains(t, ve, "code")
}

func TestCheckActivationCode_ExpiredCode(t *testing.T) {
	d := newDeps(t)
	svc := newService(t, d)
	ctx := ctxWithUser(testUserID)

	d.userRepo.EXPECT().FindByID(ctx, testUserID).Return(baseUser(), nil)
	// Redis returns ErrNotFound when the key has expired.
	d.codeStorage.EXPECT().Get(ctx, testUserID).Return(nil, model.ErrNotFound)

	err := svc.CheckActivationCode(ctx, testCode)

	var ve validation.Errors
	require.ErrorAs(t, err, &ve)
	assert.Contains(t, ve, "code")
}

func TestCheckActivationCode_WrongCode(t *testing.T) {
	d := newDeps(t)
	svc := newService(t, d)
	ctx := ctxWithUser(testUserID)
	// Hash of a different code is stored.
	codeHash := mustHashPassword(t, "XYZ789")

	d.userRepo.EXPECT().FindByID(ctx, testUserID).Return(baseUser(), nil)
	d.codeStorage.EXPECT().Get(ctx, testUserID).Return(codeHash, nil)

	err := svc.CheckActivationCode(ctx, testCode)

	var ve validation.Errors
	require.ErrorAs(t, err, &ve)
	assert.Contains(t, ve, "code")
}

func TestCheckActivationCode_UserNotFound(t *testing.T) {
	d := newDeps(t)
	svc := newService(t, d)
	ctx := ctxWithUser(testUserID)

	d.userRepo.EXPECT().FindByID(ctx, testUserID).Return(model.User{}, model.ErrNotFound)

	err := svc.CheckActivationCode(ctx, testCode)

	assert.ErrorIs(t, err, model.ErrNotFound)
}
