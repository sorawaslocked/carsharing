package publisher

import (
	"context"
	"log/slog"

	"carsharing/car-service/internal/model"
	pkglog "carsharing/shared/pkg/log"
	"carsharing/shared/pkg/utils"

	"github.com/nats-io/nats.go"
	eventcarpb "github.com/sorawaslocked/car-rental-protos/gen/event/car"
	"google.golang.org/protobuf/proto"
)

const subjectCarStatusUpdated = "car.status.updated"

type CarPublisher struct {
	log  *slog.Logger
	conn *nats.Conn
}

func NewCarPublisher(log *slog.Logger, conn *nats.Conn) *CarPublisher {
	return &CarPublisher{
		log:  pkglog.WithComponent(log, "adapter.nats.publisher.CarPublisher"),
		conn: conn,
	}
}

func (p *CarPublisher) PublishCarStatusUpdated(ctx context.Context, carID, fromStatus, toStatus string) error {
	log := pkglog.WithMetadata(pkglog.WithMethod(p.log, "PublishCarStatusUpdated"), utils.MetadataFromCtx(ctx))

	return p.publish(log, subjectCarStatusUpdated, &eventcarpb.CarStatusUpdatedEvent{
		CarId:      carID,
		FromStatus: fromStatus,
		ToStatus:   toStatus,
	})
}

func (p *CarPublisher) publish(log *slog.Logger, subject string, msg proto.Message) error {
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
