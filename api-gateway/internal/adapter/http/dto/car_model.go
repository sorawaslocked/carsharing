package dto

import (
	"strconv"
	"time"

	"carsharing/api-gateway/internal/model"
	"github.com/gin-gonic/gin"
)

type CarModelGetResponse struct {
	CarModel CarModel `json:"carModel"`
}

type CarModelsResponse struct {
	CarModels []CarModel `json:"carModels"`
}

type CarModel struct {
	ID               string    `json:"id"`
	Brand            string    `json:"brand"`
	Model            string    `json:"model"`
	Year             int16     `json:"year"`
	FuelType         string    `json:"fuelType" validate:"oneof=petrol diesel electric hybrid"`
	Transmission     string    `json:"transmission" validate:"oneof=manual auto"`
	BodyType         string    `json:"bodyType" validate:"oneof=sedan hatchback SUV crossover minivan coupe convertible pickup"`
	Class            string    `json:"class" validate:"oneof=economy compact comfort business luxury"`
	Seats            int8      `json:"seats"`
	EngineVolume     *float32  `json:"engineVolume,omitempty"`
	RangeKM          int32     `json:"rangeKm"`
	Features         []string  `json:"features,omitempty"`
	ImageStorageUrls []string  `json:"imageStorageUrls,omitempty"`
	CreatedAt        time.Time `json:"createdAt"`
	UpdatedAt        time.Time `json:"updatedAt"`
}

type CarModelCreateRequest struct {
	Brand        string   `json:"brand" binding:"required"`
	Model        string   `json:"model" binding:"required"`
	Year         int16    `json:"year" binding:"required"`
	FuelType     string   `json:"fuelType" binding:"required,oneof=petrol diesel electric hybrid"`
	Transmission string   `json:"transmission" binding:"required,oneof=manual auto"`
	BodyType     string   `json:"bodyType" binding:"required,oneof=sedan hatchback SUV crossover minivan coupe convertible pickup"`
	Class        string   `json:"class" binding:"required,oneof=economy compact comfort business luxury"`
	Seats        int8     `json:"seats" binding:"required,min=1,max=9"`
	EngineVolume *float32 `json:"engineVolume"`
	RangeKM      int32    `json:"rangeKm"`
	Features     []string `json:"features"`
}

type CarModelUpdateRequest struct {
	Brand            *string  `json:"brand"`
	Model            *string  `json:"model"`
	Year             *int16   `json:"year"`
	FuelType         *string  `json:"fuelType" validate:"omitempty,oneof=petrol diesel electric hybrid"`
	Transmission     *string  `json:"transmission" validate:"omitempty,oneof=manual auto"`
	BodyType         *string  `json:"bodyType" validate:"omitempty,oneof=sedan hatchback SUV crossover minivan coupe convertible pickup"`
	Class            *string  `json:"class" validate:"omitempty,oneof=economy compact comfort business luxury"`
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
		Brand:        req.Brand,
		Model:        req.Model,
		Year:         req.Year,
		FuelType:     req.FuelType,
		Transmission: req.Transmission,
		BodyType:     req.BodyType,
		Class:        req.Class,
		Seats:        req.Seats,
		EngineVolume: req.EngineVolume,
		RangeKM:      req.RangeKM,
		Features:     req.Features,
	}, nil
}

func FromCarModelUpdateRequest(ctx *gin.Context) (model.CarModelUpdate, error) {
	var req CarModelUpdateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		return model.CarModelUpdate{}, err
	}

	return model.CarModelUpdate{
		Brand:        req.Brand,
		Model:        req.Model,
		Year:         req.Year,
		FuelType:     req.FuelType,
		Transmission: req.Transmission,
		BodyType:     req.BodyType,
		Class:        req.Class,
		Seats:        req.Seats,
		EngineVolume: req.EngineVolume,
		RangeKM:      req.RangeKM,
		Features:     req.Features,
		ImageKeys:    req.ImageStorageKeys,
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
		ImageStorageUrls: m.ImageURLs,
		CreatedAt:        m.CreatedAt,
		UpdatedAt:        m.UpdatedAt,
	}
}
