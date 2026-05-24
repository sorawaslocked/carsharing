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

type CarChecker struct {
	log    *slog.Logger
	client carsvc.CarServiceClient
}

func NewCarChecker(log *slog.Logger, conn *grpc.ClientConn) *CarChecker {
	return &CarChecker{
		log:    pkglog.WithComponent(log, "adapter.grpc.client.CarChecker"),
		client: carsvc.NewCarServiceClient(conn),
	}
}

func (c *CarChecker) Exists(ctx context.Context, carID string) (bool, error) {
	log := pkglog.WithMetadata(pkglog.WithMethod(c.log, "Exists"), utils.MetadataFromCtx(ctx))

	_, err := c.client.GetCar(ctx, &carsvc.GetCarRequest{Id: carID})
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return false, nil
		}
		log.Error("grpc: checking car existence", slog.String("carID", carID), pkglog.Err(err))
		return false, model.ErrInternalServerError
	}

	return true, nil
}

func (c *CarChecker) GetStatus(ctx context.Context, carID string) (model.CarStatus, error) {
	log := pkglog.WithMetadata(pkglog.WithMethod(c.log, "GetStatus"), utils.MetadataFromCtx(ctx))

	resp, err := c.client.GetCar(ctx, &carsvc.GetCarRequest{Id: carID})
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return "", model.ErrCarNotFound
		}
		log.Error("grpc: getting car status", slog.String("carID", carID), pkglog.Err(err))
		return "", model.ErrInternalServerError
	}

	carStatus, ok := model.CarStatusFromString(resp.Car.Status)
	if !ok {
		log.Error("grpc: unknown car status", slog.String("carID", carID), slog.String("status", resp.Car.Status))
		return "", model.ErrInternalServerError
	}

	return carStatus, nil
}
