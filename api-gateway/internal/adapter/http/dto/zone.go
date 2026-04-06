package dto

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sorawaslocked/car-rental-api-gateway/internal/model"
)

type Zone struct {
	ID              string    `json:"id"`
	Name            string    `json:"name"`
	Type            string    `json:"type"`
	BoundaryGeoJSON string    `json:"boundary"`
	FeeAdjustment   int32     `json:"feeAdjustment"`
	IsActive        bool      `json:"isActive"`
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`
}

type ZoneCreateRequest struct {
	Name            string `json:"name"`
	Type            string `json:"type"`
	BoundaryGeoJSON string `json:"boundary"`
	FeeAdjustment   int32  `json:"feeAdjustment"`
}

type ZoneUpdateRequest struct {
	Name            *string `json:"name"`
	Type            *string `json:"type"`
	BoundaryGeoJSON *string `json:"boundary"`
	FeeAdjustment   *int32  `json:"feeAdjustment"`
	IsActive        *bool   `json:"isActive"`
}

func FromZoneCreateRequest(ctx *gin.Context) (model.ZoneCreate, error) {
	var req ZoneCreateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		return model.ZoneCreate{}, err
	}
	return model.ZoneCreate{
		Name:            req.Name,
		Type:            req.Type,
		BoundaryGeoJSON: req.BoundaryGeoJSON,
		FeeAdjustment:   req.FeeAdjustment,
	}, nil
}

func FromZoneUpdateRequest(ctx *gin.Context) (model.ZoneUpdate, error) {
	var req ZoneUpdateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		return model.ZoneUpdate{}, err
	}

	return model.ZoneUpdate{
		Name:            req.Name,
		Type:            req.Type,
		BoundaryGeoJSON: req.BoundaryGeoJSON,
		FeeAdjustment:   req.FeeAdjustment,
		IsActive:        req.IsActive,
	}, nil
}

func ZoneFilterFromCtx(ctx *gin.Context) (model.ZoneFilter, error) {
	f := model.ZoneFilter{}

	if v := ctx.Query("type"); v != "" {
		f.Type = &v
	}
	if v := ctx.Query("isActive"); v != "" {
		vBool, err := strconv.ParseBool(v)
		if err != nil {
			return model.ZoneFilter{}, model.ErrInvalidQueryParam
		}

		f.IsActive = &vBool
	}

	return f, nil
}

func ToZoneResponse(m model.Zone) Zone {
	return Zone{
		ID:              m.ID,
		Name:            m.Name,
		Type:            m.Type,
		BoundaryGeoJSON: m.BoundaryGeoJSON,
		FeeAdjustment:   m.FeeAdjustment,
		IsActive:        m.IsActive,
		CreatedAt:       m.CreatedAt,
		UpdatedAt:       m.UpdatedAt,
	}
}
