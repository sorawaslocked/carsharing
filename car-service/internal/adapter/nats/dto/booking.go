package dto

import (
	"carsharing/car-service/internal/model"

	eventbooking "github.com/sorawaslocked/car-rental-protos/gen/event/booking"
)

func BookingCreatedEventFromProto(e *eventbooking.BookingCreatedEvent) model.BookingCreatedEvent {
	event := model.BookingCreatedEvent{
		BookingID: e.GetBookingId(),
		CarID:     e.GetCarId(),
		UserID:    e.GetUserId(),
	}
	if e.GetStartsAt() != nil {
		event.StartsAt = e.GetStartsAt().AsTime()
	}
	if e.GetEndsAt() != nil {
		event.EndsAt = e.GetEndsAt().AsTime()
	}
	return event
}

func BookingCancelledEventFromProto(e *eventbooking.BookingCancelledEvent) model.BookingCancelledEvent {
	return model.BookingCancelledEvent{
		BookingID: e.GetBookingId(),
		CarID:     e.GetCarId(),
		UserID:    e.GetUserId(),
		Reason:    e.GetReason(),
	}
}

func BookingExpiredEventFromProto(e *eventbooking.BookingExpiredEvent) model.BookingExpiredEvent {
	event := model.BookingExpiredEvent{
		BookingID: e.GetBookingId(),
		CarID:     e.GetCarId(),
		UserID:    e.GetUserId(),
	}
	if e.GetExpiredAt() != nil {
		event.ExpiredAt = e.GetExpiredAt().AsTime()
	}
	return event
}

func BookingCompletedEventFromProto(e *eventbooking.BookingCompletedEvent) model.BookingCompletedEvent {
	event := model.BookingCompletedEvent{
		BookingID: e.GetBookingId(),
		CarID:     e.GetCarId(),
		UserID:    e.GetUserId(),
	}
	if e.GetCompletedAt() != nil {
		event.CompletedAt = e.GetCompletedAt().AsTime()
	}
	return event
}
