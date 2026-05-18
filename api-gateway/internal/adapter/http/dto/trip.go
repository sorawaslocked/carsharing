package dto

import (
	"time"

	"carsharing/api-gateway/internal/model"
	"github.com/gin-gonic/gin"
)

type Trip struct {
	ID        string `json:"id"`
	BookingID string `json:"bookingID"`
	UserID    string `json:"userID"`
	CarID     string `json:"carID"`
	Status    string `json:"status"`

	StartedAt      time.Time `json:"startedAt"`
	StartLocation  location  `json:"startLocation"`
	StartMileageKM int64     `json:"startMileageKM"`
	StartFuelLevel *float32  `json:"startFuelLevel,omitempty"`

	EndedAt      *time.Time `json:"endedAt,omitempty"`
	EndLocation  *location  `json:"endLocation,omitempty"`
	EndMileageKM *int64     `json:"endMileageKM,omitempty"`
	EndFuelLevel *float32   `json:"endFuelLevel,omitempty"`

	DistanceTraveledKM *float64 `json:"distanceTraveledKM,omitempty"`
	DurationSeconds    *int64   `json:"durationSeconds,omitempty"`
	FinalCostTenge     *int32   `json:"finalCostTenge,omitempty"`

	CancelReason *string `json:"cancelReason,omitempty"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type TripSummary struct {
	TripID    string    `json:"tripID"`
	BookingID string    `json:"bookingID"`
	StartedAt time.Time `json:"startedAt"`
	EndedAt   time.Time `json:"endedAt"`

	DurationSeconds    int64   `json:"durationSeconds"`
	DistanceTraveledKM float64 `json:"distanceTraveledKM"`

	PricingSnapshot   PricingSnapshot `json:"pricingSnapshot"`
	BaseCostTenge     int32           `json:"baseCostTenge"`
	DistanceCostTenge int32           `json:"distanceCostTenge"`
	OvertimeCostTenge int32           `json:"overtimeCostTenge"`
	TotalCostTenge    int32           `json:"totalCostTenge"`
}

type TripStatusReading struct {
	ID         string    `json:"id"`
	TripID     string    `json:"tripID"`
	FromStatus string    `json:"fromStatus"`
	ToStatus   string    `json:"toStatus"`
	ActorType  string    `json:"actorType"`
	ActorID    *string   `json:"actorID,omitempty"`
	Reason     *string   `json:"reason,omitempty"`
	ChangedAt  time.Time `json:"changedAt"`
}

type TripStartRequest struct {
	BookingID string `json:"bookingID"`
}

type TripCancelRequest struct {
	Reason *string `json:"reason"`
}

func FromTripStartRequest(ctx *gin.Context) (string, error) {
	var req TripStartRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		return "", err
	}

	return req.BookingID, nil
}

func FromTripCancelRequest(ctx *gin.Context) (*string, error) {
	var req TripCancelRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		return nil, err
	}

	return req.Reason, nil
}

func TripFilterFromCtx(ctx *gin.Context) (model.TripFilter, error) {
	f := model.TripFilter{}

	if v := ctx.Query("userID"); v != "" {
		f.UserID = &v
	}
	if v := ctx.Query("carID"); v != "" {
		f.CarID = &v
	}
	if v := ctx.Query("status"); v != "" {
		f.Status = &v
	}
	if v := ctx.Query("startedAfter"); v != "" {
		t, err := time.Parse(time.RFC3339, v)
		if err != nil {
			return model.TripFilter{}, model.ErrInvalidQueryParam
		}
		f.StartedAfter = &t
	}
	if v := ctx.Query("startedBefore"); v != "" {
		t, err := time.Parse(time.RFC3339, v)
		if err != nil {
			return model.TripFilter{}, model.ErrInvalidQueryParam
		}
		f.StartedBefore = &t
	}

	p, err := pagination(ctx)
	if err != nil {
		return model.TripFilter{}, model.ErrInvalidQueryParam
	}

	f.Pagination = p

	return f, nil
}

func TripStatusReadingFilterFromCtx(ctx *gin.Context) (model.TripStatusReadingFilter, error) {
	f := model.TripStatusReadingFilter{}

	if v := ctx.Query("from"); v != "" {
		t, err := time.Parse(time.RFC3339, v)
		if err != nil {
			return model.TripStatusReadingFilter{}, model.ErrInvalidQueryParam
		}
		f.From = &t
	}
	if v := ctx.Query("to"); v != "" {
		t, err := time.Parse(time.RFC3339, v)
		if err != nil {
			return model.TripStatusReadingFilter{}, model.ErrInvalidQueryParam
		}
		f.To = &t
	}

	p, err := pagination(ctx)
	if err != nil {
		return model.TripStatusReadingFilter{}, model.ErrInvalidQueryParam
	}

	f.Pagination = p

	return f, nil
}

func ToTripResponse(m model.Trip) Trip {
	t := Trip{
		ID:                 m.ID,
		BookingID:          m.BookingID,
		UserID:             m.UserID,
		CarID:              m.CarID,
		Status:             m.Status,
		StartedAt:          m.StartedAt,
		StartLocation:      location{Latitude: m.StartLocation.Latitude, Longitude: m.StartLocation.Longitude},
		StartMileageKM:     m.StartMileageKM,
		StartFuelLevel:     m.StartFuelLevel,
		EndedAt:            m.EndedAt,
		EndMileageKM:       m.EndMileageKM,
		EndFuelLevel:       m.EndFuelLevel,
		DistanceTraveledKM: m.DistanceTraveledKM,
		DurationSeconds:    m.DurationSeconds,
		FinalCostTenge:     m.FinalCostTenge,
		CancelReason:       m.CancelReason,
		CreatedAt:          m.CreatedAt,
		UpdatedAt:          m.UpdatedAt,
	}

	if m.EndLocation != nil {
		l := location{Latitude: m.EndLocation.Latitude, Longitude: m.EndLocation.Longitude}
		t.EndLocation = &l
	}

	return t
}

func ToTripSummaryResponse(m model.TripSummary) TripSummary {
	return TripSummary{
		TripID:             m.TripID,
		BookingID:          m.BookingID,
		StartedAt:          m.StartedAt,
		EndedAt:            m.EndedAt,
		DurationSeconds:    m.DurationSeconds,
		DistanceTraveledKM: m.DistanceTraveledKM,
		PricingSnapshot: PricingSnapshot{
			RateTenge:         m.PricingSnapshot.RateTenge,
			RatePerKMTenge:    m.PricingSnapshot.RatePerKMTenge,
			FreeMinutes:       m.PricingSnapshot.FreeMinutes,
			MinChargeTenge:    m.PricingSnapshot.MinChargeTenge,
			OvertimePolicy:    m.PricingSnapshot.OvertimePolicy,
			OvertimeRateTenge: m.PricingSnapshot.OvertimeRateTenge,
		},
		BaseCostTenge:     m.BaseCostTenge,
		DistanceCostTenge: m.DistanceCostTenge,
		OvertimeCostTenge: m.OvertimeCostTenge,
		TotalCostTenge:    m.TotalCostTenge,
	}
}

func ToTripStatusReadingResponse(m model.TripStatusReading) TripStatusReading {
	return TripStatusReading{
		ID:         m.ID,
		TripID:     m.TripID,
		FromStatus: m.FromStatus,
		ToStatus:   m.ToStatus,
		ActorType:  m.ActorType,
		ActorID:    m.ActorID,
		Reason:     m.Reason,
		ChangedAt:  m.ChangedAt,
	}
}
