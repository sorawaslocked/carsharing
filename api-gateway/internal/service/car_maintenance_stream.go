package service

import (
	"context"

	"carsharing/api-gateway/internal/model"
	pkglog "carsharing/shared/pkg/log"
	"carsharing/shared/pkg/utils"
)

func (s *CarMaintenanceService) StreamMaintenanceEvents(ctx context.Context, send func(model.CarMaintenanceEvent) error) error {
	log := pkglog.WithMetadata(pkglog.WithMethod(s.log, "StreamMaintenanceEvents"), utils.MetadataFromCtx(ctx))
	log.Debug("starting stream")

	if err := s.presenter.StreamMaintenanceEvents(ctx, send); err != nil {
		log.Warn("streaming maintenance events", pkglog.Err(err))

		return err
	}

	return nil
}
