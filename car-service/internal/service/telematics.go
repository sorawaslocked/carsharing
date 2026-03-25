package service

import (
	"car-rental-car-service/internal/model"
	"context"
	"log/slog"
	"time"

	"github.com/go-playground/validator/v10"
)

type TelematicsService struct {
	telematicsRepo  TelematicsRepository
	telematicsQueue TelematicsQueue
	carRepo         CarRepository

	validate *validator.Validate
	log      *slog.Logger
}

func NewTelematicsService(
	telematicsRepo TelematicsRepository,
	telematicsQueue TelematicsQueue,
	carRepo CarRepository,
	log *slog.Logger,
) *TelematicsService {
	s := &TelematicsService{
		telematicsRepo:  telematicsRepo,
		telematicsQueue: telematicsQueue,
		carRepo:         carRepo,
	}

	s.log = log.With(
		slog.Group("src",
			slog.String("component", "TelematicsService"),
		),
	)

	return s
}

func (s *TelematicsService) ProcessNextUpdate(ctx context.Context) error {
	const method = "ProcessNextUpdate"

	logger := s.log.With(slog.Group("src", slog.String("method", method)))

	update, ack, err := s.telematicsQueue.Pop(ctx)
	if err != nil {
		return handleError(logger, err)
	}

	if err = s.applyUpdate(ctx, logger, update); err != nil {
		if nackErr := ack(err); nackErr != nil {
			logger.Error("failed to nack telemetry event",
				slog.String("carID", update.CarID),
				slog.String("error", nackErr.Error()),
			)
		}
		return err
	}

	if err = ack(nil); err != nil {
		return handleError(logger, err)
	}

	return nil
}

func (s *TelematicsService) applyUpdate(ctx context.Context, logger *slog.Logger, update model.TelematicsUpdate) error {
	carID := update.CarID
	carFilter := model.CarFilter{ID: &carID}

	current, err := s.carRepo.FindOne(ctx, carFilter)
	if err != nil {
		return handleError(logger, err)
	}

	if update.OdometerKM < current.MileageKM {
		logger.Info("rejected telemetry update due to odometer regression",
			slog.String("carID", carID),
			slog.Int64("currentMileageKM", current.MileageKM),
			slog.Int64("incomingOdometerKM", update.OdometerKM),
		)
		return ErrOdometerRegression
	}

	now := time.Now()

	event := model.CarTelematicsEvent{
		CarID:        update.CarID,
		Latitude:     update.Latitude,
		Longitude:    update.Longitude,
		FuelLevel:    update.FuelLevel,
		BatteryLevel: update.BatteryLevel,
		OdometerKM:   update.OdometerKM,
		RecordedAt:   update.RecordedAt,
		ReceivedAt:   now,
	}

	if err = s.telematicsRepo.InsertEvent(ctx, event); err != nil {
		return handleError(logger, err)
	}

	carUpdate := model.CarUpdate{
		MileageKM:    &update.OdometerKM,
		FuelLevel:    update.FuelLevel,
		BatteryLevel: update.BatteryLevel,
		Location: &model.Location{
			Latitude:  update.Latitude,
			Longitude: update.Longitude,
		},
		LastSeenAt: &now,
		UpdatedAt:  now,
	}

	if err = s.carRepo.Update(ctx, carFilter, carUpdate); err != nil {
		return handleError(logger, err)
	}

	logger.Info("telemetry event applied",
		slog.String("carID", update.CarID),
		slog.Int64("odometerKM", update.OdometerKM),
		slog.Time("recordedAt", update.RecordedAt),
	)

	return nil
}
