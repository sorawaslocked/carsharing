package dto

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sorawaslocked/car-rental-api-gateway/internal/model"
)

type Car struct {
	ID               string    `json:"id"`
	ModelID          string    `json:"modelId"`
	VIN              string    `json:"vin"`
	LicensePlate     string    `json:"licensePlate"`
	Color            string    `json:"color"`
	YearManufactured int16     `json:"yearManufactured"`
	MileageKM        int64     `json:"mileageKm"`
	FuelLevel        *float32  `json:"fuelLevel,omitempty"`
	BatteryLevel     *float32  `json:"batteryLevel,omitempty"`
	Location         location  `json:"location"`
	TelematicsID     string    `json:"telematicsId"`
	FuelStatus       string    `json:"fuelStatus"`
	ZoneID           *string   `json:"zoneId,omitempty"`
	Status           string    `json:"status"`
	Notes            *string   `json:"notes,omitempty"`
	ImageStorageUrls []string  `json:"imageStorageUrls,omitempty"`
	LastSeenAt       time.Time `json:"lastSeenAt"`
	CreatedAt        time.Time `json:"createdAt"`
	UpdatedAt        time.Time `json:"updatedAt"`
}

type CarCreateRequest struct {
	ModelID          string   `json:"modelID"`
	VIN              string   `json:"vin"`
	LicensePlate     string   `json:"licensePlate"`
	Color            string   `json:"color"`
	YearManufactured int16    `json:"yearManufactured"`
	MileageKM        int64    `json:"mileageKm"`
	FuelLevel        *float32 `json:"fuelLevel"`
	BatteryLevel     *float32 `json:"batteryLevel"`
	Latitude         float64  `json:"latitude"`
	Longitude        float64  `json:"longitude"`
	TelematicsID     string   `json:"telematicsID"`
	Notes            *string  `json:"notes"`
	ImageStorageKeys []string `json:"imageStorageKeys"`
}

type CarUpdateRequest struct {
	ModelID          *string  `json:"modelId"`
	LicensePlate     *string  `json:"licensePlate"`
	Color            *string  `json:"color"`
	MileageKM        *int64   `json:"mileageKm"`
	FuelLevel        *float32 `json:"fuelLevel"`
	BatteryLevel     *float32 `json:"batteryLevel"`
	Latitude         *float64 `json:"latitude"`
	Longitude        *float64 `json:"longitude"`
	TelematicsID     *string  `json:"telematicsId"`
	ZoneID           *string  `json:"zoneId"`
	Status           *string  `json:"status"`
	Notes            *string  `json:"notes"`
	ImageStorageKeys []string `json:"imageStorageKeys"`
}

type CarStatusLogEntry struct {
	ID         string         `json:"id"`
	CarID      string         `json:"carID"`
	FromStatus string         `json:"fromStatus"`
	ToStatus   string         `json:"toStatus"`
	ActorType  string         `json:"actorType"`
	ActorID    *string        `json:"actorId,omitempty"`
	Reason     *string        `json:"reason,omitempty"`
	Metadata   map[string]any `json:"metadata,omitempty"`
	ChangedAt  time.Time      `json:"changedAt"`
}

type CarFuelReading struct {
	CarID      string    `json:"carID"`
	FuelPct    int       `json:"fuelPct"`
	RawPct     int       `json:"rawPct"`
	RecordedAt time.Time `json:"recordedAt"`
}

func FromCarCreateRequest(ctx *gin.Context) (model.CarCreate, error) {
	var req CarCreateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		return model.CarCreate{}, err
	}

	return model.CarCreate{
		ModelID:          req.ModelID,
		VIN:              req.VIN,
		LicensePlate:     req.LicensePlate,
		Color:            req.Color,
		YearManufactured: req.YearManufactured,
		MileageKM:        req.MileageKM,
		FuelLevel:        req.FuelLevel,
		BatteryLevel:     req.BatteryLevel,
		Latitude:         req.Latitude,
		Longitude:        req.Longitude,
		TelematicsID:     req.TelematicsID,
		Notes:            req.Notes,
		ImageStorageKeys: req.ImageStorageKeys,
	}, nil
}

func FromCarUpdateRequest(ctx *gin.Context) (model.CarUpdate, error) {
	var req CarUpdateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		return model.CarUpdate{}, err
	}

	return model.CarUpdate{
		ModelID:          req.ModelID,
		LicensePlate:     req.LicensePlate,
		Color:            req.Color,
		MileageKM:        req.MileageKM,
		FuelLevel:        req.FuelLevel,
		BatteryLevel:     req.BatteryLevel,
		Latitude:         req.Latitude,
		Longitude:        req.Longitude,
		TelematicsID:     req.TelematicsID,
		ZoneID:           req.ZoneID,
		Status:           req.Status,
		Notes:            req.Notes,
		ImageStorageKeys: req.ImageStorageKeys,
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
	if v := ctx.Query("latitude"); v != "" {
		vFloat, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return model.CarFilter{}, model.ErrInvalidQueryParam
		}

		f.Latitude = &vFloat
	}
	if v := ctx.Query("longitude"); v != "" {
		vFloat, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return model.CarFilter{}, model.ErrInvalidQueryParam
		}

		f.Longitude = &vFloat
	}
	if v := ctx.Query("radiusM"); v != "" {
		vInt, err := strconv.ParseInt(v, 10, 32)
		if err != nil {
			return model.CarFilter{}, model.ErrInvalidQueryParam
		}

		radiusM := int32(vInt)
		f.RadiusM = &radiusM
	}
	if v := ctx.Query("zoneId"); v != "" {
		f.ZoneID = &v
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

func CarStatusLogFilterFromCtx(ctx *gin.Context) (model.CarStatusLogFilter, error) {
	f := model.CarStatusLogFilter{}

	if v := ctx.Query("carID"); v != "" {
		f.CarID = &v
	}

	p, err := pagination(ctx)
	if err != nil {
		return model.CarStatusLogFilter{}, model.ErrInvalidQueryParam
	}

	f.Pagination = p

	return f, nil
}

func CarFuelReadingFilterFromCtx(ctx *gin.Context) (model.CarFuelReadingFilter, error) {
	f := model.CarFuelReadingFilter{}

	if v := ctx.Query("carID"); v != "" {
		f.CarID = &v
	}
	if v := ctx.Query("from"); v != "" {
		vTime, err := time.Parse("2006-01-02", v)
		if err != nil {
			return model.CarFuelReadingFilter{}, model.ErrInvalidQueryParam
		}

		f.From = &vTime
	}
	if v := ctx.Query("to"); v != "" {
		vTime, err := time.Parse("2006-01-02", v)
		if err != nil {
			return model.CarFuelReadingFilter{}, model.ErrInvalidQueryParam
		}

		f.To = &vTime
	}

	p, err := pagination(ctx)
	if err != nil {
		return model.CarFuelReadingFilter{}, model.ErrInvalidQueryParam
	}

	f.Pagination = p

	return f, nil
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
		TelematicsID:     m.TelematicsID,
		FuelStatus:       m.FuelStatus,
		ZoneID:           m.ZoneID,
		Status:           m.Status,
		Notes:            m.Notes,
		ImageStorageUrls: m.ImageStorageUrls,
		LastSeenAt:       m.LastSeenAt,
		CreatedAt:        m.CreatedAt,
		UpdatedAt:        m.UpdatedAt,
	}
}

func ToCarStatusLogEntryResponse(m model.CarStatusLogEntry) CarStatusLogEntry {
	return CarStatusLogEntry{
		ID:         m.ID,
		CarID:      m.CarID,
		FromStatus: m.FromStatus,
		ToStatus:   m.ToStatus,
		ActorType:  m.ActorType,
		ActorID:    m.ActorID,
		Reason:     m.Reason,
		Metadata:   m.Metadata,
		ChangedAt:  m.ChangedAt,
	}
}

func ToCarFuelReadingResponse(m model.CarFuelReading) CarFuelReading {
	return CarFuelReading{
		CarID:      m.CarID,
		FuelPct:    m.FuelPct,
		RawPct:     m.RawPct,
		RecordedAt: m.RecordedAt,
	}
}
