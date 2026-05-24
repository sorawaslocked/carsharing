package handler_test

import (
	"context"
	"io"
	"io/ioutil"
	"log/slog"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	tripsvc "carsharing/protos/gen/service/trip"

	"carsharing/trip-service/internal/adapter/grpc/handler"
	handlermocks "carsharing/trip-service/internal/adapter/grpc/handler/mocks"
	"carsharing/trip-service/internal/model"
)

// ── helpers ───────────────────────────────────────────────────────────────────

func discardLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(ioutil.Discard, nil))
}

func sampleTrip() model.Trip {
	return model.Trip{
		ID:        "trip-1",
		BookingID: "booking-1",
		UserID:    "user-1",
		CarID:     "car-1",
		Status:    model.TripStatusActive,
		StartedAt: time.Now(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

// fakeServerStream implements grpc.ServerStreamingServer[tripsvc.StreamTripLiveFeedResponse].
type fakeServerStream struct {
	ctx  context.Context
	sent []*tripsvc.StreamTripLiveFeedResponse
}

func (f *fakeServerStream) Send(resp *tripsvc.StreamTripLiveFeedResponse) error {
	f.sent = append(f.sent, resp)
	return nil
}

func (f *fakeServerStream) Context() context.Context     { return f.ctx }
func (f *fakeServerStream) SetHeader(metadata.MD) error  { return nil }
func (f *fakeServerStream) SendHeader(metadata.MD) error { return nil }
func (f *fakeServerStream) SetTrailer(metadata.MD)       {}
func (f *fakeServerStream) SendMsg(interface{}) error    { return nil }
func (f *fakeServerStream) RecvMsg(interface{}) error    { return nil }

// ── TripHandler ───────────────────────────────────────────────────────────────

func TestTripHandler_StartTrip_Success(t *testing.T) {
	svc := handlermocks.NewMockTripService(t)
	h := handler.NewTripHandler(discardLogger(), svc)

	svc.EXPECT().StartTrip(mock.Anything, "booking-1").Return("trip-new", nil)

	resp, err := h.StartTrip(context.Background(), &tripsvc.StartTripRequest{BookingId: "booking-1"})

	require.NoError(t, err)
	assert.Equal(t, "trip-new", resp.Id)
}

func TestTripHandler_StartTrip_ServiceError(t *testing.T) {
	svc := handlermocks.NewMockTripService(t)
	h := handler.NewTripHandler(discardLogger(), svc)

	svc.EXPECT().StartTrip(mock.Anything, "booking-1").Return("", model.ErrBookingNotCreated)

	_, err := h.StartTrip(context.Background(), &tripsvc.StartTripRequest{BookingId: "booking-1"})

	require.Error(t, err)
	assert.Equal(t, codes.FailedPrecondition, status.Code(err))
}

func TestTripHandler_GetTrip_Success(t *testing.T) {
	svc := handlermocks.NewMockTripService(t)
	h := handler.NewTripHandler(discardLogger(), svc)
	trip := sampleTrip()

	svc.EXPECT().GetTrip(mock.Anything, "trip-1").Return(trip, nil)

	resp, err := h.GetTrip(context.Background(), &tripsvc.GetTripRequest{Id: "trip-1"})

	require.NoError(t, err)
	require.NotNil(t, resp.Trip)
	assert.Equal(t, "trip-1", resp.Trip.Id)
	assert.Equal(t, "user-1", resp.Trip.UserId)
	assert.Equal(t, "active", resp.Trip.Status)
}

func TestTripHandler_GetTrip_NotFound(t *testing.T) {
	svc := handlermocks.NewMockTripService(t)
	h := handler.NewTripHandler(discardLogger(), svc)

	svc.EXPECT().GetTrip(mock.Anything, "trip-1").Return(model.Trip{}, model.ErrNotFound)

	_, err := h.GetTrip(context.Background(), &tripsvc.GetTripRequest{Id: "trip-1"})

	require.Error(t, err)
	assert.Equal(t, codes.NotFound, status.Code(err))
}

func TestTripHandler_GetTrip_PermissionDenied(t *testing.T) {
	svc := handlermocks.NewMockTripService(t)
	h := handler.NewTripHandler(discardLogger(), svc)

	svc.EXPECT().GetTrip(mock.Anything, "trip-1").Return(model.Trip{}, model.ErrInsufficientPermissions)

	_, err := h.GetTrip(context.Background(), &tripsvc.GetTripRequest{Id: "trip-1"})

	require.Error(t, err)
	assert.Equal(t, codes.PermissionDenied, status.Code(err))
}

func TestTripHandler_ListTrips_Success(t *testing.T) {
	svc := handlermocks.NewMockTripService(t)
	h := handler.NewTripHandler(discardLogger(), svc)
	trips := []model.Trip{sampleTrip(), sampleTrip()}

	svc.EXPECT().ListTrips(mock.Anything, mock.Anything).Return(trips, nil)

	resp, err := h.ListTrips(context.Background(), &tripsvc.ListTripsRequest{})

	require.NoError(t, err)
	assert.Len(t, resp.Trips, 2)
}

func TestTripHandler_EndTrip_Success(t *testing.T) {
	svc := handlermocks.NewMockTripService(t)
	h := handler.NewTripHandler(discardLogger(), svc)

	svc.EXPECT().EndTrip(mock.Anything, "trip-1").Return(nil)

	resp, err := h.EndTrip(context.Background(), &tripsvc.EndTripRequest{Id: "trip-1"})

	require.NoError(t, err)
	assert.NotNil(t, resp)
}

func TestTripHandler_EndTrip_InvalidTransition(t *testing.T) {
	svc := handlermocks.NewMockTripService(t)
	h := handler.NewTripHandler(discardLogger(), svc)

	svc.EXPECT().EndTrip(mock.Anything, "trip-1").Return(model.ErrInvalidTripStatusTransition)

	_, err := h.EndTrip(context.Background(), &tripsvc.EndTripRequest{Id: "trip-1"})

	require.Error(t, err)
	assert.Equal(t, codes.FailedPrecondition, status.Code(err))
}

func TestTripHandler_CancelTrip_Success(t *testing.T) {
	svc := handlermocks.NewMockTripService(t)
	h := handler.NewTripHandler(discardLogger(), svc)
	reason := "no longer needed"

	svc.EXPECT().CancelTrip(mock.Anything, "trip-1", &reason).Return(nil)

	resp, err := h.CancelTrip(context.Background(), &tripsvc.CancelTripRequest{Id: "trip-1", Reason: &reason})

	require.NoError(t, err)
	assert.NotNil(t, resp)
}

func TestTripHandler_CancelTrip_NoReason(t *testing.T) {
	svc := handlermocks.NewMockTripService(t)
	h := handler.NewTripHandler(discardLogger(), svc)

	svc.EXPECT().CancelTrip(mock.Anything, "trip-1", (*string)(nil)).Return(nil)

	resp, err := h.CancelTrip(context.Background(), &tripsvc.CancelTripRequest{Id: "trip-1"})

	require.NoError(t, err)
	assert.NotNil(t, resp)
}

func TestTripHandler_GetTripSummary_Success(t *testing.T) {
	svc := handlermocks.NewMockTripService(t)
	h := handler.NewTripHandler(discardLogger(), svc)
	summary := model.TripSummary{
		TripID:    "trip-1",
		BookingID: "booking-1",
		StartedAt: time.Now().Add(-1 * time.Hour),
		EndedAt:   time.Now(),
	}

	svc.EXPECT().GetTripSummary(mock.Anything, "trip-1").Return(summary, nil)

	resp, err := h.GetTripSummary(context.Background(), &tripsvc.GetTripSummaryRequest{Id: "trip-1"})

	require.NoError(t, err)
	require.NotNil(t, resp.Summary)
	assert.Equal(t, "trip-1", resp.Summary.TripId)
}

func TestTripHandler_GetTripStatusHistory_Success(t *testing.T) {
	svc := handlermocks.NewMockTripService(t)
	h := handler.NewTripHandler(discardLogger(), svc)
	readings := []model.TripStatusReading{
		{TripID: "trip-1", FromStatus: "", ToStatus: model.TripStatusActive},
		{TripID: "trip-1", FromStatus: model.TripStatusActive, ToStatus: model.TripStatusCompleted},
	}

	svc.EXPECT().GetTripStatusHistory(mock.Anything, mock.Anything).Return(readings, nil)

	resp, err := h.GetTripStatusHistory(context.Background(), &tripsvc.GetTripStatusHistoryRequest{Id: "trip-1"})

	require.NoError(t, err)
	assert.Len(t, resp.StatusHistory, 2)
}

// ── TripStreamHandler ─────────────────────────────────────────────────────────

func TestTripStreamHandler_StreamTripLiveFeed_Success(t *testing.T) {
	svc := handlermocks.NewMockTripStreamService(t)
	h := handler.NewTripStreamHandler(discardLogger(), svc)
	stream := &fakeServerStream{ctx: context.Background()}

	svc.EXPECT().
		StreamTripLiveFeed(mock.Anything, "trip-1", mock.Anything).
		Run(func(ctx context.Context, tripID string, send func(model.TripLiveFeed) error) {
			_ = send(model.TripLiveFeed{ElapsedSeconds: 60, CurrentCostTenge: 200, DistanceTraveledKM: 5.5})
		}).
		Return(nil)

	err := h.StreamTripLiveFeed(&tripsvc.StreamTripLiveFeedRequest{TripId: "trip-1"}, stream)

	require.NoError(t, err)
	require.Len(t, stream.sent, 1)
	assert.EqualValues(t, 60, stream.sent[0].ElapsedSeconds)
	assert.EqualValues(t, 200, stream.sent[0].CurrentCostTenge)
}

func TestTripStreamHandler_StreamTripLiveFeed_EOFIsCleanClose(t *testing.T) {
	svc := handlermocks.NewMockTripStreamService(t)
	h := handler.NewTripStreamHandler(discardLogger(), svc)
	stream := &fakeServerStream{ctx: context.Background()}

	svc.EXPECT().StreamTripLiveFeed(mock.Anything, "trip-1", mock.Anything).Return(io.EOF)

	err := h.StreamTripLiveFeed(&tripsvc.StreamTripLiveFeedRequest{TripId: "trip-1"}, stream)

	require.NoError(t, err) // EOF must be swallowed
}

func TestTripStreamHandler_StreamTripLiveFeed_ServiceError(t *testing.T) {
	svc := handlermocks.NewMockTripStreamService(t)
	h := handler.NewTripStreamHandler(discardLogger(), svc)
	stream := &fakeServerStream{ctx: context.Background()}

	svc.EXPECT().StreamTripLiveFeed(mock.Anything, "trip-1", mock.Anything).Return(model.ErrTripNotActive)

	err := h.StreamTripLiveFeed(&tripsvc.StreamTripLiveFeedRequest{TripId: "trip-1"}, stream)

	require.Error(t, err)
	assert.Equal(t, codes.FailedPrecondition, status.Code(err))
}
