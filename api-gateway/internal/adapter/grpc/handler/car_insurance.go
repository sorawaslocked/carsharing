package handler

import (
	"context"
	"log/slog"

	"carsharing/api-gateway/internal/adapter/grpc/dto"
	"carsharing/api-gateway/internal/model"
	pkglog "carsharing/shared/pkg/log"
	"carsharing/shared/pkg/utils"
	basepb "github.com/sorawaslocked/car-rental-protos/gen/base"
	carsvc "github.com/sorawaslocked/car-rental-protos/gen/service/car"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type CarInsuranceHandler struct {
	client carsvc.CarInsuranceServiceClient
	log    *slog.Logger
}

func NewCarInsuranceHandler(client carsvc.CarInsuranceServiceClient, logger *slog.Logger) *CarInsuranceHandler {
	return &CarInsuranceHandler{
		client: client,
		log:    pkglog.WithComponent(logger, "grpc.CarInsuranceHandler"),
	}
}

func (h *CarInsuranceHandler) Create(ctx context.Context, data model.CarInsuranceCreate) (string, error) {
	logger := pkglog.WithMethod(h.log, "Create")
	logger = pkglog.WithMetadata(logger, utils.MetadataFromCtx(ctx))

	res, err := h.client.CreateCarInsurance(ctx, &carsvc.CreateCarInsuranceRequest{
		CarId:     data.CarID,
		Type:      data.Type,
		Provider:  data.Provider,
		PolicyNum: data.PolicyNum,
		StartsAt:  timestamppb.New(data.StartsAt),
		ExpiresAt: timestamppb.New(data.ExpiresAt),
		CostTenge: data.CostTenge,
		Notes:     data.Notes,
	})
	if err != nil {
		if dto.IsSystemErr(err) {
			logger.Error("grpc call failed", pkglog.Err(err))
		}

		return "", dto.FromGrpcErr(err)
	}

	return res.GetId(), nil
}

func (h *CarInsuranceHandler) Get(ctx context.Context, id string) (model.CarInsurance, error) {
	logger := pkglog.WithMethod(h.log, "Get")
	logger = pkglog.WithMetadata(logger, utils.MetadataFromCtx(ctx))

	res, err := h.client.GetCarInsurance(ctx, &carsvc.GetCarInsuranceRequest{Id: id})
	if err != nil {
		if dto.IsSystemErr(err) {
			logger.Error("grpc call failed", pkglog.Err(err))
		}

		return model.CarInsurance{}, dto.FromGrpcErr(err)
	}

	return dto.CarInsuranceFromProto(res.GetCarInsurance()), nil
}

func (h *CarInsuranceHandler) List(ctx context.Context, filter model.CarInsuranceFilter) ([]model.CarInsurance, error) {
	logger := pkglog.WithMethod(h.log, "List")
	logger = pkglog.WithMetadata(logger, utils.MetadataFromCtx(ctx))

	req := &carsvc.ListCarInsurancesRequest{
		CarId:              filter.CarID,
		Type:               filter.Type,
		Status:             filter.Status,
		ExpiringWithinDays: filter.ExpiringWithinDays,
	}
	if filter.Pagination != nil {
		req.Pagination = &basepb.Pagination{
			Limit:  filter.Pagination.Limit,
			Offset: filter.Pagination.Offset,
		}
	}

	res, err := h.client.ListCarInsurances(ctx, req)
	if err != nil {
		if dto.IsSystemErr(err) {
			logger.Error("grpc call failed", pkglog.Err(err))
		}

		return nil, dto.FromGrpcErr(err)
	}

	insurances := make([]model.CarInsurance, len(res.GetCarInsurances()))
	for i, ins := range res.GetCarInsurances() {
		insurances[i] = dto.CarInsuranceFromProto(ins)
	}

	return insurances, nil
}

func (h *CarInsuranceHandler) Update(ctx context.Context, id string, data model.CarInsuranceUpdate) error {
	logger := pkglog.WithMethod(h.log, "Update")
	logger = pkglog.WithMetadata(logger, utils.MetadataFromCtx(ctx))

	req := &carsvc.UpdateCarInsuranceRequest{
		Id:        id,
		Provider:  data.Provider,
		PolicyNum: data.PolicyNum,
		CostTenge: data.CostTenge,
		Status:    data.Status,
		Notes:     data.Notes,
		ImageKeys: data.ImageKeys,
	}
	if data.StartsAt != nil {
		req.StartsAt = timestamppb.New(*data.StartsAt)
	}
	if data.ExpiresAt != nil {
		req.ExpiresAt = timestamppb.New(*data.ExpiresAt)
	}

	_, err := h.client.UpdateCarInsurance(ctx, req)
	if err != nil {
		if dto.IsSystemErr(err) {
			logger.Error("grpc call failed", pkglog.Err(err))
		}

		return dto.FromGrpcErr(err)
	}

	return nil
}

func (h *CarInsuranceHandler) Delete(ctx context.Context, id string) error {
	logger := pkglog.WithMethod(h.log, "Delete")
	logger = pkglog.WithMetadata(logger, utils.MetadataFromCtx(ctx))

	_, err := h.client.DeleteCarInsurance(ctx, &carsvc.DeleteCarInsuranceRequest{Id: id})
	if err != nil {
		if dto.IsSystemErr(err) {
			logger.Error("grpc call failed", pkglog.Err(err))
		}

		return dto.FromGrpcErr(err)
	}

	return nil
}

func (h *CarInsuranceHandler) GetImageUploadData(ctx context.Context) (model.ImageUploadData, error) {
	logger := pkglog.WithMethod(h.log, "GetImageUploadData")
	logger = pkglog.WithMetadata(logger, utils.MetadataFromCtx(ctx))

	res, err := h.client.GetCarInsuranceImageUploadData(ctx, &emptypb.Empty{})
	if err != nil {
		if dto.IsSystemErr(err) {
			logger.Error("grpc call failed", pkglog.Err(err))
		}

		return model.ImageUploadData{}, dto.FromGrpcErr(err)
	}

	return dto.ImageUploadDataFromProto(res.GetUploadData()), nil
}
