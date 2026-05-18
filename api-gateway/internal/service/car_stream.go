package service

import (
	"context"

	"carsharing/api-gateway/internal/model"
)

func (s *CarService) StreamCarsWithFilter(ctx context.Context, filter model.CarFilter, send func([]model.SlimCar) error) error {
	return s.presenter.StreamCarsWithFilter(ctx, filter, send)
}

func (s *CarService) StreamCarTelemetry(ctx context.Context, carID string, send func(model.CarTelemetryEvent) error) error {
	return s.presenter.StreamCarTelemetry(ctx, carID, send)
}
