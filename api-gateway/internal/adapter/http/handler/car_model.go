package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/sorawaslocked/car-rental-api-gateway/internal/adapter/http/dto"
)

type CarModelHandler struct {
	svc CarModelService
}

func NewCarModelHandler(svc CarModelService) *CarModelHandler {
	return &CarModelHandler{svc: svc}
}

// Create (CarModel) godoc
// @Summary      Create car model
// @Description  Registers a new car model (e.g. Hyundai Solaris 2022). Returns the new model's UUID.
// @Tags         car-models
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body  body      dto.CarModelCreateRequest  true  "Car model payload"
// @Success      201   {object}  map[string]any             "id"
// @Failure      400   {object}  map[string]any
// @Failure      401   {object}  map[string]any
// @Failure      409   {object}  map[string]any
// @Failure      500   {object}  map[string]any
// @Router       /car-models [post]
func (h *CarModelHandler) Create(ctx *gin.Context) {
	data, err := dto.FromCarModelCreateRequest(ctx)
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

// Get (CarModel) godoc
// @Summary      Get car model by ID
// @Tags         car-models
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string  true  "Car model UUID"
// @Success      200  {object}  map[string]any  "carModel"
// @Failure      400  {object}  map[string]any
// @Failure      404  {object}  map[string]any
// @Failure      500  {object}  map[string]any
// @Router       /car-models/{id} [get]
func (h *CarModelHandler) Get(ctx *gin.Context) {
	id, err := dto.IDParam(ctx)
	if err != nil {
		dto.FromError(ctx, err)

		return
	}

	carModel, err := h.svc.Get(ctx, id)
	if err != nil {
		dto.FromError(ctx, err)

		return
	}

	dto.Ok(ctx, gin.H{"carModel": dto.ToCarModelResponse(carModel)})
}

// GetAll (CarModel) godoc
// @Summary      List car models
// @Description  Returns a filtered, paginated list of car models.
// @Tags         car-models
// @Produce      json
// @Security     BearerAuth
// @Param        brand         query     string   false  "Brand filter (e.g. Hyundai)"
// @Param        model         query     string   false  "Model name filter"
// @Param        fuelType      query     string   false  "Fuel type (petrol, diesel, electric, hybrid)"
// @Param        transmission  query     string   false  "Transmission (automatic, manual)"
// @Param        bodyType      query     string   false  "Body type (sedan, suv, hatchback…)"
// @Param        class         query     string   false  "Class (economy, comfort, business)"
// @Param        minSeats      query     integer  false  "Minimum seat count"
// @Param        limit         query     integer  false  "Pagination limit"
// @Param        offset        query     integer  false  "Pagination offset"
// @Success      200           {object}  map[string]any  "carModels"
// @Failure      400           {object}  map[string]any
// @Failure      500           {object}  map[string]any
// @Router       /car-models [get]
func (h *CarModelHandler) GetAll(ctx *gin.Context) {
	filter, err := dto.CarModelFilterFromCtx(ctx)
	if err != nil {
		dto.FromError(ctx, err)

		return
	}

	carModels, err := h.svc.GetAll(ctx, filter)
	if err != nil {
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
// @Success      200   {object}  map[string]any
// @Failure      400   {object}  map[string]any
// @Failure      404   {object}  map[string]any
// @Failure      500   {object}  map[string]any
// @Router       /car-models/{id} [patch]
func (h *CarModelHandler) Update(ctx *gin.Context) {
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

	err = h.svc.Update(ctx, id, data)
	if err != nil {
		dto.FromError(ctx, err)

		return
	}

	dto.Ok(ctx, nil)
}

// Delete (CarModel) godoc
// @Summary      Delete car model
// @Tags         car-models
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string  true  "Car model UUID"
// @Success      200  {object}  map[string]any
// @Failure      400  {object}  map[string]any
// @Failure      404  {object}  map[string]any
// @Failure      500  {object}  map[string]any
// @Router       /car-models/{id} [delete]
func (h *CarModelHandler) Delete(ctx *gin.Context) {
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

// GetImageUploadUrl (CarModel) godoc
// @Summary      Get pre-signed image upload URL
// @Description  Returns a pre-signed S3-compatible URL and object key for uploading a car model image.
// @Tags         car-models
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  map[string]any  "uploadData"
// @Failure      401  {object}  map[string]any
// @Failure      500  {object}  map[string]any
// @Router       /car-models/image-upload [get]
func (h *CarModelHandler) GetImageUploadUrl(ctx *gin.Context) {
	uploadData, err := h.svc.GetImageUploadData(ctx)
	if err != nil {
		dto.FromError(ctx, err)

		return
	}

	dto.Ok(ctx, gin.H{"uploadData": dto.ToImageUploadDataResponse(uploadData)})
}
