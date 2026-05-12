package nats

import (
	"context"
	"log/slog"

	"github.com/nats-io/nats.go"
	eventuserpb "github.com/sorawaslocked/car-rental-protos/gen/event/user"
	"github.com/sorawaslocked/car-rental-user-service/internal/model"
	pkglog "github.com/sorawaslocked/car-rental-user-service/internal/pkg/log"
	"google.golang.org/protobuf/proto"
)

const (
	subjectUserCreated = "user.created"
	subjectUserUpdated = "user.updated"
	subjectUserDeleted = "user.deleted"
)

type Publisher struct {
	log  *slog.Logger
	conn *nats.Conn
}

func NewPublisher(log *slog.Logger, conn *nats.Conn) *Publisher {
	return &Publisher{
		log:  pkglog.WithComponent(log, "adapter.nats.Publisher"),
		conn: conn,
	}
}

func (p *Publisher) PublishUserCreated(_ context.Context, id string) error {
	return p.publish(subjectUserCreated, &eventuserpb.CreateEvent{Id: id})
}

func (p *Publisher) PublishUserUpdated(_ context.Context, id string, isSecurityUpdate bool) error {
	return p.publish(subjectUserUpdated, &eventuserpb.UpdateEvent{Id: id, IsSecurityUpdate: isSecurityUpdate})
}

func (p *Publisher) PublishUserDeleted(_ context.Context, id string) error {
	return p.publish(subjectUserDeleted, &eventuserpb.DeleteEvent{Id: id})
}

func (p *Publisher) publish(subject string, msg proto.Message) error {
	data, err := proto.Marshal(msg)
	if err != nil {
		p.log.Error("marshalling event", slog.String("subject", subject), pkglog.Err(err))
		return model.ErrNats
	}

	if err := p.conn.Publish(subject, data); err != nil {
		p.log.Error("publishing event", slog.String("subject", subject), pkglog.Err(err))
		return model.ErrNats
	}

	return nil
}
