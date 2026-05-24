package handler

import (
	"log/slog"

	"carsharing/api-gateway/internal/adapter/http/dto"
	pkglog "carsharing/shared/pkg/log"
	"carsharing/shared/pkg/utils"
	"github.com/gin-gonic/gin"
)

type CarModelHandler struct {
	svc CarModelService
	log *slog.Logger
}

func NewCarModelHandler(svc CarModelService, log *slog.Logger) *CarModelHandler {
	return &CarModelHandler{
		svc: svc,
		log: pkglog.WithComponent(log, "http.CarModelHandler"),
	}
}

// Create (CarModel) godoc
// @Summary      Create car model
// @Description  Registers a new car model (e.g. Hyundai Solaris 2022). Returns the new model's UUID. year min 1886; seats 1–9; engineVolume 0.1–28.5 L if provided; warnPct and pullPct are proportions 0.0–1.0 (not percent), pullPct must be ≥ warnPct.
// @Tags         car-models
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body  body      dto.CarModelCreateRequest  true  "Car model payload"
// @Success      201   {object}  dto.IDResponse
// @Failure      400   {object}  dto.ErrorResponse
// @Failure      401   {object}  dto.ErrorResponse
// @Failure      409   {object}  dto.ErrorResponse
// @Failure      500   {object}  dto.ErrorResponse
// @Router       /car-models [post]
func (h *CarModelHandler) Create(ctx *gin.Context) {
	log := pkglog.WithMetadata(pkglog.WithMethod(h.log, "Create"), utils.MetadataFromCtx(ctx))

	data, err := dto.FromCarModelCreateRequest(ctx)
	if err != nil {
		dto.MalformedJson(ctx)

		return
	}

	id, err := h.svc.Create(ctx, data)
	if err != nil {
		log.Warn("creating car model", pkglog.Err(err))

		dto.FromError(ctx, err)

		return
	}

	dto.Created(ctx, gin.H{"id": id})
}

// Get (CarModel) godoc
// @Summary      Get car model by ID
// @Tags         car-models
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string  true  "Car model UUID"
// @Success      200  {object}  dto.CarModelGetResponse
// @Failure      400  {object}  dto.ErrorResponse
// @Failure      401  {object}  dto.ErrorResponse
// @Failure      404  {object}  dto.ErrorResponse
// @Failure      500  {object}  dto.ErrorResponse
// @Router       /car-models/{id} [get]
func (h *CarModelHandler) Get(ctx *gin.Context) {
	log := pkglog.WithMetadata(pkglog.WithMethod(h.log, "Get"), utils.MetadataFromCtx(ctx))

	id, err := dto.IDParam(ctx)
	if err != nil {
		dto.FromError(ctx, err)

		return
	}

	carModel, err := h.svc.Get(ctx, id)
	if err != nil {
		log.Warn("getting car model", pkglog.Err(err))

		dto.FromError(ctx, err)

		return
	}

	dto.Ok(ctx, gin.H{"carModel": dto.ToCarModelResponse(carModel)})
}

// List (CarModel) godoc
// @Summary      List car models
// @Description  Returns a filtered, paginated list of car models.
// @Tags         car-models
// @Produce      json
// @Security     BearerAuth
// @Param        brand         query     string   false  "Brand filter (e.g. Hyundai)"
// @Param        model         query     string   false  "Model name filter"
// @Param        fuelType      query     string   false  "Fuel type (petrol, diesel, electric, hybrid)"
// @Param        transmission  query     string   false  "Transmission (manual, auto)"
// @Param        bodyType      query     string   false  "Body type (sedan, hatchback, suv, crossover, minivan, coupe, convertible, pickup)"
// @Param        class         query     string   false  "Class (economy, compact, comfort, business, luxury)"
// @Param        minSeats      query     integer  false  "Minimum seat count"
// @Param        limit         query     integer  false  "Pagination limit"
// @Param        offset        query     integer  false  "Pagination offset"
// @Success      200           {object}  dto.CarModelsResponse
// @Failure      400           {object}  dto.ErrorResponse
// @Failure      401           {object}  dto.ErrorResponse
// @Failure      500           {object}  dto.ErrorResponse
// @Router       /car-models [get]
func (h *CarModelHandler) List(ctx *gin.Context) {
	log := pkglog.WithMetadata(pkglog.WithMethod(h.log, "List"), utils.MetadataFromCtx(ctx))

	filter, err := dto.CarModelFilterFromCtx(ctx)
	if err != nil {
		dto.FromError(ctx, err)

		return
	}

	carModels, err := h.svc.List(ctx, filter)
	if err != nil {
		log.Warn("listing car models", pkglog.Err(err))

		dto.FromError(ctx, err)

		return
	}

	carModelResponse := make([]dto.CarModel, len(carModels))
	for i, carModel := range carModels {
		carModelResponse[i] = dto.ToCarModelResponse(carModel)
	}

	dto.Ok(ctx, gin.H{"carModels": carModelResponse})
}

// Update (CarModel) godoc
// @Summary      Update car model
// @Tags         car-models
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id    path      string                     true  "Car model UUID"
// @Param        body  body      dto.CarModelUpdateRequest  true  "Fields to update"
// @Success      204
// @Failure      400  {object}  dto.ErrorResponse
// @Failure      401  {object}  dto.ErrorResponse
// @Failure      404  {object}  dto.ErrorResponse
// @Failure      500  {object}  dto.ErrorResponse
// @Router       /car-models/{id} [patch]
func (h *CarModelHandler) Update(ctx *gin.Context) {
	log := pkglog.WithMetadata(pkglog.WithMethod(h.log, "Update"), utils.MetadataFromCtx(ctx))

	id, err := dto.IDParam(ctx)
	if err != nil {
		dto.FromError(ctx, err)

		return
	}

	data, err := dto.FromCarModelUpdateRequest(ctx)
	if err != nil {
		dto.MalformedJson(ctx)

		return
	}

	if err = h.svc.Update(ctx, id, data); err != nil {
		log.Warn("updating car model", pkglog.Err(err))

		dto.FromError(ctx, err)

		return
	}

	dto.NoContent(ctx)
}

// Delete (CarModel) godoc
// @Summary      Delete car model
// @Tags         car-models
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string  true  "Car model UUID"
// @Success      204
// @Failure      400  {object}  dto.ErrorResponse
// @Failure      401  {object}  dto.ErrorResponse
// @Failure      404  {object}  dto.ErrorResponse
// @Failure      500  {object}  dto.ErrorResponse
// @Router       /car-models/{id} [delete]
func (h *CarModelHandler) Delete(ctx *gin.Context) {
	log := pkglog.WithMetadata(pkglog.WithMethod(h.log, "Delete"), utils.MetadataFromCtx(ctx))

	id, err := dto.IDParam(ctx)
	if err != nil {
		dto.FromError(ctx, err)

		return
	}

	if err = h.svc.Delete(ctx, id); err != nil {
		log.Warn("deleting car model", pkglog.Err(err))

		dto.FromError(ctx, err)

		return
	}

	dto.NoContent(ctx)
}

// GetImageUploadUrl (CarModel) godoc
// @Summary      Get pre-signed image upload URL
// @Description  Returns a pre-signed S3-compatible URL and object key for uploading a car model image.
// @Tags         car-models
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  dto.ImageUploadResponse
// @Failure      401  {object}  dto.ErrorResponse
// @Failure      500  {object}  dto.ErrorResponse
// @Router       /car-models/image-upload [get]
func (h *CarModelHandler) GetImageUploadUrl(ctx *gin.Context) {
	log := pkglog.WithMetadata(pkglog.WithMethod(h.log, "GetImageUploadUrl"), utils.MetadataFromCtx(ctx))

	uploadData, err := h.svc.GetImageUploadData(ctx)
	if err != nil {
		log.Warn("getting car model image upload data", pkglog.Err(err))

		dto.FromError(ctx, err)

		return
	}

	dto.Ok(ctx, gin.H{"uploadData": dto.ToImageUploadDataResponse(uploadData)})
}
