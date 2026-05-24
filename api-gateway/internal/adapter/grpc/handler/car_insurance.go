package handler

import (
	"context"
	"log/slog"

	"carsharing/api-gateway/internal/adapter/grpc/dto"
	"carsharing/api-gateway/internal/model"
	basepb "carsharing/protos/gen/base"
	carsvc "carsharing/protos/gen/service/car"
	sharedmodel "carsharing/shared/model"
	pkglog "carsharing/shared/pkg/log"
	"carsharing/shared/pkg/utils"
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
	log := pkglog.WithMetadata(pkglog.WithMethod(h.log, "Create"), utils.MetadataFromCtx(ctx))

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
		log.Warn("creating car insurance", pkglog.Err(err))

		return "", dto.FromGrpcErr(err)
	}

	return res.GetId(), nil
}

func (h *CarInsuranceHandler) Get(ctx context.Context, id string) (model.CarInsurance, error) {
	log := pkglog.WithMetadata(pkglog.WithMethod(h.log, "Get"), utils.MetadataFromCtx(ctx))

	res, err := h.client.GetCarInsurance(ctx, &carsvc.GetCarInsuranceRequest{Id: id})
	if err != nil {
		log.Warn("getting car insurance", pkglog.Err(err))

		return model.CarInsurance{}, dto.FromGrpcErr(err)
	}

	return dto.CarInsuranceFromProto(res.GetCarInsurance()), nil
}

func (h *CarInsuranceHandler) List(ctx context.Context, filter model.CarInsuranceFilter) ([]model.CarInsurance, error) {
	log := pkglog.WithMetadata(pkglog.WithMethod(h.log, "List"), utils.MetadataFromCtx(ctx))

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
		log.Warn("listing car insurances", pkglog.Err(err))

		return nil, dto.FromGrpcErr(err)
	}

	insurances := make([]model.CarInsurance, len(res.GetCarInsurances()))
	for i, ins := range res.GetCarInsurances() {
		insurances[i] = dto.CarInsuranceFromProto(ins)
	}

	return insurances, nil
}

func (h *CarInsuranceHandler) Update(ctx context.Context, id string, data model.CarInsuranceUpdate) error {
	log := pkglog.WithMetadata(pkglog.WithMethod(h.log, "Update"), utils.MetadataFromCtx(ctx))

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
		log.Warn("updating car insurance", pkglog.Err(err))

		return dto.FromGrpcErr(err)
	}

	return nil
}

func (h *CarInsuranceHandler) Delete(ctx context.Context, id string) error {
	log := pkglog.WithMetadata(pkglog.WithMethod(h.log, "Delete"), utils.MetadataFromCtx(ctx))

	_, err := h.client.DeleteCarInsurance(ctx, &carsvc.DeleteCarInsuranceRequest{Id: id})
	if err != nil {
		log.Warn("deleting car insurance", pkglog.Err(err))

		return dto.FromGrpcErr(err)
	}

	return nil
}

func (h *CarInsuranceHandler) GetImageUploadData(ctx context.Context) (sharedmodel.ImageUploadData, error) {
	log := pkglog.WithMetadata(pkglog.WithMethod(h.log, "GetImageUploadData"), utils.MetadataFromCtx(ctx))

	res, err := h.client.GetCarInsuranceImageUploadData(ctx, &emptypb.Empty{})
	if err != nil {
		log.Warn("getting car insurance image upload data", pkglog.Err(err))

		return sharedmodel.ImageUploadData{}, dto.FromGrpcErr(err)
	}

	return dto.ImageUploadDataFromProto(res.GetUploadData()), nil
}
