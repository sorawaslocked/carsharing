package service_test

import (
	"context"
	"errors"
	"io"
	"io/ioutil"
	"log/slog"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	sharedmodel "carsharing/shared/model"
	"carsharing/trip-service/internal/model"
	"carsharing/trip-service/internal/service"
	svcmocks "carsharing/trip-service/internal/service/mocks"
)

// ── helpers ──────────────────────────────────────────────────────────────────

const (
	testUserID    = "user-111"
	testBookingID = "booking-222"
	testTripID    = "trip-333"
	testCarID     = "car-444"
)

type deps struct {
	tripRepo    *svcmocks.MockTripRepository
	summaryRepo *svcmocks.MockTripSummaryRepository
	statusRepo  *svcmocks.MockTripStatusReadingRepository
	booking     *svcmocks.MockBookingClient
	telematics  *svcmocks.MockTelematicsClient
	publisher   *svcmocks.MockEventPublisher
	svc         *service.TripService
}

func newDeps(t *testing.T) deps {
	t.Helper()
	d := deps{
		tripRepo:    svcmocks.NewMockTripRepository(t),
		summaryRepo: svcmocks.NewMockTripSummaryRepository(t),
		statusRepo:  svcmocks.NewMockTripStatusReadingRepository(t),
		booking:     svcmocks.NewMockBookingClient(t),
		telematics:  svcmocks.NewMockTelematicsClient(t),
		publisher:   svcmocks.NewMockEventPublisher(t),
	}
	d.svc = service.NewTripService(
		slog.New(slog.NewTextHandler(ioutil.Discard, nil)),
		d.tripRepo,
		d.summaryRepo,
		d.statusRepo,
		d.booking,
		d.telematics,
		d.publisher,
	)
	return d
}

func ctxOwner(userID string) context.Context {
	return testCtx("req-test", "127.0.0.1", userID, nil)
}

func ctxAdmin() context.Context {
	return testCtx("req-test", "127.0.0.1", "admin-user", []sharedmodel.Role{sharedmodel.RoleAdmin})
}

func ctxManager() context.Context {
	return testCtx("req-test", "127.0.0.1", "manager-user", []sharedmodel.Role{sharedmodel.RoleBookingManager})
}

func testCtx(requestID, clientIP, userID string, roles []sharedmodel.Role) context.Context {
	ctx := context.Background()
	ctx = context.WithValue(ctx, "x-request-id", requestID)
	ctx = context.WithValue(ctx, "x-client-ip", clientIP)
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
		Status: "reserved",
		PricingSnapshot: model.PricingSnapshot{
			RateTenge: 10,
		},
	}
}

func sampleTelemetry() model.CarTelemetry {
	return model.CarTelemetry{
		CarID:      testCarID,
		Location:   model.Location{Latitude: 51.5, Longitude: -0.1},
		OdometerKM: 1050,
		RecordedAt: time.Now(),
	}
}

// ── StartTrip ─────────────────────────────────────────────────────────────────

func TestTripService_StartTrip_Success(t *testing.T) {
	d := newDeps(t)
	ctx := ctxOwner(testUserID)
	booking := sampleBooking()
	tel := sampleTelemetry()
	created := model.Trip{ID: testTripID, BookingID: testBookingID, UserID: testUserID, Status: model.TripStatusActive}

	d.booking.EXPECT().GetBooking(mock.Anything, testBookingID).Return(booking, nil)
	d.telematics.EXPECT().GetLatestTelemetry(mock.Anything, testCarID).Return(tel, nil)
	d.tripRepo.EXPECT().Create(mock.Anything, mock.MatchedBy(func(tc model.TripCreate) bool {
		return tc.BookingID == testBookingID && tc.UserID == testUserID && tc.Status == model.TripStatusActive
	})).Return(created, nil)
	d.statusRepo.EXPECT().Create(mock.Anything, mock.Anything).Return(model.TripStatusReading{}, nil)
	d.publisher.EXPECT().PublishTripStarted(mock.Anything, created).Return(nil)

	id, err := d.svc.StartTrip(ctx, testBookingID)

	require.NoError(t, err)
	assert.NotEmpty(t, id)
}

func TestTripService_StartTrip_BookingNotFound(t *testing.T) {
	d := newDeps(t)
	ctx := ctxOwner(testUserID)

	d.booking.EXPECT().GetBooking(mock.Anything, testBookingID).Return(model.Booking{}, model.ErrNotFound)

	_, err := d.svc.StartTrip(ctx, testBookingID)

	assert.ErrorIs(t, err, model.ErrNotFound)
}

func TestTripService_StartTrip_InsufficientPermissions(t *testing.T) {
	d := newDeps(t)
	ctx := ctxOwner("other-user") // caller is not the booking owner
	booking := sampleBooking()    // booking belongs to testUserID

	d.booking.EXPECT().GetBooking(mock.Anything, testBookingID).Return(booking, nil)

	_, err := d.svc.StartTrip(ctx, testBookingID)

	assert.ErrorIs(t, err, model.ErrInsufficientPermissions)
}

func TestTripService_StartTrip_AdminCanStartForAnyUser(t *testing.T) {
	d := newDeps(t)
	ctx := ctxAdmin()
	booking := sampleBooking()
	tel := sampleTelemetry()
	created := model.Trip{ID: testTripID, UserID: testUserID, Status: model.TripStatusActive}

	d.booking.EXPECT().GetBooking(mock.Anything, testBookingID).Return(booking, nil)
	d.telematics.EXPECT().GetLatestTelemetry(mock.Anything, testCarID).Return(tel, nil)
	d.tripRepo.EXPECT().Create(mock.Anything, mock.Anything).Return(created, nil)
	d.statusRepo.EXPECT().Create(mock.Anything, mock.Anything).Return(model.TripStatusReading{}, nil)
	d.publisher.EXPECT().PublishTripStarted(mock.Anything, created).Return(nil)

	id, err := d.svc.StartTrip(ctx, testBookingID)

	require.NoError(t, err)
	assert.NotEmpty(t, id)
}

func TestTripService_StartTrip_BookingNotReserved(t *testing.T) {
	d := newDeps(t)
	ctx := ctxOwner(testUserID)
	booking := sampleBooking()
	booking.Status = "completed"

	d.booking.EXPECT().GetBooking(mock.Anything, testBookingID).Return(booking, nil)

	_, err := d.svc.StartTrip(ctx, testBookingID)

	assert.ErrorIs(t, err, model.ErrBookingNotCreated)
}

func TestTripService_StartTrip_TelemetryFails(t *testing.T) {
	d := newDeps(t)
	ctx := ctxOwner(testUserID)
	booking := sampleBooking()
	infraErr := errors.New("telematics unavailable")

	d.booking.EXPECT().GetBooking(mock.Anything, testBookingID).Return(booking, nil)
	d.telematics.EXPECT().GetLatestTelemetry(mock.Anything, testCarID).Return(model.CarTelemetry{}, infraErr)

	_, err := d.svc.StartTrip(ctx, testBookingID)

	assert.ErrorIs(t, err, infraErr)
}

func TestTripService_StartTrip_PublishFailureIsNonFatal(t *testing.T) {
	d := newDeps(t)
	ctx := ctxOwner(testUserID)
	booking := sampleBooking()
	tel := sampleTelemetry()
	created := model.Trip{ID: testTripID, UserID: testUserID, Status: model.TripStatusActive}

	d.booking.EXPECT().GetBooking(mock.Anything, testBookingID).Return(booking, nil)
	d.telematics.EXPECT().GetLatestTelemetry(mock.Anything, testCarID).Return(tel, nil)
	d.tripRepo.EXPECT().Create(mock.Anything, mock.Anything).Return(created, nil)
	d.statusRepo.EXPECT().Create(mock.Anything, mock.Anything).Return(model.TripStatusReading{}, nil)
	d.publisher.EXPECT().PublishTripStarted(mock.Anything, created).Return(errors.New("nats down"))

	id, err := d.svc.StartTrip(ctx, testBookingID)

	require.NoError(t, err)
	assert.NotEmpty(t, id)
}

// ── GetTrip ───────────────────────────────────────────────────────────────────

func TestTripService_GetTrip_Success(t *testing.T) {
	d := newDeps(t)
	ctx := ctxOwner(testUserID)
	trip := sampleActiveTrip()

	d.tripRepo.EXPECT().GetByID(mock.Anything, testTripID).Return(trip, nil)

	got, err := d.svc.GetTrip(ctx, testTripID)

	require.NoError(t, err)
	assert.Equal(t, testTripID, got.ID)
}

func TestTripService_GetTrip_NotFound(t *testing.T) {
	d := newDeps(t)
	ctx := ctxOwner(testUserID)

	d.tripRepo.EXPECT().GetByID(mock.Anything, testTripID).Return(model.Trip{}, model.ErrNotFound)

	_, err := d.svc.GetTrip(ctx, testTripID)

	assert.ErrorIs(t, err, model.ErrNotFound)
}

func TestTripService_GetTrip_InsufficientPermissions(t *testing.T) {
	d := newDeps(t)
	ctx := ctxOwner("other-user")
	trip := sampleActiveTrip() // owned by testUserID

	d.tripRepo.EXPECT().GetByID(mock.Anything, testTripID).Return(trip, nil)

	_, err := d.svc.GetTrip(ctx, testTripID)

	assert.ErrorIs(t, err, model.ErrInsufficientPermissions)
}

func TestTripService_GetTrip_AdminCanAccessAnyTrip(t *testing.T) {
	d := newDeps(t)
	ctx := ctxAdmin()
	trip := sampleActiveTrip()

	d.tripRepo.EXPECT().GetByID(mock.Anything, testTripID).Return(trip, nil)

	got, err := d.svc.GetTrip(ctx, testTripID)

	require.NoError(t, err)
	assert.Equal(t, testTripID, got.ID)
}

// ── ListTrips ─────────────────────────────────────────────────────────────────

func TestTripService_ListTrips_RegularUserScopedToOwnID(t *testing.T) {
	d := newDeps(t)
	ctx := ctxOwner(testUserID)

	d.tripRepo.EXPECT().List(mock.Anything, mock.MatchedBy(func(f model.TripFilter) bool {
		return f.UserID != nil && *f.UserID == testUserID
	})).Return([]model.Trip{sampleActiveTrip()}, nil)

	trips, err := d.svc.ListTrips(ctx, model.TripFilter{}) // no UserID in filter

	require.NoError(t, err)
	assert.Len(t, trips, 1)
}

func TestTripService_ListTrips_AdminPassesFilterAsIs(t *testing.T) {
	d := newDeps(t)
	ctx := ctxAdmin()
	otherUserID := "other-999"
	filter := model.TripFilter{UserID: &otherUserID}

	d.tripRepo.EXPECT().List(mock.Anything, mock.MatchedBy(func(f model.TripFilter) bool {
		return f.UserID != nil && *f.UserID == otherUserID
	})).Return([]model.Trip{}, nil)

	_, err := d.svc.ListTrips(ctx, filter)

	require.NoError(t, err)
}

func TestTripService_ListTrips_BookingManagerPassesFilterAsIs(t *testing.T) {
	d := newDeps(t)
	ctx := ctxManager()
	otherUserID := "other-999"
	filter := model.TripFilter{UserID: &otherUserID}

	d.tripRepo.EXPECT().List(mock.Anything, mock.MatchedBy(func(f model.TripFilter) bool {
		return f.UserID != nil && *f.UserID == otherUserID
	})).Return([]model.Trip{}, nil)

	_, err := d.svc.ListTrips(ctx, filter)

	require.NoError(t, err)
}

// ── EndTrip ───────────────────────────────────────────────────────────────────

func TestTripService_EndTrip_Success(t *testing.T) {
	d := newDeps(t)
	ctx := ctxOwner(testUserID)
	trip := sampleActiveTrip()
	tel := sampleTelemetry()
	booking := sampleBooking()
	updated := trip
	updated.Status = model.TripStatusCompleted

	d.tripRepo.EXPECT().GetByID(mock.Anything, testTripID).Return(trip, nil)
	d.telematics.EXPECT().GetLatestTelemetry(mock.Anything, testCarID).Return(tel, nil)
	d.booking.EXPECT().GetBooking(mock.Anything, testBookingID).Return(booking, nil)
	d.tripRepo.EXPECT().Update(mock.Anything, testTripID, mock.MatchedBy(func(u model.TripUpdate) bool {
		return u.Status != nil && *u.Status == model.TripStatusCompleted
	})).Return(updated, nil)
	d.statusRepo.EXPECT().Create(mock.Anything, mock.Anything).Return(model.TripStatusReading{}, nil)
	d.summaryRepo.EXPECT().Create(mock.Anything, mock.Anything).Return(model.TripSummary{}, nil)
	d.publisher.EXPECT().PublishTripEnded(mock.Anything, updated).Return(nil)

	err := d.svc.EndTrip(ctx, testTripID)

	require.NoError(t, err)
}

func TestTripService_EndTrip_NotFound(t *testing.T) {
	d := newDeps(t)
	ctx := ctxOwner(testUserID)

	d.tripRepo.EXPECT().GetByID(mock.Anything, testTripID).Return(model.Trip{}, model.ErrNotFound)

	err := d.svc.EndTrip(ctx, testTripID)

	assert.ErrorIs(t, err, model.ErrNotFound)
}

func TestTripService_EndTrip_InsufficientPermissions(t *testing.T) {
	d := newDeps(t)
	ctx := ctxOwner("other-user")
	trip := sampleActiveTrip()

	d.tripRepo.EXPECT().GetByID(mock.Anything, testTripID).Return(trip, nil)

	err := d.svc.EndTrip(ctx, testTripID)

	assert.ErrorIs(t, err, model.ErrInsufficientPermissions)
}

func TestTripService_EndTrip_InvalidTransition(t *testing.T) {
	d := newDeps(t)
	ctx := ctxOwner(testUserID)
	trip := sampleActiveTrip()
	trip.Status = model.TripStatusCompleted // already ended

	d.tripRepo.EXPECT().GetByID(mock.Anything, testTripID).Return(trip, nil)

	err := d.svc.EndTrip(ctx, testTripID)

	assert.ErrorIs(t, err, model.ErrInvalidStatusTransition)
}

func TestTripService_EndTrip_TelemetryFails(t *testing.T) {
	d := newDeps(t)
	ctx := ctxOwner(testUserID)
	trip := sampleActiveTrip()
	infraErr := errors.New("telematics down")

	d.tripRepo.EXPECT().GetByID(mock.Anything, testTripID).Return(trip, nil)
	d.telematics.EXPECT().GetLatestTelemetry(mock.Anything, testCarID).Return(model.CarTelemetry{}, infraErr)

	err := d.svc.EndTrip(ctx, testTripID)

	assert.ErrorIs(t, err, infraErr)
}

func TestTripService_EndTrip_PublishFailureIsNonFatal(t *testing.T) {
	d := newDeps(t)
	ctx := ctxOwner(testUserID)
	trip := sampleActiveTrip()
	tel := sampleTelemetry()
	booking := sampleBooking()
	updated := trip
	updated.Status = model.TripStatusCompleted

	d.tripRepo.EXPECT().GetByID(mock.Anything, testTripID).Return(trip, nil)
	d.telematics.EXPECT().GetLatestTelemetry(mock.Anything, testCarID).Return(tel, nil)
	d.booking.EXPECT().GetBooking(mock.Anything, testBookingID).Return(booking, nil)
	d.tripRepo.EXPECT().Update(mock.Anything, testTripID, mock.Anything).Return(updated, nil)
	d.statusRepo.EXPECT().Create(mock.Anything, mock.Anything).Return(model.TripStatusReading{}, nil)
	d.summaryRepo.EXPECT().Create(mock.Anything, mock.Anything).Return(model.TripSummary{}, nil)
	d.publisher.EXPECT().PublishTripEnded(mock.Anything, updated).Return(errors.New("nats down"))

	err := d.svc.EndTrip(ctx, testTripID)

	require.NoError(t, err)
}

// ── CancelTrip ────────────────────────────────────────────────────────────────

func TestTripService_CancelTrip_Success(t *testing.T) {
	d := newDeps(t)
	ctx := ctxOwner(testUserID)
	trip := sampleActiveTrip()
	reason := "changed mind"
	updated := trip
	updated.Status = model.TripStatusCancelled

	d.tripRepo.EXPECT().GetByID(mock.Anything, testTripID).Return(trip, nil)
	d.tripRepo.EXPECT().Update(mock.Anything, testTripID, mock.MatchedBy(func(u model.TripUpdate) bool {
		return u.Status != nil && *u.Status == model.TripStatusCancelled
	})).Return(updated, nil)
	d.statusRepo.EXPECT().Create(mock.Anything, mock.Anything).Return(model.TripStatusReading{}, nil)
	d.publisher.EXPECT().PublishTripCancelled(mock.Anything, updated).Return(nil)

	err := d.svc.CancelTrip(ctx, testTripID, &reason)

	require.NoError(t, err)
}

func TestTripService_CancelTrip_NotFound(t *testing.T) {
	d := newDeps(t)
	ctx := ctxOwner(testUserID)

	d.tripRepo.EXPECT().GetByID(mock.Anything, testTripID).Return(model.Trip{}, model.ErrNotFound)

	err := d.svc.CancelTrip(ctx, testTripID, nil)

	assert.ErrorIs(t, err, model.ErrNotFound)
}

func TestTripService_CancelTrip_InsufficientPermissions(t *testing.T) {
	d := newDeps(t)
	ctx := ctxOwner("other-user")
	trip := sampleActiveTrip()

	d.tripRepo.EXPECT().GetByID(mock.Anything, testTripID).Return(trip, nil)

	err := d.svc.CancelTrip(ctx, testTripID, nil)

	assert.ErrorIs(t, err, model.ErrInsufficientPermissions)
}

func TestTripService_CancelTrip_InvalidTransition(t *testing.T) {
	d := newDeps(t)
	ctx := ctxOwner(testUserID)
	trip := sampleActiveTrip()
	trip.Status = model.TripStatusCompleted

	d.tripRepo.EXPECT().GetByID(mock.Anything, testTripID).Return(trip, nil)

	err := d.svc.CancelTrip(ctx, testTripID, nil)

	assert.ErrorIs(t, err, model.ErrInvalidStatusTransition)
}

func TestTripService_CancelTrip_PublishFailureIsNonFatal(t *testing.T) {
	d := newDeps(t)
	ctx := ctxOwner(testUserID)
	trip := sampleActiveTrip()
	updated := trip
	updated.Status = model.TripStatusCancelled

	d.tripRepo.EXPECT().GetByID(mock.Anything, testTripID).Return(trip, nil)
	d.tripRepo.EXPECT().Update(mock.Anything, testTripID, mock.Anything).Return(updated, nil)
	d.statusRepo.EXPECT().Create(mock.Anything, mock.Anything).Return(model.TripStatusReading{}, nil)
	d.publisher.EXPECT().PublishTripCancelled(mock.Anything, updated).Return(errors.New("nats down"))

	err := d.svc.CancelTrip(ctx, testTripID, nil)

	require.NoError(t, err)
}

// ── GetTripSummary ────────────────────────────────────────────────────────────

func TestTripService_GetTripSummary_Success(t *testing.T) {
	d := newDeps(t)
	ctx := ctxOwner(testUserID)
	trip := sampleActiveTrip()
	summary := model.TripSummary{TripID: testTripID, BookingID: testBookingID}

	d.tripRepo.EXPECT().GetByID(mock.Anything, testTripID).Return(trip, nil)
	d.summaryRepo.EXPECT().GetByTripID(mock.Anything, testTripID).Return(summary, nil)

	got, err := d.svc.GetTripSummary(ctx, testTripID)

	require.NoError(t, err)
	assert.Equal(t, testTripID, got.TripID)
}

func TestTripService_GetTripSummary_TripNotFound(t *testing.T) {
	d := newDeps(t)
	ctx := ctxOwner(testUserID)

	d.tripRepo.EXPECT().GetByID(mock.Anything, testTripID).Return(model.Trip{}, model.ErrNotFound)

	_, err := d.svc.GetTripSummary(ctx, testTripID)

	assert.ErrorIs(t, err, model.ErrNotFound)
}

func TestTripService_GetTripSummary_InsufficientPermissions(t *testing.T) {
	d := newDeps(t)
	ctx := ctxOwner("other-user")
	trip := sampleActiveTrip()

	d.tripRepo.EXPECT().GetByID(mock.Anything, testTripID).Return(trip, nil)

	_, err := d.svc.GetTripSummary(ctx, testTripID)

	assert.ErrorIs(t, err, model.ErrInsufficientPermissions)
}

// ── GetTripStatusHistory ──────────────────────────────────────────────────────

func TestTripService_GetTripStatusHistory_Success(t *testing.T) {
	d := newDeps(t)
	ctx := ctxOwner(testUserID)
	trip := sampleActiveTrip()
	filter := model.TripStatusReadingFilter{TripID: testTripID}
	readings := []model.TripStatusReading{{TripID: testTripID}}

	d.tripRepo.EXPECT().GetByID(mock.Anything, testTripID).Return(trip, nil)
	d.statusRepo.EXPECT().List(mock.Anything, filter).Return(readings, nil)

	got, err := d.svc.GetTripStatusHistory(ctx, filter)

	require.NoError(t, err)
	assert.Len(t, got, 1)
}

func TestTripService_GetTripStatusHistory_TripNotFound(t *testing.T) {
	d := newDeps(t)
	ctx := ctxOwner(testUserID)

	d.tripRepo.EXPECT().GetByID(mock.Anything, testTripID).Return(model.Trip{}, model.ErrNotFound)

	_, err := d.svc.GetTripStatusHistory(ctx, model.TripStatusReadingFilter{TripID: testTripID})

	assert.ErrorIs(t, err, model.ErrNotFound)
}

func TestTripService_GetTripStatusHistory_InsufficientPermissions(t *testing.T) {
	d := newDeps(t)
	ctx := ctxOwner("other-user")
	trip := sampleActiveTrip()

	d.tripRepo.EXPECT().GetByID(mock.Anything, testTripID).Return(trip, nil)

	_, err := d.svc.GetTripStatusHistory(ctx, model.TripStatusReadingFilter{TripID: testTripID})

	assert.ErrorIs(t, err, model.ErrInsufficientPermissions)
}

// ── StreamTripLiveFeed ────────────────────────────────────────────────────────

func TestTripService_StreamTripLiveFeed_Success(t *testing.T) {
	d := newDeps(t)
	ctx := ctxOwner(testUserID)
	trip := sampleActiveTrip()
	booking := sampleBooking()

	// Called twice: initial fetch + in-callback status poll.
	d.tripRepo.EXPECT().GetByID(mock.Anything, testTripID).Return(trip, nil).Times(2)
	d.booking.EXPECT().GetBooking(mock.Anything, testBookingID).Return(booking, nil)

	var received []model.TripLiveFeed
	d.telematics.EXPECT().StreamTelemetry(mock.Anything, testCarID, mock.Anything).
		Run(func(ctx context.Context, carID string, fn func(model.CarTelemetry) error) {
			tel := model.CarTelemetry{
				OdometerKM: 1010,
				RecordedAt: trip.StartedAt.Add(5 * time.Minute),
			}
			_ = fn(tel)
		}).Return(nil)

	err := d.svc.StreamTripLiveFeed(ctx, testTripID, func(feed model.TripLiveFeed) error {
		received = append(received, feed)
		return nil
	})

	require.NoError(t, err)
	require.Len(t, received, 1)
	assert.EqualValues(t, 10, received[0].DistanceTraveledKM) // 1010 − 1000
}

func TestTripService_StreamTripLiveFeed_TripNotFound(t *testing.T) {
	d := newDeps(t)
	ctx := ctxOwner(testUserID)

	d.tripRepo.EXPECT().GetByID(mock.Anything, testTripID).Return(model.Trip{}, model.ErrNotFound)

	err := d.svc.StreamTripLiveFeed(ctx, testTripID, func(model.TripLiveFeed) error { return nil })

	assert.ErrorIs(t, err, model.ErrNotFound)
}

func TestTripService_StreamTripLiveFeed_TripNotActive(t *testing.T) {
	d := newDeps(t)
	ctx := ctxOwner(testUserID)
	trip := sampleActiveTrip()
	trip.Status = model.TripStatusCompleted

	d.tripRepo.EXPECT().GetByID(mock.Anything, testTripID).Return(trip, nil)

	err := d.svc.StreamTripLiveFeed(ctx, testTripID, func(model.TripLiveFeed) error { return nil })

	assert.ErrorIs(t, err, model.ErrTripNotActive)
}

func TestTripService_StreamTripLiveFeed_InsufficientPermissions(t *testing.T) {
	d := newDeps(t)
	ctx := ctxOwner("other-user")
	trip := sampleActiveTrip()

	d.tripRepo.EXPECT().GetByID(mock.Anything, testTripID).Return(trip, nil)

	err := d.svc.StreamTripLiveFeed(ctx, testTripID, func(model.TripLiveFeed) error { return nil })

	assert.ErrorIs(t, err, model.ErrInsufficientPermissions)
}

func TestTripService_StreamTripLiveFeed_EOFFromTelemetryIsClean(t *testing.T) {
	d := newDeps(t)
	ctx := ctxOwner(testUserID)
	trip := sampleActiveTrip()
	booking := sampleBooking()

	d.tripRepo.EXPECT().GetByID(mock.Anything, testTripID).Return(trip, nil)
	d.booking.EXPECT().GetBooking(mock.Anything, testBookingID).Return(booking, nil)
	d.telematics.EXPECT().StreamTelemetry(mock.Anything, testCarID, mock.Anything).Return(io.EOF)

	err := d.svc.StreamTripLiveFeed(ctx, testTripID, func(model.TripLiveFeed) error { return nil })

	assert.ErrorIs(t, err, io.EOF) // service propagates io.EOF; handler converts it to clean close
}
