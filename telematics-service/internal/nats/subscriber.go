package nats

import (
	"log/slog"

	natsgo "github.com/nats-io/nats.go"
	"google.golang.org/protobuf/proto"

	tripevent "github.com/sorawaslocked/car-rental-protos/gen/event/trip"
	"github.com/sorawaslocked/car-rental-telematics/internal/config"
	"github.com/sorawaslocked/car-rental-telematics/internal/service"
)

// Subscriber listens on NATS for trip lifecycle events and drives the
// simulation service accordingly.
type Subscriber struct {
	conn   *natsgo.Conn
	simSvc *service.SimulationService
	cfg    *config.Config
}

func NewSubscriber(conn *natsgo.Conn, simSvc *service.SimulationService, cfg *config.Config) *Subscriber {
	return &Subscriber{conn: conn, simSvc: simSvc, cfg: cfg}
}

// Subscribe registers handlers for TripStartedEvent and TripEndedEvent.
func (s *Subscriber) Subscribe() error {
	if _, err := s.conn.Subscribe(s.cfg.TripStartedSubject, s.handleTripStarted); err != nil {
		return err
	}
	if _, err := s.conn.Subscribe(s.cfg.TripEndedSubject, s.handleTripEnded); err != nil {
		return err
	}
	slog.Info("subscribed to NATS trip events",
		"started_subject", s.cfg.TripStartedSubject,
		"ended_subject", s.cfg.TripEndedSubject,
	)
	return nil
}

func (s *Subscriber) handleTripStarted(msg *natsgo.Msg) {
	var ev tripevent.TripStartedEvent
	if err := proto.Unmarshal(msg.Data, &ev); err != nil {
		slog.Warn("failed to unmarshal TripStartedEvent", "error", err)
		return
	}
	slog.Info("trip started", "car_id", ev.CarId, "trip_id", ev.TripId)
	// Simulation is started by the gRPC client via StreamCarTelematicsEvents,
	// which provides the initial car state (location, fuel, battery, odometer).
}

func (s *Subscriber) handleTripEnded(msg *natsgo.Msg) {
	var ev tripevent.TripEndedEvent
	if err := proto.Unmarshal(msg.Data, &ev); err != nil {
		slog.Warn("failed to unmarshal TripEndedEvent", "error", err)
		return
	}
	slog.Info("trip ended — stopping simulation", "car_id", ev.CarId, "trip_id", ev.TripId)
	s.simSvc.StopSimulation(ev.CarId)
}
