package handler

import (
	"context"

	"github.com/sorawaslocked/car-rental-api-gateway/internal/model"
)

type CarHandler struct {
	// client fleetsvc.CarServiceClient
}

func NewCarHandler() *CarHandler {
	return &CarHandler{}
}

func (h *CarHandler) Create(_ context.Context, data model.CarCreate) (string, error) {
	_ = data
	return "", errNotImplemented("CarHandler.Create")
}

func (h *CarHandler) Get(_ context.Context, id string) (model.Car, error) {
	_ = id
	return model.Car{}, errNotImplemented("CarHandler.Get")
}

func (h *CarHandler) GetAll(_ context.Context, filter model.CarFilter) ([]model.Car, error) {
	_ = filter
	return nil, errNotImplemented("CarHandler.GetAll")
}

func (h *CarHandler) Update(_ context.Context, id string, data model.CarUpdate) error {
	_ = id
	_ = data
	return errNotImplemented("CarHandler.Update")
}

func (h *CarHandler) Delete(_ context.Context, id string) error {
	_ = id
	return errNotImplemented("CarHandler.Delete")
}

func (h *CarHandler) GetCarStatusLog(_ context.Context, filter model.CarStatusLogFilter) ([]model.CarStatusLogEntry, error) {
	_ = filter
	return nil, errNotImplemented("CarHandler.GetCarStatusLog")
}

func (h *CarHandler) GetCarFuelHistory(_ context.Context, filter model.CarFuelReadingFilter) ([]model.CarFuelReading, error) {
	_ = filter
	return nil, errNotImplemented("CarHandler.GetCarFuelHistory")
}

func (h *CarHandler) GetImageUploadData(_ context.Context) (model.ImageUploadData, error) {
	return model.ImageUploadData{}, errNotImplemented("CarHandler.GetImageUploadData")
}
