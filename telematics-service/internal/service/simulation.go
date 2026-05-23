package service

import (
	"context"
	"log/slog"
	"sync"
	"sync/atomic"
	"time"

	osrm "github.com/gojuno/go.osrm"
)

const (
	speedKmh           = 40.0
	fuelConsumPerKm    = float32(0.2)
	batteryConsumPerKm = float32(0.25)
)

// SimulationRequest carries the initial car state when registering a stream.
type SimulationRequest struct {
	CarId        string
	Latitude     float64
	Longitude    float64
	FuelLevel    *float32
	BatteryLevel *float32
	OdometerKm   int64
}

// SimulationUpdate is one telemetry snapshot emitted every interval while a trip is active.
type SimulationUpdate struct {
	CarId        string
	Latitude     float64
	Longitude    float64
	FuelLevel    *float32
	BatteryLevel *float32
	OdometerKm   int64
	RecordedAt   time.Time
}

type simEntry struct {
	cancel  context.CancelFunc
	token   uint64
	startCh chan struct{}
	stopCh  chan struct{}
}

// SimulationService manages per-car telemetry simulations.
type SimulationService struct {
	osrmClient  *osrm.OSRM
	osrmProfile string
	interval    time.Duration
	mu          sync.Mutex
	activeSims  map[string]*simEntry
	counter     atomic.Uint64
}

func NewSimulationService(osrmClient *osrm.OSRM, osrmProfile string, interval time.Duration) *SimulationService {
	return &SimulationService{
		osrmClient:  osrmClient,
		osrmProfile: osrmProfile,
		interval:    interval,
		activeSims:  make(map[string]*simEntry),
	}
}

// RegisterStream registers a car for idle streaming and returns a channel of updates.
// The car starts in the idle state; updates are only emitted while a trip is active.
// Call StartTrip / EndTrip to transition between states.
func (s *SimulationService) RegisterStream(ctx context.Context, req *SimulationRequest) (<-chan *SimulationUpdate, error) {
	simCtx, cancel := context.WithCancel(ctx)
	token := s.counter.Add(1)
	entry := &simEntry{
		cancel:  cancel,
		token:   token,
		startCh: make(chan struct{}, 1),
		stopCh:  make(chan struct{}, 1),
	}

	s.mu.Lock()
	if existing, ok := s.activeSims[req.CarId]; ok {
		existing.cancel()
	}
	s.activeSims[req.CarId] = entry
	s.mu.Unlock()

	updates := make(chan *SimulationUpdate, 1)
	go s.runStream(simCtx, req, token, entry.startCh, entry.stopCh, updates)
	return updates, nil
}

// StartTrip signals the car's simulation to begin movement.
func (s *SimulationService) StartTrip(carID string) {
	s.mu.Lock()
	entry, ok := s.activeSims[carID]
	s.mu.Unlock()
	if !ok {
		slog.Warn("StartTrip: no active stream", "car_id", carID)
		return
	}
	select {
	case entry.startCh <- struct{}{}:
	default:
	}
}

// EndTrip signals the car's simulation to stop movement and return to idle.
func (s *SimulationService) EndTrip(carID string) {
	s.mu.Lock()
	entry, ok := s.activeSims[carID]
	s.mu.Unlock()
	if !ok {
		slog.Warn("EndTrip: no active stream", "car_id", carID)
		return
	}
	select {
	case entry.stopCh <- struct{}{}:
	default:
	}
}

// UnregisterStream cancels and removes the stream for the given car.
func (s *SimulationService) UnregisterStream(carID string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if entry, ok := s.activeSims[carID]; ok {
		entry.cancel()
		delete(s.activeSims, carID)
	}
}

func (s *SimulationService) runStream(
	ctx context.Context,
	req *SimulationRequest,
	token uint64,
	startCh <-chan struct{},
	stopCh <-chan struct{},
	updates chan<- *SimulationUpdate,
) {
	defer close(updates)
	defer func() {
		s.mu.Lock()
		if entry, ok := s.activeSims[req.CarId]; ok && entry.token == token {
			delete(s.activeSims, req.CarId)
		}
		s.mu.Unlock()
	}()

	lat := req.Latitude
	lng := req.Longitude
	odometer := req.OdometerKm
	var subKmAccM float64

	var fuelLevel *float32
	if req.FuelLevel != nil {
		v := *req.FuelLevel
		fuelLevel = &v
	}
	var batteryLevel *float32
	if req.BatteryLevel != nil {
		v := *req.BatteryLevel
		batteryLevel = &v
	}

	slog.Info("stream registered (idle)", "car_id", req.CarId)

	for {
		// Idle: wait for a trip start signal or stream closure.
		select {
		case <-ctx.Done():
			slog.Info("stream closed (idle)", "car_id", req.CarId)
			return
		case <-startCh:
		}

		slog.Info("trip started, simulation running", "car_id", req.CarId)

		walker := s.fetchRoute(ctx, lat, lng)
		distPerTick := speedKmh * 1000.0 / 3600.0 * s.interval.Seconds()
		ticker := time.NewTicker(s.interval)

		tripActive := true
		for tripActive {
			select {
			case <-ctx.Done():
				ticker.Stop()
				slog.Info("stream closed (active trip)", "car_id", req.CarId)
				return

			case <-stopCh:
				ticker.Stop()
				slog.Info("trip ended, simulation idle", "car_id", req.CarId)
				tripActive = false

			case <-ticker.C:
				if walker.done() {
					walker = s.fetchRoute(ctx, lat, lng)
				}

				newLat, newLng := walker.advance(distPerTick)
				distM := haversineMeters(lat, lng, newLat, newLng)

				lat = newLat
				lng = newLng

				subKmAccM += distM
				kmDelta := int64(subKmAccM / 1000)
				subKmAccM -= float64(kmDelta) * 1000
				odometer += kmDelta

				distKm := float32(distM / 1000)

				if fuelLevel != nil {
					v := *fuelLevel - distKm*fuelConsumPerKm
					if v < 0 {
						v = 0
					}
					fuelLevel = &v
				}
				if batteryLevel != nil {
					v := *batteryLevel - distKm*batteryConsumPerKm
					if v < 0 {
						v = 0
					}
					batteryLevel = &v
				}

				update := &SimulationUpdate{
					CarId:        req.CarId,
					Latitude:     lat,
					Longitude:    lng,
					FuelLevel:    cloneFloat32Ptr(fuelLevel),
					BatteryLevel: cloneFloat32Ptr(batteryLevel),
					OdometerKm:   odometer,
					RecordedAt:   time.Now(),
				}

				select {
				case updates <- update:
				case <-ctx.Done():
					ticker.Stop()
					slog.Info("stream closed (active trip)", "car_id", req.CarId)
					return
				}
			}
		}
	}
}

func cloneFloat32Ptr(v *float32) *float32 {
	if v == nil {
		return nil
	}
	c := *v
	return &c
}
