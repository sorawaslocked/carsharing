package auth_test

import (
	"context"
	"io"
	"log/slog"
	"testing"

	usersvc "carsharing/protos/gen/service/user"
	sharedmodel "carsharing/shared/model"
	"carsharing/user-service/internal/adapter/grpc/interceptor/auth"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Fake request types that satisfy the duck-typed ownerExtract interfaces.
type fakeIDReq struct{ id string }

func (r fakeIDReq) GetId() string { return r.id }

type fakeUserIDReq struct{ userID string }

func (r fakeUserIDReq) GetUserId() string { return r.userID }

const (
	callerID = "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa"
	otherID  = "bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb"
)

func newInterceptor() *auth.Interceptor {
	return auth.NewInterceptor(slog.New(slog.NewTextHandler(io.Discard, nil)))
}

// ctxWithUser builds a context the same way BaseInterceptor + MetadataFromCtx would populate it.
func ctxWithUser(userID string, roles ...sharedmodel.Role) context.Context {
	ctx := context.Background()
	ctx = context.WithValue(ctx, "x-request-id", "test-req-id")
	ctx = context.WithValue(ctx, "x-client-ip", "127.0.0.1")
	ctx = context.WithValue(ctx, "x-user-id", userID)
	if len(roles) > 0 {
		ctx = context.WithValue(ctx, "x-user-roles", roles)
	}
	return ctx
}

func ctxAnon() context.Context {
	ctx := context.Background()
	ctx = context.WithValue(ctx, "x-request-id", "test-req-id")
	ctx = context.WithValue(ctx, "x-client-ip", "127.0.0.1")
	return ctx
}

// invoke runs the interceptor for the given method and request. Returns the gRPC
// status code and whether the downstream handler was actually called.
func invoke(t *testing.T, ctx context.Context, method string, req any) (codes.Code, bool) {
	t.Helper()
	called := false
	_, err := newInterceptor().Unary(ctx, req, &grpc.UnaryServerInfo{FullMethod: method}, func(ctx context.Context, req any) (any, error) {
		called = true
		return nil, nil
	})
	if err != nil {
		st, ok := status.FromError(err)
		require.True(t, ok, "error must be a gRPC status error")
		return st.Code(), called
	}
	return codes.OK, called
}

// --- Unknown method ---

func TestAuth_UnknownMethod_PermissionDenied(t *testing.T) {
	code, called := invoke(t, ctxWithUser(callerID, sharedmodel.RoleAdmin), "/unknown.Service/Method", nil)

	assert.Equal(t, codes.PermissionDenied, code)
	assert.False(t, called)
}

// --- Public methods ---

func TestAuth_PublicMethod_AllowsAnonymous(t *testing.T) {
	code, called := invoke(t, ctxAnon(), usersvc.UserService_Register_FullMethodName, nil)

	assert.Equal(t, codes.OK, code)
	assert.True(t, called)
}

func TestAuth_PublicMethod_AllowsAuthenticated(t *testing.T) {
	code, called := invoke(t, ctxWithUser(callerID, sharedmodel.RoleUser), usersvc.UserService_SignIn_FullMethodName, nil)

	assert.Equal(t, codes.OK, code)
	assert.True(t, called)
}

// --- Any authenticated user (empty policy) ---

func TestAuth_AuthenticatedOnly_AllowsUserRole(t *testing.T) {
	code, called := invoke(t, ctxWithUser(callerID, sharedmodel.RoleUser), usersvc.UserService_SendActivationCode_FullMethodName, nil)

	assert.Equal(t, codes.OK, code)
	assert.True(t, called)
}

func TestAuth_AuthenticatedOnly_AllowsNoRoles(t *testing.T) {
	// Authenticated but no roles assigned — still allowed for empty-policy methods.
	code, called := invoke(t, ctxWithUser(callerID), usersvc.UserService_CheckActivationCode_FullMethodName, nil)

	assert.Equal(t, codes.OK, code)
	assert.True(t, called)
}

func TestAuth_AuthenticatedOnly_RejectsAnonymous(t *testing.T) {
	code, called := invoke(t, ctxAnon(), usersvc.UserService_SendActivationCode_FullMethodName, nil)

	assert.Equal(t, codes.Unauthenticated, code)
	assert.False(t, called)
}

// --- Privileged roles only ---

func TestAuth_PrivilegedOnly_AllowsAdmin(t *testing.T) {
	code, called := invoke(t, ctxWithUser(callerID, sharedmodel.RoleAdmin), usersvc.UserService_CreateUser_FullMethodName, nil)

	assert.Equal(t, codes.OK, code)
	assert.True(t, called)
}

func TestAuth_PrivilegedOnly_AllowsUserManager(t *testing.T) {
	code, called := invoke(t, ctxWithUser(callerID, sharedmodel.RoleUserManager), usersvc.UserService_ListUsers_FullMethodName, nil)

	assert.Equal(t, codes.OK, code)
	assert.True(t, called)
}

func TestAuth_PrivilegedOnly_RejectsUserRole(t *testing.T) {
	code, called := invoke(t, ctxWithUser(callerID, sharedmodel.RoleUser), usersvc.UserService_CreateUser_FullMethodName, nil)

	assert.Equal(t, codes.PermissionDenied, code)
	assert.False(t, called)
}

func TestAuth_PrivilegedOnly_RejectsEmptyRoles(t *testing.T) {
	code, called := invoke(t, ctxWithUser(callerID), usersvc.UserService_CheckDocument_FullMethodName, nil)

	assert.Equal(t, codes.PermissionDenied, code)
	assert.False(t, called)
}

func TestAuth_PrivilegedOnly_RejectsAnonymous(t *testing.T) {
	code, called := invoke(t, ctxAnon(), usersvc.UserService_ListUsers_FullMethodName, nil)

	assert.Equal(t, codes.Unauthenticated, code)
	assert.False(t, called)
}

// --- Privileged roles OR resource owner (GetId) ---

func TestAuth_RolesOrOwner_PrivilegedRoleGrantsAccess(t *testing.T) {
	req := fakeIDReq{id: otherID} // not the caller — role should still win
	code, called := invoke(t, ctxWithUser(callerID, sharedmodel.RoleAdmin), usersvc.UserService_GetUser_FullMethodName, req)

	assert.Equal(t, codes.OK, code)
	assert.True(t, called)
}

func TestAuth_RolesOrOwner_OwnerGrantsAccess(t *testing.T) {
	req := fakeIDReq{id: callerID}
	code, called := invoke(t, ctxWithUser(callerID, sharedmodel.RoleUser), usersvc.UserService_GetUser_FullMethodName, req)

	assert.Equal(t, codes.OK, code)
	assert.True(t, called)
}

func TestAuth_RolesOrOwner_RejectsNonOwnerWithoutPrivilege(t *testing.T) {
	req := fakeIDReq{id: otherID} // different user, no privileged role
	code, called := invoke(t, ctxWithUser(callerID, sharedmodel.RoleUser), usersvc.UserService_UpdateUser_FullMethodName, req)

	assert.Equal(t, codes.PermissionDenied, code)
	assert.False(t, called)
}

func TestAuth_RolesOrOwner_RejectsAnonymous(t *testing.T) {
	req := fakeIDReq{id: callerID}
	code, called := invoke(t, ctxAnon(), usersvc.UserService_DeleteUser_FullMethodName, req)

	assert.Equal(t, codes.Unauthenticated, code)
	assert.False(t, called)
}

func TestAuth_RolesOrOwner_RejectsWhenReqLacksID(t *testing.T) {
	// Request does not implement idCarrier — ownerExtract returns false, no role → denied.
	code, called := invoke(t, ctxWithUser(callerID, sharedmodel.RoleUser), usersvc.UserService_GetUser_FullMethodName, struct{}{})

	assert.Equal(t, codes.PermissionDenied, code)
	assert.False(t, called)
}

// --- Privileged roles OR resource owner (GetUserId) ---

func TestAuth_RolesOrOwner_UserIDCarrier_OwnerGrantsAccess(t *testing.T) {
	req := fakeUserIDReq{userID: callerID}
	code, called := invoke(t, ctxWithUser(callerID, sharedmodel.RoleUser), usersvc.UserService_GetProcessedDocumentsForUser_FullMethodName, req)

	assert.Equal(t, codes.OK, code)
	assert.True(t, called)
}

func TestAuth_RolesOrOwner_UserIDCarrier_RejectsNonOwner(t *testing.T) {
	req := fakeUserIDReq{userID: otherID}
	code, called := invoke(t, ctxWithUser(callerID, sharedmodel.RoleUser), usersvc.UserService_GetProcessedDocumentsForUser_FullMethodName, req)

	assert.Equal(t, codes.PermissionDenied, code)
	assert.False(t, called)
}

func TestAuth_RolesOrOwner_UserIDCarrier_PrivilegedRoleGrantsAccess(t *testing.T) {
	req := fakeUserIDReq{userID: otherID} // not owner, but admin role wins
	code, called := invoke(t, ctxWithUser(callerID, sharedmodel.RoleAdmin), usersvc.UserService_GetProcessedDocumentsForUser_FullMethodName, req)

	assert.Equal(t, codes.OK, code)
	assert.True(t, called)
}
