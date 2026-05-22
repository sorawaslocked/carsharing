package handler

import (
	"context"

	"carsharing/car-service/internal/model"
	"carsharing/car-service/internal/validation"
	sharedmodel "carsharing/shared/model"
)

type TelematicsSubscriber interface {
	SubscribeCarStream(ctx context.Context, carID string) (<-chan model.TelematicsUpdate, error)
}

type Pinger interface {
	Ping(ctx context.Context) error
}

type CarModelService interface {
	Create(ctx context.Context, createInput validation.CarModelCreate) (string, error)
	Get(ctx context.Context, id string) (model.CarModel, error)
	GetAll(ctx context.Context, filterInput validation.CarModelFilter) ([]model.CarModel, error)
	Update(ctx context.Context, id string, updateInput validation.CarModelUpdate) error
	Delete(ctx context.Context, id string) error
	GetImageUploadData(ctx context.Context) (sharedmodel.ImageUploadData, error)
}

type CarService interface {
	Create(ctx context.Context, createInput validation.CarCreate) (string, error)
	Get(ctx context.Context, id string) (model.Car, error)
	GetAll(ctx context.Context, filterInput validation.CarFilter) ([]model.Car, error)
	Update(ctx context.Context, id string, updateInput validation.CarUpdate) error
	UpdateCarStatus(ctx context.Context, id string, statusInput validation.CarStatusUpdate) error
	UpdateCarTelemetry(ctx context.Context, id string, input model.CarTelematicsUpdateInput) error
	Delete(ctx context.Context, id string) error
	GetImageUploadData(ctx context.Context) (sharedmodel.ImageUploadData, error)
	GetCarStatusHistory(ctx context.Context, filter model.CarStatusLogFilter) ([]model.CarStatusLogEntry, error)
	GetCarFuelHistory(ctx context.Context, filter model.TelematicsEventFilter) ([]model.CarTelematicsEvent, error)
	GetCarLocationHistory(ctx context.Context, filter model.TelematicsEventFilter) ([]model.CarTelematicsEvent, error)
	GetCarBatteryHistory(ctx context.Context, filter model.TelematicsEventFilter) ([]model.CarTelematicsEvent, error)
	GetCarMileageHistory(ctx context.Context, filter model.TelematicsEventFilter) ([]model.CarTelematicsEvent, error)
}

type CarInsuranceService interface {
	Create(ctx context.Context, createInput validation.CarInsuranceCreate) (string, error)
	Get(ctx context.Context, id string) (model.CarInsurance, error)
	GetAll(ctx context.Context, filterInput validation.CarInsuranceFilter) ([]model.CarInsurance, error)
	Update(ctx context.Context, id string, updateInput validation.CarInsuranceUpdate) error
	Delete(ctx context.Context, id string) error
	GetImageUploadData(ctx context.Context) (sharedmodel.ImageUploadData, error)
}

type CarMaintenanceService interface {
	CreateTemplate(ctx context.Context, createInput validation.CarMaintenanceTemplateCreate) (string, error)
	GetTemplate(ctx context.Context, id string) (model.CarMaintenanceTemplate, error)
	GetAllTemplates(ctx context.Context, filterInput validation.CarMaintenanceTemplateFilter) ([]model.CarMaintenanceTemplate, error)
	UpdateTemplate(ctx context.Context, id string, updateInput validation.CarMaintenanceTemplateUpdate) error
	DeleteTemplate(ctx context.Context, id string) error
	GetRecords(ctx context.Context, filterInput validation.CarMaintenanceRecordFilter) ([]model.CarMaintenanceRecord, error)
	CompleteRecord(ctx context.Context, id string, completeInput validation.CarMaintenanceRecordComplete) error
	GetReceiptImageUploadData(ctx context.Context) (sharedmodel.ImageUploadData, error)
}

type ZoneService interface {
	Create(ctx context.Context, createInput validation.ZoneCreate) (string, error)
	Get(ctx context.Context, id string) (model.Zone, error)
	GetAll(ctx context.Context, filterInput validation.ZoneFilter) ([]model.Zone, error)
	Update(ctx context.Context, id string, updateInput validation.ZoneUpdate) error
	Delete(ctx context.Context, id string) error
}
