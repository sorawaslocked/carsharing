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
	// speedKmh is the simulated driving speed.
	speedKmh = 40.0
	// distPerTickM is the distance covered in one 15-second tick at speedKmh.
	distPerTickM = speedKmh * 1000.0 / 3600.0 * 15.0 // ~166.7 m

	fuelConsumPerKm    = float32(0.2)  // % fuel consumed per km (0–100 scale)
	batteryConsumPerKm = float32(0.25) // % battery consumed per km (0–100 scale)
)

// SimulationRequest carries the initial car state for a new simulation.
type SimulationRequest struct {
	CarId        string
	Latitude     float64
	Longitude    float64
	FuelLevel    *float32
	BatteryLevel *float32
	OdometerKm   int64
}

// SimulationUpdate is one telemetry snapshot emitted every 15 seconds.
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
	cancel context.CancelFunc
	token  uint64
}

// SimulationService manages per-car telematics simulations.
type SimulationService struct {
	osrmClient  *osrm.OSRM
	osrmProfile string
	mu          sync.Mutex
	activeSims  map[string]*simEntry
	counter     atomic.Uint64
}

func NewSimulationService(osrmClient *osrm.OSRM, osrmProfile string) *SimulationService {
	return &SimulationService{
		osrmClient:  osrmClient,
		osrmProfile: osrmProfile,
		activeSims:  make(map[string]*simEntry),
	}
}

// StartSimulation begins a simulation for the given car and returns a channel of
// updates that will receive one event every 15 seconds until the simulation stops.
// A pre-existing simulation for the same carId is cancelled before starting the new one.
func (s *SimulationService) StartSimulation(ctx context.Context, req *SimulationRequest) (<-chan *SimulationUpdate, error) {
	simCtx, cancel := context.WithCancel(ctx)
	token := s.counter.Add(1)

	s.mu.Lock()
	if existing, ok := s.activeSims[req.CarId]; ok {
		existing.cancel()
	}
	s.activeSims[req.CarId] = &simEntry{cancel: cancel, token: token}
	s.mu.Unlock()

	updates := make(chan *SimulationUpdate, 1)
	go s.runSimulation(simCtx, req, token, updates)
	return updates, nil
}

// StopSimulation cancels a running simulation for the given car.
func (s *SimulationService) StopSimulation(carId string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if entry, ok := s.activeSims[carId]; ok {
		entry.cancel()
		delete(s.activeSims, carId)
	}
}

func (s *SimulationService) runSimulation(
	ctx context.Context,
	req *SimulationRequest,
	token uint64,
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
	var subKmAccM float64 // sub-km distance accumulator for odometer

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

	walker := s.fetchRoute(ctx, lat, lng)

	ticker := time.NewTicker(15 * time.Second)
	defer ticker.Stop()

	slog.Info("simulation started", "car_id", req.CarId, "lat", lat, "lng", lng)

	for {
		select {
		case <-ctx.Done():
			slog.Info("simulation stopped", "car_id", req.CarId)
			return
		case <-ticker.C:
			if walker.done() {
				walker = s.fetchRoute(ctx, lat, lng)
			}

			newLat, newLng := walker.advance(distPerTickM)
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
				slog.Info("simulation stopped", "car_id", req.CarId)
				return
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
