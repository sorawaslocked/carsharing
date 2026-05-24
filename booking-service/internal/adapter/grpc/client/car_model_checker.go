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

type CarModelChecker struct {
	log    *slog.Logger
	client carsvc.CarModelServiceClient
}

func NewCarModelChecker(log *slog.Logger, conn *grpc.ClientConn) *CarModelChecker {
	return &CarModelChecker{
		log:    pkglog.WithComponent(log, "adapter.grpc.client.CarModelChecker"),
		client: carsvc.NewCarModelServiceClient(conn),
	}
}

func (c *CarModelChecker) Exists(ctx context.Context, modelID string) (bool, error) {
	log := pkglog.WithMetadata(pkglog.WithMethod(c.log, "Exists"), utils.MetadataFromCtx(ctx))

	_, err := c.client.GetCarModel(ctx, &carsvc.GetCarModelRequest{Id: modelID})
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return false, nil
		}
		log.Error("grpc: checking car model existence", slog.String("modelID", modelID), pkglog.Err(err))
		return false, model.ErrInternalServerError
	}

	return true, nil
}
