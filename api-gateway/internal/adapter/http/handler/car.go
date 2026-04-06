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

func (h *CarHandler) GetImageUploadUrl(ctx *gin.Context) {
	uploadData, err := h.svc.GetImageUploadData(ctx)
	if err != nil {
		dto.FromError(ctx, err)

		return
	}

	dto.Ok(ctx, gin.H{"uploadData": dto.ToImageUploadDataResponse(uploadData)})
}
