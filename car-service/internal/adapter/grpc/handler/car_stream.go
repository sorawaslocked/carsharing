package handler

import (
	"log/slog"
	"time"

	"carsharing/car-service/internal/adapter/grpc/dto"
	"carsharing/car-service/internal/model"

	"github.com/sorawaslocked/car-rental-protos/gen/base"
	basecar "github.com/sorawaslocked/car-rental-protos/gen/base/car"
	carsvc "github.com/sorawaslocked/car-rental-protos/gen/service/car"
	"google.golang.org/grpc"
)

type CarStreamHandler struct {
	carService           CarService
	telematicsSubscriber TelematicsSubscriber

	log *slog.Logger

	carsvc.UnimplementedCarStreamServiceServer
}

func NewCarStreamHandler(carService CarService, telematicsSubscriber TelematicsSubscriber, log *slog.Logger) *CarStreamHandler {
	h := &CarStreamHandler{
		carService:           carService,
		telematicsSubscriber: telematicsSubscriber,
	}

	h.log = log.With(
		slog.Group("src",
			slog.String("component", "CarStreamHandler"),
		),
	)

	return h
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

func (h *CarStreamHandler) sendFilteredCars(req *carsvc.StreamCarsWithFilterRequest, filterInput model.CarFilterInput, stream grpc.ServerStreamingServer[carsvc.StreamCarsWithFilterResponse]) error {
	ctx := stream.Context()

	cars, err := h.carService.GetAll(ctx, filterInput)
	if err != nil {
		return dto.FromErrorToStatusCode(err)
	}

	for _, c := range cars {
		if req.MinFuelLevel != nil && (c.FuelLevel == nil || *c.FuelLevel < *req.MinFuelLevel) {
			continue
		}

		if err := stream.Send(&carsvc.StreamCarsWithFilterResponse{
			Car: []*basecar.SlimCar{dto.ToSlimCarProto(c)},
		}); err != nil {
			h.log.Error("failed to send car stream response", slog.String("error", err.Error()))
			return err
		}
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

	ch, err := h.telematicsSubscriber.SubscribeCarStream(ctx, req.CarId)
	if err != nil {
		return dto.FromErrorToStatusCode(err)
	}

	for {
		select {
		case update, ok := <-ch:
			if !ok {
				return nil
			}
			resp := &carsvc.StreamCarTelemetryResponse{
				FuelLevel:    update.FuelLevel,
				BatteryLevel: update.BatteryLevel,
				MileageKm:    update.OdometerKM,
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

func filterInputFromStreamRequest(req *carsvc.StreamCarsWithFilterRequest) model.CarFilterInput {
	filterInput := model.CarFilterInput{}

	if req.Status != nil {
		filterInput.Status = req.Status
	}

	if req.Location != nil || req.RadiusM != nil {
		lf := &model.LocationFilter{}
		if req.Location != nil {
			lf.Location = model.Location{
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
		mf := &model.CarModelFilterInput{
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
