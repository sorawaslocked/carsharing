package service

import (
	"context"
	"log/slog"

	"carsharing/api-gateway/internal/model"
	pkglog "carsharing/shared/pkg/log"
	"carsharing/shared/pkg/utils"
)

type ZoneService struct {
	presenter ZonePresenter
	log       *slog.Logger
}

func NewZoneService(presenter ZonePresenter, log *slog.Logger) *ZoneService {
	return &ZoneService{
		presenter: presenter,
		log:       pkglog.WithComponent(log, "service.ZoneService"),
	}
}

func (s *ZoneService) Create(ctx context.Context, data model.ZoneCreate) (string, error) {
	log := pkglog.WithMetadata(pkglog.WithMethod(s.log, "Create"), utils.MetadataFromCtx(ctx))
	log.Debug("creating zone")

	id, err := s.presenter.Create(ctx, data)
	if err != nil {
		log.Warn("creating zone", pkglog.Err(err))

		return "", err
	}

	log.Debug("zone created", slog.String("id", id))

	return id, nil
}

func (s *ZoneService) Get(ctx context.Context, id string) (model.Zone, error) {
	log := pkglog.WithMetadata(pkglog.WithMethod(s.log, "Get"), utils.MetadataFromCtx(ctx))
	log.Debug("getting zone")

	zone, err := s.presenter.Get(ctx, id)
	if err != nil {
		log.Warn("getting zone", pkglog.Err(err))

		return model.Zone{}, err
	}

	return zone, nil
}

func (s *ZoneService) List(ctx context.Context, filter model.ZoneFilter) ([]model.Zone, error) {
	log := pkglog.WithMetadata(pkglog.WithMethod(s.log, "List"), utils.MetadataFromCtx(ctx))
	log.Debug("listing zones")

	zones, err := s.presenter.List(ctx, filter)
	if err != nil {
		log.Warn("listing zones", pkglog.Err(err))

		return nil, err
	}

	return zones, nil
}

func (s *ZoneService) Update(ctx context.Context, id string, data model.ZoneUpdate) error {
	log := pkglog.WithMetadata(pkglog.WithMethod(s.log, "Update"), utils.MetadataFromCtx(ctx))
	log.Debug("updating zone")

	if err := s.presenter.Update(ctx, id, data); err != nil {
		log.Warn("updating zone", pkglog.Err(err))

		return err
	}

	return nil
}

func (s *ZoneService) Delete(ctx context.Context, id string) error {
	log := pkglog.WithMetadata(pkglog.WithMethod(s.log, "Delete"), utils.MetadataFromCtx(ctx))
	log.Debug("deleting zone")

	if err := s.presenter.Delete(ctx, id); err != nil {
		log.Warn("deleting zone", pkglog.Err(err))

		return err
	}

	return nil
}
