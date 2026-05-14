package service

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/sorawaslocked/car-rental-trip-service/internal/model"
	pkglog "github.com/sorawaslocked/car-rental-trip-service/internal/pkg/log"
	"github.com/sorawaslocked/car-rental-trip-service/internal/pkg/utils"
)

var ErrMissingMetadata = errors.New("missing metadata")

type TripService struct {
	log           *slog.Logger
	tripRepo      TripRepository
	summaryRepo   TripSummaryRepository
	statusRepo    TripStatusReadingRepository
	bookingClient BookingClient
	telematics    TelematicsClient
	publisher     EventPublisher
}

func NewTripService(
	log *slog.Logger,
	tripRepo TripRepository,
	summaryRepo TripSummaryRepository,
	statusRepo TripStatusReadingRepository,
	bookingClient BookingClient,
	telematics TelematicsClient,
	publisher EventPublisher,
) *TripService {
	return &TripService{
		log:           pkglog.WithComponent(log, "service.TripService"),
		tripRepo:      tripRepo,
		summaryRepo:   summaryRepo,
		statusRepo:    statusRepo,
		bookingClient: bookingClient,
		telematics:    telematics,
		publisher:     publisher,
	}
}

func (s *TripService) StartTrip(ctx context.Context, bookingID string) (string, error) {
	log := pkglog.WithMethod(s.log, "StartTrip")

	if err := validateBookingID(bookingID); err != nil {
		return "", err
	}

	md := utils.MetadataFromCtx(ctx)
	if md.UserID == nil {
		return "", ErrMissingMetadata
	}
	userID := *md.UserID

	booking, err := s.bookingClient.GetBooking(ctx, bookingID)
	if err != nil {
		return "", err
	}

	if booking.UserID != userID && !isPrivileged(md.UserRoles) {
		return "", model.ErrInsufficientPermissions
	}
	if booking.Status != "reserved" {
		return "", model.ErrBookingNotReserved
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

	_, err = s.statusRepo.Create(ctx, model.TripStatusReadingCreate{
		TripID:     tripID,
		FromStatus: model.TripStatus(""),
		ToStatus:   model.TripStatusActive,
		ActorType:  model.ActorTypeUser,
		ActorID:    &userID,
		ChangedAt:  now,
	})
	if err != nil {
		log.Error("failed to create status reading", pkglog.Err(err))
		return "", err
	}

	if err = s.bookingClient.UpdateBookingStatus(ctx, bookingID, "active"); err != nil {
		log.Error("failed to update booking status", pkglog.Err(err))
		return "", err
	}

	if err = s.publisher.PublishTripStarted(ctx, trip); err != nil {
		log.Error("failed to publish trip started event", pkglog.Err(err))
	}

	return tripID, nil
}

func (s *TripService) GetTrip(ctx context.Context, id string) (model.Trip, error) {
	if err := validateID(id); err != nil {
		return model.Trip{}, err
	}
	return s.tripRepo.GetByID(ctx, id)
}

func (s *TripService) ListTrips(ctx context.Context, filter model.TripFilter) ([]model.Trip, error) {
	return s.tripRepo.List(ctx, filter)
}

func (s *TripService) EndTrip(ctx context.Context, id string) error {
	log := pkglog.WithMethod(s.log, "EndTrip")

	if err := validateID(id); err != nil {
		return err
	}

	md := utils.MetadataFromCtx(ctx)
	if md.UserID == nil {
		return ErrMissingMetadata
	}
	userID := *md.UserID

	trip, err := s.tripRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if trip.UserID != userID && !isPrivileged(md.UserRoles) {
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

	booking, err := s.bookingClient.GetBooking(ctx, trip.BookingID)
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

	_, err = s.statusRepo.Create(ctx, model.TripStatusReadingCreate{
		TripID:     id,
		FromStatus: model.TripStatusActive,
		ToStatus:   model.TripStatusCompleted,
		ActorType:  model.ActorTypeUser,
		ActorID:    &userID,
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

	if err = s.bookingClient.UpdateBookingStatus(ctx, trip.BookingID, "completed"); err != nil {
		log.Error("failed to update booking status", pkglog.Err(err))
		return err
	}

	if err = s.publisher.PublishTripEnded(ctx, updatedTrip); err != nil {
		log.Error("failed to publish trip ended event", pkglog.Err(err))
	}

	return nil
}

func (s *TripService) CancelTrip(ctx context.Context, id string, reason *string) error {
	log := pkglog.WithMethod(s.log, "CancelTrip")

	if err := validateID(id); err != nil {
		return err
	}

	md := utils.MetadataFromCtx(ctx)
	if md.UserID == nil {
		return ErrMissingMetadata
	}
	userID := *md.UserID

	trip, err := s.tripRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if trip.UserID != userID && !isPrivileged(md.UserRoles) {
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

	_, err = s.statusRepo.Create(ctx, model.TripStatusReadingCreate{
		TripID:     id,
		FromStatus: model.TripStatusActive,
		ToStatus:   model.TripStatusCancelled,
		ActorType:  model.ActorTypeUser,
		ActorID:    &userID,
		Reason:     reason,
		ChangedAt:  now,
	})
	if err != nil {
		log.Error("failed to create status reading", pkglog.Err(err))
		return err
	}

	if err = s.bookingClient.UpdateBookingStatus(ctx, trip.BookingID, "cancelled"); err != nil {
		log.Error("failed to update booking status", pkglog.Err(err))
		return err
	}

	if err = s.publisher.PublishTripCancelled(ctx, updatedTrip); err != nil {
		log.Error("failed to publish trip cancelled event", pkglog.Err(err))
	}

	return nil
}

func (s *TripService) GetTripSummary(ctx context.Context, tripID string) (model.TripSummary, error) {
	if err := validateID(tripID); err != nil {
		return model.TripSummary{}, err
	}
	return s.summaryRepo.GetByTripID(ctx, tripID)
}

func (s *TripService) GetTripStatusHistory(ctx context.Context, filter model.TripStatusReadingFilter) ([]model.TripStatusReading, error) {
	if err := validateID(filter.TripID); err != nil {
		return nil, err
	}
	return s.statusRepo.List(ctx, filter)
}

// StreamTripLiveFeed calls send for each telemetry event while the trip remains active.
// It returns io.EOF when the trip ends normally; any other error indicates a failure.
func (s *TripService) StreamTripLiveFeed(ctx context.Context, tripID string, send func(model.TripLiveFeed) error) error {
	if err := validateID(tripID); err != nil {
		return err
	}

	trip, err := s.tripRepo.GetByID(ctx, tripID)
	if err != nil {
		return err
	}
	if trip.Status != model.TripStatusActive {
		return model.ErrTripNotActive
	}

	booking, err := s.bookingClient.GetBooking(ctx, trip.BookingID)
	if err != nil {
		return err
	}

	return s.telematics.StreamTelemetry(ctx, trip.CarID, func(t model.CarTelemetry) error {
		current, err := s.tripRepo.GetByID(ctx, tripID)
		if err != nil {
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

func ptr[T any](v T) *T {
	return &v
}
