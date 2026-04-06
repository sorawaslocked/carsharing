package handler

import (
	"context"

	"github.com/sorawaslocked/car-rental-api-gateway/internal/model"
)

type CarModelHandler struct {
	// client fleetsvc.CarModelServiceClient
}

func NewCarModelHandler() *CarModelHandler {
	return &CarModelHandler{}
}

func (h *CarModelHandler) Create(_ context.Context, data model.CarModelCreate) (string, error) {
	_ = data
	return "", errNotImplemented("CarModelHandler.Create")
}

func (h *CarModelHandler) Get(_ context.Context, id string) (model.CarModel, error) {
	_ = id
	return model.CarModel{}, errNotImplemented("CarModelHandler.Get")
}

func (h *CarModelHandler) GetAll(_ context.Context, filter model.CarModelFilter) ([]model.CarModel, error) {
	_ = filter
	return nil, errNotImplemented("CarModelHandler.GetAll")
}

func (h *CarModelHandler) Update(_ context.Context, id string, data model.CarModelUpdate) error {
	_ = id
	_ = data
	return errNotImplemented("CarModelHandler.Update")
}

func (h *CarModelHandler) Delete(_ context.Context, id string) error {
	_ = id
	return errNotImplemented("CarModelHandler.Delete")
}

func (h *CarModelHandler) GetImageUploadData(_ context.Context) (model.ImageUploadData, error) {
	return model.ImageUploadData{}, errNotImplemented("CarModelHandler.GetImageUploadData")
}
