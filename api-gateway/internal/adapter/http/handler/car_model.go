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

func (h *CarModelHandler) GetImageUploadUrl(ctx *gin.Context) {
	uploadData, err := h.svc.GetImageUploadData(ctx)
	if err != nil {
		dto.FromError(ctx, err)

		return
	}

	dto.Ok(ctx, gin.H{"uploadData": dto.ToImageUploadDataResponse(uploadData)})
}
