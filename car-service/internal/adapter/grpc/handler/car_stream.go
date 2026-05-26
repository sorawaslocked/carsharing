package handler

import (
	"log/slog"
	"time"

	"carsharing/car-service/internal/adapter/grpc/dto"
	"carsharing/car-service/internal/validation"
	pkglog "carsharing/shared/pkg/log"
	sharedvalidation "carsharing/shared/validation"

	"carsharing/protos/gen/base"
	basecar "carsharing/protos/gen/base/car"
	carsvc "carsharing/protos/gen/service/car"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type CarStreamHandler struct {
	log                 *slog.Logger
	carService          CarService
	telemetrySubscriber TelemetrySubscriber
	statusSubscriber    StatusSubscriber

	carsvc.UnimplementedCarStreamServiceServer
}

func NewCarStreamHandler(log *slog.Logger, carService CarService, telemetrySubscriber TelemetrySubscriber, statusSubscriber StatusSubscriber) *CarStreamHandler {
	return &CarStreamHandler{
		log:                 pkglog.WithComponent(log, "adapter.grpc.handler.CarStreamHandler"),
		carService:          carService,
		telemetrySubscriber: telemetrySubscriber,
		statusSubscriber:    statusSubscriber,
	}
}

const carFilterStreamInterval = 30 * time.Second

func (h *CarStreamHandler) StreamCarsWithFilter(req *carsvc.StreamCarsWithFilterRequest, stream grpc.ServerStreamingServer[carsvc.StreamCarsWithFilterResponse]) error {
	ctx := stream.Context()
	filterInput := filterInputFromStreamRequest(req)
	ticker := time.NewTicker(carFilterStreamInterval)
	defer ticker.Stop()

	if err := h.sendFilteredCars(req, filterInput, stream); err != nil {
		return err
	}

	for {
		select {
		case <-ticker.C:
			if err := h.sendFilteredCars(req, filterInput, stream); err != nil {
				return err
			}
		case <-ctx.Done():
			return nil
		}
	}
}

func (h *CarStreamHandler) sendFilteredCars(req *carsvc.StreamCarsWithFilterRequest, filterInput validation.CarFilter, stream grpc.ServerStreamingServer[carsvc.StreamCarsWithFilterResponse]) error {
	ctx := stream.Context()

	cars, err := h.carService.List(ctx, filterInput)
	if err != nil {
		return dto.FromErrorToStatusCode(err)
	}

	batch := make([]*basecar.SlimCar, 0, len(cars))
	for _, c := range cars {
		if req.MinFuelLevel != nil && (c.FuelLevel == nil || *c.FuelLevel < *req.MinFuelLevel) {
			continue
		}
		batch = append(batch, dto.ToSlimCarProto(c))
	}

	if err := stream.Send(&carsvc.StreamCarsWithFilterResponse{Car: batch}); err != nil {
		h.log.Error("failed to send car stream response", slog.String("error", err.Error()))
		return err
	}

	return nil
}

func (h *CarStreamHandler) StreamCarTelemetry(req *carsvc.StreamCarTelemetryRequest, stream grpc.ServerStreamingServer[carsvc.StreamCarTelemetryResponse]) error {
	ctx := stream.Context()

	car, err := h.carService.Get(ctx, req.CarId)
	if err != nil {
		return dto.FromErrorToStatusCode(err)
	}

	initial := &carsvc.StreamCarTelemetryResponse{
		FuelLevel:    car.FuelLevel,
		BatteryLevel: car.BatteryLevel,
		MileageKm:    car.MileageKM,
		RecordedAt:   timestamppb.New(car.LastSeenAt),
	}
	if car.Location.Latitude != 0 || car.Location.Longitude != 0 {
		initial.Location = &base.Location{
			Latitude:  car.Location.Latitude,
			Longitude: car.Location.Longitude,
		}
	}
	if err := stream.Send(initial); err != nil {
		return err
	}

	ch, unsub := h.telemetrySubscriber.SubscribeUpdates(req.CarId)
	defer unsub()

	for {
		select {
		case update, ok := <-ch:
			if !ok {
				return nil
			}
			resp := &carsvc.StreamCarTelemetryResponse{
				FuelLevel:    update.FuelLevel,
				BatteryLevel: update.BatteryLevel,
				MileageKm:    update.MileageKM,
				RecordedAt:   timestamppb.New(update.RecordedAt),
			}
			if update.Latitude != 0 || update.Longitude != 0 {
				resp.Location = &base.Location{
					Latitude:  update.Latitude,
					Longitude: update.Longitude,
				}
			}
			if err := stream.Send(resp); err != nil {
				return err
			}
		case <-ctx.Done():
			return nil
		}
	}
}

func (h *CarStreamHandler) StreamCarStatusUpdates(req *carsvc.StreamCarStatusUpdatesRequest, stream grpc.ServerStreamingServer[carsvc.StreamCarStatusUpdatesResponse]) error {
	ctx := stream.Context()

	if _, err := h.carService.Get(ctx, req.CarId); err != nil {
		return dto.FromErrorToStatusCode(err)
	}

	ch, unsub := h.statusSubscriber.SubscribeStatusUpdates(req.CarId)
	defer unsub()

	for {
		select {
		case update, ok := <-ch:
			if !ok {
				return nil
			}
			if err := stream.Send(&carsvc.StreamCarStatusUpdatesResponse{
				FromStatus: update.FromStatus.String(),
				ToStatus:   update.ToStatus.String(),
				ChangedAt:  timestamppb.New(update.ChangedAt),
			}); err != nil {
				return err
			}
		case <-ctx.Done():
			return nil
		}
	}
}

func filterInputFromStreamRequest(req *carsvc.StreamCarsWithFilterRequest) validation.CarFilter {
	filterInput := validation.CarFilter{}

	if req.Status != nil {
		filterInput.Status = req.Status
	}

	if req.Location != nil || req.RadiusM != nil {
		lf := &sharedvalidation.LocationFilter{}
		if req.Location != nil {
			lf.Location = sharedvalidation.Location{
				Latitude:  req.Location.Latitude,
				Longitude: req.Location.Longitude,
			}
		}
		if req.RadiusM != nil {
			lf.RadiusKM = float64(*req.RadiusM) / 1000
		}
		filterInput.LocationFilter = lf
	}

	if req.Brand != nil || req.Model != nil || req.FuelType != nil || req.Transmission != nil || req.BodyType != nil || req.Class != nil || req.MinSeats != nil {
		mf := &validation.CarModelFilter{
			Brand:        req.Brand,
			Model:        req.Model,
			FuelType:     req.FuelType,
			Transmission: req.Transmission,
			BodyType:     req.BodyType,
			Class:        req.Class,
		}
		if req.MinSeats != nil {
			v := int8(*req.MinSeats)
			mf.MinSeats = &v
		}
		filterInput.ModelFilter = mf
	}

	return filterInput
}
