package service

import (
	"context"
	"log/slog"

	"carsharing/api-gateway/internal/model"
	pkglog "carsharing/shared/pkg/log"
	"carsharing/shared/pkg/utils"
)

type BookingService struct {
	presenter BookingPresenter
	log       *slog.Logger
}

func NewBookingService(presenter BookingPresenter, log *slog.Logger) *BookingService {
	return &BookingService{
		presenter: presenter,
		log:       pkglog.WithComponent(log, "service.BookingService"),
	}
}

func (s *BookingService) Create(ctx context.Context, data model.BookingCreate) (string, error) {
	log := pkglog.WithMetadata(pkglog.WithMethod(s.log, "Create"), utils.MetadataFromCtx(ctx))
	log.Debug("creating booking")

	id, err := s.presenter.Create(ctx, data)
	if err != nil {
		log.Warn("creating booking", pkglog.Err(err))

		return "", err
	}

	log.Debug("booking created", slog.String("id", id))

	return id, nil
}

func (s *BookingService) Get(ctx context.Context, id string) (model.Booking, error) {
	log := pkglog.WithMetadata(pkglog.WithMethod(s.log, "Get"), utils.MetadataFromCtx(ctx))
	log.Debug("getting booking")

	booking, err := s.presenter.Get(ctx, id)
	if err != nil {
		log.Warn("getting booking", pkglog.Err(err))

		return model.Booking{}, err
	}

	return booking, nil
}

func (s *BookingService) List(ctx context.Context, filter model.BookingFilter) ([]model.Booking, error) {
	log := pkglog.WithMetadata(pkglog.WithMethod(s.log, "List"), utils.MetadataFromCtx(ctx))
	log.Debug("listing bookings")

	bookings, err := s.presenter.List(ctx, filter)
	if err != nil {
		log.Warn("listing bookings", pkglog.Err(err))

		return nil, err
	}

	return bookings, nil
}

func (s *BookingService) Cancel(ctx context.Context, id string) error {
	log := pkglog.WithMetadata(pkglog.WithMethod(s.log, "Cancel"), utils.MetadataFromCtx(ctx))
	log.Debug("cancelling booking")

	if err := s.presenter.Cancel(ctx, id); err != nil {
		log.Warn("cancelling booking", pkglog.Err(err))

		return err
	}

	return nil
}

func (s *BookingService) UpdateStatus(ctx context.Context, id string, data model.BookingStatusUpdate) error {
	log := pkglog.WithMetadata(pkglog.WithMethod(s.log, "UpdateStatus"), utils.MetadataFromCtx(ctx))
	log.Debug("updating booking status")

	if err := s.presenter.UpdateStatus(ctx, id, data); err != nil {
		log.Warn("updating booking status", pkglog.Err(err))

		return err
	}

	return nil
}

func (s *BookingService) GetStatusHistory(ctx context.Context, id string, filter model.BookingStatusReadingFilter) ([]model.BookingStatusReading, error) {
	log := pkglog.WithMetadata(pkglog.WithMethod(s.log, "GetStatusHistory"), utils.MetadataFromCtx(ctx))
	log.Debug("getting booking status history")

	history, err := s.presenter.GetStatusHistory(ctx, id, filter)
	if err != nil {
		log.Warn("getting booking status history", pkglog.Err(err))

		return nil, err
	}

	return history, nil
}
