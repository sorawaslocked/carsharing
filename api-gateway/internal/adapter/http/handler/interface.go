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
	UpdateProfile(ctx context.Context, data model.UserProfileUpdate) error
	Delete(ctx context.Context, id string) error

	Register(ctx context.Context, data model.UserCreate) (string, error)
	SignIn(ctx context.Context, cred model.Credentials) (model.AccessToken, model.RefreshToken, error)
	RefreshToken(ctx context.Context, refreshToken string) (model.AccessToken, model.RefreshToken, error)
	SignOut(ctx context.Context) error
	GetProfile(ctx context.Context) (model.User, error)

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

	UpdateTelemetry(ctx context.Context, carID string, data model.CarTelemetryUpdate) error
	UpdateStatus(ctx context.Context, carID string, data model.CarStatusUpdate) error

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

type PricingRuleService interface {
	Create(ctx context.Context, data model.PricingRuleCreate) (string, error)
	Get(ctx context.Context, id string) (model.PricingRule, error)
	List(ctx context.Context, filter model.PricingRuleFilter) ([]model.PricingRule, error)
	Update(ctx context.Context, id string, data model.PricingRuleUpdate) error
	Delete(ctx context.Context, id string) error
}

type BookingService interface {
	Create(ctx context.Context, data model.BookingCreate) (string, error)
	Get(ctx context.Context, id string) (model.Booking, error)
	List(ctx context.Context, filter model.BookingFilter) ([]model.Booking, error)
	Start(ctx context.Context, id string) error
	Cancel(ctx context.Context, id string) error
	UpdateStatus(ctx context.Context, id string, data model.BookingStatusUpdate) error
	GetStatusHistory(ctx context.Context, id string, filter model.BookingStatusReadingFilter) ([]model.BookingStatusReading, error)
}

type TripService interface {
	Start(ctx context.Context, bookingID string) (string, error)
	Get(ctx context.Context, id string) (model.Trip, error)
	List(ctx context.Context, filter model.TripFilter) ([]model.Trip, error)
	End(ctx context.Context, id string) error
	Cancel(ctx context.Context, id string, reason *string) error
	GetSummary(ctx context.Context, id string) (model.TripSummary, error)
	GetStatusHistory(ctx context.Context, id string, filter model.TripStatusReadingFilter) ([]model.TripStatusReading, error)
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
