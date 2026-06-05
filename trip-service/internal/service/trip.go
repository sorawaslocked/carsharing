package service

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"

	sharedmodel "carsharing/shared/model"
	pkglog "carsharing/shared/pkg/log"
	pkgutils "carsharing/shared/pkg/utils"
	sharedvalidation "carsharing/shared/validation"
	"carsharing/trip-service/internal/model"
	"carsharing/trip-service/internal/validation"
)

type TripService struct {
	log         *slog.Logger
	validate    *validator.Validate
	transactor  Transactor
	tripRepo    TripRepository
	summaryRepo TripSummaryRepository
	statusRepo  TripStatusReadingRepository
	booking     BookingClient
	telematics  TelematicsClient
	zonePricing ZonePricingClient
	publisher   EventPublisher
}

func NewTripService(
	log *slog.Logger,
	validate *validator.Validate,
	transactor Transactor,
	tripRepo TripRepository,
	summaryRepo TripSummaryRepository,
	statusRepo TripStatusReadingRepository,
	booking BookingClient,
	telematics TelematicsClient,
	zonePricing ZonePricingClient,
	publisher EventPublisher,
) *TripService {
	return &TripService{
		log:         pkglog.WithComponent(log, "service.TripService"),
		validate:    validate,
		transactor:  transactor,
		tripRepo:    tripRepo,
		summaryRepo: summaryRepo,
		statusRepo:  statusRepo,
		booking:     booking,
		telematics:  telematics,
		zonePricing: zonePricing,
		publisher:   publisher,
	}
}

func (s *TripService) StartTrip(ctx context.Context, bookingID string) (string, error) {
	md := pkgutils.MetadataFromCtx(ctx)
	log := pkglog.WithMetadata(pkglog.WithMethod(s.log, "StartTrip"), md)

	if err := validation.ValidateID(s.validate, bookingID); err != nil {
		return "", err
	}

	booking, err := s.booking.GetBooking(ctx, bookingID)
	if err != nil {
		return "", err
	}

	if !isPrivileged(md.UserRoles) && (md.UserID == nil || booking.UserID != *md.UserID) {
		return "", model.ErrInsufficientPermissions
	}

	if booking.Status != model.BookingStatusCreated {
		return "", model.ErrBookingNotCreated
	}

	telemetry, err := s.telematics.GetLatestTelemetry(ctx, booking.CarID)
	if err != nil {
		log.Error("telematics: getting latest telemetry", pkglog.Err(err))
		return "", err
	}

	now := time.Now()
	tripID := uuid.New().String()
	actorID := booking.UserID

	var trip model.Trip
	if err := s.transactor.InTx(ctx, func(ctx context.Context) error {
		var e error
		trip, e = s.tripRepo.Create(ctx, model.TripCreate{
			ID:        tripID,
			BookingID: bookingID,
			UserID:    booking.UserID,
			CarID:     booking.CarID,
			Status:    model.TripStatusActive,
			StartedAt: now,
			StartLocation: sharedmodel.Location{
				Latitude:  telemetry.Location.Latitude,
				Longitude: telemetry.Location.Longitude,
			},
			StartMileageKM: telemetry.MileageKM,
			StartFuelLevel: telemetry.FuelLevel,
			CreatedAt:      now,
			UpdatedAt:      now,
		})
		if e != nil {
			return e
		}
		_, e = s.statusRepo.Create(ctx, model.TripStatusReadingCreate{
			TripID:     tripID,
			FromStatus: model.TripStatus(""),
			ToStatus:   model.TripStatusActive,
			ActorType:  sharedmodel.ActorTypeUser,
			ActorID:    &actorID,
			ChangedAt:  now,
		})
		return e
	}); err != nil {
		log.Error("repo: starting trip", pkglog.Err(err))
		return "", err
	}

	if err = s.publisher.PublishTripStarted(ctx, trip); err != nil {
		log.Error("event: publishing trip started", pkglog.Err(err))
	}

	return tripID, nil
}

func (s *TripService) GetTrip(ctx context.Context, id string) (model.Trip, error) {
	md := pkgutils.MetadataFromCtx(ctx)
	log := pkglog.WithMetadata(pkglog.WithMethod(s.log, "GetTrip"), md)

	if err := validation.ValidateID(s.validate, id); err != nil {
		return model.Trip{}, err
	}

	trip, err := s.tripRepo.GetByID(ctx, id)
	if err != nil {
		if err != model.ErrNotFound {
			log.Error("repo: getting trip", pkglog.Err(err))
		}
		return model.Trip{}, err
	}

	if !isPrivileged(md.UserRoles) && (md.UserID == nil || trip.UserID != *md.UserID) {
		return model.Trip{}, model.ErrInsufficientPermissions
	}

	return trip, nil
}

func (s *TripService) ListTrips(ctx context.Context, filter validation.TripFilter) ([]model.Trip, error) {
	md := pkgutils.MetadataFromCtx(ctx)
	log := pkglog.WithMetadata(pkglog.WithMethod(s.log, "ListTrips"), md)

	if err := validation.ValidateInput(s.validate, filter); err != nil {
		return nil, err
	}

	mf := tripFilter(filter)

	if !isPrivileged(md.UserRoles) {
		mf.UserID = md.UserID
	}

	trips, err := s.tripRepo.List(ctx, mf)
	if err != nil {
		log.Error("repo: listing trips", pkglog.Err(err))
		return nil, err
	}

	return trips, nil
}

func (s *TripService) EndTrip(ctx context.Context, id string) error {
	md := pkgutils.MetadataFromCtx(ctx)
	log := pkglog.WithMetadata(pkglog.WithMethod(s.log, "EndTrip"), md)

	if err := validation.ValidateID(s.validate, id); err != nil {
		return err
	}

	trip, err := s.tripRepo.GetByID(ctx, id)
	if err != nil {
		if err != model.ErrNotFound {
			log.Error("repo: getting trip", pkglog.Err(err))
		}
		return err
	}

	if !isPrivileged(md.UserRoles) && (md.UserID == nil || trip.UserID != *md.UserID) {
		return model.ErrInsufficientPermissions
	}

	if !trip.Status.CanTransitionTo(model.TripStatusCompleted) {
		return model.ErrInvalidTripStatusTransition
	}

	telemetry, err := s.telematics.GetLatestTelemetry(ctx, trip.CarID)
	if err != nil {
		log.Error("telematics: getting latest telemetry", pkglog.Err(err))
		return err
	}

	booking, err := s.booking.GetBooking(ctx, trip.BookingID)
	if err != nil {
		log.Error("grpc: getting booking", pkglog.Err(err))
		return err
	}

	now := time.Now()
	durationSeconds := int64(now.Sub(trip.StartedAt).Seconds())
	distanceKM := float64(telemetry.MileageKM - trip.StartMileageKM)

	endLocation := sharedmodel.Location{
		Latitude:  telemetry.Location.Latitude,
		Longitude: telemetry.Location.Longitude,
	}
	endMileage := telemetry.MileageKM

	feeAdjustment, err := s.zonePricing.GetZonePricing(ctx, endLocation.Latitude, endLocation.Longitude)
	if err != nil {
		if !errors.Is(err, model.ErrLocationInNoDropZone) {
			log.Error("zone: getting zone pricing", pkglog.Err(err))
		}
		return err
	}

	baseCost, distCost, overtimeCost := calculateCosts(booking.PricingSnapshot, booking.CommittedPeriods, durationSeconds, distanceKM)
	totalCost := baseCost + distCost + overtimeCost + feeAdjustment

	actorID := trip.UserID
	var updatedTrip model.Trip
	if err := s.transactor.InTx(ctx, func(ctx context.Context) error {
		var e error
		updatedTrip, e = s.tripRepo.Update(ctx, id, model.TripUpdate{
			Status:             ptr(model.TripStatusCompleted),
			EndedAt:            &now,
			EndLocation:        &endLocation,
			EndMileageKM:       &endMileage,
			EndFuelLevel:       telemetry.FuelLevel,
			DistanceTraveledKM: ptr(distanceKM),
			DurationSeconds:    ptr(durationSeconds),
			FinalCostTenge:     ptr(totalCost),
			UpdatedAt:          now,
			ExpectedUpdatedAt:  &trip.UpdatedAt,
		})
		if e != nil {
			return e
		}
		if _, e = s.statusRepo.Create(ctx, model.TripStatusReadingCreate{
			TripID:     id,
			FromStatus: model.TripStatusActive,
			ToStatus:   model.TripStatusCompleted,
			ActorType:  sharedmodel.ActorTypeUser,
			ActorID:    &actorID,
			ChangedAt:  now,
		}); e != nil {
			return e
		}
		_, e = s.summaryRepo.Create(ctx, model.TripSummaryCreate{
			TripID:                 id,
			BookingID:              trip.BookingID,
			StartedAt:              trip.StartedAt,
			EndedAt:                now,
			DurationSeconds:        durationSeconds,
			DistanceTraveledKM:     distanceKM,
			PricingSnapshot:        booking.PricingSnapshot,
			BaseCostTenge:          baseCost,
			DistanceCostTenge:      distCost,
			OvertimeCostTenge:      overtimeCost,
			ZoneFeeAdjustmentTenge: feeAdjustment,
			TotalCostTenge:         totalCost,
		})
		return e
	}); err != nil {
		if !errors.Is(err, model.ErrConflict) {
			log.Error("repo: ending trip", pkglog.Err(err))
		}
		return err
	}

	if err = s.publisher.PublishTripEnded(ctx, updatedTrip); err != nil {
		log.Error("event: publishing trip ended", pkglog.Err(err))
	}

	return nil
}

func (s *TripService) CancelTrip(ctx context.Context, id string, reason *string) error {
	md := pkgutils.MetadataFromCtx(ctx)
	log := pkglog.WithMetadata(pkglog.WithMethod(s.log, "CancelTrip"), md)

	if err := validation.ValidateID(s.validate, id); err != nil {
		return err
	}

	trip, err := s.tripRepo.GetByID(ctx, id)
	if err != nil {
		if err != model.ErrNotFound {
			log.Error("repo: getting trip", pkglog.Err(err))
		}
		return err
	}

	if !isPrivileged(md.UserRoles) && (md.UserID == nil || trip.UserID != *md.UserID) {
		return model.ErrInsufficientPermissions
	}

	if !trip.Status.CanTransitionTo(model.TripStatusCancelled) {
		return model.ErrInvalidTripStatusTransition
	}

	now := time.Now()
	actorID := trip.UserID

	var updatedTrip model.Trip
	if err := s.transactor.InTx(ctx, func(ctx context.Context) error {
		var e error
		updatedTrip, e = s.tripRepo.Update(ctx, id, model.TripUpdate{
			Status:            ptr(model.TripStatusCancelled),
			CancelReason:      reason,
			UpdatedAt:         now,
			ExpectedUpdatedAt: &trip.UpdatedAt,
		})
		if e != nil {
			return e
		}
		_, e = s.statusRepo.Create(ctx, model.TripStatusReadingCreate{
			TripID:     id,
			FromStatus: model.TripStatusActive,
			ToStatus:   model.TripStatusCancelled,
			ActorType:  sharedmodel.ActorTypeUser,
			ActorID:    &actorID,
			Reason:     reason,
			ChangedAt:  now,
		})
		return e
	}); err != nil {
		if !errors.Is(err, model.ErrConflict) {
			log.Error("repo: cancelling trip", pkglog.Err(err))
		}
		return err
	}

	if err = s.publisher.PublishTripCancelled(ctx, updatedTrip); err != nil {
		log.Error("event: publishing trip cancelled", pkglog.Err(err))
	}

	return nil
}

func (s *TripService) GetTripSummary(ctx context.Context, tripID string) (model.TripSummary, error) {
	md := pkgutils.MetadataFromCtx(ctx)
	log := pkglog.WithMetadata(pkglog.WithMethod(s.log, "GetTripSummary"), md)

	if err := validation.ValidateID(s.validate, tripID); err != nil {
		return model.TripSummary{}, err
	}

	trip, err := s.tripRepo.GetByID(ctx, tripID)
	if err != nil {
		if err != model.ErrNotFound {
			log.Error("repo: getting trip", pkglog.Err(err))
		}
		return model.TripSummary{}, err
	}

	if !isPrivileged(md.UserRoles) && (md.UserID == nil || trip.UserID != *md.UserID) {
		return model.TripSummary{}, model.ErrInsufficientPermissions
	}

	if trip.Status != model.TripStatusCompleted {
		return model.TripSummary{}, model.ErrTripNotCompleted
	}

	summary, err := s.summaryRepo.GetByTripID(ctx, tripID)
	if err != nil {
		log.Error("repo: getting trip summary", pkglog.Err(err))
		return model.TripSummary{}, err
	}
	return summary, nil
}

func (s *TripService) GetTripStatusHistory(ctx context.Context, filter validation.TripStatusHistoryFilter) ([]model.TripStatusReading, error) {
	md := pkgutils.MetadataFromCtx(ctx)
	log := pkglog.WithMetadata(pkglog.WithMethod(s.log, "GetTripStatusHistory"), md)

	if err := validation.ValidateInput(s.validate, filter); err != nil {
		return nil, err
	}

	trip, err := s.tripRepo.GetByID(ctx, filter.TripID)
	if err != nil {
		if err != model.ErrNotFound {
			log.Error("repo: getting trip", pkglog.Err(err))
		}
		return nil, err
	}

	if !isPrivileged(md.UserRoles) && (md.UserID == nil || trip.UserID != *md.UserID) {
		return nil, model.ErrInsufficientPermissions
	}

	return s.statusRepo.List(ctx, tripStatusHistoryFilter(filter))
}

func (s *TripService) StreamTripLiveFeed(ctx context.Context, tripID string, send func(model.TripLiveFeed) error) error {
	md := pkgutils.MetadataFromCtx(ctx)
	log := pkglog.WithMetadata(pkglog.WithMethod(s.log, "StreamTripLiveFeed"), md)

	if err := validation.ValidateID(s.validate, tripID); err != nil {
		return err
	}

	trip, err := s.tripRepo.GetByID(ctx, tripID)
	if err != nil {
		return err
	}

	if !isPrivileged(md.UserRoles) && (md.UserID == nil || trip.UserID != *md.UserID) {
		return model.ErrInsufficientPermissions
	}

	if trip.Status != model.TripStatusActive {
		return model.ErrTripNotActive
	}

	booking, err := s.booking.GetBooking(ctx, trip.BookingID)
	if err != nil {
		return err
	}

	streamCtx, cancelStream := context.WithCancel(ctx)
	defer cancelStream()

	go func() {
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-streamCtx.Done():
				return
			case <-ticker.C:
				current, err := s.tripRepo.GetByID(streamCtx, tripID)
				if err != nil || current.Status != model.TripStatusActive {
					cancelStream()
					return
				}
			}
		}
	}()

	err = s.telematics.StreamTelemetry(streamCtx, trip.CarID, func(t model.CarTelemetry) error {
		current, err := s.tripRepo.GetByID(streamCtx, tripID)
		if err != nil {
			log.Error("repo: polling trip status", pkglog.Err(err))
			return err
		}
		if current.Status != model.TripStatusActive {
			return io.EOF
		}

		elapsedSeconds := int64(t.RecordedAt.Sub(trip.StartedAt).Seconds())
		distanceKM := float64(t.MileageKM - trip.StartMileageKM)
		currentCost := calculateCurrentCost(booking.PricingSnapshot, booking.CommittedPeriods, elapsedSeconds, distanceKM)

		return send(model.TripLiveFeed{
			ElapsedSeconds:     elapsedSeconds,
			CurrentCostTenge:   currentCost,
			DistanceTraveledKM: distanceKM,
		})
	})
	if streamCtx.Err() != nil {
		return io.EOF
	}
	return err
}

func isPrivileged(roles []sharedmodel.Role) bool {
	for _, r := range roles {
		if r == sharedmodel.RoleAdmin || r == sharedmodel.RoleBookingManager {
			return true
		}
	}
	return false
}

func tripFilter(filter validation.TripFilter) model.TripFilter {
	if filter.Pagination == nil {
		filter.Pagination = sharedvalidation.DefaultPagination()
	}
	mf := model.TripFilter{
		UserID: filter.UserID,
		CarID:  filter.CarID,
	}
	if filter.Status != nil {
		s := model.TripStatus(*filter.Status)
		mf.Status = &s
	}
	if filter.TimeRange != nil {
		tr := sharedmodel.TimeRange{}
		if filter.TimeRange.From != nil {
			tr.From = *filter.TimeRange.From
		}
		if filter.TimeRange.To != nil {
			tr.To = *filter.TimeRange.To
		}
		mf.TimeRange = &tr
	}
	mf.Pagination = &sharedmodel.Pagination{
		Limit:  filter.Pagination.Limit,
		Offset: filter.Pagination.Offset,
	}
	return mf
}

func tripStatusHistoryFilter(f validation.TripStatusHistoryFilter) model.TripStatusReadingFilter {
	if f.Pagination == nil {
		f.Pagination = sharedvalidation.DefaultPagination()
	}
	filter := model.TripStatusReadingFilter{
		TripID:     f.TripID,
		Pagination: &sharedmodel.Pagination{Limit: f.Pagination.Limit, Offset: f.Pagination.Offset},
	}
	if f.TimeRange != nil && f.TimeRange.From != nil && f.TimeRange.To != nil {
		filter.TimeRange = &sharedmodel.TimeRange{From: *f.TimeRange.From, To: *f.TimeRange.To}
	}
	return filter
}

func ptr[T any](v T) *T {
	return &v
}
