package handler

import (
	"context"

	"github.com/sorawaslocked/car-rental-api-gateway/internal/model"
)

type HealthChecker interface {
	Health(ctx context.Context) (model.ServiceHealth, error)
}

type UserService interface {
	Create(ctx context.Context, data model.UserCreate) (string, error)
	Get(ctx context.Context, id string) (model.User, error)
	List(ctx context.Context, filter model.UserFilter) ([]model.User, error)
	Update(ctx context.Context, id string, data model.UserUpdate) error
	Delete(ctx context.Context, id string) error

	Register(ctx context.Context, data model.UserCreate) (string, error)
	SignIn(ctx context.Context, cred model.Credentials) (model.AccessToken, model.RefreshToken, error)
	RefreshToken(ctx context.Context, refreshToken string) (model.AccessToken, model.RefreshToken, error)
	SignOut(ctx context.Context) error
	Me(ctx context.Context) (model.User, error)

	SendActivationCode(ctx context.Context) error
	CheckActivationCode(ctx context.Context, code string) error

	CreateDocument(ctx context.Context, objectKey, imageType string) (string, error)
	GetUploadDocumentData(ctx context.Context, imageType string) (model.ImageUploadData, error)
	GetProcessedDocumentsForUser(ctx context.Context, userID string) ([]model.Document, error)
	CheckDocument(ctx context.Context, docID string, status string, documentError *string) error
}

type CarModelService interface {
	Create(ctx context.Context, data model.CarModelCreate) (string, error)
	Get(ctx context.Context, id string) (model.CarModel, error)
	List(ctx context.Context, filter model.CarModelFilter) ([]model.CarModel, error)
	Update(ctx context.Context, id string, data model.CarModelUpdate) error
	Delete(ctx context.Context, id string) error

	GetImageUploadData(ctx context.Context) (model.ImageUploadData, error)
}

type CarService interface {
	Create(ctx context.Context, data model.CarCreate) (string, error)
	Get(ctx context.Context, id string) (model.Car, error)
	List(ctx context.Context, filter model.CarFilter) ([]model.Car, error)
	Update(ctx context.Context, id string, data model.CarUpdate) error
	Delete(ctx context.Context, id string) error

	ElevatedUpdate(ctx context.Context, carID string, data model.CarElevatedUpdate) error

	GetCarStatusHistory(ctx context.Context, carID string, filter model.CarStatusReadingFilter) ([]model.CarStatusReading, error)
	GetCarFuelHistory(ctx context.Context, carID string, filter model.CarFuelReadingFilter) ([]model.CarFuelReading, error)
	GetCarLocationHistory(ctx context.Context, carID string, filter model.CarLocationReadingFilter) ([]model.CarLocationReading, error)
	GetCarBatteryHistory(ctx context.Context, carID string, filter model.CarBatteryReadingFilter) ([]model.CarBatteryReading, error)
	GetCarMileageHistory(ctx context.Context, carID string, filter model.CarMileageReadingFilter) ([]model.CarMileageReading, error)

	GetImageUploadData(ctx context.Context) (model.ImageUploadData, error)
}

type CarInsuranceService interface {
	Create(ctx context.Context, data model.CarInsuranceCreate) (string, error)
	Get(ctx context.Context, id string) (model.CarInsurance, error)
	List(ctx context.Context, filter model.CarInsuranceFilter) ([]model.CarInsurance, error)
	Update(ctx context.Context, id string, data model.CarInsuranceUpdate) error
	Delete(ctx context.Context, id string) error

	GetImageUploadData(ctx context.Context) (model.ImageUploadData, error)
}

type ZoneService interface {
	Create(ctx context.Context, data model.ZoneCreate) (string, error)
	Get(ctx context.Context, id string) (model.Zone, error)
	List(ctx context.Context, filter model.ZoneFilter) ([]model.Zone, error)
	Update(ctx context.Context, id string, data model.ZoneUpdate) error
	Delete(ctx context.Context, id string) error
}

type CarMaintenanceService interface {
	CreateTemplate(ctx context.Context, data model.CarMaintenanceTemplateCreate) (string, error)
	GetTemplate(ctx context.Context, id string) (model.CarMaintenanceTemplate, error)
	ListTemplates(ctx context.Context, filter model.CarMaintenanceTemplateFilter) ([]model.CarMaintenanceTemplate, error)
	UpdateTemplate(ctx context.Context, id string, data model.CarMaintenanceTemplateUpdate) error
	DeleteTemplate(ctx context.Context, id string) error

	ListRecords(ctx context.Context, filter model.CarMaintenanceRecordFilter) ([]model.CarMaintenanceRecord, error)
	CompleteRecord(ctx context.Context, recordID string, data model.CarMaintenanceRecordComplete) error

	GetReceiptImageUploadData(ctx context.Context) (model.ImageUploadData, error)
}
