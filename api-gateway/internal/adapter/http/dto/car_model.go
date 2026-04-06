package dto

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sorawaslocked/car-rental-api-gateway/internal/model"
)

type CarModel struct {
	ID               string    `json:"id"`
	Brand            string    `json:"brand"`
	Model            string    `json:"model"`
	Year             int16     `json:"year"`
	FuelType         string    `json:"fuelType"`
	Transmission     string    `json:"transmission"`
	BodyType         string    `json:"bodyType"`
	Class            string    `json:"class"`
	Seats            int8      `json:"seats"`
	EngineVolume     *float32  `json:"engineVolume,omitempty"`
	RangeKM          int32     `json:"rangeKm"`
	Features         []string  `json:"features,omitempty"`
	ImageStorageUrls []string  `json:"imageStorageUrls,omitempty"`
	CreatedAt        time.Time `json:"createdAt"`
	UpdatedAt        time.Time `json:"updatedAt"`
}

type CarModelCreateRequest struct {
	Brand            string   `json:"brand"`
	Model            string   `json:"model"`
	Year             int16    `json:"year"`
	FuelType         string   `json:"fuelType"`
	Transmission     string   `json:"transmission"`
	BodyType         string   `json:"bodyType"`
	Class            string   `json:"class"`
	Seats            int8     `json:"seats"`
	EngineVolume     *float32 `json:"engineVolume"`
	RangeKM          int32    `json:"rangeKm"`
	Features         []string `json:"features"`
	ImageStorageKeys []string `json:"imageStorageKeys"`
}

type CarModelUpdateRequest struct {
	Brand            *string  `json:"brand"`
	Model            *string  `json:"model"`
	Year             *int16   `json:"year"`
	FuelType         *string  `json:"fuelType"`
	Transmission     *string  `json:"transmission"`
	BodyType         *string  `json:"bodyType"`
	Class            *string  `json:"class"`
	Seats            *int8    `json:"seats"`
	EngineVolume     *float32 `json:"engineVolume"`
	RangeKM          *int32   `json:"rangeKm"`
	Features         []string `json:"features"`
	ImageStorageKeys []string `json:"imageStorageKeys"`
}

func FromCarModelCreateRequest(ctx *gin.Context) (model.CarModelCreate, error) {
	var req CarModelCreateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		return model.CarModelCreate{}, err
	}

	return model.CarModelCreate{
		Brand:            req.Brand,
		Model:            req.Model,
		Year:             req.Year,
		FuelType:         req.FuelType,
		Transmission:     req.Transmission,
		BodyType:         req.BodyType,
		Class:            req.Class,
		Seats:            req.Seats,
		EngineVolume:     req.EngineVolume,
		RangeKM:          req.RangeKM,
		Features:         req.Features,
		ImageStorageKeys: req.ImageStorageKeys,
	}, nil
}

func FromCarModelUpdateRequest(ctx *gin.Context) (model.CarModelUpdate, error) {
	var req CarModelUpdateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		return model.CarModelUpdate{}, err
	}

	return model.CarModelUpdate{
		Brand:            req.Brand,
		Model:            req.Model,
		Year:             req.Year,
		FuelType:         req.FuelType,
		Transmission:     req.Transmission,
		BodyType:         req.BodyType,
		Class:            req.Class,
		Seats:            req.Seats,
		EngineVolume:     req.EngineVolume,
		RangeKM:          req.RangeKM,
		Features:         req.Features,
		ImageStorageKeys: req.ImageStorageKeys,
	}, nil
}

func CarModelFilterFromCtx(ctx *gin.Context) (model.CarModelFilter, error) {
	f := model.CarModelFilter{}

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
		vInt, err := strconv.Atoi(v)
		if err != nil {
			return model.CarModelFilter{}, model.ErrInvalidQueryParam
		}

		minSeats := int8(vInt)
		f.MinSeats = &minSeats
	}

	p, err := pagination(ctx)
	if err != nil {
		return model.CarModelFilter{}, model.ErrInvalidQueryParam
	}

	f.Pagination = p

	return f, nil
}

func ToCarModelResponse(m model.CarModel) CarModel {
	return CarModel{
		ID:               m.ID,
		Brand:            m.Brand,
		Model:            m.Model,
		Year:             m.Year,
		FuelType:         m.FuelType,
		Transmission:     m.Transmission,
		BodyType:         m.BodyType,
		Class:            m.Class,
		Seats:            m.Seats,
		EngineVolume:     m.EngineVolume,
		RangeKM:          m.RangeKM,
		Features:         m.Features,
		ImageStorageUrls: m.ImageStorageUrls,
		CreatedAt:        m.CreatedAt,
		UpdatedAt:        m.UpdatedAt,
	}
}
