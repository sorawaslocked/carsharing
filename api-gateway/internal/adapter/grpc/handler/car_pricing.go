package handler

import (
	"context"
	"log/slog"

	"github.com/sorawaslocked/car-rental-api-gateway/internal/adapter/grpc/dto"
	"github.com/sorawaslocked/car-rental-api-gateway/internal/model"
	pkglog "github.com/sorawaslocked/car-rental-api-gateway/internal/pkg/log"
	basepb "github.com/sorawaslocked/car-rental-protos/gen/base"
	carsvc "github.com/sorawaslocked/car-rental-protos/gen/service/car"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type CarPricingRuleHandler struct {
	client carsvc.CarPricingServiceClient
	log    *slog.Logger
}

func NewCarPricingRuleHandler(client carsvc.CarPricingServiceClient, logger *slog.Logger) *CarPricingRuleHandler {
	return &CarPricingRuleHandler{
		client: client,
		log:    pkglog.WithComponent(logger, "grpc.CarPricingRuleHandler"),
	}
}

func (h *CarPricingRuleHandler) Create(ctx context.Context, data model.CarPricingRuleCreate) (string, error) {
	logger := pkglog.WithMethod(h.log, "Create")

	req := &carsvc.CreatePricingRuleRequest{
		ModelId:            data.ModelID,
		ZoneId:             data.ZoneID,
		Class:              data.Class,
		RatePerMinuteTenge: data.RatePerMinuteTenge,
		RatePerKmTenge:     data.RatePerKMTenge,
		FreeMinutes:        data.FreeMinutes,
		MinChargeTenge:     data.MinChargeTenge,
	}
	if data.StartsAt != nil {
		req.StartsAt = timestamppb.New(*data.StartsAt)
	}
	if data.ExpiresAt != nil {
		req.ExpiresAt = timestamppb.New(*data.ExpiresAt)
	}

	res, err := h.client.CreatePricingRule(ctx, req)
	if err != nil {
		if dto.IsSystemErr(err) {
			logger.Error("grpc call failed", pkglog.Err(err))
		}

		return "", dto.FromGrpcErr(err)
	}

	return res.GetId(), nil
}

func (h *CarPricingRuleHandler) Get(ctx context.Context, id string) (model.CarPricingRule, error) {
	logger := pkglog.WithMethod(h.log, "Get")

	res, err := h.client.GetPricingRule(ctx, &carsvc.GetPricingRuleRequest{Id: id})
	if err != nil {
		if dto.IsSystemErr(err) {
			logger.Error("grpc call failed", pkglog.Err(err))
		}

		return model.CarPricingRule{}, dto.FromGrpcErr(err)
	}

	return dto.CarPricingRuleFromProto(res.GetRule()), nil
}

func (h *CarPricingRuleHandler) List(ctx context.Context, filter model.CarPricingRuleFilter) ([]model.CarPricingRule, error) {
	logger := pkglog.WithMethod(h.log, "List")

	req := &carsvc.ListPricingRulesRequest{
		ModelId:  filter.ModelID,
		ZoneId:   filter.ZoneID,
		Class:    filter.Class,
		IsActive: filter.IsActive,
	}
	if filter.Pagination != nil {
		req.Pagination = &basepb.Pagination{
			Limit:  filter.Pagination.Limit,
			Offset: filter.Pagination.Offset,
		}
	}

	res, err := h.client.ListPricingRules(ctx, req)
	if err != nil {
		if dto.IsSystemErr(err) {
			logger.Error("grpc call failed", pkglog.Err(err))
		}

		return nil, dto.FromGrpcErr(err)
	}

	rules := make([]model.CarPricingRule, len(res.GetRules()))
	for i, r := range res.GetRules() {
		rules[i] = dto.CarPricingRuleFromProto(r)
	}

	return rules, nil
}

func (h *CarPricingRuleHandler) Update(ctx context.Context, id string, data model.CarPricingRuleUpdate) error {
	logger := pkglog.WithMethod(h.log, "Update")

	req := &carsvc.UpdatePricingRuleRequest{
		Id:                 id,
		ModelId:            data.ModelID,
		ZoneId:             data.ZoneID,
		Class:              data.Class,
		RatePerMinuteTenge: data.RatePerMinuteTenge,
		RatePerKmTenge:     data.RatePerKMTenge,
		FreeMinutes:        data.FreeMinutes,
		MinChargeTenge:     data.MinChargeTenge,
		IsActive:           data.IsActive,
	}
	if data.StartsAt != nil {
		req.StartsAt = timestamppb.New(*data.StartsAt)
	}
	if data.ExpiresAt != nil {
		req.ExpiresAt = timestamppb.New(*data.ExpiresAt)
	}

	_, err := h.client.UpdatePricingRule(ctx, req)
	if err != nil {
		if dto.IsSystemErr(err) {
			logger.Error("grpc call failed", pkglog.Err(err))
		}

		return dto.FromGrpcErr(err)
	}

	return nil
}

func (h *CarPricingRuleHandler) Delete(ctx context.Context, id string) error {
	logger := pkglog.WithMethod(h.log, "Delete")

	_, err := h.client.DeletePricingRule(ctx, &carsvc.DeletePricingRuleRequest{Id: id})
	if err != nil {
		if dto.IsSystemErr(err) {
			logger.Error("grpc call failed", pkglog.Err(err))
		}

		return dto.FromGrpcErr(err)
	}

	return nil
}
