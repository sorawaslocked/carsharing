package handler

import (
	"context"

	"github.com/sorawaslocked/car-rental-api-gateway/internal/model"
)

type CarInsuranceHandler struct {
	// client fleetsvc.InsuranceServiceClient
}

func NewInsuranceHandler() *CarInsuranceHandler {
	return &CarInsuranceHandler{}
}

func (h *CarInsuranceHandler) Create(_ context.Context, data model.CarInsuranceCreate) (string, error) {
	_ = data
	return "", errNotImplemented("CarInsuranceHandler.Create")
}

func (h *CarInsuranceHandler) Get(_ context.Context, id string) (model.CarInsurance, error) {
	_ = id
	return model.CarInsurance{}, errNotImplemented("CarInsuranceHandler.Get")
}

func (h *CarInsuranceHandler) GetAll(_ context.Context, filter model.CarInsuranceFilter) ([]model.CarInsurance, error) {
	_ = filter
	return nil, errNotImplemented("CarInsuranceHandler.GetAll")
}

func (h *CarInsuranceHandler) Update(_ context.Context, id string, data model.CarInsuranceUpdate) error {
	_ = id
	_ = data
	return errNotImplemented("CarInsuranceHandler.Update")
}

func (h *CarInsuranceHandler) Delete(_ context.Context, id string) error {
	_ = id
	return errNotImplemented("CarInsuranceHandler.Delete")
}

func (h *CarInsuranceHandler) GetImageUploadData(_ context.Context) (model.ImageUploadData, error) {
	return model.ImageUploadData{}, errNotImplemented("GetImageUploadUrl")
}
