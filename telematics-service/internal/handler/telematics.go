package handler

import (
	"database/sql"
	"errors"
	"log/slog"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"carsharing/telematics-service/internal/db"
	"carsharing/telematics-service/internal/service"
	telematicspb "github.com/sorawaslocked/car-rental-protos/gen/service/telematics"
)

type TelematicsHandler struct {
	telematicspb.UnimplementedCarTelematicsStreamServiceServer
	simSvc  *service.SimulationService
	carRepo *db.CarRepository
}

func NewTelematicsHandler(simSvc *service.SimulationService, carRepo *db.CarRepository) *TelematicsHandler {
	return &TelematicsHandler{simSvc: simSvc, carRepo: carRepo}
}

func (h *TelematicsHandler) StreamCarTelematicsEvents(
	req *telematicspb.StreamCarTelematicsEventsRequest,
	stream grpc.ServerStreamingServer[telematicspb.StreamCarTelematicsEventsResponse],
) error {
	slog.Info("telematics stream opened", "car_id", req.CarId)

	telemetry, err := h.carRepo.GetCarTelemetry(req.CarId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return status.Errorf(codes.NotFound, "car %s not found", req.CarId)
		}
		return status.Errorf(codes.Internal, "fetch car telemetry: %v", err)
	}

	// Send the current idle state immediately so the client has an initial position.
	if err := stream.Send(&telematicspb.StreamCarTelematicsEventsResponse{
		CarId:        req.CarId,
		Latitude:     telemetry.Latitude,
		Longitude:    telemetry.Longitude,
		FuelLevel:    telemetry.FuelLevel,
		BatteryLevel: telemetry.BatteryLevel,
		OdometerKm:   telemetry.MileageKm,
		RecordedAt:   timestamppb.Now(),
	}); err != nil {
		return err
	}

	updates, err := h.simSvc.RegisterStream(stream.Context(), &service.SimulationRequest{
		CarId:        req.CarId,
		Latitude:     telemetry.Latitude,
		Longitude:    telemetry.Longitude,
		FuelLevel:    telemetry.FuelLevel,
		BatteryLevel: telemetry.BatteryLevel,
		OdometerKm:   telemetry.MileageKm,
	})
	if err != nil {
		return status.Errorf(codes.Internal, "register stream: %v", err)
	}

	for update := range updates {
		if err := stream.Send(&telematicspb.StreamCarTelematicsEventsResponse{
			CarId:        update.CarId,
			Latitude:     update.Latitude,
			Longitude:    update.Longitude,
			FuelLevel:    update.FuelLevel,
			BatteryLevel: update.BatteryLevel,
			OdometerKm:   update.OdometerKm,
			RecordedAt:   timestamppb.New(update.RecordedAt),
		}); err != nil {
			slog.Warn("stream send failed", "car_id", req.CarId, "error", err)
			return err
		}
	}

	slog.Info("telematics stream closed", "car_id", req.CarId)
	return nil
}
