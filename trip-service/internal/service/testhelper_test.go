package service_test

import (
	"context"
	"io"
	"log/slog"
	"testing"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/require"

	sharedmodel "carsharing/shared/model"
	"carsharing/trip-service/internal/model"
	"carsharing/trip-service/internal/service"
	mocks "carsharing/trip-service/internal/service/mocks"
	"carsharing/trip-service/internal/validation"
)

const (
	testUserID    = "11111111-1111-4111-8111-111111111111"
	testBookingID = "22222222-2222-4222-8222-222222222222"
	testTripID    = "33333333-3333-4333-8333-333333333333"
	testCarID     = "44444444-4444-4444-8444-444444444444"
	testOtherID   = "55555555-5555-4555-8555-555555555555"
)

type deps struct {
	tripRepo    *mocks.MockTripRepository
	summaryRepo *mocks.MockTripSummaryRepository
	statusRepo  *mocks.MockTripStatusReadingRepository
	booking     *mocks.MockBookingClient
	telematics  *mocks.MockTelematicsClient
	publisher   *mocks.MockEventPublisher
}

func newDeps(t *testing.T) *deps {
	t.Helper()
	return &deps{
		tripRepo:    mocks.NewMockTripRepository(t),
		summaryRepo: mocks.NewMockTripSummaryRepository(t),
		statusRepo:  mocks.NewMockTripStatusReadingRepository(t),
		booking:     mocks.NewMockBookingClient(t),
		telematics:  mocks.NewMockTelematicsClient(t),
		publisher:   mocks.NewMockEventPublisher(t),
	}
}

func newService(t *testing.T, d *deps) *service.TripService {
	t.Helper()
	v := validator.New()
	require.NoError(t, validation.RegisterCustomValidators(v, slog.New(slog.NewTextHandler(io.Discard, nil))))
	return service.NewTripService(
		slog.New(slog.NewTextHandler(io.Discard, nil)),
		v,
		d.tripRepo,
		d.summaryRepo,
		d.statusRepo,
		d.booking,
		d.telematics,
		d.publisher,
	)
}

func ctxOwner(userID string) context.Context {
	return testCtx(userID, nil)
}

func ctxAdmin() context.Context {
	return testCtx("admin-user", []sharedmodel.Role{sharedmodel.RoleAdmin})
}

func ctxManager() context.Context {
	return testCtx("manager-user", []sharedmodel.Role{sharedmodel.RoleBookingManager})
}

func testCtx(userID string, roles []sharedmodel.Role) context.Context {
	ctx := context.Background()
	ctx = context.WithValue(ctx, "x-request-id", "test-req-id")
	ctx = context.WithValue(ctx, "x-client-ip", "127.0.0.1")
	if userID != "" {
		ctx = context.WithValue(ctx, "x-user-id", userID)
	}
	if len(roles) > 0 {
		ctx = context.WithValue(ctx, "x-user-roles", roles)
	}
	return ctx
}

func sampleActiveTrip() model.Trip {
	return model.Trip{
		ID:             testTripID,
		BookingID:      testBookingID,
		UserID:         testUserID,
		CarID:          testCarID,
		Status:         model.TripStatusActive,
		StartedAt:      time.Now().Add(-10 * time.Minute),
		StartMileageKM: 1000,
	}
}

func sampleBooking() model.Booking {
	return model.Booking{
		ID:     testBookingID,
		UserID: testUserID,
		CarID:  testCarID,
		Status: "created",
		PricingSnapshot: model.PricingSnapshot{
			RateTenge: 10,
		},
	}
}

func sampleTelemetry() model.CarTelemetry {
	return model.CarTelemetry{
		CarID:      testCarID,
		Location:   sharedmodel.Location{Latitude: 51.5, Longitude: -0.1},
		MileageKM:  1050,
		RecordedAt: time.Now(),
	}
}
