package service

import (
	"context"

	"carsharing/car-service/internal/model"
	sharedmodel "carsharing/shared/model"
)

type CarModelRepository interface {
	Insert(ctx context.Context, carModel model.CarModel) (string, error)
	FindByID(ctx context.Context, id string) (model.CarModel, error)
	Find(ctx context.Context, filter model.CarModelFilter) ([]model.CarModel, error)
	Update(ctx context.Context, id string, update model.CarModelUpdate) error
	Delete(ctx context.Context, id string) error
}

type CarRepository interface {
	Insert(ctx context.Context, car model.Car) (string, error)
	FindByID(ctx context.Context, id string) (model.Car, error)
	Find(ctx context.Context, filter model.CarFilter) ([]model.Car, error)
	Update(ctx context.Context, id string, update model.CarUpdate) error
	Delete(ctx context.Context, id string) error
}

type CarStatusReadingRepository interface {
	Insert(ctx context.Context, entry model.CarStatusReading) error
	Find(ctx context.Context, filter model.CarStatusReadingFilter) ([]model.CarStatusReading, error)
}

type ZoneRepository interface {
	Insert(ctx context.Context, zone model.Zone) (string, error)
	FindByID(ctx context.Context, id string) (model.Zone, error)
	FindByLocation(ctx context.Context, lat, lng float64) (*model.Zone, error)
	Find(ctx context.Context, filter model.ZoneFilter) ([]model.Zone, error)
	Update(ctx context.Context, id string, update model.ZoneUpdate) error
	Delete(ctx context.Context, id string) error
}

type CarInsuranceRepository interface {
	Insert(ctx context.Context, insurance model.CarInsurance) (string, error)
	FindByID(ctx context.Context, id string) (model.CarInsurance, error)
	Find(ctx context.Context, filter model.CarInsuranceFilter) ([]model.CarInsurance, error)
	Update(ctx context.Context, id string, update model.CarInsuranceUpdate) error
	Delete(ctx context.Context, id string) error
}

type CarMaintenanceTemplateRepository interface {
	Insert(ctx context.Context, template model.CarMaintenanceTemplate) (string, error)
	FindByID(ctx context.Context, id string) (model.CarMaintenanceTemplate, error)
	Find(ctx context.Context, filter model.CarMaintenanceTemplateFilter) ([]model.CarMaintenanceTemplate, error)
	Update(ctx context.Context, id string, update model.CarMaintenanceTemplateUpdate) error
	Delete(ctx context.Context, id string) error
}

type CarMaintenanceRecordRepository interface {
	Insert(ctx context.Context, record model.CarMaintenanceRecord) (string, error)
	FindByID(ctx context.Context, id string) (model.CarMaintenanceRecord, error)
	Find(ctx context.Context, filter model.CarMaintenanceRecordFilter) ([]model.CarMaintenanceRecord, error)
	Update(ctx context.Context, id string, update model.CarMaintenanceRecordUpdate) error
	UpdateWithServiceState(ctx context.Context, id string, update model.CarMaintenanceRecordUpdate, state model.CarServiceState) error
}

type CarServiceStateRepository interface {
	Upsert(ctx context.Context, state model.CarServiceState) error
	FindAll(ctx context.Context, filter model.CarServiceStateFilter) ([]model.CarServiceState, error)
}

type TelemetryReadingRepository interface {
	Insert(ctx context.Context, reading model.TelemetryReading) error
	Find(ctx context.Context, filter model.TelemetryReadingFilter) ([]model.TelemetryReading, error)
}

type TelemetryStreamClient interface {
	Subscribe(ctx context.Context, car model.Car) (<-chan model.TelemetryUpdate, error)
}

type CarCreatedNotifier interface {
	OnCarCreated(car model.Car)
}

type EventPublisher interface {
	PublishCarStatusUpdated(ctx context.Context, carID, fromStatus, toStatus string) error
}

type ObjectStorage interface {
	GetCarImageUploadData(ctx context.Context) (sharedmodel.ImageUploadData, error)
	GetCarModelImageUploadData(ctx context.Context) (sharedmodel.ImageUploadData, error)
	GetInsuranceImageUploadData(ctx context.Context) (sharedmodel.ImageUploadData, error)
	GetMaintenanceReceiptImageUploadData(ctx context.Context) (sharedmodel.ImageUploadData, error)
	GetPresignedURL(ctx context.Context, key string) (string, error)
}
