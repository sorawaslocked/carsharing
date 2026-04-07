package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/sorawaslocked/car-rental-api-gateway/internal/adapter/http/dto"
)

type CarHandler struct {
	svc CarService
}

func NewCarHandler(svc CarService) *CarHandler {
	return &CarHandler{svc: svc}
}

// Create (Car) godoc
// @Summary      Register a car in the fleet
// @Description  Creates a new physical car record linked to a car model and telematics device.
// @Tags         cars
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body  body      dto.CarCreateRequest  true  "Car create payload"
// @Success      201   {object}  map[string]any        "id"
// @Failure      400   {object}  map[string]any
// @Failure      401   {object}  map[string]any
// @Failure      409   {object}  map[string]any
// @Failure      500   {object}  map[string]any
// @Router       /cars [post]
func (h *CarHandler) Create(ctx *gin.Context) {
	data, err := dto.FromCarCreateRequest(ctx)
	if err != nil {
		dto.MalformedJson(ctx)

		return
	}

	id, err := h.svc.Create(ctx, data)
	if err != nil {
		dto.FromError(ctx, err)

		return
	}

	dto.Created(ctx, gin.H{"id": id})
}

// Get (Car) godoc
// @Summary      Get car by ID
// @Tags         cars
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string  true  "Car UUID"
// @Success      200  {object}  dto.Car
// @Failure      400  {object}  map[string]any
// @Failure      404  {object}  map[string]any
// @Failure      500  {object}  map[string]any
// @Router       /cars/{id} [get]
func (h *CarHandler) Get(ctx *gin.Context) {
	id, err := dto.IDParam(ctx)
	if err != nil {
		dto.FromError(ctx, err)

		return
	}

	car, err := h.svc.Get(ctx, id)
	if err != nil {
		dto.FromError(ctx, err)

		return
	}

	dto.Ok(ctx, dto.ToCarResponse(car))
}

// GetAll (Car) godoc
// @Summary      List available cars
// @Description  Returns cars filtered by status, location radius, fuel level, zone, and model attributes. Pass latitude+longitude+radiusM to get nearby cars sorted by distance.
// @Tags         cars
// @Produce      json
// @Security     BearerAuth
// @Param        brand         query     string   false  "Brand filter"
// @Param        model         query     string   false  "Model name filter"
// @Param        fuelType      query     string   false  "Fuel type"
// @Param        transmission  query     string   false  "Transmission"
// @Param        bodyType      query     string   false  "Body type"
// @Param        class         query     string   false  "Class"
// @Param        minSeats      query     integer  false  "Minimum seats"
// @Param        latitude      query     number   false  "User latitude (for proximity search)"
// @Param        longitude     query     number   false  "User longitude (for proximity search)"
// @Param        radiusM       query     integer  false  "Search radius in metres"
// @Param        zoneId        query     string   false  "Filter by zone UUID"
// @Param        minFuelLevel  query     number   false  "Minimum fuel level (0–100)"
// @Param        status        query     string   false  "Car status (available, reserved, in_trip, …)"
// @Param        limit         query     integer  false  "Pagination limit"
// @Param        offset        query     integer  false  "Pagination offset"
// @Success      200           {object}  map[string]any  "cars"
// @Failure      400           {object}  map[string]any
// @Failure      500           {object}  map[string]any
// @Router       /cars [get]
func (h *CarHandler) GetAll(ctx *gin.Context) {
	filter, err := dto.CarFilterFromCtx(ctx)
	if err != nil {
		dto.FromError(ctx, err)

		return
	}

	cars, err := h.svc.GetAll(ctx, filter)
	if err != nil {
		dto.FromError(ctx, err)

		return
	}

	carResponse := make([]dto.Car, len(cars))
	for i, car := range cars {
		carResponse[i] = dto.ToCarResponse(car)
	}

	dto.Ok(ctx, gin.H{"cars": carResponse})
}

// Update (Car) godoc
// @Summary      Update car
// @Tags         cars
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id    path      string                true  "Car UUID"
// @Param        body  body      dto.CarUpdateRequest  true  "Fields to update"
// @Success      200   {object}  map[string]any
// @Failure      400   {object}  map[string]any
// @Failure      404   {object}  map[string]any
// @Failure      500   {object}  map[string]any
// @Router       /cars/{id} [patch]
func (h *CarHandler) Update(ctx *gin.Context) {
	id, err := dto.IDParam(ctx)
	if err != nil {
		dto.FromError(ctx, err)

		return
	}

	data, err := dto.FromCarUpdateRequest(ctx)
	if err != nil {
		dto.MalformedJson(ctx)

		return
	}

	err = h.svc.Update(ctx, id, data)
	if err != nil {
		dto.FromError(ctx, err)

		return
	}

	dto.Ok(ctx, nil)
}

// Delete (Car) godoc
// @Summary      Delete car
// @Tags         cars
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string  true  "Car UUID"
// @Success      200  {object}  map[string]any
// @Failure      400  {object}  map[string]any
// @Failure      404  {object}  map[string]any
// @Failure      500  {object}  map[string]any
// @Router       /cars/{id} [delete]
func (h *CarHandler) Delete(ctx *gin.Context) {
	id, err := dto.IDParam(ctx)
	if err != nil {
		dto.FromError(ctx, err)

		return
	}

	err = h.svc.Delete(ctx, id)
	if err != nil {
		dto.FromError(ctx, err)

		return
	}

	dto.Ok(ctx, nil)
}

// GetCarStatusLog godoc
// @Summary      Get car status change log
// @Description  Returns the immutable audit trail of every car state transition (available → reserved → in_trip → …).
// @Tags         cars
// @Produce      json
// @Security     BearerAuth
// @Param        carID   query     string   false  "Filter by car UUID"
// @Param        limit   query     integer  false  "Pagination limit"
// @Param        offset  query     integer  false  "Pagination offset"
// @Success      200     {object}  map[string]any  "logs"
// @Failure      400     {object}  map[string]any
// @Failure      500     {object}  map[string]any
// @Router       /cars/status-log [get]
func (h *CarHandler) GetCarStatusLog(ctx *gin.Context) {
	filter, err := dto.CarStatusLogFilterFromCtx(ctx)
	if err != nil {
		dto.FromError(ctx, err)

		return
	}

	logs, err := h.svc.GetCarStatusLog(ctx, filter)
	if err != nil {
		dto.FromError(ctx, err)

		return
	}

	logResponse := make([]dto.CarStatusLogEntry, len(logs))
	for i, le := range logs {
		logResponse[i] = dto.ToCarStatusLogEntryResponse(le)
	}

	dto.Ok(ctx, gin.H{"logs": logResponse})
}

// GetCarFuelHistory godoc
// @Summary      Get car fuel reading history
// @Description  Returns time-series fuel readings for a car. Useful for detecting theft (sudden drops) or verifying refuels.
// @Tags         cars
// @Produce      json
// @Security     BearerAuth
// @Param        carID   query     string  false  "Filter by car UUID"
// @Param        from    query     string  false  "Start date (YYYY-MM-DD)"
// @Param        to      query     string  false  "End date (YYYY-MM-DD)"
// @Param        limit   query     integer false  "Pagination limit"
// @Param        offset  query     integer false  "Pagination offset"
// @Success      200     {object}  map[string]any  "fuelHistory"
// @Failure      400     {object}  map[string]any
// @Failure      500     {object}  map[string]any
// @Router       /cars/fuel-history [get]
func (h *CarHandler) GetCarFuelHistory(ctx *gin.Context) {
	filter, err := dto.CarFuelReadingFilterFromCtx(ctx)
	if err != nil {
		dto.FromError(ctx, err)

		return
	}

	fuelHistory, err := h.svc.GetCarFuelHistory(ctx, filter)
	if err != nil {
		dto.FromError(ctx, err)

		return
	}

	fuelHistoryResponse := make([]dto.CarFuelReading, len(fuelHistory))
	for i, fh := range fuelHistory {
		fuelHistoryResponse[i] = dto.ToCarFuelReadingResponse(fh)
	}

	dto.Ok(ctx, gin.H{"fuelHistory": fuelHistoryResponse})
}

// GetImageUploadUrl (Car) godoc
// @Summary      Get pre-signed image upload URL for car photos
// @Description  Returns a pre-signed S3-compatible URL and object key for uploading a car or trip photo.
// @Tags         cars
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  map[string]any  "uploadData"
// @Failure      401  {object}  map[string]any
// @Failure      500  {object}  map[string]any
// @Router       /cars/image-upload [get]
func (h *CarHandler) GetImageUploadUrl(ctx *gin.Context) {
	uploadData, err := h.svc.GetImageUploadData(ctx)
	if err != nil {
		dto.FromError(ctx, err)

		return
	}

	dto.Ok(ctx, gin.H{"uploadData": dto.ToImageUploadDataResponse(uploadData)})
}
