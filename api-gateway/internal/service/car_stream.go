package service

import (
	"context"

	"carsharing/api-gateway/internal/model"
	pkglog "carsharing/shared/pkg/log"
	"carsharing/shared/pkg/utils"
)

func (s *CarService) StreamCarsWithFilter(ctx context.Context, filter model.CarFilter, send func([]model.SlimCar) error) error {
	log := pkglog.WithMetadata(pkglog.WithMethod(s.log, "StreamCarsWithFilter"), utils.MetadataFromCtx(ctx))
	log.Debug("starting stream")

	if err := s.presenter.StreamCarsWithFilter(ctx, filter, send); err != nil {
		log.Warn("streaming cars with filter", pkglog.Err(err))

		return err
	}

	return nil
}

func (s *CarService) StreamCarTelemetry(ctx context.Context, carID string, send func(model.CarTelemetryEvent) error) error {
	log := pkglog.WithMetadata(pkglog.WithMethod(s.log, "StreamCarTelemetry"), utils.MetadataFromCtx(ctx))
	log.Debug("starting stream")

	if err := s.presenter.StreamCarTelemetry(ctx, carID, send); err != nil {
		log.Warn("streaming car telemetry", pkglog.Err(err))

		return err
	}

	return nil
}

func (s *CarService) StreamCarStatusUpdates(ctx context.Context, carID string, send func(model.CarStatusEvent) error) error {
	log := pkglog.WithMetadata(pkglog.WithMethod(s.log, "StreamCarStatusUpdates"), utils.MetadataFromCtx(ctx))
	log.Debug("starting stream")

	if err := s.presenter.StreamCarStatusUpdates(ctx, carID, send); err != nil {
		log.Warn("streaming car status updates", pkglog.Err(err))

		return err
	}

	return nil
}
