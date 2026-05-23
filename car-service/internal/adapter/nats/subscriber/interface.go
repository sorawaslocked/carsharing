package subscriber

import (
	"context"

	"carsharing/car-service/internal/model"
)

type BookingEventHandler interface {
	OnBookingCreated(ctx context.Context, event model.BookingCreatedEvent) error
	OnBookingCancelled(ctx context.Context, event model.BookingCancelledEvent) error
	OnBookingExpired(ctx context.Context, event model.BookingExpiredEvent) error
	OnBookingCompleted(ctx context.Context, event model.BookingCompletedEvent) error
}

type TripEventHandler interface {
	OnTripStarted(ctx context.Context, event model.TripStartedEvent) error
	OnTripEnded(ctx context.Context, event model.TripEndedEvent) error
	OnTripCancelled(ctx context.Context, event model.TripCancelledEvent) error
}
