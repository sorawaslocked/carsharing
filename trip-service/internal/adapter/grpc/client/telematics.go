package client

import (
	"context"
	"io"
	"log/slog"
	"time"

	"google.golang.org/grpc"

	carsvc "github.com/sorawaslocked/car-rental-protos/gen/service/car"

	pkglog "carsharing/shared/pkg/log"
	pkgutils "carsharing/shared/pkg/utils"
	"carsharing/trip-service/internal/model"
)

type TelematicsClient struct {
	log          *slog.Logger
	carClient    carsvc.CarServiceClient
	streamClient carsvc.CarStreamServiceClient
}

func NewTelematicsClient(log *slog.Logger, carConn, streamConn *grpc.ClientConn) *TelematicsClient {
	return &TelematicsClient{
		log:          pkglog.WithComponent(log, "client.TelematicsClient"),
		carClient:    carsvc.NewCarServiceClient(carConn),
		streamClient: carsvc.NewCarStreamServiceClient(streamConn),
	}
}

func (c *TelematicsClient) GetLatestTelemetry(ctx context.Context, carID string) (model.CarTelemetry, error) {
	log := pkglog.WithMethod(c.log, "GetLatestTelemetry")
	log = pkglog.WithMetadata(log, pkgutils.MetadataFromCtx(ctx))

	resp, err := c.carClient.GetCar(ctx, &carsvc.GetCarRequest{Id: carID})
	if err != nil {
		log.Error("failed to get car", pkglog.Err(err))
		return model.CarTelemetry{}, err
	}

	car := resp.Car
	t := model.CarTelemetry{
		CarID:      car.Id,
		MileageKM:  car.MileageKm,
		FuelLevel:  car.FuelLevel,
		RecordedAt: time.Now(),
	}
	if car.Location != nil {
		t.Location = model.Location{
			Latitude:  car.Location.Latitude,
			Longitude: car.Location.Longitude,
		}
	}
	if car.LastSeenAt != nil {
		t.RecordedAt = car.LastSeenAt.AsTime()
	}
	return t, nil
}

func (c *TelematicsClient) StreamTelemetry(ctx context.Context, carID string, fn func(model.CarTelemetry) error) error {
	log := pkglog.WithMethod(c.log, "StreamTelemetry")
	log = pkglog.WithMetadata(log, pkgutils.MetadataFromCtx(ctx))

	stream, err := c.streamClient.StreamCarTelemetry(ctx, &carsvc.StreamCarTelemetryRequest{CarId: carID})
	if err != nil {
		log.Error("failed to open telemetry stream", pkglog.Err(err))
		return err
	}

	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			return io.EOF
		}
		if err != nil {
			log.Error("telemetry stream error", pkglog.Err(err))
			return err
		}

		t := model.CarTelemetry{
			CarID:     carID,
			MileageKM: resp.MileageKm,
			FuelLevel: resp.FuelLevel,
		}
		if resp.Location != nil {
			t.Location = model.Location{
				Latitude:  resp.Location.Latitude,
				Longitude: resp.Location.Longitude,
			}
		}
		if resp.RecordedAt != nil {
			t.RecordedAt = resp.RecordedAt.AsTime()
		}
		if err = fn(t); err != nil {
			return err
		}
	}
}
