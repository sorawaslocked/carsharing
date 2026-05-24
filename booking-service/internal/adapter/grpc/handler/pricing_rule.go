package handler

import (
	"context"
	"log/slog"

	"carsharing/booking-service/internal/adapter/grpc/dto"
	basebookingpb "carsharing/protos/gen/base/booking"
	servicebookingpb "carsharing/protos/gen/service/booking"
	pkglog "carsharing/shared/pkg/log"
	"carsharing/shared/pkg/utils"
	"google.golang.org/protobuf/types/known/emptypb"
)

type PricingRuleHandler struct {
	servicebookingpb.UnimplementedPricingRuleServiceServer
	log *slog.Logger
	svc PricingRuleService
}

func NewPricingRuleHandler(log *slog.Logger, svc PricingRuleService) *PricingRuleHandler {
	return &PricingRuleHandler{
		log: pkglog.WithComponent(log, "grpc.PricingRuleHandler"),
		svc: svc,
	}
}

func (h *PricingRuleHandler) CreatePricingRule(ctx context.Context, req *servicebookingpb.CreatePricingRuleRequest) (*servicebookingpb.CreatePricingRuleResponse, error) {
	log := pkglog.WithMethod(h.log, "CreatePricingRule")
	log = pkglog.WithMetadata(log, utils.MetadataFromCtx(ctx))

	id, err := h.svc.Create(ctx, dto.PricingRuleCreateFromProto(req))
	if err != nil {
		return nil, dto.ToStatusError(err)
	}

	log.Info("pricing rule created", slog.String("id", id))

	return &servicebookingpb.CreatePricingRuleResponse{Id: id}, nil
}

func (h *PricingRuleHandler) GetPricingRule(ctx context.Context, req *servicebookingpb.GetPricingRuleRequest) (*servicebookingpb.GetPricingRuleResponse, error) {
	log := pkglog.WithMethod(h.log, "GetPricingRule")
	log = pkglog.WithMetadata(log, utils.MetadataFromCtx(ctx))

	rule, err := h.svc.GetByID(ctx, req.Id)
	if err != nil {
		return nil, dto.ToStatusError(err)
	}

	log.Info("pricing rule retrieved", slog.String("id", req.Id))

	return &servicebookingpb.GetPricingRuleResponse{Rule: dto.PricingRuleToProto(rule)}, nil
}

func (h *PricingRuleHandler) ListPricingRules(ctx context.Context, req *servicebookingpb.ListPricingRulesRequest) (*servicebookingpb.ListPricingRulesResponse, error) {
	log := pkglog.WithMethod(h.log, "ListPricingRules")
	log = pkglog.WithMetadata(log, utils.MetadataFromCtx(ctx))

	rules, err := h.svc.List(ctx, dto.PricingRuleListFilterFromProto(req))
	if err != nil {
		return nil, dto.ToStatusError(err)
	}

	log.Info("pricing rules listed", slog.Int("count", len(rules)))

	pbRules := make([]*basebookingpb.PricingRule, 0, len(rules))
	for _, r := range rules {
		pbRules = append(pbRules, dto.PricingRuleToProto(r))
	}

	return &servicebookingpb.ListPricingRulesResponse{Rules: pbRules}, nil
}

func (h *PricingRuleHandler) UpdatePricingRule(ctx context.Context, req *servicebookingpb.UpdatePricingRuleRequest) (*emptypb.Empty, error) {
	log := pkglog.WithMethod(h.log, "UpdatePricingRule")
	log = pkglog.WithMetadata(log, utils.MetadataFromCtx(ctx))

	if err := h.svc.Update(ctx, req.Id, dto.PricingRuleUpdateFromProto(req)); err != nil {
		return nil, dto.ToStatusError(err)
	}

	log.Info("pricing rule updated", slog.String("id", req.Id))

	return &emptypb.Empty{}, nil
}

func (h *PricingRuleHandler) DeletePricingRule(ctx context.Context, req *servicebookingpb.DeletePricingRuleRequest) (*emptypb.Empty, error) {
	log := pkglog.WithMethod(h.log, "DeletePricingRule")
	log = pkglog.WithMetadata(log, utils.MetadataFromCtx(ctx))

	if err := h.svc.Delete(ctx, req.Id); err != nil {
		return nil, dto.ToStatusError(err)
	}

	log.Info("pricing rule deleted", slog.String("id", req.Id))

	return &emptypb.Empty{}, nil
}
