package service

import (
	"context"
	"log/slog"
	"sync"
	"sync/atomic"
	"time"

	"carsharing/car-service/internal/model"
	sharedmodel "carsharing/shared/model"
	pkglog "carsharing/shared/pkg/log"

	"github.com/go-playground/validator/v10"
)

const telemetryReconnectDelay = 5 * time.Second

type TelemetryService struct {
	log      *slog.Logger
	validate *validator.Validate

	streamClient         TelemetryStreamClient
	telemetryReadingRepo TelemetryReadingRepository
	carRepo              CarRepository

	mu  sync.Mutex
	ctx context.Context
	wg  sync.WaitGroup

	subsMu sync.RWMutex
	subs   map[string]map[uint64]chan model.TelemetryUpdate
	subSeq atomic.Uint64
}

func NewTelemetryService(
	log *slog.Logger,
	validate *validator.Validate,
	streamClient TelemetryStreamClient,
	telemetryReadingRepo TelemetryReadingRepository,
	carRepo CarRepository,
) *TelemetryService {
	return &TelemetryService{
		log:                  pkglog.WithComponent(log, "service.TelemetryService"),
		validate:             validate,
		streamClient:         streamClient,
		telemetryReadingRepo: telemetryReadingRepo,
		carRepo:              carRepo,
		subs:                 make(map[string]map[uint64]chan model.TelemetryUpdate),
	}
}

// Start loads all cars and subscribes to each car's telemetry stream in a
// dedicated goroutine. Goroutines run until ctx is cancelled.
func (s *TelemetryService) Start(ctx context.Context) error {
	log := pkglog.WithMethod(s.log, "Start")

	cars, err := s.carRepo.Find(ctx, model.CarFilter{})
	if err != nil {
		log.Error("repo: listing cars", pkglog.Err(err))
		return err
	}

	s.mu.Lock()
	s.ctx = ctx
	s.mu.Unlock()

	log.Info("subscribing to telemetry streams", slog.Int("cars", len(cars)))

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
func (s *TelemetryService) Stop() {
	s.wg.Wait()
}

// OnCarCreated starts a telemetry stream goroutine for a newly created car.
func (s *TelemetryService) OnCarCreated(car model.Car) {
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

// SubscribeUpdates registers a listener for live telemetry updates for carID.
// Updates are fanned out from the single background stream maintained by TelemetryService.
// The returned unsubscribe func must be called when the caller is done.
func (s *TelemetryService) SubscribeUpdates(carID string) (<-chan model.TelemetryUpdate, func()) {
	id := s.subSeq.Add(1)
	ch := make(chan model.TelemetryUpdate, 16)

	s.subsMu.Lock()
	if s.subs[carID] == nil {
		s.subs[carID] = make(map[uint64]chan model.TelemetryUpdate)
	}
	s.subs[carID][id] = ch
	s.subsMu.Unlock()

	return ch, func() {
		s.subsMu.Lock()
		if m, ok := s.subs[carID]; ok {
			delete(m, id)
			if len(m) == 0 {
				delete(s.subs, carID)
			}
		}
		s.subsMu.Unlock()
	}
}

func (s *TelemetryService) fanOut(carID string, update model.TelemetryUpdate) {
	s.subsMu.RLock()
	defer s.subsMu.RUnlock()
	for _, ch := range s.subs[carID] {
		select {
		case ch <- update:
		default:
		}
	}
}

func (s *TelemetryService) subscribeToCarStream(ctx context.Context, car model.Car) {
	log := pkglog.WithMethod(s.log, "subscribeToCarStream")
	log = log.With(slog.String("carID", car.ID))

	for {
		ch, err := s.streamClient.Subscribe(ctx, car)
		if err != nil {
			if ctx.Err() != nil {
				return
			}
			log.Error("failed to subscribe to telemetry stream", pkglog.Err(err))
			select {
			case <-ctx.Done():
				return
			case <-time.After(telemetryReconnectDelay):
			}
			continue
		}

		for update := range ch {
			if err := s.applyUpdate(ctx, log, update); err != nil {
				log.Error("failed to apply telemetry update", pkglog.Err(err))
			}
			s.fanOut(update.CarID, update)
		}

		if ctx.Err() != nil {
			return
		}

		log.Info("telemetry stream closed, reconnecting",
			slog.Duration("in", telemetryReconnectDelay),
		)
		select {
		case <-ctx.Done():
			return
		case <-time.After(telemetryReconnectDelay):
		}
	}
}

func (s *TelemetryService) applyUpdate(ctx context.Context, log *slog.Logger, update model.TelemetryUpdate) error {
	current, err := s.carRepo.FindByID(ctx, update.CarID)
	if err != nil {
		log.Error("repo: finding car by id", pkglog.Err(err))
		return err
	}

	if update.MileageKM < current.MileageKM {
		log.Info("rejected telemetry update: mileage regression",
			slog.Int64("current", current.MileageKM),
			slog.Int64("incoming", update.MileageKM),
		)
		return model.ErrMileageRegression
	}

	loc := sharedmodel.Location{Latitude: update.Latitude, Longitude: update.Longitude}
	mileage := update.MileageKM

	if err = s.telemetryReadingRepo.Insert(ctx, model.TelemetryReading{
		CarID:        update.CarID,
		Location:     &loc,
		FuelPct:      update.FuelLevel,
		BatteryLevel: update.BatteryLevel,
		MileageKM:    &mileage,
		ActorType:    sharedmodel.ActorTypeTelemetry,
		RecordedAt:   update.RecordedAt,
	}); err != nil {
		log.Error("repo: inserting telemetry event", pkglog.Err(err))
		return err
	}

	now := time.Now()
	if err = s.carRepo.Update(ctx, update.CarID, model.CarUpdate{
		MileageKM:    &update.MileageKM,
		FuelLevel:    update.FuelLevel,
		BatteryLevel: update.BatteryLevel,
		Location:     &loc,
		LastSeenAt:   &now,
		UpdatedAt:    now,
	}); err != nil {
		log.Error("repo: updating car", pkglog.Err(err))
		return err
	}

	log.Info("telemetry update applied",
		slog.Int64("mileageKM", update.MileageKM),
		slog.Time("recordedAt", update.RecordedAt),
	)

	return nil
}
