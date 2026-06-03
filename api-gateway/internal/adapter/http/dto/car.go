package dto

import (
	"strconv"
	"time"

	"carsharing/api-gateway/internal/model"
	sharedmodel "carsharing/shared/model"

	"github.com/gin-gonic/gin"
)

type CarResponse struct {
	Car Car `json:"car"`
}

type CarsResponse struct {
	Cars []Car `json:"cars"`
}

type CarStatusHistoryResponse struct {
	StatusHistory []CarStatusReading `json:"statusHistory"`
}

type CarTelemetryHistoryResponse struct {
	TelemetryHistory []CarTelemetryReading `json:"telemetryHistory"`
}

type Car struct {
	ID               string    `json:"id"`
	ModelID          string    `json:"modelID"`
	VIN              string    `json:"vin"`
	LicensePlate     string    `json:"licensePlate"`
	Color            string    `json:"color"`
	YearManufactured int16     `json:"yearManufactured"`
	MileageKM        int64     `json:"mileageKm"`
	FuelLevel        *float32  `json:"fuelLevel,omitempty"`
	BatteryLevel     *float32  `json:"batteryLevel,omitempty"`
	Location         location  `json:"location"`
	TelemetryID      string    `json:"telemetryId"`
	FuelStatus       string    `json:"fuelStatus"`
	Status           string    `json:"status" validate:"oneof=available reserved in_use maintenance out_of_service"`
	IsRetired        bool      `json:"isRetired"`
	Notes            *string   `json:"notes,omitempty"`
	ImageStorageUrls []string  `json:"imageStorageUrls,omitempty"`
	LastSeenAt       time.Time `json:"lastSeenAt"`
	CreatedAt        time.Time `json:"createdAt"`
	UpdatedAt        time.Time `json:"updatedAt"`
}

type CarStatusReading struct {
	ID         string         `json:"id"`
	CarID      string         `json:"carID"`
	FromStatus string         `json:"fromStatus" validate:"oneof=available reserved in_use maintenance out_of_service"`
	ToStatus   string         `json:"toStatus" validate:"oneof=available reserved in_use maintenance out_of_service"`
	ActorType  string         `json:"actorType" validate:"oneof=user system telemetry"`
	ActorID    *string        `json:"actorID,omitempty"`
	Reason     *string        `json:"reason,omitempty"`
	Metadata   map[string]any `json:"metadata,omitempty"`
	RecordedAt time.Time      `json:"recordedAt"`
}

type CarTelemetryReading struct {
	ID           string         `json:"id"`
	CarID        string         `json:"carID"`
	FuelPct      *float32       `json:"fuelPct,omitempty"`
	FuelRawPct   *float32       `json:"fuelRawPct,omitempty"`
	BatteryLevel *float32       `json:"batteryLevel,omitempty"`
	MileageKM    *int64         `json:"mileageKm,omitempty"`
	Location     *location      `json:"location,omitempty"`
	ActorType    string         `json:"actorType" validate:"oneof=user system telemetry"`
	ActorID      *string        `json:"actorID,omitempty"`
	Reason       *string        `json:"reason,omitempty"`
	Metadata     map[string]any `json:"metadata,omitempty"`
	RecordedAt   time.Time      `json:"recordedAt"`
}

type CarCreateRequest struct {
	ModelID          string   `json:"modelID" binding:"required"`
	VIN              string   `json:"vin" binding:"required"`
	LicensePlate     string   `json:"licensePlate" binding:"required"`
	Color            string   `json:"color" binding:"required"`
	YearManufactured int16    `json:"yearManufactured" binding:"required"`
	TelemetryID      string   `json:"telemetryId" binding:"required"`
	MileageKM        *int64   `json:"mileageKm"`
	FuelLevel        *float32 `json:"fuelLevel"`
	BatteryLevel     *float32 `json:"batteryLevel"`
	Latitude         *float64 `json:"latitude"`
	Longitude        *float64 `json:"longitude"`
	Notes            *string  `json:"notes"`
}

type CarUpdateRequest struct {
	ModelID      *string  `json:"modelID"`
	LicensePlate *string  `json:"licensePlate"`
	Color        *string  `json:"color"`
	TelemetryID  *string  `json:"telemetryId"`
	IsRetired    *bool    `json:"isRetired"`
	Notes        *string  `json:"notes"`
	ImageKeys    []string `json:"imageKeys"`
}

type CarTelemetryUpdateRequest struct {
	MileageKM    *int64         `json:"mileageKm"`
	FuelLevel    *float32       `json:"fuelLevel"`
	BatteryLevel *float32       `json:"batteryLevel"`
	Latitude     *float64       `json:"latitude"`
	Longitude    *float64       `json:"longitude"`
	Reason       string         `json:"reason" binding:"required"`
	Metadata     map[string]any `json:"metadata"`
}

type CarStatusUpdateRequest struct {
	Status   string         `json:"status" binding:"required,oneof=available reserved in_use maintenance out_of_service"`
	Reason   string         `json:"reason" binding:"required"`
	Metadata map[string]any `json:"metadata"`
}

func FromCarCreateRequest(ctx *gin.Context) (model.CarCreate, error) {
	var req CarCreateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		return model.CarCreate{}, err
	}

	create := model.CarCreate{
		ModelID:          req.ModelID,
		VIN:              req.VIN,
		LicensePlate:     req.LicensePlate,
		Color:            req.Color,
		YearManufactured: req.YearManufactured,
		TelemetryID:      req.TelemetryID,
		MileageKM:        req.MileageKM,
		FuelLevel:        req.FuelLevel,
		BatteryLevel:     req.BatteryLevel,
		Notes:            req.Notes,
	}
	if req.Latitude != nil && req.Longitude != nil {
		create.Location = &sharedmodel.Location{
			Latitude:  *req.Latitude,
			Longitude: *req.Longitude,
		}
	}
	return create, nil
}

func FromCarUpdateRequest(ctx *gin.Context) (model.CarUpdate, error) {
	var req CarUpdateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		return model.CarUpdate{}, err
	}

	return model.CarUpdate{
		ModelID:      req.ModelID,
		LicensePlate: req.LicensePlate,
		Color:        req.Color,
		TelemetryID:  req.TelemetryID,
		IsRetired:    req.IsRetired,
		Notes:        req.Notes,
		ImageKeys:    req.ImageKeys,
	}, nil
}

func FromCarTelemetryUpdateRequest(ctx *gin.Context) (model.CarTelemetryUpdate, error) {
	var req CarTelemetryUpdateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		return model.CarTelemetryUpdate{}, err
	}

	update := model.CarTelemetryUpdate{
		MileageKM:    req.MileageKM,
		FuelLevel:    req.FuelLevel,
		BatteryLevel: req.BatteryLevel,
		Reason:       req.Reason,
		Metadata:     req.Metadata,
	}
	if req.Latitude != nil && req.Longitude != nil {
		update.Location = &sharedmodel.Location{
			Latitude:  *req.Latitude,
			Longitude: *req.Longitude,
		}
	}

	return update, nil
}

func FromCarStatusUpdateRequest(ctx *gin.Context) (model.CarStatusUpdate, error) {
	var req CarStatusUpdateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		return model.CarStatusUpdate{}, err
	}

	return model.CarStatusUpdate{
		Status:   req.Status,
		Reason:   req.Reason,
		Metadata: req.Metadata,
	}, nil
}

func CarFilterFromCtx(ctx *gin.Context) (model.CarFilter, error) {
	f := model.CarFilter{}

	if v := ctx.Query("brand"); v != "" {
		f.Brand = &v
	}
	if v := ctx.Query("model"); v != "" {
		f.Model = &v
	}
	if v := ctx.Query("fuelType"); v != "" {
		f.FuelType = &v
	}
	if v := ctx.Query("transmission"); v != "" {
		f.Transmission = &v
	}
	if v := ctx.Query("bodyType"); v != "" {
		f.BodyType = &v
	}
	if v := ctx.Query("class"); v != "" {
		f.Class = &v
	}
	if v := ctx.Query("minSeats"); v != "" {
		vInt, err := strconv.ParseInt(v, 10, 8)
		if err != nil {
			return model.CarFilter{}, model.ErrInvalidQueryParam
		}

		minSeats := int8(vInt)
		f.MinSeats = &minSeats
	}
	if lat := ctx.Query("latitude"); lat != "" {
		if lng := ctx.Query("longitude"); lng != "" {
			latF, err := strconv.ParseFloat(lat, 64)
			if err != nil {
				return model.CarFilter{}, model.ErrInvalidQueryParam
			}
			lngF, err := strconv.ParseFloat(lng, 64)
			if err != nil {
				return model.CarFilter{}, model.ErrInvalidQueryParam
			}
			f.Location = &sharedmodel.Location{Latitude: latF, Longitude: lngF}
		}
	}
	if v := ctx.Query("radiusM"); v != "" {
		vInt, err := strconv.ParseInt(v, 10, 32)
		if err != nil {
			return model.CarFilter{}, model.ErrInvalidQueryParam
		}

		radiusM := int32(vInt)
		f.RadiusM = &radiusM
	}
	if v := ctx.Query("minFuelLevel"); v != "" {
		vFloat, err := strconv.ParseFloat(v, 32)
		if err != nil {
			return model.CarFilter{}, model.ErrInvalidQueryParam
		}

		minFuelLevel := float32(vFloat)
		f.MinFuelLevel = &minFuelLevel
	}
	if v := ctx.Query("status"); v != "" {
		f.Status = &v
	}

	p, err := pagination(ctx)
	if err != nil {
		return model.CarFilter{}, model.ErrInvalidQueryParam
	}

	f.Pagination = p

	return f, nil
}

func carTimeRangeFilter(ctx *gin.Context) (tr *sharedmodel.TimeRange, p *sharedmodel.Pagination, err error) {
	var timeRange sharedmodel.TimeRange
	hasRange := false
	if v := ctx.Query("from"); v != "" {
		t, parseErr := time.Parse("2006-01-02", v)
		if parseErr != nil {
			return nil, nil, model.ErrInvalidQueryParam
		}
		timeRange.From = t
		hasRange = true
	}
	if v := ctx.Query("to"); v != "" {
		t, parseErr := time.Parse("2006-01-02", v)
		if parseErr != nil {
			return nil, nil, model.ErrInvalidQueryParam
		}
		timeRange.To = t
		hasRange = true
	}
	if hasRange {
		tr = &timeRange
	}

	p, err = pagination(ctx)
	if err != nil {
		return nil, nil, model.ErrInvalidQueryParam
	}

	return tr, p, nil
}

func CarStatusReadingFilterFromCtx(ctx *gin.Context) (model.CarStatusReadingFilter, error) {
	tr, p, err := carTimeRangeFilter(ctx)
	if err != nil {
		return model.CarStatusReadingFilter{}, err
	}

	return model.CarStatusReadingFilter{TimeRange: tr, Pagination: p}, nil
}

func CarTelemetryReadingFilterFromCtx(ctx *gin.Context) (model.CarTelemetryReadingFilter, error) {
	tr, p, err := carTimeRangeFilter(ctx)
	if err != nil {
		return model.CarTelemetryReadingFilter{}, err
	}

	return model.CarTelemetryReadingFilter{TimeRange: tr, Pagination: p}, nil
}

func ToCarResponse(m model.Car) Car {
	return Car{
		ID:               m.ID,
		ModelID:          m.ModelID,
		VIN:              m.VIN,
		LicensePlate:     m.LicensePlate,
		Color:            m.Color,
		YearManufactured: m.YearManufactured,
		MileageKM:        m.MileageKM,
		FuelLevel:        m.FuelLevel,
		BatteryLevel:     m.BatteryLevel,
		Location: location{
			Latitude:  m.Location.Latitude,
			Longitude: m.Location.Longitude,
		},
		TelemetryID:      m.TelemetryID,
		FuelStatus:       m.FuelStatus,
		Status:           m.Status,
		IsRetired:        m.IsRetired,
		Notes:            m.Notes,
		ImageStorageUrls: m.ImageURLs,
		LastSeenAt:       m.LastSeenAt,
		CreatedAt:        m.CreatedAt,
		UpdatedAt:        m.UpdatedAt,
	}
}

func ToCarStatusReadingResponse(m model.CarStatusReading) CarStatusReading {
	return CarStatusReading{
		ID:         m.ID,
		CarID:      m.CarID,
		FromStatus: m.FromStatus,
		ToStatus:   m.ToStatus,
		ActorType:  m.ActorType,
		ActorID:    m.ActorID,
		Reason:     m.Reason,
		Metadata:   m.Metadata,
		RecordedAt: m.RecordedAt,
	}
}

func ToCarTelemetryReadingResponse(m model.CarTelemetryReading) CarTelemetryReading {
	r := CarTelemetryReading{
		ID:           m.ID,
		CarID:        m.CarID,
		FuelPct:      m.FuelPct,
		FuelRawPct:   m.FuelRawPct,
		BatteryLevel: m.BatteryLevel,
		MileageKM:    m.MileageKM,
		ActorType:    m.ActorType,
		ActorID:      m.ActorID,
		Reason:       m.Reason,
		Metadata:     m.Metadata,
		RecordedAt:   m.RecordedAt,
	}
	if m.Location != nil {
		r.Location = &location{
			Latitude:  m.Location.Latitude,
			Longitude: m.Location.Longitude,
		}
	}
	return r
}
