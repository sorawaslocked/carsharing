package client

import (
	"context"
	"log/slog"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	carsvc "carsharing/protos/gen/service/car"
	pkglog "carsharing/shared/pkg/log"
	pkgutils "carsharing/shared/pkg/utils"
	"carsharing/trip-service/internal/model"
)

type ZoneClient struct {
	log    *slog.Logger
	client carsvc.ZoneServiceClient
}

func NewZoneClient(log *slog.Logger, conn *grpc.ClientConn) *ZoneClient {
	return &ZoneClient{
		log:    pkglog.WithComponent(log, "adapter.grpc.client.ZoneClient"),
		client: carsvc.NewZoneServiceClient(conn),
	}
}

func (c *ZoneClient) GetZonePricing(ctx context.Context, lat, lng float64) (int32, error) {
	log := pkglog.WithMetadata(pkglog.WithMethod(c.log, "GetZonePricing"), pkgutils.MetadataFromCtx(ctx))

	resp, err := c.client.GetZonePricing(ctx, &carsvc.GetZonePricingRequest{
		Latitude:  lat,
		Longitude: lng,
	})
	if err != nil {
		if status.Code(err) == codes.FailedPrecondition {
			return 0, model.ErrLocationInNoDropZone
		}
		log.Error("grpc: getting zone pricing", pkglog.Err(err))
		return 0, err
	}

	return resp.FeeAdjustment, nil
}
