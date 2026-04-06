package handler

import (
	"context"

	"github.com/sorawaslocked/car-rental-api-gateway/internal/model"
)

type ZoneHandler struct {
	// client fleetsvc.ZoneServiceClient
}

func NewZoneHandler() *ZoneHandler {
	return &ZoneHandler{}
}

func (h *ZoneHandler) Create(_ context.Context, data model.ZoneCreate) (string, error) {
	_ = data
	return "", errNotImplemented("ZoneHandler.Create")
}

func (h *ZoneHandler) Get(_ context.Context, id string) (model.Zone, error) {
	_ = id
	return model.Zone{}, errNotImplemented("ZoneHandler.Get")
}

func (h *ZoneHandler) GetAll(_ context.Context, filter model.ZoneFilter) ([]model.Zone, error) {
	_ = filter
	return nil, errNotImplemented("ZoneHandler.GetAll")
}

func (h *ZoneHandler) Update(_ context.Context, id string, data model.ZoneUpdate) error {
	_ = id
	_ = data
	return errNotImplemented("ZoneHandler.Update")
}

func (h *ZoneHandler) Delete(_ context.Context, id string) error {
	_ = id
	return errNotImplemented("ZoneHandler.Delete")
}
