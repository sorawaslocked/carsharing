package handler

import (
	"context"
	"log/slog"

	"carsharing/api-gateway/internal/adapter/grpc/dto"
	"carsharing/api-gateway/internal/model"
	basepb "carsharing/protos/gen/base"
	bookingsvc "carsharing/protos/gen/service/booking"
	pkglog "carsharing/shared/pkg/log"
	"carsharing/shared/pkg/utils"
)

type PricingRuleHandler struct {
	client bookingsvc.PricingRuleServiceClient
	log    *slog.Logger
}

func NewPricingRuleHandler(client bookingsvc.PricingRuleServiceClient, logger *slog.Logger) *PricingRuleHandler {
	return &PricingRuleHandler{
		client: client,
		log:    pkglog.WithComponent(logger, "grpc.PricingRuleHandler"),
	}
}

func (h *PricingRuleHandler) Create(ctx context.Context, data model.PricingRuleCreate) (string, error) {
	logger := pkglog.WithMethod(h.log, "Create")
	logger = pkglog.WithMetadata(logger, utils.MetadataFromCtx(ctx))

	req := &bookingsvc.CreatePricingRuleRequest{
		ModelId:           data.ModelID,
		ZoneId:            data.ZoneID,
		Class:             data.Class,
		Type:              data.Type,
		RateTenge:         data.RateTenge,
		RatePerKmTenge:    data.RatePerKMTenge,
		FreeMinutes:       data.FreeMinutes,
		MinChargeTenge:    data.MinChargeTenge,
		OvertimePolicy:    data.OvertimePolicy,
		OvertimeRateTenge: data.OvertimeRateTenge,
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

func (h *PricingRuleHandler) Get(ctx context.Context, id string) (model.PricingRule, error) {
	logger := pkglog.WithMethod(h.log, "Get")
	logger = pkglog.WithMetadata(logger, utils.MetadataFromCtx(ctx))

	res, err := h.client.GetPricingRule(ctx, &bookingsvc.GetPricingRuleRequest{Id: id})
	if err != nil {
		if dto.IsSystemErr(err) {
			logger.Error("grpc call failed", pkglog.Err(err))
		}

		return model.PricingRule{}, dto.FromGrpcErr(err)
	}

	return dto.PricingRuleFromProto(res.GetRule()), nil
}

func (h *PricingRuleHandler) List(ctx context.Context, filter model.PricingRuleFilter) ([]model.PricingRule, error) {
	logger := pkglog.WithMethod(h.log, "List")
	logger = pkglog.WithMetadata(logger, utils.MetadataFromCtx(ctx))

	req := &bookingsvc.ListPricingRulesRequest{
		ModelId:  filter.ModelID,
		ZoneId:   filter.ZoneID,
		Class:    filter.Class,
		Type:     filter.Type,
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

	rules := make([]model.PricingRule, len(res.GetRules()))
	for i, r := range res.GetRules() {
		rules[i] = dto.PricingRuleFromProto(r)
	}

	return rules, nil
}

func (h *PricingRuleHandler) Update(ctx context.Context, id string, data model.PricingRuleUpdate) error {
	logger := pkglog.WithMethod(h.log, "Update")
	logger = pkglog.WithMetadata(logger, utils.MetadataFromCtx(ctx))

	req := &bookingsvc.UpdatePricingRuleRequest{
		Id:                id,
		ModelId:           data.ModelID,
		ZoneId:            data.ZoneID,
		Class:             data.Class,
		Type:              data.Type,
		RateTenge:         data.RateTenge,
		RatePerKmTenge:    data.RatePerKMTenge,
		FreeMinutes:       data.FreeMinutes,
		MinChargeTenge:    data.MinChargeTenge,
		OvertimePolicy:    data.OvertimePolicy,
		OvertimeRateTenge: data.OvertimeRateTenge,
		IsActive:          data.IsActive,
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

func (h *PricingRuleHandler) Delete(ctx context.Context, id string) error {
	logger := pkglog.WithMethod(h.log, "Delete")
	logger = pkglog.WithMetadata(logger, utils.MetadataFromCtx(ctx))

	_, err := h.client.DeletePricingRule(ctx, &bookingsvc.DeletePricingRuleRequest{Id: id})
	if err != nil {
		if dto.IsSystemErr(err) {
			logger.Error("grpc call failed", pkglog.Err(err))
		}

		return dto.FromGrpcErr(err)
	}

	return nil
}
