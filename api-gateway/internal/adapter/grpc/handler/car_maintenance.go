package handler

import (
	"context"

	"github.com/sorawaslocked/car-rental-api-gateway/internal/model"
)

type CarMaintenanceHandler struct {
	// client fleetsvc.CarMaintenanceServiceClient
}

func NewCarMaintenanceHandler() *CarMaintenanceHandler {
	return &CarMaintenanceHandler{}
}

func (h *CarMaintenanceHandler) CreateTemplate(_ context.Context, data model.CarMaintenanceTemplateCreate) (string, error) {
	_ = data
	return "", errNotImplemented("CarMaintenanceHandler.CreateTemplate")
}

func (h *CarMaintenanceHandler) GetTemplate(_ context.Context, id string) (model.CarMaintenanceTemplate, error) {
	_ = id
	return model.CarMaintenanceTemplate{}, errNotImplemented("CarMaintenanceHandler.GetTemplate")
}

func (h *CarMaintenanceHandler) GetAllTemplates(_ context.Context, filter model.CarMaintenanceTemplateFilter) ([]model.CarMaintenanceTemplate, error) {
	_ = filter
	return nil, errNotImplemented("CarMaintenanceHandler.GetAllTemplates")
}

func (h *CarMaintenanceHandler) UpdateTemplate(_ context.Context, id string, data model.CarMaintenanceTemplateUpdate) error {
	_ = id
	_ = data
	return errNotImplemented("CarMaintenanceHandler.UploadTemplate")
}

func (h *CarMaintenanceHandler) DeleteTemplate(_ context.Context, id string) error {
	_ = id
	return errNotImplemented("CarMaintenanceHandler.DeleteTemplate")
}

func (h *CarMaintenanceHandler) GetRecords(_ context.Context, filter model.CarMaintenanceRecordFilter) ([]model.CarMaintenanceRecord, error) {
	_ = filter
	return nil, errNotImplemented("CarMaintenanceHandler.GetRecords")
}

func (h *CarMaintenanceHandler) CompleteRecord(_ context.Context, recordID string, data model.CarMaintenanceRecordComplete) error {
	_ = recordID
	_ = data
	return errNotImplemented("CarMaintenanceHandler.CompleteRecord")
}

func (h *CarMaintenanceHandler) GetReceiptImageUploadData(_ context.Context) (model.ImageUploadData, error) {
	return model.ImageUploadData{}, errNotImplemented("CarMaintenanceHandler.GetReceiptImageUploadData")
}
