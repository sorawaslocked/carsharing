package handler

import (
	"car-rental-car-service/internal/adapter/grpc/dto"
	"context"
	"log/slog"

	carsvc "github.com/sorawaslocked/car-rental-protos/gen/service/car"
)

type CarHandler struct {
	carService CarService

	log *slog.Logger

	carsvc.UnimplementedCarServiceServer
}

func NewCarHandler(carService CarService, log *slog.Logger) *CarHandler {
	s := &CarHandler{
		carService: carService,
	}

	s.log = log.With(
		slog.Group("src",
			slog.String("component", "CarHandler"),
		),
	)

	return s
}

func (h *CarHandler) CreateCar(ctx context.Context, req *carsvc.CreateCarRequest) (*carsvc.CreateCarResponse, error) {
	createInput := dto.FromCreateCarRequest(req)

	id, err := h.carService.Create(ctx, createInput)
	if err != nil {
		return nil, dto.FromErrorToStatusCode(err)
	}

	return &carsvc.CreateCarResponse{
		ID: id,
	}, nil
}

func (h *CarHandler) GetCar(ctx context.Context, req *carsvc.GetCarRequest) (*carsvc.GetCarResponse, error) {
	filterInput := dto.FromGetCarRequest(req)

	car, err := h.carService.Get(ctx, filterInput)
	if err != nil {
		return nil, dto.FromErrorToStatusCode(err)
	}

	return &carsvc.GetCarResponse{
		Car: dto.ToCarProto(car),
	}, nil
}

func (h *CarHandler) GetCars(ctx context.Context, req *carsvc.GetCarsRequest) (*carsvc.GetCarsResponse, error) {
	filterInput := dto.FromGetCarsRequest(req)

	cars, err := h.carService.GetAll(ctx, filterInput)
	if err != nil {
		return nil, dto.FromErrorToStatusCode(err)
	}

	return &carsvc.GetCarsResponse{
		Cars: dto.ToCarProtos(cars),
	}, nil
}

func (h *CarHandler) GetAvailableCars(ctx context.Context, req *carsvc.GetAvailableCarsRequest) (*carsvc.GetAvailableCarsResponse, error) {
	filterInput := dto.FromGetAvailableCars(req)

	cars, err := h.carService.GetAvailableCars(ctx, filterInput)
	if err != nil {
		return nil, dto.FromErrorToStatusCode(err)
	}

	return &carsvc.GetAvailableCarsResponse{
		Cars: dto.ToCarProtos(cars),
	}, nil
}

func (h *CarHandler) UpdateCar(ctx context.Context, req *carsvc.UpdateCarRequest) (*carsvc.UpdateCarResponse, error) {
	filterInput, updateInput := dto.FromUpdateCarRequest(req)

	err := h.carService.Update(ctx, filterInput, updateInput)
	if err != nil {
		return nil, dto.FromErrorToStatusCode(err)
	}

	return &carsvc.UpdateCarResponse{}, nil
}

func (h *CarHandler) UpdateCarStatus(ctx context.Context, req *carsvc.UpdateCarStatusRequest) (*carsvc.UpdateCarStatusResponse, error) {
	filterInput, statusUpdateInput := dto.FromUpdateCarStatusRequest(req)

	err := h.carService.UpdateCarStatus(ctx, filterInput, statusUpdateInput)
	if err != nil {
		return nil, dto.FromErrorToStatusCode(err)
	}

	return &carsvc.UpdateCarStatusResponse{}, nil
}

func (h *CarHandler) DeleteCar(ctx context.Context, req *carsvc.DeleteCarRequest) (*carsvc.DeleteCarResponse, error) {
	filterInput := dto.FromDeleteCarRequest(req)

	err := h.carService.Delete(ctx, filterInput)
	if err != nil {
		return nil, dto.FromErrorToStatusCode(err)
	}

	return &carsvc.DeleteCarResponse{}, nil
}
