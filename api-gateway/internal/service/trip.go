package service

import (
	"context"
	"log/slog"

	"carsharing/api-gateway/internal/model"
	pkglog "carsharing/shared/pkg/log"
	"carsharing/shared/pkg/utils"
)

type TripService struct {
	presenter TripPresenter
	log       *slog.Logger
}

func NewTripService(presenter TripPresenter, log *slog.Logger) *TripService {
	return &TripService{
		presenter: presenter,
		log:       pkglog.WithComponent(log, "service.TripService"),
	}
}

func (s *TripService) Start(ctx context.Context, bookingID string) (string, error) {
	log := pkglog.WithMetadata(pkglog.WithMethod(s.log, "Start"), utils.MetadataFromCtx(ctx))
	log.Debug("starting trip")

	id, err := s.presenter.Start(ctx, bookingID)
	if err != nil {
		log.Warn("starting trip", pkglog.Err(err))

		return "", err
	}

	log.Debug("trip created", slog.String("id", id))

	return id, nil
}

func (s *TripService) Get(ctx context.Context, id string) (model.Trip, error) {
	log := pkglog.WithMetadata(pkglog.WithMethod(s.log, "Get"), utils.MetadataFromCtx(ctx))
	log.Debug("getting trip")

	trip, err := s.presenter.Get(ctx, id)
	if err != nil {
		log.Warn("getting trip", pkglog.Err(err))

		return model.Trip{}, err
	}

	return trip, nil
}

func (s *TripService) List(ctx context.Context, filter model.TripFilter) ([]model.Trip, error) {
	log := pkglog.WithMetadata(pkglog.WithMethod(s.log, "List"), utils.MetadataFromCtx(ctx))
	log.Debug("listing trips")

	trips, err := s.presenter.List(ctx, filter)
	if err != nil {
		log.Warn("listing trips", pkglog.Err(err))

		return nil, err
	}

	return trips, nil
}

func (s *TripService) End(ctx context.Context, id string) error {
	log := pkglog.WithMetadata(pkglog.WithMethod(s.log, "End"), utils.MetadataFromCtx(ctx))
	log.Debug("ending trip")

	if err := s.presenter.End(ctx, id); err != nil {
		log.Warn("ending trip", pkglog.Err(err))

		return err
	}

	return nil
}

func (s *TripService) Cancel(ctx context.Context, id string, reason *string) error {
	log := pkglog.WithMetadata(pkglog.WithMethod(s.log, "Cancel"), utils.MetadataFromCtx(ctx))
	log.Debug("cancelling trip")

	if err := s.presenter.Cancel(ctx, id, reason); err != nil {
		log.Warn("cancelling trip", pkglog.Err(err))

		return err
	}

	return nil
}

func (s *TripService) GetSummary(ctx context.Context, id string) (model.TripSummary, error) {
	log := pkglog.WithMetadata(pkglog.WithMethod(s.log, "GetSummary"), utils.MetadataFromCtx(ctx))
	log.Debug("getting trip summary")

	summary, err := s.presenter.GetSummary(ctx, id)
	if err != nil {
		log.Warn("getting trip summary", pkglog.Err(err))

		return model.TripSummary{}, err
	}

	return summary, nil
}

func (s *TripService) GetStatusHistory(ctx context.Context, id string, filter model.TripStatusReadingFilter) ([]model.TripStatusReading, error) {
	log := pkglog.WithMetadata(pkglog.WithMethod(s.log, "GetStatusHistory"), utils.MetadataFromCtx(ctx))
	log.Debug("getting trip status history")

	history, err := s.presenter.GetStatusHistory(ctx, id, filter)
	if err != nil {
		log.Warn("getting trip status history", pkglog.Err(err))

		return nil, err
	}

	return history, nil
}
