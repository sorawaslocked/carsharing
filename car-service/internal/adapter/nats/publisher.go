package nats

import (
	"context"
	"fmt"

	"github.com/nats-io/nats.go"
	eventpb "github.com/sorawaslocked/car-rental-protos/gen/event/car"
	"google.golang.org/protobuf/proto"
)

const subjectCarStatusUpdated = "car.status.updated"

type Publisher struct {
	conn *nats.Conn
}

func NewPublisher(conn *nats.Conn) *Publisher {
	return &Publisher{conn: conn}
}

func (p *Publisher) PublishCarStatusUpdated(_ context.Context, carID, fromStatus, toStatus string) error {
	event := &eventpb.CarStatusUpdatedEvent{
		CarId:      carID,
		FromStatus: fromStatus,
		ToStatus:   toStatus,
	}

	data, err := proto.Marshal(event)
	if err != nil {
		return fmt.Errorf("marshal car status updated event: %w", err)
	}

	return p.conn.Publish(subjectCarStatusUpdated, data)
}
