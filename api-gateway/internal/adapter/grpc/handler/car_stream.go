package handler

import (
	"context"
	"errors"
	"io"

	"carsharing/api-gateway/internal/adapter/grpc/dto"
	"carsharing/api-gateway/internal/model"
	basepb "carsharing/protos/gen/base"
	carsvc "carsharing/protos/gen/service/car"
	pkglog "carsharing/shared/pkg/log"
	"carsharing/shared/pkg/utils"
)

func (h *CarHandler) StreamCarsWithFilter(ctx context.Context, filter model.CarFilter, send func([]model.SlimCar) error) error {
	logger := pkglog.WithMethod(h.log, "StreamCarsWithFilter")
	logger = pkglog.WithMetadata(logger, utils.MetadataFromCtx(ctx))

	req := &carsvc.StreamCarsWithFilterRequest{
		Brand:        filter.Brand,
		Model:        filter.Model,
		FuelType:     filter.FuelType,
		Transmission: filter.Transmission,
		BodyType:     filter.BodyType,
		Class:        filter.Class,
		ZoneId:       filter.ZoneID,
		Status:       filter.Status,
		RadiusM:      filter.RadiusM,
		MinFuelLevel: filter.MinFuelLevel,
	}
	if filter.MinSeats != nil {
		v := int32(*filter.MinSeats)
		req.MinSeats = &v
	}
	if filter.Location != nil {
		req.Location = &basepb.Location{
			Latitude:  filter.Location.Latitude,
			Longitude: filter.Location.Longitude,
		}
	}

	stream, err := h.streamClient.StreamCarsWithFilter(ctx, req)
	if err != nil {
		if dto.IsSystemErr(err) {
			logger.Error("grpc call failed", pkglog.Err(err))
		}
		return dto.FromGrpcErr(err)
	}

	for {
		msg, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			return nil
		}
		if err != nil {
			if ctx.Err() != nil {
				return nil
			}
			if dto.IsSystemErr(err) {
				logger.Error("stream recv failed", pkglog.Err(err))
			}
			return dto.FromGrpcErr(err)
		}

		cars := make([]model.SlimCar, len(msg.GetCar()))
		for i, c := range msg.GetCar() {
			cars[i] = model.SlimCar{
				ID:           c.GetId(),
				ModelID:      c.GetModelId(),
				LicensePlate: c.GetLicensePlate(),
				Color:        c.GetColor(),
				Location:     dto.LocationFromProto(c.GetLocation()),
				FuelLevel:    c.GetFuelLevel(),
				Status:       c.GetStatus(),
			}
		}

		if err = send(cars); err != nil {
			return err
		}
	}
}

func (h *CarHandler) StreamCarTelemetry(ctx context.Context, carID string, send func(model.CarTelemetryEvent) error) error {
	logger := pkglog.WithMethod(h.log, "StreamCarTelemetry")
	logger = pkglog.WithMetadata(logger, utils.MetadataFromCtx(ctx))

	stream, err := h.streamClient.StreamCarTelemetry(ctx, &carsvc.StreamCarTelemetryRequest{CarId: carID})
	if err != nil {
		if dto.IsSystemErr(err) {
			logger.Error("grpc call failed", pkglog.Err(err))
		}
		return dto.FromGrpcErr(err)
	}

	for {
		msg, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			return nil
		}
		if err != nil {
			if ctx.Err() != nil {
				return nil
			}
			if dto.IsSystemErr(err) {
				logger.Error("stream recv failed", pkglog.Err(err))
			}
			return dto.FromGrpcErr(err)
		}

		event := model.CarTelemetryEvent{
			FuelLevel:    msg.GetFuelLevel(),
			BatteryLevel: msg.GetBatteryLevel(),
			MileageKM:    msg.GetMileageKm(),
			Location:     dto.LocationFromProto(msg.GetLocation()),
		}
		if t := msg.GetRecordedAt(); t != nil {
			event.RecordedAt = t.AsTime()
		}

		if err = send(event); err != nil {
			return err
		}
	}
}
