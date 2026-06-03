package service

import (
	"context"

	"carsharing/trip-service/internal/model"
)

type Transactor interface {
	InTx(ctx context.Context, fn func(ctx context.Context) error) error
}

type TripRepository interface {
	Create(ctx context.Context, trip model.TripCreate) (model.Trip, error)
	GetByID(ctx context.Context, id string) (model.Trip, error)
	List(ctx context.Context, filter model.TripFilter) ([]model.Trip, error)
	Update(ctx context.Context, id string, update model.TripUpdate) (model.Trip, error)
}

type TripSummaryRepository interface {
	Create(ctx context.Context, summary model.TripSummaryCreate) (model.TripSummary, error)
	GetByTripID(ctx context.Context, tripID string) (model.TripSummary, error)
}

type TripStatusReadingRepository interface {
	Create(ctx context.Context, reading model.TripStatusReadingCreate) (model.TripStatusReading, error)
	List(ctx context.Context, filter model.TripStatusReadingFilter) ([]model.TripStatusReading, error)
}

type BookingClient interface {
	GetBooking(ctx context.Context, id string) (model.Booking, error)
}

type TelematicsClient interface {
	GetLatestTelemetry(ctx context.Context, carID string) (model.CarTelemetry, error)
	StreamTelemetry(ctx context.Context, carID string, fn func(model.CarTelemetry) error) error
}

type ZonePricingClient interface {
	GetZonePricing(ctx context.Context, lat, lng float64) (int32, error)
}

type EventPublisher interface {
	PublishTripStarted(ctx context.Context, trip model.Trip) error
	PublishTripEnded(ctx context.Context, trip model.Trip) error
	PublishTripCancelled(ctx context.Context, trip model.Trip) error
}
