package dto

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sorawaslocked/car-rental-api-gateway/internal/model"
)

type BookingPricingSnapshot struct {
	RateTenge         int32   `json:"rateTenge"`
	RatePerKMTenge    *int32  `json:"ratePerKMTenge,omitempty"`
	FreeMinutes       *int32  `json:"freeMinutes,omitempty"`
	MinChargeTenge    *int32  `json:"minChargeTenge,omitempty"`
	OvertimePolicy    *string `json:"overtimePolicy,omitempty"`
	OvertimeRateTenge *int32  `json:"overtimeRateTenge,omitempty"`
}

type Booking struct {
	ID               string                 `json:"id"`
	UserID           string                 `json:"userID"`
	CarID            string                 `json:"carID"`
	CommittedPeriods *int32                 `json:"committedPeriods,omitempty"`
	Status           string                 `json:"status"`
	PricingRuleID    string                 `json:"pricingRuleID"`
	PricingSnapshot  BookingPricingSnapshot `json:"pricingSnapshot"`
	CreatedAt        time.Time              `json:"createdAt"`
	UpdatedAt        time.Time              `json:"updatedAt"`
}

type BookingStatusReading struct {
	ID         string    `json:"id"`
	BookingID  string    `json:"bookingID"`
	FromStatus string    `json:"fromStatus"`
	ToStatus   string    `json:"toStatus"`
	ActorType  string    `json:"actorType"`
	ActorID    *string   `json:"actorID,omitempty"`
	Reason     *string   `json:"reason,omitempty"`
	ChangedAt  time.Time `json:"changedAt"`
}

type BookingCreateRequest struct {
	CarID            string `json:"carID"`
	PricingRuleID    string `json:"pricingRuleID"`
	CommittedPeriods *int32 `json:"committedPeriods"`
}

type BookingStatusUpdateRequest struct {
	Status string  `json:"status"`
	Reason *string `json:"reason"`
}

func FromBookingCreateRequest(ctx *gin.Context) (model.BookingCreate, error) {
	userID := ctx.GetString("x-user-id")

	var req BookingCreateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		return model.BookingCreate{}, err
	}

	return model.BookingCreate{
		UserID:           userID,
		CarID:            req.CarID,
		PricingRuleID:    req.PricingRuleID,
		CommittedPeriods: req.CommittedPeriods,
	}, nil
}

func FromBookingStatusUpdateRequest(ctx *gin.Context) (model.BookingStatusUpdate, error) {
	var req BookingStatusUpdateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		return model.BookingStatusUpdate{}, err
	}

	return model.BookingStatusUpdate{
		Status: req.Status,
		Reason: req.Reason,
	}, nil
}

func BookingFilterFromCtx(ctx *gin.Context) (model.BookingFilter, error) {
	f := model.BookingFilter{}

	if v := ctx.Query("userID"); v != "" {
		f.UserID = &v
	}
	if v := ctx.Query("carID"); v != "" {
		f.CarID = &v
	}
	if v := ctx.Query("status"); v != "" {
		f.Status = &v
	}
	if v := ctx.Query("pricingRuleID"); v != "" {
		f.PricingRuleID = &v
	}

	p, err := pagination(ctx)
	if err != nil {
		return model.BookingFilter{}, model.ErrInvalidQueryParam
	}

	f.Pagination = p

	return f, nil
}

func BookingStatusReadingFilterFromCtx(ctx *gin.Context) (model.BookingStatusReadingFilter, error) {
	f := model.BookingStatusReadingFilter{}

	if v := ctx.Query("from"); v != "" {
		t, err := time.Parse(time.RFC3339, v)
		if err != nil {
			return model.BookingStatusReadingFilter{}, model.ErrInvalidQueryParam
		}
		f.From = &t
	}
	if v := ctx.Query("to"); v != "" {
		t, err := time.Parse(time.RFC3339, v)
		if err != nil {
			return model.BookingStatusReadingFilter{}, model.ErrInvalidQueryParam
		}
		f.To = &t
	}

	p, err := pagination(ctx)
	if err != nil {
		return model.BookingStatusReadingFilter{}, model.ErrInvalidQueryParam
	}

	f.Pagination = p

	return f, nil
}

func ToBookingResponse(m model.Booking) Booking {
	return Booking{
		ID:               m.ID,
		UserID:           m.UserID,
		CarID:            m.CarID,
		CommittedPeriods: m.CommittedPeriods,
		Status:           m.Status,
		PricingRuleID:    m.PricingRuleID,
		PricingSnapshot: BookingPricingSnapshot{
			RateTenge:         m.PricingSnapshot.RateTenge,
			RatePerKMTenge:    m.PricingSnapshot.RatePerKMTenge,
			FreeMinutes:       m.PricingSnapshot.FreeMinutes,
			MinChargeTenge:    m.PricingSnapshot.MinChargeTenge,
			OvertimePolicy:    m.PricingSnapshot.OvertimePolicy,
			OvertimeRateTenge: m.PricingSnapshot.OvertimeRateTenge,
		},
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}

func ToBookingStatusReadingResponse(m model.BookingStatusReading) BookingStatusReading {
	return BookingStatusReading{
		ID:         m.ID,
		BookingID:  m.BookingID,
		FromStatus: m.FromStatus,
		ToStatus:   m.ToStatus,
		ActorType:  m.ActorType,
		ActorID:    m.ActorID,
		Reason:     m.Reason,
		ChangedAt:  m.ChangedAt,
	}
}
