package service

import (
	"context"

	"carsharing/api-gateway/internal/model"
	pkglog "carsharing/shared/pkg/log"
	"carsharing/shared/pkg/utils"
)

func (s *TripService) StreamTripLiveFeed(ctx context.Context, tripID string, send func(model.TripLiveFeed) error) error {
	log := pkglog.WithMetadata(pkglog.WithMethod(s.log, "StreamTripLiveFeed"), utils.MetadataFromCtx(ctx))
	log.Debug("starting stream")

	if err := s.presenter.StreamTripLiveFeed(ctx, tripID, send); err != nil {
		log.Warn("streaming trip live feed", pkglog.Err(err))

		return err
	}

	return nil
}
