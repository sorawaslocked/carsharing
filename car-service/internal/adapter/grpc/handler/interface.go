package handler

import (
	"context"

	"carsharing/car-service/internal/model"
	"carsharing/car-service/internal/validation"
	sharedmodel "carsharing/shared/model"
)

type TelemetrySubscriber interface {
	SubscribeUpdates(carID string) (<-chan model.TelemetryUpdate, func())
}

type StatusSubscriber interface {
	SubscribeStatusUpdates(carID string) (<-chan model.CarStatusUpdate, func())
}

type Pinger interface {
	Ping(ctx context.Context) error
}

type CarModelService interface {
	Create(ctx context.Context, createInput validation.CarModelCreate) (string, error)
	Get(ctx context.Context, id string) (model.CarModel, error)
	List(ctx context.Context, filterInput validation.CarModelFilter) ([]model.CarModel, error)
	Update(ctx context.Context, id string, updateInput validation.CarModelUpdate) error
	Delete(ctx context.Context, id string) error
	GetImageUploadData(ctx context.Context) (sharedmodel.ImageUploadData, error)
}

type CarService interface {
	Create(ctx context.Context, createInput validation.CarCreate) (string, error)
	Get(ctx context.Context, id string) (model.Car, error)
	List(ctx context.Context, filterInput validation.CarFilter) ([]model.Car, error)
	Update(ctx context.Context, id string, updateInput validation.CarUpdate) error
	UpdateCarStatus(ctx context.Context, id string, statusInput validation.CarStatusUpdate) error
	UpdateCarTelemetry(ctx context.Context, id string, data validation.CarTelemetryUpdate) error
	Delete(ctx context.Context, id string) error
	GetImageUploadData(ctx context.Context) (sharedmodel.ImageUploadData, error)
	ListCarStatusHistory(ctx context.Context, filter validation.CarStatusReadingFilter) ([]model.CarStatusReading, error)
	ListCarTelemetryHistory(ctx context.Context, filter validation.TelemetryReadingFilter) ([]model.TelemetryReading, error)
}

type CarInsuranceService interface {
	Create(ctx context.Context, createInput validation.CarInsuranceCreate) (string, error)
	Get(ctx context.Context, id string) (model.CarInsurance, error)
	List(ctx context.Context, filterInput validation.CarInsuranceFilter) ([]model.CarInsurance, error)
	Update(ctx context.Context, id string, updateInput validation.CarInsuranceUpdate) error
	Delete(ctx context.Context, id string) error
	GetImageUploadData(ctx context.Context) (sharedmodel.ImageUploadData, error)
}

type CarMaintenanceService interface {
	CreateTemplate(ctx context.Context, createInput validation.CarMaintenanceTemplateCreate) (string, error)
	GetTemplate(ctx context.Context, id string) (model.CarMaintenanceTemplate, error)
	ListTemplates(ctx context.Context, filterInput validation.CarMaintenanceTemplateFilter) ([]model.CarMaintenanceTemplate, error)
	UpdateTemplate(ctx context.Context, id string, updateInput validation.CarMaintenanceTemplateUpdate) error
	DeleteTemplate(ctx context.Context, id string) error
	AssignCarTemplate(ctx context.Context, data validation.CarTemplateAssign) error
	ListRecords(ctx context.Context, filterInput validation.CarMaintenanceRecordFilter) ([]model.CarMaintenanceRecord, error)
	CompleteRecord(ctx context.Context, id string, completeInput validation.CarMaintenanceRecordComplete) error
	GetReceiptImageUploadData(ctx context.Context) (sharedmodel.ImageUploadData, error)
}

type MaintenanceEventSubscriber interface {
	SubscribeMaintenanceEvents() (<-chan model.CarMaintenanceEvent, func())
}

type ZoneService interface {
	Create(ctx context.Context, createInput validation.ZoneCreate) (string, error)
	Get(ctx context.Context, id string) (model.Zone, error)
	List(ctx context.Context, filterInput validation.ZoneFilter) ([]model.Zone, error)
	Update(ctx context.Context, id string, updateInput validation.ZoneUpdate) error
	Delete(ctx context.Context, id string) error
	GetZonePricing(ctx context.Context, data validation.ZoneGetPricing) (int32, error)
}
