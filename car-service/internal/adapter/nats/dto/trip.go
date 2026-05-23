package dto

import (
	"carsharing/car-service/internal/model"

	eventtrip "github.com/sorawaslocked/car-rental-protos/gen/event/trip"
)

func TripStartedEventFromProto(e *eventtrip.TripStartedEvent) model.TripStartedEvent {
	event := model.TripStartedEvent{
		TripID:    e.GetTripId(),
		BookingID: e.GetBookingId(),
		CarID:     e.GetCarId(),
		UserID:    e.GetUserId(),
	}
	if e.GetStartedAt() != nil {
		event.StartedAt = e.GetStartedAt().AsTime()
	}
	return event
}

func TripEndedEventFromProto(e *eventtrip.TripEndedEvent) model.TripEndedEvent {
	event := model.TripEndedEvent{
		TripID:    e.GetTripId(),
		BookingID: e.GetBookingId(),
		CarID:     e.GetCarId(),
		UserID:    e.GetUserId(),
	}
	if e.GetEndedAt() != nil {
		event.EndedAt = e.GetEndedAt().AsTime()
	}
	return event
}

func TripCancelledEventFromProto(e *eventtrip.TripCancelledEvent) model.TripCancelledEvent {
	event := model.TripCancelledEvent{
		TripID:    e.GetTripId(),
		BookingID: e.GetBookingId(),
		CarID:     e.GetCarId(),
		UserID:    e.GetUserId(),
		Reason:    e.GetReason(),
	}
	if e.GetCancelledAt() != nil {
		event.CancelledAt = e.GetCancelledAt().AsTime()
	}
	return event
}
