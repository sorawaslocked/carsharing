package service

import (
	"context"
	"io"
	"log/slog"
	"time"

	sharedmodel "carsharing/shared/model"
	pkglog "carsharing/shared/pkg/log"
	pkgutils "carsharing/shared/pkg/utils"
	"carsharing/trip-service/internal/model"

	"github.com/google/uuid"
)

type TripService struct {
	log         *slog.Logger
	tripRepo    TripRepository
	summaryRepo TripSummaryRepository
	statusRepo  TripStatusReadingRepository
	booking     BookingClient
	telematics  TelematicsClient
	publisher   EventPublisher
}

func NewTripService(
	log *slog.Logger,
	tripRepo TripRepository,
	summaryRepo TripSummaryRepository,
	statusRepo TripStatusReadingRepository,
	booking BookingClient,
	telematics TelematicsClient,
	publisher EventPublisher,
) *TripService {
	return &TripService{
		log:         pkglog.WithComponent(log, "service.TripService"),
		tripRepo:    tripRepo,
		summaryRepo: summaryRepo,
		statusRepo:  statusRepo,
		booking:     booking,
		telematics:  telematics,
		publisher:   publisher,
	}
}

func (s *TripService) StartTrip(ctx context.Context, bookingID string) (string, error) {
	md := pkgutils.MetadataFromCtx(ctx)
	log := pkglog.WithMethod(s.log, "StartTrip")
	log = pkglog.WithMetadata(log, md)

	booking, err := s.booking.GetBooking(ctx, bookingID)
	if err != nil {
		return "", err
	}

	if !canAccess(md, booking.UserID) {
		return "", model.ErrInsufficientPermissions
	}

	if booking.Status != "created" {
		return "", model.ErrBookingNotCreated
	}

	telemetry, err := s.telematics.GetLatestTelemetry(ctx, booking.CarID)
	if err != nil {
		log.Error("failed to get telemetry", pkglog.Err(err))
		return "", err
	}

	now := time.Now()
	tripID := uuid.New().String()

	trip, err := s.tripRepo.Create(ctx, model.TripCreate{
		ID:        tripID,
		BookingID: bookingID,
		UserID:    booking.UserID,
		CarID:     booking.CarID,
		Status:    model.TripStatusActive,
		StartedAt: now,
		StartLocation: model.Location{
			Latitude:  telemetry.Location.Latitude,
			Longitude: telemetry.Location.Longitude,
		},
		StartMileageKM: telemetry.OdometerKM,
		StartFuelLevel: telemetry.FuelLevel,
	})
	if err != nil {
		log.Error("failed to create trip", pkglog.Err(err))
		return "", err
	}

	actorID := booking.UserID
	_, err = s.statusRepo.Create(ctx, model.TripStatusReadingCreate{
		TripID:     tripID,
		FromStatus: model.TripStatus(""),
		ToStatus:   model.TripStatusActive,
		ActorType:  model.ActorTypeUser,
		ActorID:    &actorID,
		ChangedAt:  now,
	})
	if err != nil {
		log.Error("failed to create status reading", pkglog.Err(err))
		return "", err
	}

	if err = s.publisher.PublishTripStarted(ctx, trip); err != nil {
		log.Error("failed to publish trip started event", pkglog.Err(err))
	}

	return tripID, nil
}

func (s *TripService) GetTrip(ctx context.Context, id string) (model.Trip, error) {
	trip, err := s.tripRepo.GetByID(ctx, id)
	if err != nil {
		return model.Trip{}, err
	}
	if !canAccess(pkgutils.MetadataFromCtx(ctx), trip.UserID) {
		return model.Trip{}, model.ErrInsufficientPermissions
	}
	return trip, nil
}

func (s *TripService) ListTrips(ctx context.Context, filter model.TripFilter) ([]model.Trip, error) {
	md := pkgutils.MetadataFromCtx(ctx)
	for _, r := range md.UserRoles {
		if r == sharedmodel.RoleAdmin || r == sharedmodel.RoleBookingManager {
			return s.tripRepo.List(ctx, filter)
		}
	}
	filter.UserID = md.UserID
	return s.tripRepo.List(ctx, filter)
}

func (s *TripService) EndTrip(ctx context.Context, id string) error {
	md := pkgutils.MetadataFromCtx(ctx)
	log := pkglog.WithMethod(s.log, "EndTrip")
	log = pkglog.WithMetadata(log, md)

	trip, err := s.tripRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if !canAccess(md, trip.UserID) {
		return model.ErrInsufficientPermissions
	}

	if !trip.Status.CanTransitionTo(model.TripStatusCompleted) {
		return model.ErrInvalidStatusTransition
	}

	telemetry, err := s.telematics.GetLatestTelemetry(ctx, trip.CarID)
	if err != nil {
		log.Error("failed to get telemetry", pkglog.Err(err))
		return err
	}

	booking, err := s.booking.GetBooking(ctx, trip.BookingID)
	if err != nil {
		log.Error("failed to get booking", pkglog.Err(err))
		return err
	}

	now := time.Now()
	durationSeconds := int64(now.Sub(trip.StartedAt).Seconds())
	distanceKM := float64(telemetry.OdometerKM - trip.StartMileageKM)

	baseCost, distCost, overtimeCost := calculateCosts(booking.PricingSnapshot, booking.CommittedPeriods, durationSeconds, distanceKM)
	totalCost := baseCost + distCost + overtimeCost

	endLocation := model.Location{
		Latitude:  telemetry.Location.Latitude,
		Longitude: telemetry.Location.Longitude,
	}
	endMileage := telemetry.OdometerKM

	updatedTrip, err := s.tripRepo.Update(ctx, id, model.TripUpdate{
		Status:             ptr(model.TripStatusCompleted),
		EndedAt:            &now,
		EndLocation:        &endLocation,
		EndMileageKM:       &endMileage,
		EndFuelLevel:       telemetry.FuelLevel,
		DistanceTraveledKM: ptr(distanceKM),
		DurationSeconds:    ptr(durationSeconds),
		FinalCostTenge:     ptr(totalCost),
		UpdatedAt:          now,
	})
	if err != nil {
		log.Error("failed to update trip", pkglog.Err(err))
		return err
	}

	actorID := trip.UserID
	_, err = s.statusRepo.Create(ctx, model.TripStatusReadingCreate{
		TripID:     id,
		FromStatus: model.TripStatusActive,
		ToStatus:   model.TripStatusCompleted,
		ActorType:  model.ActorTypeUser,
		ActorID:    &actorID,
		ChangedAt:  now,
	})
	if err != nil {
		log.Error("failed to create status reading", pkglog.Err(err))
		return err
	}

	_, err = s.summaryRepo.Create(ctx, model.TripSummaryCreate{
		TripID:             id,
		BookingID:          trip.BookingID,
		StartedAt:          trip.StartedAt,
		EndedAt:            now,
		DurationSeconds:    durationSeconds,
		DistanceTraveledKM: distanceKM,
		PricingSnapshot:    booking.PricingSnapshot,
		BaseCostTenge:      baseCost,
		DistanceCostTenge:  distCost,
		OvertimeCostTenge:  overtimeCost,
		TotalCostTenge:     totalCost,
	})
	if err != nil {
		log.Error("failed to create trip summary", pkglog.Err(err))
		return err
	}

	if err = s.publisher.PublishTripEnded(ctx, updatedTrip); err != nil {
		log.Error("failed to publish trip ended event", pkglog.Err(err))
	}

	return nil
}

func (s *TripService) CancelTrip(ctx context.Context, id string, reason *string) error {
	md := pkgutils.MetadataFromCtx(ctx)
	log := pkglog.WithMethod(s.log, "CancelTrip")
	log = pkglog.WithMetadata(log, md)

	trip, err := s.tripRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if !canAccess(md, trip.UserID) {
		return model.ErrInsufficientPermissions
	}

	if !trip.Status.CanTransitionTo(model.TripStatusCancelled) {
		return model.ErrInvalidStatusTransition
	}

	now := time.Now()

	updatedTrip, err := s.tripRepo.Update(ctx, id, model.TripUpdate{
		Status:       ptr(model.TripStatusCancelled),
		CancelReason: reason,
		UpdatedAt:    now,
	})
	if err != nil {
		log.Error("failed to update trip", pkglog.Err(err))
		return err
	}

	actorID := trip.UserID
	_, err = s.statusRepo.Create(ctx, model.TripStatusReadingCreate{
		TripID:     id,
		FromStatus: model.TripStatusActive,
		ToStatus:   model.TripStatusCancelled,
		ActorType:  model.ActorTypeUser,
		ActorID:    &actorID,
		Reason:     reason,
		ChangedAt:  now,
	})
	if err != nil {
		log.Error("failed to create status reading", pkglog.Err(err))
		return err
	}

	if err = s.publisher.PublishTripCancelled(ctx, updatedTrip); err != nil {
		log.Error("failed to publish trip cancelled event", pkglog.Err(err))
	}

	return nil
}

func (s *TripService) GetTripSummary(ctx context.Context, tripID string) (model.TripSummary, error) {
	trip, err := s.tripRepo.GetByID(ctx, tripID)
	if err != nil {
		return model.TripSummary{}, err
	}
	if !canAccess(pkgutils.MetadataFromCtx(ctx), trip.UserID) {
		return model.TripSummary{}, model.ErrInsufficientPermissions
	}
	return s.summaryRepo.GetByTripID(ctx, tripID)
}

func (s *TripService) GetTripStatusHistory(ctx context.Context, filter model.TripStatusReadingFilter) ([]model.TripStatusReading, error) {
	trip, err := s.tripRepo.GetByID(ctx, filter.TripID)
	if err != nil {
		return nil, err
	}
	if !canAccess(pkgutils.MetadataFromCtx(ctx), trip.UserID) {
		return nil, model.ErrInsufficientPermissions
	}
	return s.statusRepo.List(ctx, filter)
}

// StreamTripLiveFeed calls send for each telemetry event while the trip remains active.
// It returns io.EOF when the trip ends normally; any other error indicates a failure.
func (s *TripService) StreamTripLiveFeed(ctx context.Context, tripID string, send func(model.TripLiveFeed) error) error {
	log := pkglog.WithMethod(s.log, "StreamTripLiveFeed")
	log = pkglog.WithMetadata(log, pkgutils.MetadataFromCtx(ctx))

	trip, err := s.tripRepo.GetByID(ctx, tripID)
	if err != nil {
		return err
	}
	if trip.Status != model.TripStatusActive {
		return model.ErrTripNotActive
	}
	if !canAccess(pkgutils.MetadataFromCtx(ctx), trip.UserID) {
		return model.ErrInsufficientPermissions
	}

	booking, err := s.booking.GetBooking(ctx, trip.BookingID)
	if err != nil {
		return err
	}

	return s.telematics.StreamTelemetry(ctx, trip.CarID, func(t model.CarTelemetry) error {
		current, err := s.tripRepo.GetByID(ctx, tripID)
		if err != nil {
			log.Error("failed to poll trip status", pkglog.Err(err))
			return err
		}
		if current.Status != model.TripStatusActive {
			return io.EOF
		}

		elapsedSeconds := int64(t.RecordedAt.Sub(trip.StartedAt).Seconds())
		distanceKM := float64(t.OdometerKM - trip.StartMileageKM)
		currentCost := calculateCurrentCost(booking.PricingSnapshot, booking.CommittedPeriods, elapsedSeconds, distanceKM)

		return send(model.TripLiveFeed{
			ElapsedSeconds:     elapsedSeconds,
			CurrentCostTenge:   currentCost,
			DistanceTraveledKM: distanceKM,
		})
	})
}

func canAccess(md pkgutils.Metadata, ownerID string) bool {
	if md.UserID != nil && *md.UserID == ownerID {
		return true
	}
	for _, r := range md.UserRoles {
		if r == sharedmodel.RoleAdmin || r == sharedmodel.RoleBookingManager {
			return true
		}
	}
	return false
}

func ptr[T any](v T) *T {
	return &v
}
