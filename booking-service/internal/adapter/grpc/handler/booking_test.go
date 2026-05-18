package handler_test

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"testing"

	"carsharing/booking-service/internal/adapter/grpc/handler"
	"carsharing/booking-service/internal/adapter/grpc/handler/mocks"
	"carsharing/booking-service/internal/model"
	servicebookingpb "github.com/sorawaslocked/car-rental-protos/gen/service/booking"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func discardLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}

func assertCode(t *testing.T, err error, want codes.Code) {
	t.Helper()
	require.Error(t, err)
	st, ok := status.FromError(err)
	require.True(t, ok, "expected a gRPC status error")
	assert.Equal(t, want, st.Code())
}

// --- CreateBooking ---

func TestBookingHandler_CreateBooking_HappyPath(t *testing.T) {
	svc := mocks.NewMockBookingService(t)
	svc.EXPECT().Create(context.Background(), model.BookingCreate{
		UserID:        "u-1",
		CarID:         "c-1",
		PricingRuleID: "r-1",
	}).Return("b-1", nil)

	h := handler.NewBookingHandler(discardLogger(), svc)
	resp, err := h.CreateBooking(context.Background(), &servicebookingpb.CreateBookingRequest{
		UserId:        "u-1",
		CarId:         "c-1",
		PricingRuleId: "r-1",
	})

	require.NoError(t, err)
	assert.Equal(t, "b-1", resp.Id)
}

func TestBookingHandler_CreateBooking_ServiceError(t *testing.T) {
	svc := mocks.NewMockBookingService(t)
	svc.EXPECT().Create(context.Background(), model.BookingCreate{PricingRuleID: "missing"}).
		Return("", model.ErrNotFound)

	h := handler.NewBookingHandler(discardLogger(), svc)
	_, err := h.CreateBooking(context.Background(), &servicebookingpb.CreateBookingRequest{PricingRuleId: "missing"})

	assertCode(t, err, codes.NotFound)
}

// --- GetBooking ---

func TestBookingHandler_GetBooking_HappyPath(t *testing.T) {
	booking := model.Booking{ID: "b-1", UserID: "u-1", Status: model.BookingStatusCreated}

	svc := mocks.NewMockBookingService(t)
	svc.EXPECT().GetByID(context.Background(), "b-1").Return(booking, nil)

	h := handler.NewBookingHandler(discardLogger(), svc)
	resp, err := h.GetBooking(context.Background(), &servicebookingpb.GetBookingRequest{Id: "b-1"})

	require.NoError(t, err)
	assert.Equal(t, "b-1", resp.Booking.Id)
	assert.Equal(t, "u-1", resp.Booking.UserId)
	assert.Equal(t, string(model.BookingStatusCreated), resp.Booking.Status)
}

func TestBookingHandler_GetBooking_NotFound(t *testing.T) {
	svc := mocks.NewMockBookingService(t)
	svc.EXPECT().GetByID(context.Background(), "missing").Return(model.Booking{}, model.ErrNotFound)

	h := handler.NewBookingHandler(discardLogger(), svc)
	_, err := h.GetBooking(context.Background(), &servicebookingpb.GetBookingRequest{Id: "missing"})

	assertCode(t, err, codes.NotFound)
}

// --- ListBookings ---

func TestBookingHandler_ListBookings_HappyPath(t *testing.T) {
	bookings := []model.Booking{{ID: "b-1"}, {ID: "b-2"}}
	filter := model.BookingListFilter{}

	svc := mocks.NewMockBookingService(t)
	svc.EXPECT().List(context.Background(), filter).Return(bookings, nil)

	h := handler.NewBookingHandler(discardLogger(), svc)
	resp, err := h.ListBookings(context.Background(), &servicebookingpb.ListBookingsRequest{})

	require.NoError(t, err)
	assert.Len(t, resp.Bookings, 2)
	assert.Equal(t, "b-1", resp.Bookings[0].Id)
	assert.Equal(t, "b-2", resp.Bookings[1].Id)
}

func TestBookingHandler_ListBookings_ServiceError(t *testing.T) {
	svc := mocks.NewMockBookingService(t)
	svc.EXPECT().List(context.Background(), model.BookingListFilter{}).
		Return(nil, errors.New("db error"))

	h := handler.NewBookingHandler(discardLogger(), svc)
	_, err := h.ListBookings(context.Background(), &servicebookingpb.ListBookingsRequest{})

	assertCode(t, err, codes.Internal)
}

// --- CancelBooking ---

func TestBookingHandler_CancelBooking_HappyPath(t *testing.T) {
	svc := mocks.NewMockBookingService(t)
	svc.EXPECT().Cancel(context.Background(), "b-1", (*string)(nil)).Return(nil)

	h := handler.NewBookingHandler(discardLogger(), svc)
	resp, err := h.CancelBooking(context.Background(), &servicebookingpb.CancelBookingRequest{Id: "b-1"})

	require.NoError(t, err)
	assert.NotNil(t, resp)
}

func TestBookingHandler_CancelBooking_NotFound(t *testing.T) {
	svc := mocks.NewMockBookingService(t)
	svc.EXPECT().Cancel(context.Background(), "missing", (*string)(nil)).Return(model.ErrNotFound)

	h := handler.NewBookingHandler(discardLogger(), svc)
	_, err := h.CancelBooking(context.Background(), &servicebookingpb.CancelBookingRequest{Id: "missing"})

	assertCode(t, err, codes.NotFound)
}

func TestBookingHandler_CancelBooking_InvalidTransition(t *testing.T) {
	svc := mocks.NewMockBookingService(t)
	svc.EXPECT().Cancel(context.Background(), "b-1", (*string)(nil)).Return(model.ErrInvalidTransition)

	h := handler.NewBookingHandler(discardLogger(), svc)
	_, err := h.CancelBooking(context.Background(), &servicebookingpb.CancelBookingRequest{Id: "b-1"})

	assertCode(t, err, codes.InvalidArgument)
}

// --- UpdateBookingStatus ---

func TestBookingHandler_UpdateBookingStatus_HappyPath(t *testing.T) {
	svc := mocks.NewMockBookingService(t)
	svc.EXPECT().UpdateStatus(context.Background(), "b-1", "cancelled", (*string)(nil)).Return(nil)

	h := handler.NewBookingHandler(discardLogger(), svc)
	resp, err := h.UpdateBookingStatus(context.Background(), &servicebookingpb.UpdateBookingStatusRequest{
		Id:     "b-1",
		Status: "cancelled",
	})

	require.NoError(t, err)
	assert.NotNil(t, resp)
}

func TestBookingHandler_UpdateBookingStatus_InvalidStatus(t *testing.T) {
	svc := mocks.NewMockBookingService(t)
	svc.EXPECT().UpdateStatus(context.Background(), "b-1", "BOGUS", (*string)(nil)).Return(model.ErrInvalidStatus)

	h := handler.NewBookingHandler(discardLogger(), svc)
	_, err := h.UpdateBookingStatus(context.Background(), &servicebookingpb.UpdateBookingStatusRequest{
		Id:     "b-1",
		Status: "BOGUS",
	})

	assertCode(t, err, codes.InvalidArgument)
}

func TestBookingHandler_UpdateBookingStatus_InvalidTransition(t *testing.T) {
	svc := mocks.NewMockBookingService(t)
	svc.EXPECT().UpdateStatus(context.Background(), "b-1", "cancelled", (*string)(nil)).Return(model.ErrInvalidTransition)

	h := handler.NewBookingHandler(discardLogger(), svc)
	_, err := h.UpdateBookingStatus(context.Background(), &servicebookingpb.UpdateBookingStatusRequest{
		Id:     "b-1",
		Status: "cancelled",
	})

	assertCode(t, err, codes.InvalidArgument)
}

// --- GetBookingStatusHistory ---

func TestBookingHandler_GetBookingStatusHistory_HappyPath(t *testing.T) {
	history := []model.BookingStatusReading{
		{ID: "h-1", BookingID: "b-1", FromStatus: "created", ToStatus: "cancelled"},
		{ID: "h-2", BookingID: "b-1", FromStatus: "cancelled", ToStatus: "expired"},
	}
	filter := model.BookingStatusHistoryFilter{BookingID: "b-1"}

	svc := mocks.NewMockBookingService(t)
	svc.EXPECT().GetStatusHistory(context.Background(), filter).Return(history, nil)

	h := handler.NewBookingHandler(discardLogger(), svc)
	resp, err := h.GetBookingStatusHistory(context.Background(), &servicebookingpb.GetBookingStatusHistoryRequest{Id: "b-1"})

	require.NoError(t, err)
	assert.Len(t, resp.StatusHistory, 2)
	assert.Equal(t, "h-1", resp.StatusHistory[0].Id)
	assert.Equal(t, "created", resp.StatusHistory[0].FromStatus)
	assert.Equal(t, "cancelled", resp.StatusHistory[0].ToStatus)
}

func TestBookingHandler_GetBookingStatusHistory_ServiceError(t *testing.T) {
	svc := mocks.NewMockBookingService(t)
	svc.EXPECT().GetStatusHistory(context.Background(), model.BookingStatusHistoryFilter{BookingID: "b-1"}).
		Return(nil, errors.New("db error"))

	h := handler.NewBookingHandler(discardLogger(), svc)
	_, err := h.GetBookingStatusHistory(context.Background(), &servicebookingpb.GetBookingStatusHistoryRequest{Id: "b-1"})

	assertCode(t, err, codes.Internal)
}
