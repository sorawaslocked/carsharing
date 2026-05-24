package handler

import (
	"log/slog"

	"carsharing/api-gateway/internal/adapter/http/dto"
	pkglog "carsharing/shared/pkg/log"
	"carsharing/shared/pkg/utils"
	"github.com/gin-gonic/gin"
)

type CarHandler struct {
	svc CarService
	log *slog.Logger
}

func NewCarHandler(svc CarService, log *slog.Logger) *CarHandler {
	return &CarHandler{
		svc: svc,
		log: pkglog.WithComponent(log, "http.CarHandler"),
	}
}

// Create (Car) godoc
// @Summary      Register a car in the fleet
// @Description  Creates a new physical car record linked to a car model and telemetry device. VIN must be exactly 17 alphanumeric characters. yearManufactured minimum 1886. fuelLevel and batteryLevel 0–100 if provided.
// @Tags         cars
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body  body      dto.CarCreateRequest  true  "Car create payload"
// @Success      201   {object}  dto.IDResponse
// @Failure      400   {object}  dto.ErrorResponse
// @Failure      401   {object}  dto.ErrorResponse
// @Failure      409   {object}  dto.ErrorResponse
// @Failure      500   {object}  dto.ErrorResponse
// @Router       /cars [post]
func (h *CarHandler) Create(ctx *gin.Context) {
	log := pkglog.WithMetadata(pkglog.WithMethod(h.log, "Create"), utils.MetadataFromCtx(ctx))

	data, err := dto.FromCarCreateRequest(ctx)
	if err != nil {
		dto.MalformedJson(ctx)

		return
	}

	id, err := h.svc.Create(ctx, data)
	if err != nil {
		log.Warn("creating car", pkglog.Err(err))

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
// @Success      200  {object}  dto.CarResponse
// @Failure      400  {object}  dto.ErrorResponse
// @Failure      401  {object}  dto.ErrorResponse
// @Failure      404  {object}  dto.ErrorResponse
// @Failure      500  {object}  dto.ErrorResponse
// @Router       /cars/{id} [get]
func (h *CarHandler) Get(ctx *gin.Context) {
	log := pkglog.WithMetadata(pkglog.WithMethod(h.log, "Get"), utils.MetadataFromCtx(ctx))

	id, err := dto.IDParam(ctx)
	if err != nil {
		dto.FromError(ctx, err)

		return
	}

	car, err := h.svc.Get(ctx, id)
	if err != nil {
		log.Warn("getting car", pkglog.Err(err))

		dto.FromError(ctx, err)

		return
	}

	dto.Ok(ctx, dto.ToCarResponse(car))
}

// List (Car) godoc
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
// @Param        latitude      query     number   false  "Latitude for proximity search"
// @Param        longitude     query     number   false  "Longitude for proximity search"
// @Param        radiusM       query     integer  false  "Search radius in metres"
// @Param        zoneId        query     string   false  "Filter by zone UUID"
// @Param        minFuelLevel  query     number   false  "Minimum fuel level (0–100)"
// @Param        status        query     string   false  "Car status (available, reserved, in_use, maintenance, out_of_service)"
// @Param        limit         query     integer  false  "Pagination limit"
// @Param        offset        query     integer  false  "Pagination offset"
// @Success      200           {object}  dto.CarsResponse
// @Failure      400           {object}  dto.ErrorResponse
// @Failure      401           {object}  dto.ErrorResponse
// @Failure      500           {object}  dto.ErrorResponse
// @Router       /cars [get]
func (h *CarHandler) List(ctx *gin.Context) {
	log := pkglog.WithMetadata(pkglog.WithMethod(h.log, "List"), utils.MetadataFromCtx(ctx))

	filter, err := dto.CarFilterFromCtx(ctx)
	if err != nil {
		dto.FromError(ctx, err)

		return
	}

	cars, err := h.svc.List(ctx, filter)
	if err != nil {
		log.Warn("listing cars", pkglog.Err(err))

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
// @Success      204
// @Failure      400  {object}  dto.ErrorResponse
// @Failure      401  {object}  dto.ErrorResponse
// @Failure      404  {object}  dto.ErrorResponse
// @Failure      500  {object}  dto.ErrorResponse
// @Router       /cars/{id} [patch]
func (h *CarHandler) Update(ctx *gin.Context) {
	log := pkglog.WithMetadata(pkglog.WithMethod(h.log, "Update"), utils.MetadataFromCtx(ctx))

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

	if err = h.svc.Update(ctx, id, data); err != nil {
		log.Warn("updating car", pkglog.Err(err))

		dto.FromError(ctx, err)

		return
	}

	dto.NoContent(ctx)
}

// Delete (Car) godoc
// @Summary      Delete car
// @Tags         cars
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string  true  "Car UUID"
// @Success      204
// @Failure      400  {object}  dto.ErrorResponse
// @Failure      401  {object}  dto.ErrorResponse
// @Failure      404  {object}  dto.ErrorResponse
// @Failure      500  {object}  dto.ErrorResponse
// @Router       /cars/{id} [delete]
func (h *CarHandler) Delete(ctx *gin.Context) {
	log := pkglog.WithMetadata(pkglog.WithMethod(h.log, "Delete"), utils.MetadataFromCtx(ctx))

	id, err := dto.IDParam(ctx)
	if err != nil {
		dto.FromError(ctx, err)

		return
	}

	if err = h.svc.Delete(ctx, id); err != nil {
		log.Warn("deleting car", pkglog.Err(err))

		dto.FromError(ctx, err)

		return
	}

	dto.NoContent(ctx)
}

// UpdateTelemetry (Car) godoc
// @Summary      Update car telemetry (admin)
// @Description  Records audited sensor readings (mileage, fuel, battery, location) for a car. mileageKm min 0; fuelLevel and batteryLevel 0–100 if provided; reason optional (1–500 chars).
// @Tags         cars
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id    path      string                          true  "Car UUID"
// @Param        body  body      dto.CarTelemetryUpdateRequest  true  "Telemetry update payload"
// @Success      204
// @Failure      400  {object}  dto.ErrorResponse
// @Failure      401  {object}  dto.ErrorResponse
// @Failure      404  {object}  dto.ErrorResponse
// @Failure      500  {object}  dto.ErrorResponse
// @Router       /cars/{id}/telemetry [patch]
func (h *CarHandler) UpdateTelemetry(ctx *gin.Context) {
	log := pkglog.WithMetadata(pkglog.WithMethod(h.log, "UpdateTelemetry"), utils.MetadataFromCtx(ctx))

	id, err := dto.IDParam(ctx)
	if err != nil {
		dto.FromError(ctx, err)

		return
	}

	data, err := dto.FromCarTelemetryUpdateRequest(ctx)
	if err != nil {
		dto.MalformedJson(ctx)

		return
	}

	if err = h.svc.UpdateTelemetry(ctx, id, data); err != nil {
		log.Warn("updating car telemetry", pkglog.Err(err))

		dto.FromError(ctx, err)

		return
	}

	dto.NoContent(ctx)
}

// UpdateStatus (Car) godoc
// @Summary      Update car status (admin)
// @Description  Records an audited operational status transition for a car. status must be one of: available, reserved, in_use, maintenance, out_of_service. reason optional (1–500 chars).
// @Tags         cars
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id    path      string                       true  "Car UUID"
// @Param        body  body      dto.CarStatusUpdateRequest  true  "Status update payload"
// @Success      204
// @Failure      400  {object}  dto.ErrorResponse
// @Failure      401  {object}  dto.ErrorResponse
// @Failure      404  {object}  dto.ErrorResponse
// @Failure      500  {object}  dto.ErrorResponse
// @Router       /cars/{id}/status [patch]
func (h *CarHandler) UpdateStatus(ctx *gin.Context) {
	log := pkglog.WithMetadata(pkglog.WithMethod(h.log, "UpdateStatus"), utils.MetadataFromCtx(ctx))

	id, err := dto.IDParam(ctx)
	if err != nil {
		dto.FromError(ctx, err)

		return
	}

	data, err := dto.FromCarStatusUpdateRequest(ctx)
	if err != nil {
		dto.MalformedJson(ctx)

		return
	}

	if err = h.svc.UpdateStatus(ctx, id, data); err != nil {
		log.Warn("updating car status", pkglog.Err(err))

		dto.FromError(ctx, err)

		return
	}

	dto.NoContent(ctx)
}

// GetCarStatusHistory godoc
// @Summary      Get car status change history
// @Description  Returns the immutable audit trail of every car state transition.
// @Tags         cars
// @Produce      json
// @Security     BearerAuth
// @Param        id      path      string   true   "Car UUID"
// @Param        from    query     string   false  "Start date (YYYY-MM-DD)"
// @Param        to      query     string   false  "End date (YYYY-MM-DD)"
// @Param        limit   query     integer  false  "Pagination limit"
// @Param        offset  query     integer  false  "Pagination offset"
// @Success      200     {object}  dto.CarStatusHistoryResponse
// @Failure      400     {object}  dto.ErrorResponse
// @Failure      401     {object}  dto.ErrorResponse
// @Failure      404     {object}  dto.ErrorResponse
// @Failure      500     {object}  dto.ErrorResponse
// @Router       /cars/{id}/status-history [get]
func (h *CarHandler) GetCarStatusHistory(ctx *gin.Context) {
	log := pkglog.WithMetadata(pkglog.WithMethod(h.log, "GetCarStatusHistory"), utils.MetadataFromCtx(ctx))

	carID, err := dto.IDParam(ctx)
	if err != nil {
		dto.FromError(ctx, err)

		return
	}

	filter, err := dto.CarStatusReadingFilterFromCtx(ctx)
	if err != nil {
		dto.FromError(ctx, err)

		return
	}

	history, err := h.svc.GetCarStatusHistory(ctx, carID, filter)
	if err != nil {
		log.Warn("getting car status history", pkglog.Err(err))

		dto.FromError(ctx, err)

		return
	}

	response := make([]dto.CarStatusReading, len(history))
	for i, r := range history {
		response[i] = dto.ToCarStatusReadingResponse(r)
	}

	dto.Ok(ctx, gin.H{"statusHistory": response})
}

// GetCarTelemetryHistory godoc
// @Summary      Get car telemetry history
// @Description  Returns time-series telemetry readings for a car.
// @Tags         cars
// @Produce      json
// @Security     BearerAuth
// @Param        id      path      string   true   "Car UUID"
// @Param        from    query     string   false  "Start date (YYYY-MM-DD)"
// @Param        to      query     string   false  "End date (YYYY-MM-DD)"
// @Param        limit   query     integer  false  "Pagination limit"
// @Param        offset  query     integer  false  "Pagination offset"
// @Success      200     {object}  dto.CarTelemetryHistoryResponse
// @Failure      400     {object}  dto.ErrorResponse
// @Failure      401     {object}  dto.ErrorResponse
// @Failure      404     {object}  dto.ErrorResponse
// @Failure      500     {object}  dto.ErrorResponse
// @Router       /cars/{id}/telemetry-history [get]
func (h *CarHandler) GetCarTelemetryHistory(ctx *gin.Context) {
	log := pkglog.WithMetadata(pkglog.WithMethod(h.log, "GetCarTelemetryHistory"), utils.MetadataFromCtx(ctx))

	carID, err := dto.IDParam(ctx)
	if err != nil {
		dto.FromError(ctx, err)

		return
	}

	filter, err := dto.CarTelemetryReadingFilterFromCtx(ctx)
	if err != nil {
		dto.FromError(ctx, err)

		return
	}

	history, err := h.svc.GetCarTelemetryHistory(ctx, carID, filter)
	if err != nil {
		log.Warn("getting car telemetry history", pkglog.Err(err))

		dto.FromError(ctx, err)

		return
	}

	response := make([]dto.CarTelemetryReading, len(history))
	for i, r := range history {
		response[i] = dto.ToCarTelemetryReadingResponse(r)
	}

	dto.Ok(ctx, gin.H{"telemetryHistory": response})
}

// GetImageUploadUrl (Car) godoc
// @Summary      Get pre-signed image upload URL for car photos
// @Description  Returns a pre-signed S3-compatible URL and object key for uploading a car photo.
// @Tags         cars
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  dto.ImageUploadResponse
// @Failure      401  {object}  dto.ErrorResponse
// @Failure      500  {object}  dto.ErrorResponse
// @Router       /cars/image-upload [get]
func (h *CarHandler) GetImageUploadUrl(ctx *gin.Context) {
	log := pkglog.WithMetadata(pkglog.WithMethod(h.log, "GetImageUploadUrl"), utils.MetadataFromCtx(ctx))

	uploadData, err := h.svc.GetImageUploadData(ctx)
	if err != nil {
		log.Warn("getting car image upload data", pkglog.Err(err))

		dto.FromError(ctx, err)

		return
	}

	dto.Ok(ctx, gin.H{"uploadData": dto.ToImageUploadDataResponse(uploadData)})
}
