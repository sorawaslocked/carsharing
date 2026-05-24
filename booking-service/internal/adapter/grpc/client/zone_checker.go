package client

import (
	"context"
	"log/slog"

	"carsharing/booking-service/internal/model"
	carsvc "carsharing/protos/gen/service/car"
	pkglog "carsharing/shared/pkg/log"
	"carsharing/shared/pkg/utils"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ZoneChecker struct {
	log    *slog.Logger
	client carsvc.ZoneServiceClient
}

func NewZoneChecker(log *slog.Logger, conn *grpc.ClientConn) *ZoneChecker {
	return &ZoneChecker{
		log:    pkglog.WithComponent(log, "adapter.grpc.client.ZoneChecker"),
		client: carsvc.NewZoneServiceClient(conn),
	}
}

func (c *ZoneChecker) Exists(ctx context.Context, zoneID string) (bool, error) {
	log := pkglog.WithMetadata(pkglog.WithMethod(c.log, "Exists"), utils.MetadataFromCtx(ctx))

	_, err := c.client.GetZone(ctx, &carsvc.GetZoneRequest{Id: zoneID})
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return false, nil
		}
		log.Error("grpc: checking zone existence", slog.String("zoneID", zoneID), pkglog.Err(err))
		return false, model.ErrInternalServerError
	}

	return true, nil
}
