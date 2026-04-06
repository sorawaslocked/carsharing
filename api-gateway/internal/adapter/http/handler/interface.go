package handler

import (
	"context"

	"github.com/sorawaslocked/car-rental-api-gateway/internal/model"
)

type AuthService interface {
	Register(ctx context.Context, data model.UserCreateData) (uint64, error)
	Login(ctx context.Context, cred model.Credentials) (model.Token, error)
	RefreshToken(ctx context.Context, refreshToken string) (model.Token, error)
	Logout(ctx context.Context, refreshToken string) error
}

type UserService interface {
	Insert(ctx context.Context, data model.UserCreateData) (uint64, error)
	FindOne(ctx context.Context, filter model.UserFilter) (model.User, error)
	Find(ctx context.Context, filter model.UserFilter) ([]model.User, error)
	Update(ctx context.Context, filter model.UserFilter, data model.UserUpdateData) error
	Delete(ctx context.Context, filter model.UserFilter) error
	Me(ctx context.Context) (model.User, error)
	SendActivationCode(ctx context.Context) error
	CheckActivationCode(ctx context.Context, code string) error
}

type CarModelService interface {
	Create(ctx context.Context, data model.CarModelCreate) (string, error)
	Get(ctx context.Context, id string) (model.CarModel, error)
	GetAll(ctx context.Context, filter model.CarModelFilter) ([]model.CarModel, error)
	Update(ctx context.Context, id string, data model.CarModelUpdate) error
	Delete(ctx context.Context, id string) error

	GetImageUploadData(ctx context.Context) (model.ImageUploadData, error)
}

type CarService interface {
	Create(ctx context.Context, data model.CarCreate) (string, error)
	Get(ctx context.Context, id string) (model.Car, error)
	GetAll(ctx context.Context, filter model.CarFilter) ([]model.Car, error)
	Update(ctx context.Context, id string, data model.CarUpdate) error
	Delete(ctx context.Context, id string) error

	GetCarStatusLog(ctx context.Context, filter model.CarStatusLogFilter) ([]model.CarStatusLogEntry, error)
	GetCarFuelHistory(ctx context.Context, filter model.CarFuelReadingFilter) ([]model.CarFuelReading, error)

	GetImageUploadData(ctx context.Context) (model.ImageUploadData, error)
}

type CarInsuranceService interface {
	Create(ctx context.Context, data model.CarInsuranceCreate) (string, error)
	Get(ctx context.Context, id string) (model.CarInsurance, error)
	GetAll(ctx context.Context, filter model.CarInsuranceFilter) ([]model.CarInsurance, error)
	Update(ctx context.Context, id string, data model.CarInsuranceUpdate) error
	Delete(ctx context.Context, id string) error

	GetImageUploadData(ctx context.Context) (model.ImageUploadData, error)
}

type ZoneService interface {
	Create(ctx context.Context, data model.ZoneCreate) (string, error)
	Get(ctx context.Context, id string) (model.Zone, error)
	GetAll(ctx context.Context, filter model.ZoneFilter) ([]model.Zone, error)
	Update(ctx context.Context, id string, data model.ZoneUpdate) error
	Delete(ctx context.Context, id string) error
}

type CarMaintenanceService interface {
	CreateTemplate(ctx context.Context, data model.CarMaintenanceTemplateCreate) (string, error)
	GetTemplate(ctx context.Context, id string) (model.CarMaintenanceTemplate, error)
	GetAllTemplates(ctx context.Context, filter model.CarMaintenanceTemplateFilter) ([]model.CarMaintenanceTemplate, error)
	UpdateTemplate(ctx context.Context, id string, data model.CarMaintenanceTemplateUpdate) error
	DeleteTemplate(ctx context.Context, id string) error

	GetRecords(ctx context.Context, filter model.CarMaintenanceRecordFilter) ([]model.CarMaintenanceRecord, error)
	CompleteRecord(ctx context.Context, recordID string, data model.CarMaintenanceRecordComplete) error

	GetReceiptImageUploadData(ctx context.Context) (model.ImageUploadData, error)
}
