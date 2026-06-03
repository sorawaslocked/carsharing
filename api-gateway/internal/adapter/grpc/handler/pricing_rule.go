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
	log := pkglog.WithMetadata(pkglog.WithMethod(h.log, "Create"), utils.MetadataFromCtx(ctx))

	req := &bookingsvc.CreatePricingRuleRequest{
		ModelId:           data.ModelID,
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
		log.Warn("creating pricing rule", pkglog.Err(err))

		return "", dto.FromGrpcErr(err)
	}

	return res.GetId(), nil
}

func (h *PricingRuleHandler) Get(ctx context.Context, id string) (model.PricingRule, error) {
	log := pkglog.WithMetadata(pkglog.WithMethod(h.log, "Get"), utils.MetadataFromCtx(ctx))

	res, err := h.client.GetPricingRule(ctx, &bookingsvc.GetPricingRuleRequest{Id: id})
	if err != nil {
		log.Warn("getting pricing rule", pkglog.Err(err))

		return model.PricingRule{}, dto.FromGrpcErr(err)
	}

	return dto.PricingRuleFromProto(res.GetRule()), nil
}

func (h *PricingRuleHandler) List(ctx context.Context, filter model.PricingRuleFilter) ([]model.PricingRule, error) {
	log := pkglog.WithMetadata(pkglog.WithMethod(h.log, "List"), utils.MetadataFromCtx(ctx))

	req := &bookingsvc.ListPricingRulesRequest{
		ModelId:  filter.ModelID,
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
		log.Warn("listing pricing rules", pkglog.Err(err))

		return nil, dto.FromGrpcErr(err)
	}

	rules := make([]model.PricingRule, len(res.GetRules()))
	for i, r := range res.GetRules() {
		rules[i] = dto.PricingRuleFromProto(r)
	}

	return rules, nil
}

func (h *PricingRuleHandler) Update(ctx context.Context, id string, data model.PricingRuleUpdate) error {
	log := pkglog.WithMetadata(pkglog.WithMethod(h.log, "Update"), utils.MetadataFromCtx(ctx))

	_, err := h.client.UpdatePricingRule(ctx, &bookingsvc.UpdatePricingRuleRequest{
		Id:                id,
		ModelId:           data.ModelID,
		Class:             data.Class,
		Type:              data.Type,
		RateTenge:         data.RateTenge,
		RatePerKmTenge:    data.RatePerKMTenge,
		FreeMinutes:       data.FreeMinutes,
		MinChargeTenge:    data.MinChargeTenge,
		OvertimePolicy:    data.OvertimePolicy,
		OvertimeRateTenge: data.OvertimeRateTenge,
		IsActive:          data.IsActive,
	})
	if err != nil {
		log.Warn("updating pricing rule", pkglog.Err(err))

		return dto.FromGrpcErr(err)
	}

	return nil
}

func (h *PricingRuleHandler) Delete(ctx context.Context, id string) error {
	log := pkglog.WithMetadata(pkglog.WithMethod(h.log, "Delete"), utils.MetadataFromCtx(ctx))

	_, err := h.client.DeletePricingRule(ctx, &bookingsvc.DeletePricingRuleRequest{Id: id})
	if err != nil {
		log.Warn("deleting pricing rule", pkglog.Err(err))

		return dto.FromGrpcErr(err)
	}

	return nil
}
