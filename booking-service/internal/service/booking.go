package service

import (
	"context"
	"log/slog"
	"time"

	"github.com/sorawaslocked/car-rental-booking-service/internal/model"
	pkglog "github.com/sorawaslocked/car-rental-booking-service/internal/pkg/log"
	"github.com/sorawaslocked/car-rental-booking-service/internal/pkg/utils"
)

const (
	defaultExpiryDuration = 24 * time.Hour
	expiryPollInterval    = 30 * time.Second
)

type BookingService struct {
	log         *slog.Logger
	bookingRepo BookingRepository
	ruleRepo    PricingRuleRepository
	publisher   EventPublisher
}

func NewBookingService(log *slog.Logger, bookingRepo BookingRepository, ruleRepo PricingRuleRepository, publisher EventPublisher) *BookingService {
	return &BookingService{
		log:         pkglog.WithComponent(log, "service.BookingService"),
		bookingRepo: bookingRepo,
		ruleRepo:    ruleRepo,
		publisher:   publisher,
	}
}

func (s *BookingService) Create(ctx context.Context, data model.BookingCreate) (string, error) {
	log := pkglog.WithMethod(s.log, "Create")
	md := utils.MetadataFromCtx(ctx)
	log = pkglog.WithMetadata(log, md)

	if _, err := s.ruleRepo.GetByID(ctx, data.PricingRuleID); err != nil {
		return "", err
	}

	expiresAt := time.Now().Add(defaultExpiryDuration)

	id, err := s.bookingRepo.Create(ctx, data, expiresAt)
	if err != nil {
		log.Error("failed to create booking", pkglog.Err(err))
		return "", err
	}

	booking, err := s.bookingRepo.GetByID(ctx, id)
	if err != nil {
		log.Error("failed to fetch created booking for event", pkglog.Err(err))
		return id, nil
	}

	if err := s.publisher.PublishBookingCreated(ctx, booking); err != nil {
		log.Error("failed to publish booking.created", pkglog.Err(err))
	}

	return id, nil
}

func (s *BookingService) GetByID(ctx context.Context, id string) (model.Booking, error) {
	log := pkglog.WithMethod(s.log, "GetByID")
	md := utils.MetadataFromCtx(ctx)
	log = pkglog.WithMetadata(log, md)

	booking, err := s.bookingRepo.GetByID(ctx, id)
	if err != nil {
		if err != model.ErrNotFound {
			log.Error("failed to get booking", pkglog.Err(err))
		}
		return model.Booking{}, err
	}

	return booking, nil
}

func (s *BookingService) List(ctx context.Context, filter model.BookingListFilter) ([]model.Booking, error) {
	log := pkglog.WithMethod(s.log, "List")
	md := utils.MetadataFromCtx(ctx)
	log = pkglog.WithMetadata(log, md)

	bookings, err := s.bookingRepo.List(ctx, filter)
	if err != nil {
		log.Error("failed to list bookings", pkglog.Err(err))
		return nil, err
	}

	return bookings, nil
}

func (s *BookingService) Cancel(ctx context.Context, id string, reason *string) error {
	log := pkglog.WithMethod(s.log, "Cancel")
	md := utils.MetadataFromCtx(ctx)
	log = pkglog.WithMetadata(log, md)

	booking, err := s.bookingRepo.GetByID(ctx, id)
	if err != nil {
		if err != model.ErrNotFound {
			log.Error("failed to get booking for cancel", pkglog.Err(err))
		}
		return err
	}

	if err := model.ValidateTransition(booking.Status, model.BookingStatusCancelled); err != nil {
		return err
	}

	actorType := "user"
	if err := s.bookingRepo.UpdateStatus(ctx, id, string(model.BookingStatusCancelled), actorType, md.UserID, reason); err != nil {
		log.Error("failed to cancel booking", pkglog.Err(err))
		return err
	}

	reasonStr := ""
	if reason != nil {
		reasonStr = *reason
	}

	if err := s.publisher.PublishBookingCancelled(ctx, booking, reasonStr); err != nil {
		log.Error("failed to publish booking.cancelled", pkglog.Err(err))
	}

	return nil
}

func (s *BookingService) Complete(ctx context.Context, id string) error {
	log := pkglog.WithMethod(s.log, "Complete")

	booking, err := s.bookingRepo.GetByID(ctx, id)
	if err != nil {
		if err != model.ErrNotFound {
			log.Error("failed to get booking for complete", pkglog.Err(err))
		}
		return err
	}

	if err := model.ValidateTransition(booking.Status, model.BookingStatusCompleted); err != nil {
		return err
	}

	actorType := "system"
	if err := s.bookingRepo.UpdateStatus(ctx, id, string(model.BookingStatusCompleted), actorType, nil, nil); err != nil {
		log.Error("failed to complete booking", pkglog.Err(err))
		return err
	}

	if err := s.publisher.PublishBookingCompleted(ctx, booking); err != nil {
		log.Error("failed to publish booking.completed", pkglog.Err(err))
	}

	return nil
}

func (s *BookingService) UpdateStatus(ctx context.Context, id, rawStatus string, reason *string) error {
	log := pkglog.WithMethod(s.log, "UpdateStatus")
	md := utils.MetadataFromCtx(ctx)
	log = pkglog.WithMetadata(log, md)

	toStatus, err := model.ParseBookingStatus(rawStatus)
	if err != nil {
		return err
	}

	booking, err := s.bookingRepo.GetByID(ctx, id)
	if err != nil {
		if err != model.ErrNotFound {
			log.Error("failed to get booking for status update", pkglog.Err(err))
		}
		return err
	}

	if err := model.ValidateTransition(booking.Status, toStatus); err != nil {
		return err
	}

	actorType := "system"
	if err := s.bookingRepo.UpdateStatus(ctx, id, string(toStatus), actorType, md.UserID, reason); err != nil {
		if err != model.ErrNotFound {
			log.Error("failed to update booking status", pkglog.Err(err))
		}
		return err
	}

	return nil
}

func (s *BookingService) GetStatusHistory(ctx context.Context, filter model.BookingStatusHistoryFilter) ([]model.BookingStatusReading, error) {
	log := pkglog.WithMethod(s.log, "GetStatusHistory")
	md := utils.MetadataFromCtx(ctx)
	log = pkglog.WithMetadata(log, md)

	history, err := s.bookingRepo.GetStatusHistory(ctx, filter)
	if err != nil {
		log.Error("failed to get booking status history", pkglog.Err(err))
		return nil, err
	}

	return history, nil
}

func (s *BookingService) StartExpiryWatcher(ctx context.Context) {
	log := pkglog.WithMethod(s.log, "StartExpiryWatcher")
	ticker := time.NewTicker(expiryPollInterval)
	defer ticker.Stop()

	log.Info("expiry watcher started")

	for {
		select {
		case <-ctx.Done():
			log.Info("expiry watcher stopped")
			return
		case <-ticker.C:
			s.expireBookings(ctx)
		}
	}
}

func (s *BookingService) expireBookings(ctx context.Context) {
	log := pkglog.WithMethod(s.log, "expireBookings")

	bookings, err := s.bookingRepo.ListCreatedExpired(ctx, time.Now())
	if err != nil {
		log.Error("failed to list expired bookings", pkglog.Err(err))
		return
	}

	for _, b := range bookings {
		actorType := "system"
		if err := s.bookingRepo.UpdateStatus(ctx, b.ID, string(model.BookingStatusExpired), actorType, nil, nil); err != nil {
			log.Error("failed to expire booking", slog.String("bookingID", b.ID), pkglog.Err(err))
			continue
		}

		log.Info("booking expired", slog.String("bookingID", b.ID))

		if err := s.publisher.PublishBookingExpired(context.Background(), b); err != nil {
			log.Error("failed to publish booking.expired", slog.String("bookingID", b.ID), pkglog.Err(err))
		}
	}
}
