package service

import (
	"context"

	"github.com/sorawaslocked/car-rental-car-service/internal/model"
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

type CarStatusLogRepository interface {
	Insert(ctx context.Context, entry model.CarStatusLogEntry) error
	Find(ctx context.Context, filter model.CarStatusLogFilter) ([]model.CarStatusLogEntry, error)
}

type ZoneRepository interface {
	Insert(ctx context.Context, zone model.Zone) (string, error)
	FindByID(ctx context.Context, id string) (model.Zone, error)
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
}

type CarServiceStateRepository interface {
	Upsert(ctx context.Context, state model.CarServiceState) error
	FindAll(ctx context.Context, filter model.CarServiceStateFilter) ([]model.CarServiceState, error)
}

type TelematicsRepository interface {
	InsertEvent(ctx context.Context, event model.CarTelematicsEvent) error
	FindEvents(ctx context.Context, filter model.TelematicsEventFilter) ([]model.CarTelematicsEvent, error)
}

type TelematicsStreamClient interface {
	Subscribe(ctx context.Context, car model.Car) (<-chan model.TelematicsUpdate, error)
}

type CarCreatedNotifier interface {
	OnCarCreated(car model.Car)
}

type EventPublisher interface {
	PublishCarStatusUpdated(ctx context.Context, carID, fromStatus, toStatus string) error
}

type ObjectStorage interface {
	GetCarImageUploadData(ctx context.Context) (model.ImageUploadData, error)
	GetCarModelImageUploadData(ctx context.Context) (model.ImageUploadData, error)
	GetInsuranceImageUploadData(ctx context.Context) (model.ImageUploadData, error)
	GetMaintenanceReceiptImageUploadData(ctx context.Context) (model.ImageUploadData, error)
	GetPresignedURL(ctx context.Context, key string) (string, error)
}
