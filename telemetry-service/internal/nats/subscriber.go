package nats

import (
	"fmt"
	"log/slog"

	"carsharing/telematics-service/internal/service"
	"github.com/nats-io/nats.go"
	tripevent "github.com/sorawaslocked/car-rental-protos/gen/event/trip"
	"google.golang.org/protobuf/proto"
)

// TripSubscriber listens for trip lifecycle events and drives the simulation state machine.
type TripSubscriber struct {
	conn             *nats.Conn
	simSvc           *service.SimulationService
	subs             []*nats.Subscription
	startedSubject   string
	endedSubject     string
	cancelledSubject string
}

func NewTripSubscriber(url, startedSubject, endedSubject, cancelledSubject string, simSvc *service.SimulationService) (*TripSubscriber, error) {
	conn, err := nats.Connect(url)
	if err != nil {
		return nil, fmt.Errorf("connect nats: %w", err)
	}
	return &TripSubscriber{
		conn:             conn,
		simSvc:           simSvc,
		startedSubject:   startedSubject,
		endedSubject:     endedSubject,
		cancelledSubject: cancelledSubject,
	}, nil
}

func (s *TripSubscriber) Subscribe() error {
	startSub, err := s.conn.Subscribe(s.startedSubject, func(msg *nats.Msg) {
		var ev tripevent.TripStartedEvent
		if err := proto.Unmarshal(msg.Data, &ev); err != nil {
			slog.Warn("invalid trip started payload", "subject", s.startedSubject, "error", err)
			return
		}
		slog.Info("trip started event received", "car_id", ev.GetCarId(), "trip_id", ev.GetTripId())
		s.simSvc.StartTrip(ev.GetCarId())
	})
	if err != nil {
		return fmt.Errorf("subscribe %s: %w", s.startedSubject, err)
	}

	endSub, err := s.conn.Subscribe(s.endedSubject, func(msg *nats.Msg) {
		var ev tripevent.TripEndedEvent
		if err := proto.Unmarshal(msg.Data, &ev); err != nil {
			slog.Warn("invalid trip ended payload", "subject", s.endedSubject, "error", err)
			return
		}
		slog.Info("trip ended event received", "car_id", ev.GetCarId(), "trip_id", ev.GetTripId())
		s.simSvc.EndTrip(ev.GetCarId())
	})
	if err != nil {
		startSub.Unsubscribe()
		return fmt.Errorf("subscribe %s: %w", s.endedSubject, err)
	}

	cancelSub, err := s.conn.Subscribe(s.cancelledSubject, func(msg *nats.Msg) {
		var ev tripevent.TripCancelledEvent
		if err := proto.Unmarshal(msg.Data, &ev); err != nil {
			slog.Warn("invalid trip cancelled payload", "subject", s.cancelledSubject, "error", err)
			return
		}
		slog.Info("trip cancelled event received", "car_id", ev.GetCarId(), "trip_id", ev.GetTripId(), "reason", ev.GetReason())
		s.simSvc.EndTrip(ev.GetCarId())
	})
	if err != nil {
		startSub.Unsubscribe()
		endSub.Unsubscribe()
		return fmt.Errorf("subscribe %s: %w", s.cancelledSubject, err)
	}

	s.subs = []*nats.Subscription{startSub, endSub, cancelSub}
	slog.Info("NATS trip subscriptions active",
		"started", s.startedSubject,
		"ended", s.endedSubject,
		"cancelled", s.cancelledSubject,
	)
	return nil
}

func (s *TripSubscriber) Close() {
	for _, sub := range s.subs {
		_ = sub.Unsubscribe()
	}
	s.conn.Close()
}
