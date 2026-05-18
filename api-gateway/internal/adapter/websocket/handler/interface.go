package handler

import (
	"context"

	"carsharing/api-gateway/internal/model"
)

type CarStreamService interface {
	StreamCarsWithFilter(ctx context.Context, filter model.CarFilter, send func([]model.SlimCar) error) error
	StreamCarTelemetry(ctx context.Context, carID string, send func(model.CarTelemetryEvent) error) error
}

type TripStreamService interface {
	StreamTripLiveFeed(ctx context.Context, tripID string, send func(model.TripLiveFeed) error) error
}
