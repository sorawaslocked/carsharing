package service

import (
	"context"
	"log/slog"
	"sync"
	"time"

	"carsharing/car-service/internal/model"
	pkglog "carsharing/shared/pkg/log"
)

const telematicsReconnectDelay = 5 * time.Second

type TelematicsService struct {
	streamClient   TelematicsStreamClient
	telematicsRepo TelematicsRepository
	carRepo        CarRepository

	mu  sync.Mutex
	ctx context.Context
	wg  sync.WaitGroup

	log *slog.Logger
}

func NewTelematicsService(
	streamClient TelematicsStreamClient,
	telematicsRepo TelematicsRepository,
	carRepo CarRepository,
	log *slog.Logger,
) *TelematicsService {
	s := &TelematicsService{
		streamClient:   streamClient,
		telematicsRepo: telematicsRepo,
		carRepo:        carRepo,
	}

	s.log = pkglog.WithComponent(log, "service.TelematicsService")

	return s
}

// Start loads all cars and subscribes to each car's telemetry stream in a
// dedicated goroutine. Goroutines run until ctx is cancelled.
func (s *TelematicsService) Start(ctx context.Context) error {
	logger := pkglog.WithMethod(s.log, "Start")

	cars, err := s.carRepo.Find(ctx, model.CarFilter{})
	if err != nil {
		return handleError(logger, err)
	}

	s.mu.Lock()
	s.ctx = ctx
	s.mu.Unlock()

	logger.Info("subscribing to telemetry streams", slog.Int("cars", len(cars)))

	for _, car := range cars {
		s.wg.Add(1)
		go func(c model.Car) {
			defer s.wg.Done()
			s.subscribeToCarStream(ctx, c)
		}(car)
	}

	return nil
}

// Stop waits for all telemetry goroutines to finish after the context is cancelled.
func (s *TelematicsService) Stop() {
	s.wg.Wait()
}

// OnCarCreated starts a telemetry stream goroutine for a newly created car.
func (s *TelematicsService) OnCarCreated(car model.Car) {
	s.mu.Lock()
	ctx := s.ctx
	s.mu.Unlock()

	if ctx == nil {
		s.log.Warn("OnCarCreated called before Start",
			slog.String("carID", car.ID),
		)
		return
	}

	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		s.subscribeToCarStream(ctx, car)
	}()
}

// SubscribeCarStream opens a live telemetry stream for a single car and returns
// a channel of updates. Used by the gRPC streaming handler.
func (s *TelematicsService) SubscribeCarStream(ctx context.Context, carID string) (<-chan model.TelematicsUpdate, error) {
	logger := pkglog.WithMethod(s.log, "SubscribeCarStream")

	car, err := s.carRepo.FindByID(ctx, carID)
	if err != nil {
		return nil, handleError(logger, err)
	}

	return s.streamClient.Subscribe(ctx, car)
}

func (s *TelematicsService) subscribeToCarStream(ctx context.Context, car model.Car) {
	logger := pkglog.WithMethod(s.log, "subscribeToCarStream")
	logger = logger.With(slog.String("carID", car.ID))

	for {
		ch, err := s.streamClient.Subscribe(ctx, car)
		if err != nil {
			if ctx.Err() != nil {
				return
			}
			logger.Error("failed to subscribe to telemetry stream", pkglog.Err(err))
			select {
			case <-ctx.Done():
				return
			case <-time.After(telematicsReconnectDelay):
			}
			continue
		}

		for update := range ch {
			if err := s.applyUpdate(ctx, logger, update); err != nil {
				logger.Error("failed to apply telemetry update", pkglog.Err(err))
			}
		}

		if ctx.Err() != nil {
			return
		}

		logger.Info("telemetry stream closed, reconnecting",
			slog.Duration("in", telematicsReconnectDelay),
		)
		select {
		case <-ctx.Done():
			return
		case <-time.After(telematicsReconnectDelay):
		}
	}
}

func (s *TelematicsService) applyUpdate(ctx context.Context, logger *slog.Logger, update model.TelematicsUpdate) error {
	current, err := s.carRepo.FindByID(ctx, update.CarID)
	if err != nil {
		return handleError(logger, err)
	}

	if update.OdometerKM < current.MileageKM {
		logger.Info("rejected telemetry update: odometer regression",
			slog.Int64("current", current.MileageKM),
			slog.Int64("incoming", update.OdometerKM),
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
		ActorType:    string(model.CarStatusActorTelematics),
		RecordedAt:   update.RecordedAt,
		ReceivedAt:   now,
	}

	if err = s.telematicsRepo.InsertEvent(ctx, event); err != nil {
		return handleError(logger, err)
	}

	if err = s.carRepo.Update(ctx, update.CarID, model.CarUpdate{
		MileageKM:    &update.OdometerKM,
		FuelLevel:    update.FuelLevel,
		BatteryLevel: update.BatteryLevel,
		Location: &model.Location{
			Latitude:  update.Latitude,
			Longitude: update.Longitude,
		},
		LastSeenAt: &now,
		UpdatedAt:  now,
	}); err != nil {
		return handleError(logger, err)
	}

	logger.Info("telemetry update applied",
		slog.Int64("odometerKM", update.OdometerKM),
		slog.Time("recordedAt", update.RecordedAt),
	)

	return nil
}
