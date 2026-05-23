package handler

import (
	"context"
	"log/slog"

	"carsharing/car-service/internal/adapter/grpc/dto"
	pkglog "carsharing/shared/pkg/log"

	carsvc "carsharing/protos/gen/service/car"

	"google.golang.org/protobuf/types/known/emptypb"
)

type CarInsuranceHandler struct {
	log              *slog.Logger
	insuranceService CarInsuranceService

	carsvc.UnimplementedCarInsuranceServiceServer
}

func NewCarInsuranceHandler(log *slog.Logger, insuranceService CarInsuranceService) *CarInsuranceHandler {
	return &CarInsuranceHandler{
		log:              pkglog.WithComponent(log, "adapter.grpc.handler.CarInsuranceHandler"),
		insuranceService: insuranceService,
	}
}

func (h *CarInsuranceHandler) CreateCarInsurance(ctx context.Context, req *carsvc.CreateCarInsuranceRequest) (*carsvc.CreateCarInsuranceResponse, error) {
	createInput := dto.FromCreateCarInsuranceRequest(req)

	id, err := h.insuranceService.Create(ctx, createInput)
	if err != nil {
		return nil, dto.FromErrorToStatusCode(err)
	}

	return &carsvc.CreateCarInsuranceResponse{Id: id}, nil
}

func (h *CarInsuranceHandler) GetCarInsurance(ctx context.Context, req *carsvc.GetCarInsuranceRequest) (*carsvc.GetCarInsuranceResponse, error) {
	insurance, err := h.insuranceService.Get(ctx, req.Id)
	if err != nil {
		return nil, dto.FromErrorToStatusCode(err)
	}

	return &carsvc.GetCarInsuranceResponse{CarInsurance: dto.ToCarInsuranceProto(insurance)}, nil
}

func (h *CarInsuranceHandler) ListCarInsurances(ctx context.Context, req *carsvc.ListCarInsurancesRequest) (*carsvc.ListCarInsurancesResponse, error) {
	filterInput := dto.FromListCarInsurancesRequest(req)

	insurances, err := h.insuranceService.List(ctx, filterInput)
	if err != nil {
		return nil, dto.FromErrorToStatusCode(err)
	}

	return &carsvc.ListCarInsurancesResponse{CarInsurances: dto.ToCarInsuranceProtos(insurances)}, nil
}

func (h *CarInsuranceHandler) UpdateCarInsurance(ctx context.Context, req *carsvc.UpdateCarInsuranceRequest) (*emptypb.Empty, error) {
	updateInput := dto.FromUpdateCarInsuranceRequest(req)

	if err := h.insuranceService.Update(ctx, req.Id, updateInput); err != nil {
		return nil, dto.FromErrorToStatusCode(err)
	}

	return &emptypb.Empty{}, nil
}

func (h *CarInsuranceHandler) DeleteCarInsurance(ctx context.Context, req *carsvc.DeleteCarInsuranceRequest) (*emptypb.Empty, error) {
	if err := h.insuranceService.Delete(ctx, req.Id); err != nil {
		return nil, dto.FromErrorToStatusCode(err)
	}

	return &emptypb.Empty{}, nil
}

func (h *CarInsuranceHandler) GetCarInsuranceImageUploadData(ctx context.Context, _ *emptypb.Empty) (*carsvc.GetCarInsuranceImageUploadDataResponse, error) {
	data, err := h.insuranceService.GetImageUploadData(ctx)
	if err != nil {
		return nil, dto.FromErrorToStatusCode(err)
	}

	return &carsvc.GetCarInsuranceImageUploadDataResponse{
		UploadData: dto.ToImageUploadData(data),
	}, nil
}
