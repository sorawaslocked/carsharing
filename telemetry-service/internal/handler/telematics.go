package handler

import (
	"database/sql"
	"errors"
	"log/slog"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	telemetrypb "carsharing/protos/gen/service/telemetry"
	"carsharing/telematics-service/internal/db"
	"carsharing/telematics-service/internal/service"
)

type TelematicsHandler struct {
	telemetrypb.UnimplementedCarTelemetryStreamServiceServer
	simSvc  *service.SimulationService
	carRepo *db.CarRepository
}

func NewTelematicsHandler(simSvc *service.SimulationService, carRepo *db.CarRepository) *TelematicsHandler {
	return &TelematicsHandler{simSvc: simSvc, carRepo: carRepo}
}

func (h *TelematicsHandler) StreamCarTelemetryEvents(
	req *telemetrypb.StreamCarTelemetryEventsRequest,
	stream grpc.ServerStreamingServer[telemetrypb.StreamCarTelemetryEventsResponse],
) error {
	slog.Info("telemetry stream opened", "car_id", req.CarId)

	telemetry, err := h.carRepo.GetCarTelemetry(req.CarId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return status.Errorf(codes.NotFound, "car %s not found", req.CarId)
		}
		return status.Errorf(codes.Internal, "fetch car telemetry: %v", err)
	}

	// Send the current idle state immediately so the client has an initial position.
	if err := stream.Send(&telemetrypb.StreamCarTelemetryEventsResponse{
		CarId:        req.CarId,
		Latitude:     telemetry.Latitude,
		Longitude:    telemetry.Longitude,
		FuelLevel:    telemetry.FuelLevel,
		BatteryLevel: telemetry.BatteryLevel,
		MileageKm:    telemetry.MileageKm,
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
		MileageKM:    telemetry.MileageKm,
	})
	if err != nil {
		return status.Errorf(codes.Internal, "register stream: %v", err)
	}

	for update := range updates {
		if err := stream.Send(&telemetrypb.StreamCarTelemetryEventsResponse{
			CarId:        update.CarId,
			Latitude:     update.Latitude,
			Longitude:    update.Longitude,
			FuelLevel:    update.FuelLevel,
			BatteryLevel: update.BatteryLevel,
			MileageKm:    update.MileageKM,
			RecordedAt:   timestamppb.New(update.RecordedAt),
		}); err != nil {
			slog.Warn("stream send failed", "car_id", req.CarId, "error", err)
			return err
		}
	}

	slog.Info("telemetry stream closed", "car_id", req.CarId)
	return nil
}
