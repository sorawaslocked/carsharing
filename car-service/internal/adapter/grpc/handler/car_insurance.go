package handler

import (
	"context"
	"log/slog"

	"github.com/sorawaslocked/car-rental-car-service/internal/adapter/grpc/dto"

	carsvc "github.com/sorawaslocked/car-rental-protos/gen/service/car"
	"google.golang.org/protobuf/types/known/emptypb"
)

type CarInsuranceHandler struct {
	insuranceService CarInsuranceService

	log *slog.Logger

	carsvc.UnimplementedCarInsuranceServiceServer
}

func NewCarInsuranceHandler(insuranceService CarInsuranceService, log *slog.Logger) *CarInsuranceHandler {
	h := &CarInsuranceHandler{
		insuranceService: insuranceService,
	}

	h.log = log.With(
		slog.Group("src",
			slog.String("component", "CarInsuranceHandler"),
		),
	)

	return h
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

	insurances, err := h.insuranceService.GetAll(ctx, filterInput)
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
		UploadData: dto.ToImageUploadData(data.URL, data.ObjectKey),
	}, nil
}
