package handler

import (
	"context"
	"github.com/sorawaslocked/car-rental-car-service/internal/adapter/grpc/dto"
	"log/slog"

	carsvc "github.com/sorawaslocked/car-rental-protos/gen/service/car"
)

type CarModelHandler struct {
	carModelService CarModelService

	log *slog.Logger

	carsvc.UnimplementedCarModelServiceServer
}

func NewCarModelHandler(carModelService CarModelService, log *slog.Logger) *CarModelHandler {
	s := &CarModelHandler{
		carModelService: carModelService,
	}

	s.log = log.With(
		slog.Group("src",
			slog.String("component", "CarModelHandler"),
		),
	)

	return s
}

func (h *CarModelHandler) CreateCarModel(ctx context.Context, req *carsvc.CreateCarModelRequest) (*carsvc.CreateCarModelResponse, error) {
	createInput := dto.FromCreateCarModelRequest(req)

	id, err := h.carModelService.Create(ctx, createInput)
	if err != nil {
		return nil, dto.FromErrorToStatusCode(err)
	}

	return &carsvc.CreateCarModelResponse{
		ID: id,
	}, nil
}

func (h *CarModelHandler) GetCarModel(ctx context.Context, req *carsvc.GetCarModelRequest) (*carsvc.GetCarModelResponse, error) {
	filterInput := dto.FromGetCarModelRequest(req)

	carModel, err := h.carModelService.Get(ctx, filterInput)
	if err != nil {
		return nil, dto.FromErrorToStatusCode(err)
	}

	return &carsvc.GetCarModelResponse{
		CarModel: dto.ToCarModelProto(carModel),
	}, nil
}

func (h *CarModelHandler) GetCarModels(ctx context.Context, req *carsvc.GetCarModelsRequest) (*carsvc.GetCarModelsResponse, error) {
	filterInput := dto.FromGetCarModelsRequest(req)

	carModels, err := h.carModelService.GetAll(ctx, filterInput)
	if err != nil {
		return nil, dto.FromErrorToStatusCode(err)
	}

	return &carsvc.GetCarModelsResponse{
		CarModels: dto.ToCarModelProtos(carModels),
	}, nil
}

func (h *CarModelHandler) UpdateCarModel(ctx context.Context, req *carsvc.UpdateCarModelRequest) (*carsvc.UpdateCarModelResponse, error) {
	filterInput, updateInput := dto.FromUpdateCarModelRequest(req)

	err := h.carModelService.Update(ctx, filterInput, updateInput)
	if err != nil {
		return nil, dto.FromErrorToStatusCode(err)
	}

	return &carsvc.UpdateCarModelResponse{}, nil
}

func (h *CarModelHandler) DeleteCarModel(ctx context.Context, req *carsvc.DeleteCarModelRequest) (*carsvc.DeleteCarModelResponse, error) {
	filterInput := dto.FromDeleteCarModelRequest(req)

	err := h.carModelService.Delete(ctx, filterInput)
	if err != nil {
		return nil, dto.FromErrorToStatusCode(err)
	}

	return &carsvc.DeleteCarModelResponse{}, nil
}
