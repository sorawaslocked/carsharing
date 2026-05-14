package client

import (
	"context"
	"io"
	"log/slog"
	"time"

	"google.golang.org/grpc"

	carsvc "github.com/sorawaslocked/car-rental-protos/gen/service/car"

	"github.com/sorawaslocked/car-rental-trip-service/internal/model"
	pkglog "github.com/sorawaslocked/car-rental-trip-service/internal/pkg/log"
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
	resp, err := c.carClient.GetCar(ctx, &carsvc.GetCarRequest{Id: carID})
	if err != nil {
		return model.CarTelemetry{}, err
	}
	car := resp.Car
	t := model.CarTelemetry{
		CarID:      car.Id,
		OdometerKM: car.MileageKm,
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
	stream, err := c.streamClient.StreamCarTelemetry(ctx, &carsvc.StreamCarTelemetryRequest{CarId: carID})
	if err != nil {
		return err
	}
	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			return io.EOF
		}
		if err != nil {
			return err
		}
		t := model.CarTelemetry{
			CarID:      carID,
			OdometerKM: resp.MileageKm,
			FuelLevel:  resp.FuelLevel,
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
