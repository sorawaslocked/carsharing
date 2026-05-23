package publisher

import (
	"context"
	"log/slog"

	pkglog "carsharing/shared/pkg/log"
	"carsharing/shared/pkg/utils"
	"carsharing/user-service/internal/model"

	eventuserpb "carsharing/protos/gen/event/user"
	"github.com/nats-io/nats.go"
	"google.golang.org/protobuf/proto"
)

const (
	subjectUserCreated = "user.created"
	subjectUserUpdated = "user.updated"
	subjectUserDeleted = "user.deleted"
)

type UserPublisher struct {
	log  *slog.Logger
	conn *nats.Conn
}

func NewUserPublisher(log *slog.Logger, conn *nats.Conn) *UserPublisher {
	return &UserPublisher{
		log:  pkglog.WithComponent(log, "adapter.nats.publisher.UserPublisher"),
		conn: conn,
	}
}

func (p *UserPublisher) PublishUserCreated(ctx context.Context, id string) error {
	log := pkglog.WithMetadata(pkglog.WithMethod(p.log, "PublishUserCreated"), utils.MetadataFromCtx(ctx))

	return p.publish(log, subjectUserCreated, &eventuserpb.UserCreatedEvent{Id: id})
}

func (p *UserPublisher) PublishUserUpdated(ctx context.Context, id string, isSecurityUpdate bool) error {
	log := pkglog.WithMetadata(pkglog.WithMethod(p.log, "PublishUserUpdated"), utils.MetadataFromCtx(ctx))

	return p.publish(log, subjectUserUpdated, &eventuserpb.UserUpdatedEvent{Id: id, IsSecurityUpdate: isSecurityUpdate})
}

func (p *UserPublisher) PublishUserDeleted(ctx context.Context, id string) error {
	log := pkglog.WithMetadata(pkglog.WithMethod(p.log, "PublishUserDeleted"), utils.MetadataFromCtx(ctx))

	return p.publish(log, subjectUserDeleted, &eventuserpb.UserDeletedEvent{Id: id})
}

func (p *UserPublisher) publish(log *slog.Logger, subject string, msg proto.Message) error {
	data, err := proto.Marshal(msg)
	if err != nil {
		log.Error("marshalling event", slog.String("subject", subject), pkglog.Err(err))

		return model.ErrNats
	}

	if err := p.conn.Publish(subject, data); err != nil {
		log.Error("publishing event", slog.String("subject", subject), pkglog.Err(err))
		return model.ErrNats
	}

	return nil
}
