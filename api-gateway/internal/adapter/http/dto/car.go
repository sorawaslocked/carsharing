package dto

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sorawaslocked/car-rental-api-gateway/internal/model"
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

type CarFuelHistoryResponse struct {
	FuelHistory []CarFuelReading `json:"fuelHistory"`
}

type CarLocationHistoryResponse struct {
	LocationHistory []CarLocationReading `json:"locationHistory"`
}

type CarBatteryHistoryResponse struct {
	BatteryHistory []CarBatteryReading `json:"batteryHistory"`
}

type CarMileageHistoryResponse struct {
	MileageHistory []CarMileageReading `json:"mileageHistory"`
}

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
	IsRetired        bool      `json:"isRetired"`
	Notes            *string   `json:"notes,omitempty"`
	ImageStorageUrls []string  `json:"imageStorageUrls,omitempty"`
	LastSeenAt       time.Time `json:"lastSeenAt"`
	CreatedAt        time.Time `json:"createdAt"`
	UpdatedAt        time.Time `json:"updatedAt"`
}

type CarStatusReading struct {
	ID         string         `json:"id"`
	CarID      string         `json:"carId"`
	FromStatus string         `json:"fromStatus"`
	ToStatus   string         `json:"toStatus"`
	ActorType  string         `json:"actorType"`
	ActorID    *string        `json:"actorId,omitempty"`
	Reason     *string        `json:"reason,omitempty"`
	Metadata   map[string]any `json:"metadata,omitempty"`
	ChangedAt  time.Time      `json:"changedAt"`
}

type CarFuelReading struct {
	ID         string         `json:"id"`
	CarID      string         `json:"carId"`
	FuelPct    float32        `json:"fuelPct"`
	RawPct     float32        `json:"rawPct"`
	ActorType  string         `json:"actorType"`
	ActorID    *string        `json:"actorId,omitempty"`
	Reason     *string        `json:"reason,omitempty"`
	Metadata   map[string]any `json:"metadata,omitempty"`
	RecordedAt time.Time      `json:"recordedAt"`
}

type CarLocationReading struct {
	ID         string         `json:"id"`
	CarID      string         `json:"carId"`
	Location   location       `json:"location"`
	ActorType  string         `json:"actorType"`
	ActorID    *string        `json:"actorId,omitempty"`
	Reason     *string        `json:"reason,omitempty"`
	Metadata   map[string]any `json:"metadata,omitempty"`
	RecordedAt time.Time      `json:"recordedAt"`
}

type CarBatteryReading struct {
	ID           string         `json:"id"`
	CarID        string         `json:"carId"`
	BatteryLevel float32        `json:"batteryLevel"`
	ActorType    string         `json:"actorType"`
	ActorID      *string        `json:"actorId,omitempty"`
	Reason       *string        `json:"reason,omitempty"`
	Metadata     map[string]any `json:"metadata,omitempty"`
	RecordedAt   time.Time      `json:"recordedAt"`
}

type CarMileageReading struct {
	ID         string         `json:"id"`
	CarID      string         `json:"carId"`
	MileageKM  int64          `json:"mileageKm"`
	ActorType  string         `json:"actorType"`
	ActorID    *string        `json:"actorId,omitempty"`
	Reason     *string        `json:"reason,omitempty"`
	Metadata   map[string]any `json:"metadata,omitempty"`
	RecordedAt time.Time      `json:"recordedAt"`
}

type CarCreateRequest struct {
	ModelID          string  `json:"modelId"`
	VIN              string  `json:"vin"`
	LicensePlate     string  `json:"licensePlate"`
	Color            string  `json:"color"`
	YearManufactured int16   `json:"yearManufactured"`
	TelematicsID     string  `json:"telematicsId"`
	Notes            *string `json:"notes"`
}

type CarUpdateRequest struct {
	ModelID      *string  `json:"modelId"`
	LicensePlate *string  `json:"licensePlate"`
	Color        *string  `json:"color"`
	MileageKM    *int64   `json:"mileageKm"`
	FuelLevel    *float32 `json:"fuelLevel"`
	BatteryLevel *float32 `json:"batteryLevel"`
	Latitude     *float64 `json:"latitude"`
	Longitude    *float64 `json:"longitude"`
	TelematicsID *string  `json:"telematicsId"`
	ZoneID       *string  `json:"zoneId"`
	Status       *string  `json:"status"`
	IsRetired    *bool    `json:"isRetired"`
	Notes        *string  `json:"notes"`
	ImageKeys    []string `json:"imageKeys"`
}

type CarElevatedUpdateRequest struct {
	Status       *string        `json:"status"`
	MileageKM    *int64         `json:"mileageKm"`
	FuelLevel    *float32       `json:"fuelLevel"`
	BatteryLevel *float32       `json:"batteryLevel"`
	Latitude     *float64       `json:"latitude"`
	Longitude    *float64       `json:"longitude"`
	Reason       string         `json:"reason"`
	Metadata     map[string]any `json:"metadata"`
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
		TelematicsID:     req.TelematicsID,
		Notes:            req.Notes,
	}, nil
}

func FromCarUpdateRequest(ctx *gin.Context) (model.CarUpdate, error) {
	var req CarUpdateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		return model.CarUpdate{}, err
	}

	update := model.CarUpdate{
		ModelID:      req.ModelID,
		LicensePlate: req.LicensePlate,
		Color:        req.Color,
		MileageKM:    req.MileageKM,
		FuelLevel:    req.FuelLevel,
		BatteryLevel: req.BatteryLevel,
		TelematicsID: req.TelematicsID,
		ZoneID:       req.ZoneID,
		Status:       req.Status,
		IsRetired:    req.IsRetired,
		Notes:        req.Notes,
		ImageKeys:    req.ImageKeys,
	}
	if req.Latitude != nil && req.Longitude != nil {
		update.Location = &model.Location{
			Latitude:  *req.Latitude,
			Longitude: *req.Longitude,
		}
	}

	return update, nil
}

func FromCarElevatedUpdateRequest(ctx *gin.Context) (model.CarElevatedUpdate, error) {
	var req CarElevatedUpdateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		return model.CarElevatedUpdate{}, err
	}

	update := model.CarElevatedUpdate{
		Status:       req.Status,
		MileageKM:    req.MileageKM,
		FuelLevel:    req.FuelLevel,
		BatteryLevel: req.BatteryLevel,
		Reason:       req.Reason,
		Metadata:     req.Metadata,
	}
	if req.Latitude != nil && req.Longitude != nil {
		update.Location = &model.Location{
			Latitude:  *req.Latitude,
			Longitude: *req.Longitude,
		}
	}

	return update, nil
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
			f.Location = &model.Location{Latitude: latF, Longitude: lngF}
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

func carTimeRangeFilter(ctx *gin.Context) (from, to *time.Time, p *model.Pagination, err error) {
	if v := ctx.Query("from"); v != "" {
		t, parseErr := time.Parse("2006-01-02", v)
		if parseErr != nil {
			return nil, nil, nil, model.ErrInvalidQueryParam
		}
		from = &t
	}
	if v := ctx.Query("to"); v != "" {
		t, parseErr := time.Parse("2006-01-02", v)
		if parseErr != nil {
			return nil, nil, nil, model.ErrInvalidQueryParam
		}
		to = &t
	}

	p, err = pagination(ctx)
	if err != nil {
		return nil, nil, nil, model.ErrInvalidQueryParam
	}

	return from, to, p, nil
}

func CarStatusReadingFilterFromCtx(ctx *gin.Context) (model.CarStatusReadingFilter, error) {
	from, to, p, err := carTimeRangeFilter(ctx)
	if err != nil {
		return model.CarStatusReadingFilter{}, err
	}

	return model.CarStatusReadingFilter{From: from, To: to, Pagination: p}, nil
}

func CarFuelReadingFilterFromCtx(ctx *gin.Context) (model.CarFuelReadingFilter, error) {
	from, to, p, err := carTimeRangeFilter(ctx)
	if err != nil {
		return model.CarFuelReadingFilter{}, err
	}

	return model.CarFuelReadingFilter{From: from, To: to, Pagination: p}, nil
}

func CarLocationReadingFilterFromCtx(ctx *gin.Context) (model.CarLocationReadingFilter, error) {
	from, to, p, err := carTimeRangeFilter(ctx)
	if err != nil {
		return model.CarLocationReadingFilter{}, err
	}

	return model.CarLocationReadingFilter{From: from, To: to, Pagination: p}, nil
}

func CarBatteryReadingFilterFromCtx(ctx *gin.Context) (model.CarBatteryReadingFilter, error) {
	from, to, p, err := carTimeRangeFilter(ctx)
	if err != nil {
		return model.CarBatteryReadingFilter{}, err
	}

	return model.CarBatteryReadingFilter{From: from, To: to, Pagination: p}, nil
}

func CarMileageReadingFilterFromCtx(ctx *gin.Context) (model.CarMileageReadingFilter, error) {
	from, to, p, err := carTimeRangeFilter(ctx)
	if err != nil {
		return model.CarMileageReadingFilter{}, err
	}

	return model.CarMileageReadingFilter{From: from, To: to, Pagination: p}, nil
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
		ChangedAt:  m.ChangedAt,
	}
}

func ToCarFuelReadingResponse(m model.CarFuelReading) CarFuelReading {
	return CarFuelReading{
		ID:         m.ID,
		CarID:      m.CarID,
		FuelPct:    m.FuelPct,
		RawPct:     m.RawPct,
		ActorType:  m.ActorType,
		ActorID:    m.ActorID,
		Reason:     m.Reason,
		Metadata:   m.Metadata,
		RecordedAt: m.RecordedAt,
	}
}

func ToCarLocationReadingResponse(m model.CarLocationReading) CarLocationReading {
	return CarLocationReading{
		ID:    m.ID,
		CarID: m.CarID,
		Location: location{
			Latitude:  m.Location.Latitude,
			Longitude: m.Location.Longitude,
		},
		ActorType:  m.ActorType,
		ActorID:    m.ActorID,
		Reason:     m.Reason,
		Metadata:   m.Metadata,
		RecordedAt: m.RecordedAt,
	}
}

func ToCarBatteryReadingResponse(m model.CarBatteryReading) CarBatteryReading {
	return CarBatteryReading{
		ID:           m.ID,
		CarID:        m.CarID,
		BatteryLevel: m.BatteryLevel,
		ActorType:    m.ActorType,
		ActorID:      m.ActorID,
		Reason:       m.Reason,
		Metadata:     m.Metadata,
		RecordedAt:   m.RecordedAt,
	}
}

func ToCarMileageReadingResponse(m model.CarMileageReading) CarMileageReading {
	return CarMileageReading{
		ID:         m.ID,
		CarID:      m.CarID,
		MileageKM:  m.MileageKM,
		ActorType:  m.ActorType,
		ActorID:    m.ActorID,
		Reason:     m.Reason,
		Metadata:   m.Metadata,
		RecordedAt: m.RecordedAt,
	}
}
