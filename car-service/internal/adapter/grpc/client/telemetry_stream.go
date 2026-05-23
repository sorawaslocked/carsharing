package client

import (
	"context"
	"log/slog"
	"time"

	"carsharing/car-service/internal/model"
	pkglog "carsharing/shared/pkg/log"

	telemetrypb "carsharing/protos/gen/service/telemetry"
	"google.golang.org/grpc"
)

type TelemetryStreamClient struct {
	client telemetrypb.CarTelemetryStreamServiceClient
	log    *slog.Logger
}

func NewTelemetryStreamClient(conn *grpc.ClientConn, log *slog.Logger) *TelemetryStreamClient {
	return &TelemetryStreamClient{
		client: telemetrypb.NewCarTelemetryStreamServiceClient(conn),
		log:    pkglog.WithComponent(log, "grpc.client.TelemetryStreamClient"),
	}
}

func (c *TelemetryStreamClient) Subscribe(ctx context.Context, car model.Car) (<-chan model.TelemetryUpdate, error) {
	req := &telemetrypb.StreamCarTelemetryEventsRequest{
		CarId: car.ID,
	}

	stream, err := c.client.StreamCarTelemetryEvents(ctx, req)
	if err != nil {
		return nil, err
	}

	ch := make(chan model.TelemetryUpdate, 16)

	go func() {
		defer close(ch)
		for {
			resp, err := stream.Recv()
			if err != nil {
				if ctx.Err() == nil {
					c.log.Error("telemetry stream recv error",
						pkglog.Err(err),
						slog.String("carID", car.ID),
					)
				}
				return
			}

			update := model.TelemetryUpdate{
				CarID:        resp.GetCarId(),
				Latitude:     resp.GetLatitude(),
				Longitude:    resp.GetLongitude(),
				MileageKM:    resp.GetMileageKm(),
				FuelLevel:    resp.FuelLevel,
				BatteryLevel: resp.BatteryLevel,
			}
			if ts := resp.GetRecordedAt(); ts != nil {
				update.RecordedAt = ts.AsTime()
			} else {
				update.RecordedAt = time.Now()
			}

			select {
			case ch <- update:
			case <-ctx.Done():
				return
			}
		}
	}()

	return ch, nil
}
