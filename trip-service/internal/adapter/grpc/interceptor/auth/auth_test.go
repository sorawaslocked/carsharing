package auth_test

import (
	"context"
	"io"
	"log/slog"
	"testing"

	tripsvc "carsharing/protos/gen/service/trip"
	sharedmodel "carsharing/shared/model"
	"carsharing/trip-service/internal/adapter/grpc/interceptor/auth"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	callerID = "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa"
)

func newInterceptor() *auth.Interceptor {
	return auth.NewInterceptor(slog.New(slog.NewTextHandler(io.Discard, nil)))
}

// ctxWithUser builds a context the way BaseInterceptor + MetadataFromCtx populate it.
func ctxWithUser(userID string, roles ...sharedmodel.Role) context.Context {
	ctx := context.Background()
	ctx = context.WithValue(ctx, "x-user-id", userID)
	if len(roles) > 0 {
		ctx = context.WithValue(ctx, "x-user-roles", roles)
	}
	return ctx
}

func ctxAnon() context.Context { return context.Background() }

// invoke runs the Unary interceptor and returns the gRPC code and whether the handler was called.
func invoke(t *testing.T, ctx context.Context, method string) (codes.Code, bool) {
	t.Helper()
	called := false
	_, err := newInterceptor().Unary(ctx, nil, &grpc.UnaryServerInfo{FullMethod: method},
		func(ctx context.Context, req any) (any, error) {
			called = true
			return nil, nil
		},
	)
	if err != nil {
		st, ok := status.FromError(err)
		require.True(t, ok)
		return st.Code(), called
	}
	return codes.OK, called
}

// invokeStream runs the Stream interceptor and returns the gRPC code and whether the handler was called.
func invokeStream(t *testing.T, ctx context.Context, method string) (codes.Code, bool) {
	t.Helper()
	called := false
	ss := &fakeStream{ctx: ctx}
	err := newInterceptor().Stream(nil, ss, &grpc.StreamServerInfo{FullMethod: method},
		func(srv any, stream grpc.ServerStream) error {
			called = true
			return nil
		},
	)
	if err != nil {
		st, ok := status.FromError(err)
		require.True(t, ok)
		return st.Code(), called
	}
	return codes.OK, called
}

type fakeStream struct {
	grpc.ServerStream
	ctx context.Context
}

func (f *fakeStream) Context() context.Context { return f.ctx }

// --- Unknown method ---

func TestAuth_UnknownMethod_PermissionDenied(t *testing.T) {
	code, called := invoke(t, ctxWithUser(callerID, sharedmodel.RoleAdmin), "/unknown.Service/Method")

	assert.Equal(t, codes.PermissionDenied, code)
	assert.False(t, called)
}

// --- Public method (Health) ---

func TestAuth_PublicMethod_AllowsAnonymous(t *testing.T) {
	code, called := invoke(t, ctxAnon(), tripsvc.HealthService_Health_FullMethodName)

	assert.Equal(t, codes.OK, code)
	assert.True(t, called)
}

func TestAuth_PublicMethod_AllowsAuthenticated(t *testing.T) {
	code, called := invoke(t, ctxWithUser(callerID), tripsvc.HealthService_Health_FullMethodName)

	assert.Equal(t, codes.OK, code)
	assert.True(t, called)
}

// --- Any authenticated user ---

func TestAuth_AuthenticatedOnly_AllowsUser(t *testing.T) {
	code, called := invoke(t, ctxWithUser(callerID), tripsvc.TripService_StartTrip_FullMethodName)

	assert.Equal(t, codes.OK, code)
	assert.True(t, called)
}

func TestAuth_AuthenticatedOnly_RejectsAnonymous(t *testing.T) {
	code, called := invoke(t, ctxAnon(), tripsvc.TripService_GetTrip_FullMethodName)

	assert.Equal(t, codes.Unauthenticated, code)
	assert.False(t, called)
}

// --- Manager-restricted (GetTripStatusHistory) ---

func TestAuth_ManagerOnly_AllowsAdmin(t *testing.T) {
	code, called := invoke(t, ctxWithUser(callerID, sharedmodel.RoleAdmin), tripsvc.TripService_GetTripStatusHistory_FullMethodName)

	assert.Equal(t, codes.OK, code)
	assert.True(t, called)
}

func TestAuth_ManagerOnly_AllowsBookingManager(t *testing.T) {
	code, called := invoke(t, ctxWithUser(callerID, sharedmodel.RoleBookingManager), tripsvc.TripService_GetTripStatusHistory_FullMethodName)

	assert.Equal(t, codes.OK, code)
	assert.True(t, called)
}

func TestAuth_ManagerOnly_RejectsRegularUser(t *testing.T) {
	code, called := invoke(t, ctxWithUser(callerID, sharedmodel.RoleUser), tripsvc.TripService_GetTripStatusHistory_FullMethodName)

	assert.Equal(t, codes.PermissionDenied, code)
	assert.False(t, called)
}

func TestAuth_ManagerOnly_RejectsNoRoles(t *testing.T) {
	code, called := invoke(t, ctxWithUser(callerID), tripsvc.TripService_GetTripStatusHistory_FullMethodName)

	assert.Equal(t, codes.PermissionDenied, code)
	assert.False(t, called)
}

func TestAuth_ManagerOnly_RejectsAnonymous(t *testing.T) {
	code, called := invoke(t, ctxAnon(), tripsvc.TripService_GetTripStatusHistory_FullMethodName)

	assert.Equal(t, codes.Unauthenticated, code)
	assert.False(t, called)
}

// --- Stream interceptor ---

func TestAuth_Stream_AuthenticatedOnly_AllowsUser(t *testing.T) {
	code, called := invokeStream(t, ctxWithUser(callerID), tripsvc.TripStreamService_StreamTripLiveFeed_FullMethodName)

	assert.Equal(t, codes.OK, code)
	assert.True(t, called)
}

func TestAuth_Stream_AuthenticatedOnly_RejectsAnonymous(t *testing.T) {
	code, called := invokeStream(t, ctxAnon(), tripsvc.TripStreamService_StreamTripLiveFeed_FullMethodName)

	assert.Equal(t, codes.Unauthenticated, code)
	assert.False(t, called)
}

func TestAuth_Stream_UnknownMethod_PermissionDenied(t *testing.T) {
	code, called := invokeStream(t, ctxWithUser(callerID), "/unknown.Service/Stream")

	assert.Equal(t, codes.PermissionDenied, code)
	assert.False(t, called)
}
