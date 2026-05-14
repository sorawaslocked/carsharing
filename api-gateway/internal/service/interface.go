package service

import (
	"context"
	"time"

	"github.com/sorawaslocked/car-rental-api-gateway/internal/model"
)

type TokenManager interface {
	GenerateAccessToken(userID string) (token string, exp time.Time, err error)
	GenerateRefreshToken(userID string) (token string, exp time.Time, err error)
	ParseToken(token string) (userID string, exp time.Time, err error)
}

type UserSessionCache interface {
	IsSignedIn(ctx context.Context, userID, deviceID string) (bool, error)
	SetSignedIn(ctx context.Context, userID, deviceID string, loggedIn bool) error
}

type UserPresenter interface {
	Create(ctx context.Context, data model.UserCreate) (string, error)
	Get(ctx context.Context, id string) (model.User, error)
	List(ctx context.Context, filter model.UserFilter) ([]model.User, error)
	Update(ctx context.Context, id string, data model.UserUpdate) error
	Delete(ctx context.Context, id string) error

	Register(ctx context.Context, data model.UserCreate) (string, error)
	SignIn(ctx context.Context, creds model.Credentials) (id string, err error)

	SendActivationCode(ctx context.Context) error
	CheckActivationCode(ctx context.Context, code string) error

	CreateDocument(ctx context.Context, objectKey, imageType string) (string, error)
	GetDocumentImageUploadData(ctx context.Context, imageType string) (model.ImageUploadData, error)
	GetProfileImageUploadData(ctx context.Context) (model.ImageUploadData, error)
	GetProcessedDocumentsForUser(ctx context.Context, userID string) ([]model.Document, error)
	CheckDocument(ctx context.Context, docID string, status string, documentError *string) error
}

type CarModelPresenter interface {
	Create(ctx context.Context, data model.CarModelCreate) (string, error)
	Get(ctx context.Context, id string) (model.CarModel, error)
	List(ctx context.Context, filter model.CarModelFilter) ([]model.CarModel, error)
	Update(ctx context.Context, id string, data model.CarModelUpdate) error
	Delete(ctx context.Context, id string) error

	GetImageUploadData(ctx context.Context) (model.ImageUploadData, error)
}

type CarPresenter interface {
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

	StreamCarsWithFilter(ctx context.Context, filter model.CarFilter, send func([]model.SlimCar) error) error
	StreamCarTelemetry(ctx context.Context, carID string, send func(model.CarTelemetryEvent) error) error
}

type PricingRulePresenter interface {
	Create(ctx context.Context, data model.PricingRuleCreate) (string, error)
	Get(ctx context.Context, id string) (model.PricingRule, error)
	List(ctx context.Context, filter model.PricingRuleFilter) ([]model.PricingRule, error)
	Update(ctx context.Context, id string, data model.PricingRuleUpdate) error
	Delete(ctx context.Context, id string) error
}

type CarInsurancePresenter interface {
	Create(ctx context.Context, data model.CarInsuranceCreate) (string, error)
	Get(ctx context.Context, id string) (model.CarInsurance, error)
	List(ctx context.Context, filter model.CarInsuranceFilter) ([]model.CarInsurance, error)
	Update(ctx context.Context, id string, data model.CarInsuranceUpdate) error
	Delete(ctx context.Context, id string) error

	GetImageUploadData(ctx context.Context) (model.ImageUploadData, error)
}

type ZonePresenter interface {
	Create(ctx context.Context, data model.ZoneCreate) (string, error)
	Get(ctx context.Context, id string) (model.Zone, error)
	List(ctx context.Context, filter model.ZoneFilter) ([]model.Zone, error)
	Update(ctx context.Context, id string, data model.ZoneUpdate) error
	Delete(ctx context.Context, id string) error
}

type BookingPresenter interface {
	Create(ctx context.Context, data model.BookingCreate) (string, error)
	Get(ctx context.Context, id string) (model.Booking, error)
	List(ctx context.Context, filter model.BookingFilter) ([]model.Booking, error)
	Cancel(ctx context.Context, id string) error
	UpdateStatus(ctx context.Context, id string, data model.BookingStatusUpdate) error
	GetStatusHistory(ctx context.Context, id string, filter model.BookingStatusReadingFilter) ([]model.BookingStatusReading, error)
}

type TripPresenter interface {
	Start(ctx context.Context, bookingID string) (string, error)
	Get(ctx context.Context, id string) (model.Trip, error)
	List(ctx context.Context, filter model.TripFilter) ([]model.Trip, error)
	End(ctx context.Context, id string) error
	Cancel(ctx context.Context, id string, reason *string) error
	GetSummary(ctx context.Context, id string) (model.TripSummary, error)
	GetStatusHistory(ctx context.Context, id string, filter model.TripStatusReadingFilter) ([]model.TripStatusReading, error)

	StreamTripLiveFeed(ctx context.Context, tripID string, send func(model.TripLiveFeed) error) error
}

type CarMaintenancePresenter interface {
	CreateTemplate(ctx context.Context, data model.CarMaintenanceTemplateCreate) (string, error)
	GetTemplate(ctx context.Context, id string) (model.CarMaintenanceTemplate, error)
	ListTemplates(ctx context.Context, filter model.CarMaintenanceTemplateFilter) ([]model.CarMaintenanceTemplate, error)
	UpdateTemplate(ctx context.Context, id string, data model.CarMaintenanceTemplateUpdate) error
	DeleteTemplate(ctx context.Context, id string) error

	ListRecords(ctx context.Context, filter model.CarMaintenanceRecordFilter) ([]model.CarMaintenanceRecord, error)
	CompleteRecord(ctx context.Context, recordID string, data model.CarMaintenanceRecordComplete) error

	GetReceiptImageUploadData(ctx context.Context) (model.ImageUploadData, error)
}
