package handler

import (
	"context"

	"carsharing/trip-service/internal/model"
	"carsharing/trip-service/internal/validation"
)

type TripService interface {
	StartTrip(ctx context.Context, bookingID string) (string, error)
	GetTrip(ctx context.Context, id string) (model.Trip, error)
	ListTrips(ctx context.Context, filter validation.TripFilter) ([]model.Trip, error)
	EndTrip(ctx context.Context, id string) error
	CancelTrip(ctx context.Context, id string, reason *string) error
	GetTripSummary(ctx context.Context, tripID string) (model.TripSummary, error)
	GetTripStatusHistory(ctx context.Context, filter validation.TripStatusHistoryFilter) ([]model.TripStatusReading, error)
	StreamTripLiveFeed(ctx context.Context, tripID string, send func(model.TripLiveFeed) error) error
}
