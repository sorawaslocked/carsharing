package handler

import (
	"context"

	"carsharing/car-service/internal/model"
)

type TelematicsSubscriber interface {
	SubscribeCarStream(ctx context.Context, carID string) (<-chan model.TelematicsUpdate, error)
}

type Pinger interface {
	Ping(ctx context.Context) error
}

type CarModelService interface {
	Create(ctx context.Context, createInput model.CarModelCreateInput) (string, error)
	Get(ctx context.Context, id string) (model.CarModel, error)
	GetAll(ctx context.Context, filterInput model.CarModelFilterInput) ([]model.CarModel, error)
	Update(ctx context.Context, id string, updateInput model.CarModelUpdateInput) error
	Delete(ctx context.Context, id string) error
	GetImageUploadData(ctx context.Context) (model.ImageUploadData, error)
}

type CarService interface {
	Create(ctx context.Context, createInput model.CarCreateInput) (string, error)
	Get(ctx context.Context, id string) (model.Car, error)
	GetAll(ctx context.Context, filterInput model.CarFilterInput) ([]model.Car, error)
	Update(ctx context.Context, id string, updateInput model.CarUpdateInput) error
	UpdateCarStatus(ctx context.Context, id string, statusInput model.CarStatusUpdateInput) error
	UpdateCarTelemetry(ctx context.Context, id string, input model.CarTelematicsUpdateInput) error
	Delete(ctx context.Context, id string) error
	GetImageUploadData(ctx context.Context) (model.ImageUploadData, error)
	GetCarStatusHistory(ctx context.Context, filter model.CarStatusLogFilter) ([]model.CarStatusLogEntry, error)
	GetCarFuelHistory(ctx context.Context, filter model.TelematicsEventFilter) ([]model.CarTelematicsEvent, error)
	GetCarLocationHistory(ctx context.Context, filter model.TelematicsEventFilter) ([]model.CarTelematicsEvent, error)
	GetCarBatteryHistory(ctx context.Context, filter model.TelematicsEventFilter) ([]model.CarTelematicsEvent, error)
	GetCarMileageHistory(ctx context.Context, filter model.TelematicsEventFilter) ([]model.CarTelematicsEvent, error)
}

type CarInsuranceService interface {
	Create(ctx context.Context, createInput model.CarInsuranceCreateInput) (string, error)
	Get(ctx context.Context, id string) (model.CarInsurance, error)
	GetAll(ctx context.Context, filterInput model.CarInsuranceFilterInput) ([]model.CarInsurance, error)
	Update(ctx context.Context, id string, updateInput model.CarInsuranceUpdateInput) error
	Delete(ctx context.Context, id string) error
	GetImageUploadData(ctx context.Context) (model.ImageUploadData, error)
}

type CarMaintenanceService interface {
	CreateTemplate(ctx context.Context, createInput model.CarMaintenanceTemplateCreateInput) (string, error)
	GetTemplate(ctx context.Context, id string) (model.CarMaintenanceTemplate, error)
	GetAllTemplates(ctx context.Context, filterInput model.CarMaintenanceTemplateFilterInput) ([]model.CarMaintenanceTemplate, error)
	UpdateTemplate(ctx context.Context, id string, updateInput model.CarMaintenanceTemplateUpdateInput) error
	DeleteTemplate(ctx context.Context, id string) error
	GetRecords(ctx context.Context, filterInput model.CarMaintenanceRecordFilterInput) ([]model.CarMaintenanceRecord, error)
	CompleteRecord(ctx context.Context, id string, completeInput model.CarMaintenanceRecordCompleteInput) error
	GetReceiptImageUploadData(ctx context.Context) (model.ImageUploadData, error)
}

type ZoneService interface {
	Create(ctx context.Context, createInput model.ZoneCreateInput) (string, error)
	Get(ctx context.Context, id string) (model.Zone, error)
	GetAll(ctx context.Context, filterInput model.ZoneFilterInput) ([]model.Zone, error)
	Update(ctx context.Context, id string, updateInput model.ZoneUpdateInput) error
	Delete(ctx context.Context, id string) error
}
