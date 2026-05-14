package handler

import (
	"log/slog"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	telematicspb "github.com/sorawaslocked/car-rental-protos/gen/service/telematics"
	"github.com/sorawaslocked/car-rental-telematics/internal/service"
)

// TelematicsHandler implements the CarTelematicsStreamServiceServer interface.
type TelematicsHandler struct {
	telematicspb.UnimplementedCarTelematicsStreamServiceServer
	simSvc *service.SimulationService
}

func NewTelematicsHandler(simSvc *service.SimulationService) *TelematicsHandler {
	return &TelematicsHandler{simSvc: simSvc}
}

// StreamCarTelematicsEvents starts a telematics simulation for the requested car
// and streams one update every 15 seconds until the trip ends or the client disconnects.
func (h *TelematicsHandler) StreamCarTelematicsEvents(
	req *telematicspb.StreamCarTelematicsEventsRequest,
	stream grpc.ServerStreamingServer[telematicspb.StreamCarTelematicsEventsResponse],
) error {
	slog.Info("telematics stream opened", "car_id", req.CarId)

	updates, err := h.simSvc.StartSimulation(stream.Context(), &service.SimulationRequest{
		CarId:        req.CarId,
		Latitude:     req.CurrentLocationLatitude,
		Longitude:    req.CurrentLocationLongitude,
		FuelLevel:    req.FuelLevel,
		BatteryLevel: req.BatteryLevel,
		OdometerKm:   req.OdometerKm,
	})
	if err != nil {
		return status.Errorf(codes.Internal, "start simulation: %v", err)
	}

	for update := range updates {
		resp := &telematicspb.StreamCarTelematicsEventsResponse{
			CarId:        update.CarId,
			Latitude:     update.Latitude,
			Longitude:    update.Longitude,
			FuelLevel:    update.FuelLevel,
			BatteryLevel: update.BatteryLevel,
			OdometerKm:   update.OdometerKm,
			RecordedAt:   timestamppb.New(update.RecordedAt),
		}
		if err := stream.Send(resp); err != nil {
			slog.Warn("stream send failed", "car_id", req.CarId, "error", err)
			return err
		}
	}

	slog.Info("telematics stream closed", "car_id", req.CarId)
	return nil
}
