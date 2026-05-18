package grpc

import (
	"context"
	"log/slog"
	"time"

	"carsharing/car-service/internal/model"
	pkglog "carsharing/car-service/internal/pkg/log"
	telematicspb "github.com/sorawaslocked/car-rental-protos/gen/service/telematics"
	"google.golang.org/grpc"
)

type TelematicsStreamClient struct {
	client telematicspb.CarTelematicsStreamServiceClient
	log    *slog.Logger
}

func NewTelematicsStreamClient(conn *grpc.ClientConn, log *slog.Logger) *TelematicsStreamClient {
	return &TelematicsStreamClient{
		client: telematicspb.NewCarTelematicsStreamServiceClient(conn),
		log:    pkglog.WithComponent(log, "adapter.TelematicsStreamClient"),
	}
}

func (c *TelematicsStreamClient) Subscribe(ctx context.Context, car model.Car) (<-chan model.TelematicsUpdate, error) {
	req := &telematicspb.StreamCarTelematicsEventsRequest{
		CarId: car.ID,
	}

	stream, err := c.client.StreamCarTelematicsEvents(ctx, req)
	if err != nil {
		return nil, err
	}

	ch := make(chan model.TelematicsUpdate, 16)

	go func() {
		defer close(ch)
		for {
			resp, err := stream.Recv()
			if err != nil {
				if ctx.Err() == nil {
					c.log.Error("telematics stream recv error",
						pkglog.Err(err),
						slog.String("carID", car.ID),
					)
				}
				return
			}

			update := model.TelematicsUpdate{
				CarID:        resp.GetCarId(),
				Latitude:     resp.GetLatitude(),
				Longitude:    resp.GetLongitude(),
				OdometerKM:   resp.GetOdometerKm(),
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
