package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/sorawaslocked/car-rental-api-gateway/internal/adapter/http/dto"
)

type CarInsuranceHandler struct {
	svc CarInsuranceService
}

func NewCarInsuranceHandler(svc CarInsuranceService) *CarInsuranceHandler {
	return &CarInsuranceHandler{svc: svc}
}

// Create (CarInsurance) godoc
// @Summary      Create car insurance record
// @Description  Records an ОСАГО or КАСКО policy for a specific car.
// @Tags         car-insurances
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body  body      dto.CarInsuranceCreateRequest  true  "Insurance payload"
// @Success      200   {object}  map[string]any                 "id"
// @Failure      400   {object}  map[string]any
// @Failure      401   {object}  map[string]any
// @Failure      409   {object}  map[string]any
// @Failure      500   {object}  map[string]any
// @Router       /car-insurances [post]
func (h *CarInsuranceHandler) Create(ctx *gin.Context) {
	data, err := dto.FromCarInsuranceCreateRequest(ctx)
	if err != nil {
		dto.MalformedJson(ctx)

		return
	}

	id, err := h.svc.Create(ctx, data)
	if err != nil {
		dto.FromError(ctx, err)

		return
	}

	dto.Ok(ctx, gin.H{"id": id})
}

// Get (CarInsurance) godoc
// @Summary      Get insurance record by ID
// @Tags         car-insurances
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string  true  "Insurance UUID"
// @Success      200  {object}  map[string]any  "insurance"
// @Failure      400  {object}  map[string]any
// @Failure      404  {object}  map[string]any
// @Failure      500  {object}  map[string]any
// @Router       /car-insurances/{id} [get]
func (h *CarInsuranceHandler) Get(ctx *gin.Context) {
	id, err := dto.IDParam(ctx)
	if err != nil {
		dto.FromError(ctx, err)

		return
	}

	insurance, err := h.svc.Get(ctx, id)
	if err != nil {
		dto.FromError(ctx, err)

		return
	}

	dto.Ok(ctx, gin.H{"insurance": dto.ToCarInsuranceResponse(insurance)})
}

// GetAll (CarInsurance) godoc
// @Summary      List insurance records
// @Description  Returns insurance records filtered by car, type, status, or expiry window.
// @Tags         car-insurances
// @Produce      json
// @Security     BearerAuth
// @Param        carID               query     string   false  "Filter by car UUID"
// @Param        type                query     string   false  "Insurance type (osago, kasko)"
// @Param        status              query     string   false  "Status (active, expired, cancelled)"
// @Param        expiringWithinDays  query     integer  false  "Return policies expiring within N days"
// @Param        limit               query     integer  false  "Pagination limit"
// @Param        offset              query     integer  false  "Pagination offset"
// @Success      200                 {object}  map[string]any  "insurances"
// @Failure      400                 {object}  map[string]any
// @Failure      500                 {object}  map[string]any
// @Router       /car-insurances [get]
func (h *CarInsuranceHandler) GetAll(ctx *gin.Context) {
	filter, err := dto.CarInsuranceFilterFromCtx(ctx)
	if err != nil {
		dto.FromError(ctx, err)

		return
	}

	insurances, err := h.svc.GetAll(ctx, filter)
	if err != nil {
		dto.FromError(ctx, err)

		return
	}

	insuranceResponse := make([]dto.CarInsurance, len(insurances))
	for i, insurance := range insurances {
		insuranceResponse[i] = dto.ToCarInsuranceResponse(insurance)
	}

	dto.Ok(ctx, gin.H{"insurances": insuranceResponse})
}

// Update (CarInsurance) godoc
// @Summary      Update insurance record
// @Tags         car-insurances
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id    path      string                         true  "Insurance UUID"
// @Param        body  body      dto.CarInsuranceUpdateRequest  true  "Fields to update"
// @Success      200   {object}  map[string]any
// @Failure      400   {object}  map[string]any
// @Failure      404   {object}  map[string]any
// @Failure      500   {object}  map[string]any
// @Router       /car-insurances/{id} [patch]
func (h *CarInsuranceHandler) Update(ctx *gin.Context) {
	id, err := dto.IDParam(ctx)
	if err != nil {
		dto.FromError(ctx, err)

		return
	}

	data, err := dto.FromCarInsuranceUpdateRequest(ctx)
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

// Delete (CarInsurance) godoc
// @Summary      Delete insurance record
// @Tags         car-insurances
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string  true  "Insurance UUID"
// @Success      200  {object}  map[string]any
// @Failure      400  {object}  map[string]any
// @Failure      404  {object}  map[string]any
// @Failure      500  {object}  map[string]any
// @Router       /car-insurances/{id} [delete]
func (h *CarInsuranceHandler) Delete(ctx *gin.Context) {
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

// GetImageUploadUrl (CarInsurance) godoc
// @Summary      Get pre-signed upload URL for insurance document images
// @Tags         car-insurances
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  map[string]any  "uploadData"
// @Failure      401  {object}  map[string]any
// @Failure      500  {object}  map[string]any
// @Router       /car-insurances/image-upload [get]
func (h *CarInsuranceHandler) GetImageUploadUrl(ctx *gin.Context) {
	uploadData, err := h.svc.GetImageUploadData(ctx)
	if err != nil {
		dto.FromError(ctx, err)

		return
	}

	dto.Ok(ctx, gin.H{"uploadData": dto.ToImageUploadDataResponse(uploadData)})
}
