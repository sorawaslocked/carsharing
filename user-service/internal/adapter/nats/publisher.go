package nats

import (
	"context"
	"log/slog"

	"github.com/nats-io/nats.go"
	eventuserpb "github.com/sorawaslocked/car-rental-protos/gen/event/user"
	"github.com/sorawaslocked/car-rental-user-service/internal/model"
	pkglog "github.com/sorawaslocked/car-rental-user-service/internal/pkg/log"
	"github.com/sorawaslocked/car-rental-user-service/internal/pkg/utils"
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

func (p *Publisher) PublishUserCreated(ctx context.Context, id string) error {
	logger := pkglog.WithMetadata(pkglog.WithMethod(p.log, "PublishUserCreated"), utils.MetadataFromCtx(ctx))
	return p.publish(logger, subjectUserCreated, &eventuserpb.CreateEvent{Id: id})
}

func (p *Publisher) PublishUserUpdated(ctx context.Context, id string, isSecurityUpdate bool) error {
	logger := pkglog.WithMetadata(pkglog.WithMethod(p.log, "PublishUserUpdated"), utils.MetadataFromCtx(ctx))
	return p.publish(logger, subjectUserUpdated, &eventuserpb.UpdateEvent{Id: id, IsSecurityUpdate: isSecurityUpdate})
}

func (p *Publisher) PublishUserDeleted(ctx context.Context, id string) error {
	logger := pkglog.WithMetadata(pkglog.WithMethod(p.log, "PublishUserDeleted"), utils.MetadataFromCtx(ctx))
	return p.publish(logger, subjectUserDeleted, &eventuserpb.DeleteEvent{Id: id})
}

func (p *Publisher) publish(logger *slog.Logger, subject string, msg proto.Message) error {
	data, err := proto.Marshal(msg)
	if err != nil {
		logger.Error("marshalling event", slog.String("subject", subject), pkglog.Err(err))
		return model.ErrNats
	}

	if err := p.conn.Publish(subject, data); err != nil {
		logger.Error("publishing event", slog.String("subject", subject), pkglog.Err(err))
		return model.ErrNats
	}

	return nil
}
