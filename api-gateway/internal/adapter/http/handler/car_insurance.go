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

func (h *CarInsuranceHandler) GetImageUploadUrl(ctx *gin.Context) {
	uploadData, err := h.svc.GetImageUploadData(ctx)
	if err != nil {
		dto.FromError(ctx, err)

		return
	}

	dto.Ok(ctx, gin.H{"uploadData": dto.ToImageUploadDataResponse(uploadData)})
}
