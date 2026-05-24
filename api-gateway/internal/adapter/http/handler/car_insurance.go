package handler

import (
	"log/slog"

	"carsharing/api-gateway/internal/adapter/http/dto"
	pkglog "carsharing/shared/pkg/log"
	"carsharing/shared/pkg/utils"
	"github.com/gin-gonic/gin"
)

type CarInsuranceHandler struct {
	svc CarInsuranceService
	log *slog.Logger
}

func NewCarInsuranceHandler(svc CarInsuranceService, log *slog.Logger) *CarInsuranceHandler {
	return &CarInsuranceHandler{
		svc: svc,
		log: pkglog.WithComponent(log, "http.CarInsuranceHandler"),
	}
}

// Create (CarInsurance) godoc
// @Summary      Create car insurance record
// @Description  Records an ОСАГО or КАСКО policy for a specific car. type must be osago or kasko; expiresAt must be after startsAt; costTenge min 0; provider and policyNum 1–100 chars.
// @Tags         car-insurances
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body  body      dto.CarInsuranceCreateRequest  true  "Insurance payload"
// @Success      200   {object}  dto.IDResponse
// @Failure      400   {object}  dto.ErrorResponse
// @Failure      401   {object}  dto.ErrorResponse
// @Failure      409   {object}  dto.ErrorResponse
// @Failure      500   {object}  dto.ErrorResponse
// @Router       /car-insurances [post]
func (h *CarInsuranceHandler) Create(ctx *gin.Context) {
	log := pkglog.WithMetadata(pkglog.WithMethod(h.log, "Create"), utils.MetadataFromCtx(ctx))

	data, err := dto.FromCarInsuranceCreateRequest(ctx)
	if err != nil {
		dto.MalformedJson(ctx)

		return
	}

	id, err := h.svc.Create(ctx, data)
	if err != nil {
		log.Warn("creating car insurance", pkglog.Err(err))

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
// @Success      200  {object}  dto.CarInsuranceResponse
// @Failure      400  {object}  dto.ErrorResponse
// @Failure      401  {object}  dto.ErrorResponse
// @Failure      404  {object}  dto.ErrorResponse
// @Failure      500  {object}  dto.ErrorResponse
// @Router       /car-insurances/{id} [get]
func (h *CarInsuranceHandler) Get(ctx *gin.Context) {
	log := pkglog.WithMetadata(pkglog.WithMethod(h.log, "Get"), utils.MetadataFromCtx(ctx))

	id, err := dto.IDParam(ctx)
	if err != nil {
		dto.FromError(ctx, err)

		return
	}

	insurance, err := h.svc.Get(ctx, id)
	if err != nil {
		log.Warn("getting car insurance", pkglog.Err(err))

		dto.FromError(ctx, err)

		return
	}

	dto.Ok(ctx, gin.H{"insurance": dto.ToCarInsuranceResponse(insurance)})
}

// List (CarInsurance) godoc
// @Summary      List insurance records
// @Description  Returns insurance records filtered by car, type, status, or expiry window.
// @Tags         car-insurances
// @Produce      json
// @Security     BearerAuth
// @Param        carID               query     string   false  "Filter by car UUID"
// @Param        type                query     string   false  "Insurance type"  Enums(osago, kasko)
// @Param        status              query     string   false  "Status"          Enums(active, expired, cancelled)
// @Param        expiringWithinDays  query     integer  false  "Return policies expiring within N days (1–365)"
// @Param        limit               query     integer  false  "Pagination limit"
// @Param        offset              query     integer  false  "Pagination offset"
// @Success      200                 {object}  dto.CarInsurancesResponse
// @Failure      400                 {object}  dto.ErrorResponse
// @Failure      401                 {object}  dto.ErrorResponse
// @Failure      500                 {object}  dto.ErrorResponse
// @Router       /car-insurances [get]
func (h *CarInsuranceHandler) List(ctx *gin.Context) {
	log := pkglog.WithMetadata(pkglog.WithMethod(h.log, "List"), utils.MetadataFromCtx(ctx))

	filter, err := dto.CarInsuranceFilterFromCtx(ctx)
	if err != nil {
		dto.FromError(ctx, err)

		return
	}

	insurances, err := h.svc.List(ctx, filter)
	if err != nil {
		log.Warn("listing car insurances", pkglog.Err(err))

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
// @Success      204
// @Failure      400  {object}  dto.ErrorResponse
// @Failure      401  {object}  dto.ErrorResponse
// @Failure      404  {object}  dto.ErrorResponse
// @Failure      500  {object}  dto.ErrorResponse
// @Router       /car-insurances/{id} [patch]
func (h *CarInsuranceHandler) Update(ctx *gin.Context) {
	log := pkglog.WithMetadata(pkglog.WithMethod(h.log, "Update"), utils.MetadataFromCtx(ctx))

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

	if err = h.svc.Update(ctx, id, data); err != nil {
		log.Warn("updating car insurance", pkglog.Err(err))

		dto.FromError(ctx, err)

		return
	}

	dto.NoContent(ctx)
}

// Delete (CarInsurance) godoc
// @Summary      Delete insurance record
// @Tags         car-insurances
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string  true  "Insurance UUID"
// @Success      204
// @Failure      400  {object}  dto.ErrorResponse
// @Failure      401  {object}  dto.ErrorResponse
// @Failure      404  {object}  dto.ErrorResponse
// @Failure      500  {object}  dto.ErrorResponse
// @Router       /car-insurances/{id} [delete]
func (h *CarInsuranceHandler) Delete(ctx *gin.Context) {
	log := pkglog.WithMetadata(pkglog.WithMethod(h.log, "Delete"), utils.MetadataFromCtx(ctx))

	id, err := dto.IDParam(ctx)
	if err != nil {
		dto.FromError(ctx, err)

		return
	}

	if err = h.svc.Delete(ctx, id); err != nil {
		log.Warn("deleting car insurance", pkglog.Err(err))

		dto.FromError(ctx, err)

		return
	}

	dto.NoContent(ctx)
}

// GetImageUploadUrl (CarInsurance) godoc
// @Summary      Get pre-signed upload URL for insurance document images
// @Tags         car-insurances
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  dto.ImageUploadResponse
// @Failure      401  {object}  dto.ErrorResponse
// @Failure      500  {object}  dto.ErrorResponse
// @Router       /car-insurances/image-upload [get]
func (h *CarInsuranceHandler) GetImageUploadUrl(ctx *gin.Context) {
	log := pkglog.WithMetadata(pkglog.WithMethod(h.log, "GetImageUploadUrl"), utils.MetadataFromCtx(ctx))

	uploadData, err := h.svc.GetImageUploadData(ctx)
	if err != nil {
		log.Warn("getting car insurance image upload data", pkglog.Err(err))

		dto.FromError(ctx, err)

		return
	}

	dto.Ok(ctx, gin.H{"uploadData": dto.ToImageUploadDataResponse(uploadData)})
}
