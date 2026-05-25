package service

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"carsharing/booking-service/internal/model"
	"carsharing/booking-service/internal/validation"
	sharedmodel "carsharing/shared/model"
	pkglog "carsharing/shared/pkg/log"
	"carsharing/shared/pkg/utils"
	sharedvalidation "carsharing/shared/validation"

	"github.com/go-playground/validator/v10"
)

const (
	defaultExpiryDuration = 24 * time.Hour
	expiryPollInterval    = 30 * time.Second
)

type BookingService struct {
	log         *slog.Logger
	validate    *validator.Validate
	bookingRepo BookingRepository
	ruleRepo    PricingRuleRepository
	publisher   EventPublisher
	carChecker  CarChecker
}

func NewBookingService(
	log *slog.Logger,
	validate *validator.Validate,
	bookingRepo BookingRepository,
	ruleRepo PricingRuleRepository,
	publisher EventPublisher,
	carChecker CarChecker,
) *BookingService {
	return &BookingService{
		log:         pkglog.WithComponent(log, "service.BookingService"),
		validate:    validate,
		bookingRepo: bookingRepo,
		ruleRepo:    ruleRepo,
		publisher:   publisher,
		carChecker:  carChecker,
	}
}

func (s *BookingService) Create(ctx context.Context, data validation.BookingCreate) (string, error) {
	log := pkglog.WithMetadata(pkglog.WithMethod(s.log, "Create"), utils.MetadataFromCtx(ctx))
	md := utils.MetadataFromCtx(ctx)

	if err := validation.ValidateInput(s.validate, data); err != nil {
		return "", err
	}

	if err := checkOwnerAccess(md, data.UserID); err != nil {
		return "", err
	}

	carExists, err := s.carChecker.Exists(ctx, data.CarID)
	if err != nil {
		log.Error("checking car existence", pkglog.Err(err))
		return "", err
	}
	if !carExists {
		return "", model.ErrCarNotFound
	}

	carStatus, err := s.carChecker.GetStatus(ctx, data.CarID)
	if err != nil {
		log.Error("getting car status", pkglog.Err(err))
		return "", err
	}
	if carStatus != model.CarStatusAvailable {
		return "", model.ErrCarNotAvailable
	}

	if _, err := s.ruleRepo.GetByID(ctx, data.PricingRuleID); err != nil {
		if errors.Is(err, model.ErrPricingRuleNotFound) {
			return "", model.ErrPricingRuleNotFound
		}
		log.Error("getting pricing rule", pkglog.Err(err))
		return "", err
	}

	expiresAt := time.Now().Add(defaultExpiryDuration)

	bookingData := model.BookingCreate{
		UserID:           data.UserID,
		CarID:            data.CarID,
		CommittedPeriods: data.CommittedPeriods,
		PricingRuleID:    data.PricingRuleID,
	}

	id, err := s.bookingRepo.Create(ctx, bookingData, expiresAt)
	if err != nil {
		log.Error("repo: creating booking", pkglog.Err(err))
		return "", err
	}

	booking, err := s.bookingRepo.GetByID(ctx, id)
	if err != nil {
		log.Error("repo: fetching created booking for event", pkglog.Err(err))
		return id, nil
	}

	if err := s.publisher.PublishBookingCreated(ctx, booking); err != nil {
		log.Error("event: publishing booking.created", pkglog.Err(err))
	}

	return id, nil
}

func (s *BookingService) GetByID(ctx context.Context, id string) (model.Booking, error) {
	log := pkglog.WithMetadata(pkglog.WithMethod(s.log, "GetByID"), utils.MetadataFromCtx(ctx))
	md := utils.MetadataFromCtx(ctx)

	if err := validation.ValidateID(s.validate, id); err != nil {
		return model.Booking{}, err
	}

	booking, err := s.bookingRepo.GetByID(ctx, id)
	if err != nil {
		if !errors.Is(err, model.ErrBookingNotFound) {
			log.Error("repo: getting booking", pkglog.Err(err))
		}
		return model.Booking{}, err
	}

	if err := checkOwnerAccess(md, booking.UserID); err != nil {
		return model.Booking{}, err
	}

	return booking, nil
}

func (s *BookingService) List(ctx context.Context, filter validation.BookingListFilter) ([]model.Booking, error) {
	log := pkglog.WithMetadata(pkglog.WithMethod(s.log, "List"), utils.MetadataFromCtx(ctx))
	md := utils.MetadataFromCtx(ctx)

	if err := validation.ValidateInput(s.validate, filter); err != nil {
		return nil, err
	}

	repoFilter := bookingListFilter(filter)

	if !isPrivileged(md.UserRoles) {
		repoFilter.UserID = md.UserID
	}

	bookings, err := s.bookingRepo.List(ctx, repoFilter)
	if err != nil {
		log.Error("repo: listing bookings", pkglog.Err(err))
		return nil, err
	}

	return bookings, nil
}

func (s *BookingService) Cancel(ctx context.Context, id string, reason *string) error {
	log := pkglog.WithMetadata(pkglog.WithMethod(s.log, "Cancel"), utils.MetadataFromCtx(ctx))
	md := utils.MetadataFromCtx(ctx)

	if err := validation.ValidateID(s.validate, id); err != nil {
		return err
	}

	booking, err := s.bookingRepo.GetByID(ctx, id)
	if err != nil {
		if !errors.Is(err, model.ErrBookingNotFound) {
			log.Error("repo: getting booking for cancel", pkglog.Err(err))
		}
		return err
	}

	if err := checkOwnerAccess(md, booking.UserID); err != nil {
		return err
	}

	if err := model.ValidateTransition(booking.Status, model.BookingStatusCancelled); err != nil {
		return err
	}

	if err := s.bookingRepo.UpdateStatus(ctx, id, model.BookingStatusCancelled, sharedmodel.ActorTypeUser, md.UserID, reason); err != nil {
		log.Error("repo: cancelling booking", pkglog.Err(err))
		return err
	}

	reasonStr := ""
	if reason != nil {
		reasonStr = *reason
	}

	if err := s.publisher.PublishBookingCancelled(ctx, booking, reasonStr); err != nil {
		log.Error("event: publishing booking.cancelled", pkglog.Err(err))
	}

	return nil
}

func (s *BookingService) Complete(ctx context.Context, id string) error {
	log := pkglog.WithMetadata(pkglog.WithMethod(s.log, "Complete"), utils.MetadataFromCtx(ctx))

	booking, err := s.bookingRepo.GetByID(ctx, id)
	if err != nil {
		if !errors.Is(err, model.ErrBookingNotFound) {
			log.Error("repo: getting booking for complete", pkglog.Err(err))
		}
		return err
	}

	if err := model.ValidateTransition(booking.Status, model.BookingStatusCompleted); err != nil {
		return err
	}

	if err := s.bookingRepo.UpdateStatus(ctx, id, model.BookingStatusCompleted, sharedmodel.ActorTypeSystem, nil, nil); err != nil {
		log.Error("repo: completing booking", pkglog.Err(err))
		return err
	}

	if err := s.publisher.PublishBookingCompleted(ctx, booking); err != nil {
		log.Error("event: publishing booking.completed", pkglog.Err(err))
	}

	return nil
}

func (s *BookingService) UpdateStatus(ctx context.Context, id string, data validation.BookingStatusUpdate) error {
	log := pkglog.WithMetadata(pkglog.WithMethod(s.log, "UpdateStatus"), utils.MetadataFromCtx(ctx))
	md := utils.MetadataFromCtx(ctx)

	if err := validation.ValidateID(s.validate, id); err != nil {
		return err
	}

	if err := validation.ValidateInput(s.validate, data); err != nil {
		return err
	}

	if !isPrivileged(md.UserRoles) {
		return model.ErrInsufficientPermissions
	}

	toStatus, _ := model.ParseBookingStatus(data.Status)

	booking, err := s.bookingRepo.GetByID(ctx, id)
	if err != nil {
		if !errors.Is(err, model.ErrBookingNotFound) {
			log.Error("repo: getting booking for status update", pkglog.Err(err))
		}
		return err
	}

	if err := model.ValidateTransition(booking.Status, toStatus); err != nil {
		return err
	}

	if err := s.bookingRepo.UpdateStatus(ctx, id, toStatus, sharedmodel.ActorTypeUser, md.UserID, data.Reason); err != nil {
		log.Error("repo: updating booking status", pkglog.Err(err))
		return err
	}

	return nil
}

func (s *BookingService) GetStatusHistory(ctx context.Context, filter validation.BookingStatusHistoryFilter) ([]model.BookingStatusReading, error) {
	log := pkglog.WithMetadata(pkglog.WithMethod(s.log, "GetStatusHistory"), utils.MetadataFromCtx(ctx))
	md := utils.MetadataFromCtx(ctx)

	if err := validation.ValidateInput(s.validate, filter); err != nil {
		return nil, err
	}

	booking, err := s.bookingRepo.GetByID(ctx, filter.BookingID)
	if err != nil {
		if !errors.Is(err, model.ErrBookingNotFound) {
			log.Error("repo: getting booking for history", pkglog.Err(err))
		}
		return nil, err
	}

	if err := checkOwnerAccess(md, booking.UserID); err != nil {
		return nil, err
	}

	history, err := s.bookingRepo.GetStatusHistory(ctx, bookingStatusHistoryFilter(filter))
	if err != nil {
		log.Error("repo: getting booking status history", pkglog.Err(err))
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

func isPrivileged(roles []sharedmodel.Role) bool {
	for _, r := range roles {
		if r == sharedmodel.RoleAdmin || r == sharedmodel.RoleBookingManager {
			return true
		}
	}
	return false
}

func checkOwnerAccess(md utils.Metadata, ownerID string) error {
	if !isPrivileged(md.UserRoles) && (md.UserID == nil || ownerID != *md.UserID) {
		return model.ErrInsufficientPermissions
	}
	return nil
}

func (s *BookingService) expireBookings(ctx context.Context) {
	log := pkglog.WithMethod(s.log, "expireBookings")

	bookings, err := s.bookingRepo.ListCreatedExpired(ctx, time.Now())
	if err != nil {
		log.Error("repo: listing expired bookings", pkglog.Err(err))
		return
	}

	for _, b := range bookings {
		if err := s.bookingRepo.UpdateStatus(ctx, b.ID, model.BookingStatusExpired, sharedmodel.ActorTypeSystem, nil, nil); err != nil {
			log.Error("repo: expiring booking", slog.String("bookingID", b.ID), pkglog.Err(err))
			continue
		}

		log.Info("booking expired", slog.String("bookingID", b.ID))

		if err := s.publisher.PublishBookingExpired(ctx, b); err != nil {
			log.Error("event: publishing booking.expired", slog.String("bookingID", b.ID), pkglog.Err(err))
		}
	}
}

func bookingListFilter(f validation.BookingListFilter) model.BookingListFilter {
	if f.Pagination == nil {
		f.Pagination = sharedvalidation.DefaultPagination()
	}
	return model.BookingListFilter{
		UserID:        f.UserID,
		CarID:         f.CarID,
		Status:        f.Status,
		PricingRuleID: f.PricingRuleID,
		Pagination:    sharedmodel.Pagination{Limit: f.Pagination.Limit, Offset: f.Pagination.Offset},
	}
}

func bookingStatusHistoryFilter(f validation.BookingStatusHistoryFilter) model.BookingStatusHistoryFilter {
	if f.Pagination == nil {
		f.Pagination = sharedvalidation.DefaultPagination()
	}
	filter := model.BookingStatusHistoryFilter{
		BookingID:  f.BookingID,
		Pagination: sharedmodel.Pagination{Limit: f.Pagination.Limit, Offset: f.Pagination.Offset},
	}
	if f.TimeRange != nil && f.TimeRange.From != nil && f.TimeRange.To != nil {
		filter.TimeRange = &sharedmodel.TimeRange{From: *f.TimeRange.From, To: *f.TimeRange.To}
	}
	return filter
}
