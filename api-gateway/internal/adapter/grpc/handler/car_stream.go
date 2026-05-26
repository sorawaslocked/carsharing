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
	log := pkglog.WithMetadata(pkglog.WithMethod(h.log, "StreamCarsWithFilter"), utils.MetadataFromCtx(ctx))
	log.Debug("starting stream")

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
		if ctx.Err() != nil {
			return nil
		}
		log.Warn("streaming cars with filter", pkglog.Err(err))
		return dto.FromGrpcErr(err)
	}
	log.Debug("stream opened")

	for {
		msg, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			return nil
		}
		if err != nil {
			if ctx.Err() != nil {
				return nil
			}
			log.Warn("receiving car stream", pkglog.Err(err))
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

func (h *CarHandler) StreamCarStatusUpdates(ctx context.Context, carID string, send func(model.CarStatusEvent) error) error {
	log := pkglog.WithMetadata(pkglog.WithMethod(h.log, "StreamCarStatusUpdates"), utils.MetadataFromCtx(ctx))
	log.Debug("starting stream")

	req := &carsvc.StreamCarStatusUpdatesRequest{CarId: carID}

	stream, err := h.streamClient.StreamCarStatusUpdates(ctx, req)
	if err != nil {
		if ctx.Err() != nil {
			return nil
		}
		log.Warn("streaming car status updates", pkglog.Err(err))
		return dto.FromGrpcErr(err)
	}
	log.Debug("stream opened")

	for {
		msg, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			return nil
		}
		if err != nil {
			if ctx.Err() != nil {
				return nil
			}
			log.Warn("receiving car status updates stream", pkglog.Err(err))
			return dto.FromGrpcErr(err)
		}

		event := model.CarStatusEvent{
			FromStatus: msg.GetFromStatus(),
			ToStatus:   msg.GetToStatus(),
		}
		if t := msg.GetChangedAt(); t != nil {
			event.ChangedAt = t.AsTime()
		}

		if err = send(event); err != nil {
			return err
		}
	}
}

func (h *CarHandler) StreamCarTelemetry(ctx context.Context, carID string, send func(model.CarTelemetryEvent) error) error {
	log := pkglog.WithMetadata(pkglog.WithMethod(h.log, "StreamCarTelemetry"), utils.MetadataFromCtx(ctx))
	log.Debug("starting stream")

	req := &carsvc.StreamCarTelemetryRequest{CarId: carID}

	stream, err := h.streamClient.StreamCarTelemetry(ctx, req)
	if err != nil {
		if ctx.Err() != nil {
			return nil
		}
		log.Warn("streaming car telemetry", pkglog.Err(err))
		return dto.FromGrpcErr(err)
	}
	log.Debug("stream opened")

	for {
		msg, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			return nil
		}
		if err != nil {
			if ctx.Err() != nil {
				return nil
			}
			log.Warn("receiving car telemetry stream", pkglog.Err(err))
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
