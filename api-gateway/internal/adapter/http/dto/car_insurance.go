package dto

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sorawaslocked/car-rental-api-gateway/internal/model"
)

type CarInsurance struct {
	ID               string    `json:"id"`
	CarID            string    `json:"carId"`
	Type             string    `json:"type"`
	Provider         string    `json:"provider"`
	PolicyNum        string    `json:"policyNum"`
	StartsAt         time.Time `json:"startsAt"`
	ExpiresAt        time.Time `json:"expiresAt"`
	CostTenge        int32     `json:"costTenge"`
	Status           string    `json:"status"`
	ImageStorageUrls []string  `json:"imageStorageUrls,omitempty"`
	Notes            *string   `json:"notes,omitempty"`
	CreatedAt        time.Time `json:"createdAt"`
	UpdatedAt        time.Time `json:"updatedAt"`
}

type CarInsuranceCreateRequest struct {
	CarID            string    `json:"carId"`
	Type             string    `json:"type"`
	Provider         string    `json:"provider"`
	PolicyNum        string    `json:"policyNum"`
	StartsAt         time.Time `json:"startsAt"`
	ExpiresAt        time.Time `json:"expiresAt"`
	CostTenge        int32     `json:"costTenge"`
	ImageStorageKeys []string  `json:"imageStorageKeys"`
	Notes            *string   `json:"notes"`
}

type CarInsuranceUpdateRequest struct {
	Provider         *string    `json:"provider"`
	PolicyNum        *string    `json:"policyNum"`
	StartsAt         *time.Time `json:"startsAt"`
	ExpiresAt        *time.Time `json:"expiresAt"`
	CostTenge        *int32     `json:"costTenge"`
	Status           *string    `json:"status"`
	ImageStorageKeys []string   `json:"imageStorageKeys"`
	Notes            *string    `json:"notes"`
}

func FromCarInsuranceCreateRequest(ctx *gin.Context) (model.CarInsuranceCreate, error) {
	var req CarInsuranceCreateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		return model.CarInsuranceCreate{}, err
	}
	return model.CarInsuranceCreate{
		CarID:            req.CarID,
		Type:             req.Type,
		Provider:         req.Provider,
		PolicyNum:        req.PolicyNum,
		StartsAt:         req.StartsAt,
		ExpiresAt:        req.ExpiresAt,
		CostTenge:        req.CostTenge,
		ImageStorageKeys: req.ImageStorageKeys,
		Notes:            req.Notes,
	}, nil
}

func FromCarInsuranceUpdateRequest(ctx *gin.Context) (model.CarInsuranceUpdate, error) {
	var req CarInsuranceUpdateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		return model.CarInsuranceUpdate{}, err
	}
	return model.CarInsuranceUpdate{
		Provider:         req.Provider,
		PolicyNum:        req.PolicyNum,
		StartsAt:         req.StartsAt,
		ExpiresAt:        req.ExpiresAt,
		CostTenge:        req.CostTenge,
		Status:           req.Status,
		ImageStorageKeys: req.ImageStorageKeys,
		Notes:            req.Notes,
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
		ID:               m.ID,
		CarID:            m.CarID,
		Type:             m.Type,
		Provider:         m.Provider,
		PolicyNum:        m.PolicyNum,
		StartsAt:         m.StartsAt,
		ExpiresAt:        m.ExpiresAt,
		CostTenge:        m.CostTenge,
		Status:           m.Status,
		ImageStorageUrls: m.ImageStorageUrls,
		Notes:            m.Notes,
		CreatedAt:        m.CreatedAt,
		UpdatedAt:        m.UpdatedAt,
	}
}
