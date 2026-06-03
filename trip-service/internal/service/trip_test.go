package service_test

import (
	"context"
	"errors"
	"io"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"carsharing/trip-service/internal/model"
	"carsharing/trip-service/internal/validation"
)

// ── StartTrip ─────────────────────────────────────────────────────────────────

func TestTripService_StartTrip_Success(t *testing.T) {
	d := newDeps(t)
	svc := newService(t, d)
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

	id, err := svc.StartTrip(ctx, testBookingID)

	require.NoError(t, err)
	assert.NotEmpty(t, id)
}

func TestTripService_StartTrip_InvalidID(t *testing.T) {
	d := newDeps(t)
	svc := newService(t, d)
	ctx := ctxOwner(testUserID)

	_, err := svc.StartTrip(ctx, "not-a-uuid")

	assert.Error(t, err)
}

func TestTripService_StartTrip_BookingNotFound(t *testing.T) {
	d := newDeps(t)
	svc := newService(t, d)
	ctx := ctxOwner(testUserID)

	d.booking.EXPECT().GetBooking(mock.Anything, testBookingID).Return(model.Booking{}, model.ErrNotFound)

	_, err := svc.StartTrip(ctx, testBookingID)

	assert.ErrorIs(t, err, model.ErrNotFound)
}

func TestTripService_StartTrip_InsufficientPermissions(t *testing.T) {
	d := newDeps(t)
	svc := newService(t, d)
	ctx := ctxOwner(testOtherID)
	booking := sampleBooking()

	d.booking.EXPECT().GetBooking(mock.Anything, testBookingID).Return(booking, nil)

	_, err := svc.StartTrip(ctx, testBookingID)

	assert.ErrorIs(t, err, model.ErrInsufficientPermissions)
}

func TestTripService_StartTrip_AdminCanStartForAnyUser(t *testing.T) {
	d := newDeps(t)
	svc := newService(t, d)
	ctx := ctxAdmin()
	booking := sampleBooking()
	tel := sampleTelemetry()
	created := model.Trip{ID: testTripID, UserID: testUserID, Status: model.TripStatusActive}

	d.booking.EXPECT().GetBooking(mock.Anything, testBookingID).Return(booking, nil)
	d.telematics.EXPECT().GetLatestTelemetry(mock.Anything, testCarID).Return(tel, nil)
	d.tripRepo.EXPECT().Create(mock.Anything, mock.Anything).Return(created, nil)
	d.statusRepo.EXPECT().Create(mock.Anything, mock.Anything).Return(model.TripStatusReading{}, nil)
	d.publisher.EXPECT().PublishTripStarted(mock.Anything, created).Return(nil)

	id, err := svc.StartTrip(ctx, testBookingID)

	require.NoError(t, err)
	assert.NotEmpty(t, id)
}

func TestTripService_StartTrip_BookingNotCreated(t *testing.T) {
	d := newDeps(t)
	svc := newService(t, d)
	ctx := ctxOwner(testUserID)
	booking := sampleBooking()
	booking.Status = "completed"

	d.booking.EXPECT().GetBooking(mock.Anything, testBookingID).Return(booking, nil)

	_, err := svc.StartTrip(ctx, testBookingID)

	assert.ErrorIs(t, err, model.ErrBookingNotCreated)
}

func TestTripService_StartTrip_TelemetryFails(t *testing.T) {
	d := newDeps(t)
	svc := newService(t, d)
	ctx := ctxOwner(testUserID)
	booking := sampleBooking()
	infraErr := errors.New("telemetry unavailable")

	d.booking.EXPECT().GetBooking(mock.Anything, testBookingID).Return(booking, nil)
	d.telematics.EXPECT().GetLatestTelemetry(mock.Anything, testCarID).Return(model.CarTelemetry{}, infraErr)

	_, err := svc.StartTrip(ctx, testBookingID)

	assert.ErrorIs(t, err, infraErr)
}

func TestTripService_StartTrip_PublishFailureIsNonFatal(t *testing.T) {
	d := newDeps(t)
	svc := newService(t, d)
	ctx := ctxOwner(testUserID)
	booking := sampleBooking()
	tel := sampleTelemetry()
	created := model.Trip{ID: testTripID, UserID: testUserID, Status: model.TripStatusActive}

	d.booking.EXPECT().GetBooking(mock.Anything, testBookingID).Return(booking, nil)
	d.telematics.EXPECT().GetLatestTelemetry(mock.Anything, testCarID).Return(tel, nil)
	d.tripRepo.EXPECT().Create(mock.Anything, mock.Anything).Return(created, nil)
	d.statusRepo.EXPECT().Create(mock.Anything, mock.Anything).Return(model.TripStatusReading{}, nil)
	d.publisher.EXPECT().PublishTripStarted(mock.Anything, created).Return(errors.New("nats down"))

	id, err := svc.StartTrip(ctx, testBookingID)

	require.NoError(t, err)
	assert.NotEmpty(t, id)
}

// ── GetTrip ───────────────────────────────────────────────────────────────────

func TestTripService_GetTrip_Success(t *testing.T) {
	d := newDeps(t)
	svc := newService(t, d)
	ctx := ctxOwner(testUserID)
	trip := sampleActiveTrip()

	d.tripRepo.EXPECT().GetByID(mock.Anything, testTripID).Return(trip, nil)

	got, err := svc.GetTrip(ctx, testTripID)

	require.NoError(t, err)
	assert.Equal(t, testTripID, got.ID)
}

func TestTripService_GetTrip_InvalidID(t *testing.T) {
	d := newDeps(t)
	svc := newService(t, d)

	_, err := svc.GetTrip(ctxOwner(testUserID), "not-a-uuid")

	assert.Error(t, err)
}

func TestTripService_GetTrip_NotFound(t *testing.T) {
	d := newDeps(t)
	svc := newService(t, d)
	ctx := ctxOwner(testUserID)

	d.tripRepo.EXPECT().GetByID(mock.Anything, testTripID).Return(model.Trip{}, model.ErrNotFound)

	_, err := svc.GetTrip(ctx, testTripID)

	assert.ErrorIs(t, err, model.ErrNotFound)
}

func TestTripService_GetTrip_InsufficientPermissions(t *testing.T) {
	d := newDeps(t)
	svc := newService(t, d)
	ctx := ctxOwner(testOtherID)
	trip := sampleActiveTrip()

	d.tripRepo.EXPECT().GetByID(mock.Anything, testTripID).Return(trip, nil)

	_, err := svc.GetTrip(ctx, testTripID)

	assert.ErrorIs(t, err, model.ErrInsufficientPermissions)
}

func TestTripService_GetTrip_AdminCanAccessAnyTrip(t *testing.T) {
	d := newDeps(t)
	svc := newService(t, d)
	ctx := ctxAdmin()
	trip := sampleActiveTrip()

	d.tripRepo.EXPECT().GetByID(mock.Anything, testTripID).Return(trip, nil)

	got, err := svc.GetTrip(ctx, testTripID)

	require.NoError(t, err)
	assert.Equal(t, testTripID, got.ID)
}

// ── ListTrips ─────────────────────────────────────────────────────────────────

func TestTripService_ListTrips_RegularUserScopedToOwnID(t *testing.T) {
	d := newDeps(t)
	svc := newService(t, d)
	ctx := ctxOwner(testUserID)

	d.tripRepo.EXPECT().List(mock.Anything, mock.MatchedBy(func(f model.TripFilter) bool {
		return f.UserID != nil && *f.UserID == testUserID
	})).Return([]model.Trip{sampleActiveTrip()}, nil)

	trips, err := svc.ListTrips(ctx, validation.TripFilter{})

	require.NoError(t, err)
	assert.Len(t, trips, 1)
}

func TestTripService_ListTrips_AdminPassesFilterAsIs(t *testing.T) {
	d := newDeps(t)
	svc := newService(t, d)
	ctx := ctxAdmin()
	otherUserID := testOtherID

	d.tripRepo.EXPECT().List(mock.Anything, mock.MatchedBy(func(f model.TripFilter) bool {
		return f.UserID != nil && *f.UserID == otherUserID
	})).Return([]model.Trip{}, nil)

	_, err := svc.ListTrips(ctx, validation.TripFilter{UserID: &otherUserID})

	require.NoError(t, err)
}

func TestTripService_ListTrips_BookingManagerPassesFilterAsIs(t *testing.T) {
	d := newDeps(t)
	svc := newService(t, d)
	ctx := ctxManager()
	otherUserID := testOtherID

	d.tripRepo.EXPECT().List(mock.Anything, mock.MatchedBy(func(f model.TripFilter) bool {
		return f.UserID != nil && *f.UserID == otherUserID
	})).Return([]model.Trip{}, nil)

	_, err := svc.ListTrips(ctx, validation.TripFilter{UserID: &otherUserID})

	require.NoError(t, err)
}

// ── EndTrip ───────────────────────────────────────────────────────────────────

func TestTripService_EndTrip_Success(t *testing.T) {
	d := newDeps(t)
	svc := newService(t, d)
	ctx := ctxOwner(testUserID)
	trip := sampleActiveTrip()
	tel := sampleTelemetry()
	booking := sampleBooking()
	updated := trip
	updated.Status = model.TripStatusCompleted

	d.tripRepo.EXPECT().GetByID(mock.Anything, testTripID).Return(trip, nil)
	d.telematics.EXPECT().GetLatestTelemetry(mock.Anything, testCarID).Return(tel, nil)
	d.booking.EXPECT().GetBooking(mock.Anything, testBookingID).Return(booking, nil)
	d.zonePricing.EXPECT().GetZonePricing(mock.Anything, tel.Location.Latitude, tel.Location.Longitude).Return(int32(0), nil)
	d.tripRepo.EXPECT().Update(mock.Anything, testTripID, mock.MatchedBy(func(u model.TripUpdate) bool {
		return u.Status != nil && *u.Status == model.TripStatusCompleted
	})).Return(updated, nil)
	d.statusRepo.EXPECT().Create(mock.Anything, mock.Anything).Return(model.TripStatusReading{}, nil)
	d.summaryRepo.EXPECT().Create(mock.Anything, mock.Anything).Return(model.TripSummary{}, nil)
	d.publisher.EXPECT().PublishTripEnded(mock.Anything, updated).Return(nil)

	err := svc.EndTrip(ctx, testTripID)

	require.NoError(t, err)
}

func TestTripService_EndTrip_ZoneInNoDropZone(t *testing.T) {
	d := newDeps(t)
	svc := newService(t, d)
	ctx := ctxOwner(testUserID)
	trip := sampleActiveTrip()
	tel := sampleTelemetry()
	booking := sampleBooking()

	d.tripRepo.EXPECT().GetByID(mock.Anything, testTripID).Return(trip, nil)
	d.telematics.EXPECT().GetLatestTelemetry(mock.Anything, testCarID).Return(tel, nil)
	d.booking.EXPECT().GetBooking(mock.Anything, testBookingID).Return(booking, nil)
	d.zonePricing.EXPECT().GetZonePricing(mock.Anything, tel.Location.Latitude, tel.Location.Longitude).Return(int32(0), model.ErrLocationInNoDropZone)

	err := svc.EndTrip(ctx, testTripID)

	assert.ErrorIs(t, err, model.ErrLocationInNoDropZone)
}

func TestTripService_EndTrip_ZonePricingError(t *testing.T) {
	d := newDeps(t)
	svc := newService(t, d)
	ctx := ctxOwner(testUserID)
	trip := sampleActiveTrip()
	tel := sampleTelemetry()
	booking := sampleBooking()
	infraErr := errors.New("car service unavailable")

	d.tripRepo.EXPECT().GetByID(mock.Anything, testTripID).Return(trip, nil)
	d.telematics.EXPECT().GetLatestTelemetry(mock.Anything, testCarID).Return(tel, nil)
	d.booking.EXPECT().GetBooking(mock.Anything, testBookingID).Return(booking, nil)
	d.zonePricing.EXPECT().GetZonePricing(mock.Anything, tel.Location.Latitude, tel.Location.Longitude).Return(int32(0), infraErr)

	err := svc.EndTrip(ctx, testTripID)

	assert.ErrorIs(t, err, infraErr)
}

func TestTripService_EndTrip_InvalidID(t *testing.T) {
	d := newDeps(t)
	svc := newService(t, d)

	err := svc.EndTrip(ctxOwner(testUserID), "not-a-uuid")

	assert.Error(t, err)
}

func TestTripService_EndTrip_NotFound(t *testing.T) {
	d := newDeps(t)
	svc := newService(t, d)
	ctx := ctxOwner(testUserID)

	d.tripRepo.EXPECT().GetByID(mock.Anything, testTripID).Return(model.Trip{}, model.ErrNotFound)

	err := svc.EndTrip(ctx, testTripID)

	assert.ErrorIs(t, err, model.ErrNotFound)
}

func TestTripService_EndTrip_InsufficientPermissions(t *testing.T) {
	d := newDeps(t)
	svc := newService(t, d)
	ctx := ctxOwner(testOtherID)
	trip := sampleActiveTrip()

	d.tripRepo.EXPECT().GetByID(mock.Anything, testTripID).Return(trip, nil)

	err := svc.EndTrip(ctx, testTripID)

	assert.ErrorIs(t, err, model.ErrInsufficientPermissions)
}

func TestTripService_EndTrip_InvalidTransition(t *testing.T) {
	d := newDeps(t)
	svc := newService(t, d)
	ctx := ctxOwner(testUserID)
	trip := sampleActiveTrip()
	trip.Status = model.TripStatusCompleted

	d.tripRepo.EXPECT().GetByID(mock.Anything, testTripID).Return(trip, nil)

	err := svc.EndTrip(ctx, testTripID)

	assert.ErrorIs(t, err, model.ErrInvalidTripStatusTransition)
}

func TestTripService_EndTrip_TelemetryFails(t *testing.T) {
	d := newDeps(t)
	svc := newService(t, d)
	ctx := ctxOwner(testUserID)
	trip := sampleActiveTrip()
	infraErr := errors.New("telemetry down")

	d.tripRepo.EXPECT().GetByID(mock.Anything, testTripID).Return(trip, nil)
	d.telematics.EXPECT().GetLatestTelemetry(mock.Anything, testCarID).Return(model.CarTelemetry{}, infraErr)

	err := svc.EndTrip(ctx, testTripID)

	assert.ErrorIs(t, err, infraErr)
}

func TestTripService_EndTrip_PublishFailureIsNonFatal(t *testing.T) {
	d := newDeps(t)
	svc := newService(t, d)
	ctx := ctxOwner(testUserID)
	trip := sampleActiveTrip()
	tel := sampleTelemetry()
	booking := sampleBooking()
	updated := trip
	updated.Status = model.TripStatusCompleted

	d.tripRepo.EXPECT().GetByID(mock.Anything, testTripID).Return(trip, nil)
	d.telematics.EXPECT().GetLatestTelemetry(mock.Anything, testCarID).Return(tel, nil)
	d.booking.EXPECT().GetBooking(mock.Anything, testBookingID).Return(booking, nil)
	d.zonePricing.EXPECT().GetZonePricing(mock.Anything, tel.Location.Latitude, tel.Location.Longitude).Return(int32(0), nil)
	d.tripRepo.EXPECT().Update(mock.Anything, testTripID, mock.Anything).Return(updated, nil)
	d.statusRepo.EXPECT().Create(mock.Anything, mock.Anything).Return(model.TripStatusReading{}, nil)
	d.summaryRepo.EXPECT().Create(mock.Anything, mock.Anything).Return(model.TripSummary{}, nil)
	d.publisher.EXPECT().PublishTripEnded(mock.Anything, updated).Return(errors.New("nats down"))

	err := svc.EndTrip(ctx, testTripID)

	require.NoError(t, err)
}

// ── CancelTrip ────────────────────────────────────────────────────────────────

func TestTripService_CancelTrip_Success(t *testing.T) {
	d := newDeps(t)
	svc := newService(t, d)
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

	err := svc.CancelTrip(ctx, testTripID, &reason)

	require.NoError(t, err)
}

func TestTripService_CancelTrip_InvalidID(t *testing.T) {
	d := newDeps(t)
	svc := newService(t, d)

	err := svc.CancelTrip(ctxOwner(testUserID), "not-a-uuid", nil)

	assert.Error(t, err)
}

func TestTripService_CancelTrip_NotFound(t *testing.T) {
	d := newDeps(t)
	svc := newService(t, d)
	ctx := ctxOwner(testUserID)

	d.tripRepo.EXPECT().GetByID(mock.Anything, testTripID).Return(model.Trip{}, model.ErrNotFound)

	err := svc.CancelTrip(ctx, testTripID, nil)

	assert.ErrorIs(t, err, model.ErrNotFound)
}

func TestTripService_CancelTrip_InsufficientPermissions(t *testing.T) {
	d := newDeps(t)
	svc := newService(t, d)
	ctx := ctxOwner(testOtherID)
	trip := sampleActiveTrip()

	d.tripRepo.EXPECT().GetByID(mock.Anything, testTripID).Return(trip, nil)

	err := svc.CancelTrip(ctx, testTripID, nil)

	assert.ErrorIs(t, err, model.ErrInsufficientPermissions)
}

func TestTripService_CancelTrip_InvalidTransition(t *testing.T) {
	d := newDeps(t)
	svc := newService(t, d)
	ctx := ctxOwner(testUserID)
	trip := sampleActiveTrip()
	trip.Status = model.TripStatusCompleted

	d.tripRepo.EXPECT().GetByID(mock.Anything, testTripID).Return(trip, nil)

	err := svc.CancelTrip(ctx, testTripID, nil)

	assert.ErrorIs(t, err, model.ErrInvalidTripStatusTransition)
}

func TestTripService_CancelTrip_PublishFailureIsNonFatal(t *testing.T) {
	d := newDeps(t)
	svc := newService(t, d)
	ctx := ctxOwner(testUserID)
	trip := sampleActiveTrip()
	updated := trip
	updated.Status = model.TripStatusCancelled

	d.tripRepo.EXPECT().GetByID(mock.Anything, testTripID).Return(trip, nil)
	d.tripRepo.EXPECT().Update(mock.Anything, testTripID, mock.Anything).Return(updated, nil)
	d.statusRepo.EXPECT().Create(mock.Anything, mock.Anything).Return(model.TripStatusReading{}, nil)
	d.publisher.EXPECT().PublishTripCancelled(mock.Anything, updated).Return(errors.New("nats down"))

	err := svc.CancelTrip(ctx, testTripID, nil)

	require.NoError(t, err)
}

// ── GetTripSummary ────────────────────────────────────────────────────────────

func TestTripService_GetTripSummary_Success(t *testing.T) {
	d := newDeps(t)
	svc := newService(t, d)
	ctx := ctxOwner(testUserID)
	trip := sampleActiveTrip()
	trip.Status = model.TripStatusCompleted
	summary := model.TripSummary{TripID: testTripID, BookingID: testBookingID}

	d.tripRepo.EXPECT().GetByID(mock.Anything, testTripID).Return(trip, nil)
	d.summaryRepo.EXPECT().GetByTripID(mock.Anything, testTripID).Return(summary, nil)

	got, err := svc.GetTripSummary(ctx, testTripID)

	require.NoError(t, err)
	assert.Equal(t, testTripID, got.TripID)
}

func TestTripService_GetTripSummary_InvalidID(t *testing.T) {
	d := newDeps(t)
	svc := newService(t, d)

	_, err := svc.GetTripSummary(ctxOwner(testUserID), "not-a-uuid")

	assert.Error(t, err)
}

func TestTripService_GetTripSummary_TripNotFound(t *testing.T) {
	d := newDeps(t)
	svc := newService(t, d)
	ctx := ctxOwner(testUserID)

	d.tripRepo.EXPECT().GetByID(mock.Anything, testTripID).Return(model.Trip{}, model.ErrNotFound)

	_, err := svc.GetTripSummary(ctx, testTripID)

	assert.ErrorIs(t, err, model.ErrNotFound)
}

func TestTripService_GetTripSummary_InsufficientPermissions(t *testing.T) {
	d := newDeps(t)
	svc := newService(t, d)
	ctx := ctxOwner(testOtherID)
	trip := sampleActiveTrip()

	d.tripRepo.EXPECT().GetByID(mock.Anything, testTripID).Return(trip, nil)

	_, err := svc.GetTripSummary(ctx, testTripID)

	assert.ErrorIs(t, err, model.ErrInsufficientPermissions)
}

// ── GetTripStatusHistory ──────────────────────────────────────────────────────

func TestTripService_GetTripStatusHistory_Success(t *testing.T) {
	d := newDeps(t)
	svc := newService(t, d)
	ctx := ctxOwner(testUserID)
	trip := sampleActiveTrip()
	readings := []model.TripStatusReading{{TripID: testTripID}}

	d.tripRepo.EXPECT().GetByID(mock.Anything, testTripID).Return(trip, nil)
	d.statusRepo.EXPECT().List(mock.Anything, mock.Anything).Return(readings, nil)

	got, err := svc.GetTripStatusHistory(ctx, validation.TripStatusHistoryFilter{TripID: testTripID})

	require.NoError(t, err)
	assert.Len(t, got, 1)
}

func TestTripService_GetTripStatusHistory_InvalidID(t *testing.T) {
	d := newDeps(t)
	svc := newService(t, d)

	_, err := svc.GetTripStatusHistory(ctxOwner(testUserID), validation.TripStatusHistoryFilter{TripID: "not-a-uuid"})

	assert.Error(t, err)
}

func TestTripService_GetTripStatusHistory_TripNotFound(t *testing.T) {
	d := newDeps(t)
	svc := newService(t, d)
	ctx := ctxOwner(testUserID)

	d.tripRepo.EXPECT().GetByID(mock.Anything, testTripID).Return(model.Trip{}, model.ErrNotFound)

	_, err := svc.GetTripStatusHistory(ctx, validation.TripStatusHistoryFilter{TripID: testTripID})

	assert.ErrorIs(t, err, model.ErrNotFound)
}

func TestTripService_GetTripStatusHistory_InsufficientPermissions(t *testing.T) {
	d := newDeps(t)
	svc := newService(t, d)
	ctx := ctxOwner(testOtherID)
	trip := sampleActiveTrip()

	d.tripRepo.EXPECT().GetByID(mock.Anything, testTripID).Return(trip, nil)

	_, err := svc.GetTripStatusHistory(ctx, validation.TripStatusHistoryFilter{TripID: testTripID})

	assert.ErrorIs(t, err, model.ErrInsufficientPermissions)
}

// ── StreamTripLiveFeed ────────────────────────────────────────────────────────

func TestTripService_StreamTripLiveFeed_Success(t *testing.T) {
	d := newDeps(t)
	svc := newService(t, d)
	ctx := ctxOwner(testUserID)
	trip := sampleActiveTrip()
	booking := sampleBooking()

	d.tripRepo.EXPECT().GetByID(mock.Anything, testTripID).Return(trip, nil).Times(2)
	d.booking.EXPECT().GetBooking(mock.Anything, testBookingID).Return(booking, nil)

	var received []model.TripLiveFeed
	d.telematics.EXPECT().StreamTelemetry(mock.Anything, testCarID, mock.Anything).
		Run(func(_ context.Context, _ string, fn func(model.CarTelemetry) error) {
			tel := model.CarTelemetry{
				MileageKM:  1010,
				RecordedAt: trip.StartedAt.Add(5 * time.Minute),
			}
			_ = fn(tel)
		}).Return(nil)

	err := svc.StreamTripLiveFeed(ctx, testTripID, func(feed model.TripLiveFeed) error {
		received = append(received, feed)
		return nil
	})

	require.NoError(t, err)
	require.Len(t, received, 1)
	assert.EqualValues(t, 10, received[0].DistanceTraveledKM)
}

func TestTripService_StreamTripLiveFeed_InvalidID(t *testing.T) {
	d := newDeps(t)
	svc := newService(t, d)

	err := svc.StreamTripLiveFeed(ctxOwner(testUserID), "not-a-uuid", func(model.TripLiveFeed) error { return nil })

	assert.Error(t, err)
}

func TestTripService_StreamTripLiveFeed_TripNotFound(t *testing.T) {
	d := newDeps(t)
	svc := newService(t, d)
	ctx := ctxOwner(testUserID)

	d.tripRepo.EXPECT().GetByID(mock.Anything, testTripID).Return(model.Trip{}, model.ErrNotFound)

	err := svc.StreamTripLiveFeed(ctx, testTripID, func(model.TripLiveFeed) error { return nil })

	assert.ErrorIs(t, err, model.ErrNotFound)
}

func TestTripService_StreamTripLiveFeed_TripNotActive(t *testing.T) {
	d := newDeps(t)
	svc := newService(t, d)
	ctx := ctxOwner(testUserID)
	trip := sampleActiveTrip()
	trip.Status = model.TripStatusCompleted

	d.tripRepo.EXPECT().GetByID(mock.Anything, testTripID).Return(trip, nil)

	err := svc.StreamTripLiveFeed(ctx, testTripID, func(model.TripLiveFeed) error { return nil })

	assert.ErrorIs(t, err, model.ErrTripNotActive)
}

func TestTripService_StreamTripLiveFeed_InsufficientPermissions(t *testing.T) {
	d := newDeps(t)
	svc := newService(t, d)
	ctx := ctxOwner(testOtherID)
	trip := sampleActiveTrip()

	d.tripRepo.EXPECT().GetByID(mock.Anything, testTripID).Return(trip, nil)

	err := svc.StreamTripLiveFeed(ctx, testTripID, func(model.TripLiveFeed) error { return nil })

	assert.ErrorIs(t, err, model.ErrInsufficientPermissions)
}

func TestTripService_StreamTripLiveFeed_EOFFromTelemetryIsClean(t *testing.T) {
	d := newDeps(t)
	svc := newService(t, d)
	ctx := ctxOwner(testUserID)
	trip := sampleActiveTrip()
	booking := sampleBooking()

	d.tripRepo.EXPECT().GetByID(mock.Anything, testTripID).Return(trip, nil)
	d.booking.EXPECT().GetBooking(mock.Anything, testBookingID).Return(booking, nil)
	d.telematics.EXPECT().StreamTelemetry(mock.Anything, testCarID, mock.Anything).Return(io.EOF)

	err := svc.StreamTripLiveFeed(ctx, testTripID, func(model.TripLiveFeed) error { return nil })

	assert.ErrorIs(t, err, io.EOF)
}

func TestTripService_StreamTripLiveFeed_TripCompletedMidStream(t *testing.T) {
	d := newDeps(t)
	svc := newService(t, d)
	ctx := ctxOwner(testUserID)
	trip := sampleActiveTrip()
	completed := trip
	completed.Status = model.TripStatusCompleted
	booking := sampleBooking()

	d.tripRepo.EXPECT().GetByID(mock.Anything, testTripID).Return(trip, nil).Once()
	d.tripRepo.EXPECT().GetByID(mock.Anything, testTripID).Return(completed, nil).Once()
	d.booking.EXPECT().GetBooking(mock.Anything, testBookingID).Return(booking, nil)
	d.telematics.EXPECT().StreamTelemetry(mock.Anything, testCarID, mock.Anything).
		RunAndReturn(func(_ context.Context, _ string, fn func(model.CarTelemetry) error) error {
			tel := model.CarTelemetry{MileageKM: 1010, RecordedAt: trip.StartedAt.Add(5 * time.Minute)}
			return fn(tel)
		})

	err := svc.StreamTripLiveFeed(ctx, testTripID, func(model.TripLiveFeed) error { return nil })

	assert.ErrorIs(t, err, io.EOF)
}

// ── EndTrip (conflict) ────────────────────────────────────────────────────────

func TestTripService_EndTrip_Conflict(t *testing.T) {
	d := newDeps(t)
	svc := newService(t, d)
	ctx := ctxOwner(testUserID)
	trip := sampleActiveTrip()
	tel := sampleTelemetry()
	booking := sampleBooking()

	d.tripRepo.EXPECT().GetByID(mock.Anything, testTripID).Return(trip, nil)
	d.telematics.EXPECT().GetLatestTelemetry(mock.Anything, testCarID).Return(tel, nil)
	d.booking.EXPECT().GetBooking(mock.Anything, testBookingID).Return(booking, nil)
	d.zonePricing.EXPECT().GetZonePricing(mock.Anything, tel.Location.Latitude, tel.Location.Longitude).Return(int32(0), nil)
	d.tripRepo.EXPECT().Update(mock.Anything, testTripID, mock.Anything).Return(model.Trip{}, model.ErrConflict)

	err := svc.EndTrip(ctx, testTripID)

	assert.ErrorIs(t, err, model.ErrConflict)
}

// ── CancelTrip (conflict) ─────────────────────────────────────────────────────

func TestTripService_CancelTrip_Conflict(t *testing.T) {
	d := newDeps(t)
	svc := newService(t, d)
	ctx := ctxOwner(testUserID)
	trip := sampleActiveTrip()

	d.tripRepo.EXPECT().GetByID(mock.Anything, testTripID).Return(trip, nil)
	d.tripRepo.EXPECT().Update(mock.Anything, testTripID, mock.Anything).Return(model.Trip{}, model.ErrConflict)

	err := svc.CancelTrip(ctx, testTripID, nil)

	assert.ErrorIs(t, err, model.ErrConflict)
}

// ── GetTripSummary (not completed) ────────────────────────────────────────────

func TestTripService_GetTripSummary_TripNotCompleted(t *testing.T) {
	d := newDeps(t)
	svc := newService(t, d)
	ctx := ctxOwner(testUserID)
	trip := sampleActiveTrip()

	d.tripRepo.EXPECT().GetByID(mock.Anything, testTripID).Return(trip, nil)

	_, err := svc.GetTripSummary(ctx, testTripID)

	assert.ErrorIs(t, err, model.ErrTripNotCompleted)
}
