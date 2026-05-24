package dto

import (
	"strconv"
	"time"

	"carsharing/api-gateway/internal/model"
	"github.com/gin-gonic/gin"
)

type CarInsuranceResponse struct {
	Insurance CarInsurance `json:"insurance"`
}

type CarInsurancesResponse struct {
	Insurances []CarInsurance `json:"insurances"`
}

type CarInsurance struct {
	ID        string    `json:"id"`
	CarID     string    `json:"carID"`
	Type      string    `json:"type" validate:"oneof=osago kasko"`
	Provider  string    `json:"provider"`
	PolicyNum string    `json:"policyNum"`
	StartsAt  time.Time `json:"startsAt"`
	ExpiresAt time.Time `json:"expiresAt"`
	CostTenge int32     `json:"costTenge" validate:"min=0"`
	Status    string    `json:"status" validate:"oneof=active expired cancelled"`
	ImageURLs []string  `json:"imageURLs,omitempty"`
	Notes     *string   `json:"notes,omitempty"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type CarInsuranceCreateRequest struct {
	CarID     string    `json:"carID" binding:"required"`
	Type      string    `json:"type" binding:"required,oneof=osago kasko"`
	Provider  string    `json:"provider" binding:"required"`
	PolicyNum string    `json:"policyNum" binding:"required"`
	StartsAt  time.Time `json:"startsAt" binding:"required"`
	ExpiresAt time.Time `json:"expiresAt" binding:"required"`
	CostTenge int32     `json:"costTenge" validate:"min=0"`
	Notes     *string   `json:"notes"`
}

type CarInsuranceUpdateRequest struct {
	Provider  *string    `json:"provider"`
	PolicyNum *string    `json:"policyNum"`
	StartsAt  *time.Time `json:"startsAt"`
	ExpiresAt *time.Time `json:"expiresAt"`
	CostTenge *int32     `json:"costTenge" validate:"omitempty,min=0"`
	Status    *string    `json:"status" validate:"omitempty,oneof=active expired cancelled"`
	ImageKeys []string   `json:"imageKeys"`
	Notes     *string    `json:"notes"`
}

func FromCarInsuranceCreateRequest(ctx *gin.Context) (model.CarInsuranceCreate, error) {
	var req CarInsuranceCreateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		return model.CarInsuranceCreate{}, err
	}
	return model.CarInsuranceCreate{
		CarID:     req.CarID,
		Type:      req.Type,
		Provider:  req.Provider,
		PolicyNum: req.PolicyNum,
		StartsAt:  req.StartsAt,
		ExpiresAt: req.ExpiresAt,
		CostTenge: req.CostTenge,
		Notes:     req.Notes,
	}, nil
}

func FromCarInsuranceUpdateRequest(ctx *gin.Context) (model.CarInsuranceUpdate, error) {
	var req CarInsuranceUpdateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		return model.CarInsuranceUpdate{}, err
	}
	return model.CarInsuranceUpdate{
		Provider:  req.Provider,
		PolicyNum: req.PolicyNum,
		StartsAt:  req.StartsAt,
		ExpiresAt: req.ExpiresAt,
		CostTenge: req.CostTenge,
		Status:    req.Status,
		ImageKeys: req.ImageKeys,
		Notes:     req.Notes,
	}, nil
}

func CarInsuranceFilterFromCtx(ctx *gin.Context) (model.CarInsuranceFilter, error) {
	f := model.CarInsuranceFilter{}

	if v := ctx.Query("carID"); v != "" {
		f.CarID = &v
	}
	if v := ctx.Query("type"); v != "" {
		f.Type = &v
	}
	if v := ctx.Query("status"); v != "" {
		f.Status = &v
	}
	if v := ctx.Query("expiringWithinDays"); v != "" {
		vInt, err := strconv.ParseInt(v, 10, 32)
		if err != nil {
			return model.CarInsuranceFilter{}, model.ErrInvalidQueryParam
		}

		days := int32(vInt)
		f.ExpiringWithinDays = &days
	}

	p, err := pagination(ctx)
	if err != nil {
		return model.CarInsuranceFilter{}, model.ErrInvalidQueryParam
	}

	f.Pagination = p

	return f, nil
}

func ToCarInsuranceResponse(m model.CarInsurance) CarInsurance {
	return CarInsurance{
		ID:        m.ID,
		CarID:     m.CarID,
		Type:      m.Type,
		Provider:  m.Provider,
		PolicyNum: m.PolicyNum,
		StartsAt:  m.StartsAt,
		ExpiresAt: m.ExpiresAt,
		CostTenge: m.CostTenge,
		Status:    m.Status,
		ImageURLs: m.ImageURLs,
		Notes:     m.Notes,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}
