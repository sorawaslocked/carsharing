package auth_test

import (
	"context"
	"io"
	"log/slog"
	"testing"

	"carsharing/car-service/internal/adapter/grpc/interceptor/auth"
	carsvc "carsharing/protos/gen/service/car"
	sharedmodel "carsharing/shared/model"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const callerID = "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa"

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

// invoke runs the interceptor for the given method. Returns the gRPC status code
// and whether the downstream handler was actually called.
func invoke(t *testing.T, ctx context.Context, method string) (codes.Code, bool) {
	t.Helper()
	called := false
	_, err := newInterceptor().Unary(ctx, nil, &grpc.UnaryServerInfo{FullMethod: method}, func(ctx context.Context, req any) (any, error) {
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
	code, called := invoke(t, ctxWithUser(callerID, sharedmodel.RoleAdmin), "/unknown.Service/Method")

	assert.Equal(t, codes.PermissionDenied, code)
	assert.False(t, called)
}

// --- Public methods ---

func TestAuth_PublicMethod_AllowsAnonymous(t *testing.T) {
	code, called := invoke(t, ctxAnon(), carsvc.HealthService_Health_FullMethodName)

	assert.Equal(t, codes.OK, code)
	assert.True(t, called)
}

func TestAuth_PublicMethod_AllowsAuthenticated(t *testing.T) {
	code, called := invoke(t, ctxWithUser(callerID, sharedmodel.RoleUser), carsvc.HealthService_Health_FullMethodName)

	assert.Equal(t, codes.OK, code)
	assert.True(t, called)
}

// --- Any authenticated user (read methods) ---

func TestAuth_AuthenticatedOnly_AllowsUserRole(t *testing.T) {
	code, called := invoke(t, ctxWithUser(callerID, sharedmodel.RoleUser), carsvc.CarService_GetCar_FullMethodName)

	assert.Equal(t, codes.OK, code)
	assert.True(t, called)
}

func TestAuth_AuthenticatedOnly_AllowsNoRoles(t *testing.T) {
	// Authenticated but no roles assigned — still allowed for empty-policy methods.
	code, called := invoke(t, ctxWithUser(callerID), carsvc.CarModelService_ListCarModels_FullMethodName)

	assert.Equal(t, codes.OK, code)
	assert.True(t, called)
}

func TestAuth_AuthenticatedOnly_RejectsAnonymous(t *testing.T) {
	code, called := invoke(t, ctxAnon(), carsvc.CarService_ListCars_FullMethodName)

	assert.Equal(t, codes.Unauthenticated, code)
	assert.False(t, called)
}

// --- Fleet roles only (write methods) ---

func TestAuth_FleetOnly_AllowsAdmin(t *testing.T) {
	code, called := invoke(t, ctxWithUser(callerID, sharedmodel.RoleAdmin), carsvc.CarService_CreateCar_FullMethodName)

	assert.Equal(t, codes.OK, code)
	assert.True(t, called)
}

func TestAuth_FleetOnly_AllowsFleetManager(t *testing.T) {
	code, called := invoke(t, ctxWithUser(callerID, sharedmodel.RoleFleetManager), carsvc.CarModelService_UpdateCarModel_FullMethodName)

	assert.Equal(t, codes.OK, code)
	assert.True(t, called)
}

func TestAuth_FleetOnly_RejectsUserRole(t *testing.T) {
	code, called := invoke(t, ctxWithUser(callerID, sharedmodel.RoleUser), carsvc.CarService_DeleteCar_FullMethodName)

	assert.Equal(t, codes.PermissionDenied, code)
	assert.False(t, called)
}

func TestAuth_FleetOnly_RejectsEmptyRoles(t *testing.T) {
	code, called := invoke(t, ctxWithUser(callerID), carsvc.CarMaintenanceService_CompleteMaintenanceRecord_FullMethodName)

	assert.Equal(t, codes.PermissionDenied, code)
	assert.False(t, called)
}

func TestAuth_FleetOnly_RejectsAnonymous(t *testing.T) {
	code, called := invoke(t, ctxAnon(), carsvc.ZoneService_CreateZone_FullMethodName)

	assert.Equal(t, codes.Unauthenticated, code)
	assert.False(t, called)
}
