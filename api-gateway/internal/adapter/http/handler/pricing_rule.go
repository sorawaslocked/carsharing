package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/sorawaslocked/car-rental-api-gateway/internal/adapter/http/dto"
)

type PricingRuleHandler struct {
	svc PricingRuleService
}

func NewPricingRuleHandler(svc PricingRuleService) *PricingRuleHandler {
	return &PricingRuleHandler{svc: svc}
}

// Create (PricingRule) godoc
// @Summary      Create a pricing rule
// @Description  Defines a new pricing rule. Nil scope fields (modelID, zoneID, class) act as wildcards and apply to all values of that dimension.
// @Tags         pricing-rules
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body  body      dto.PricingRuleCreateRequest  true  "Pricing rule create payload"
// @Success      201   {object}  dto.IDResponse
// @Failure      400   {object}  dto.ErrorResponse
// @Failure      401   {object}  dto.ErrorResponse
// @Failure      409   {object}  dto.ErrorResponse
// @Failure      500   {object}  dto.ErrorResponse
// @Router       /pricing-rules [post]
func (h *PricingRuleHandler) Create(ctx *gin.Context) {
	data, err := dto.FromPricingRuleCreateRequest(ctx)
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

// Get (PricingRule) godoc
// @Summary      Get pricing rule by ID
// @Tags         pricing-rules
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string  true  "Pricing rule UUID"
// @Success      200  {object}  dto.PricingRuleResponse
// @Failure      400  {object}  dto.ErrorResponse
// @Failure      401  {object}  dto.ErrorResponse
// @Failure      404  {object}  dto.ErrorResponse
// @Failure      500  {object}  dto.ErrorResponse
// @Router       /pricing-rules/{id} [get]
func (h *PricingRuleHandler) Get(ctx *gin.Context) {
	id, err := dto.IDParam(ctx)
	if err != nil {
		dto.FromError(ctx, err)

		return
	}

	rule, err := h.svc.Get(ctx, id)
	if err != nil {
		dto.FromError(ctx, err)

		return
	}

	dto.Ok(ctx, gin.H{"rule": dto.ToPricingRuleResponse(rule)})
}

// List (PricingRule) godoc
// @Summary      List pricing rules
// @Description  Returns pricing rules filtered by optional scope dimensions.
// @Tags         pricing-rules
// @Produce      json
// @Security     BearerAuth
// @Param        modelID   query     string   false  "Filter by car model UUID"
// @Param        zoneID    query     string   false  "Filter by zone UUID"
// @Param        class     query     string   false  "Filter by car class"
// @Param        type      query     string   false  "Filter by pricing type"
// @Param        isActive  query     boolean  false  "Filter by active status"
// @Param        limit     query     integer  false  "Pagination limit"
// @Param        offset    query     integer  false  "Pagination offset"
// @Success      200       {object}  dto.PricingRulesResponse
// @Failure      400       {object}  dto.ErrorResponse
// @Failure      401       {object}  dto.ErrorResponse
// @Failure      500       {object}  dto.ErrorResponse
// @Router       /pricing-rules [get]
func (h *PricingRuleHandler) List(ctx *gin.Context) {
	filter, err := dto.PricingRuleFilterFromCtx(ctx)
	if err != nil {
		dto.FromError(ctx, err)

		return
	}

	rules, err := h.svc.List(ctx, filter)
	if err != nil {
		dto.FromError(ctx, err)

		return
	}

	response := make([]dto.PricingRule, len(rules))
	for i, rule := range rules {
		response[i] = dto.ToPricingRuleResponse(rule)
	}

	dto.Ok(ctx, gin.H{"rules": response})
}

// Update (PricingRule) godoc
// @Summary      Update pricing rule
// @Description  Partially updates a pricing rule. Only provided fields are changed.
// @Tags         pricing-rules
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id    path      string                       true  "Pricing rule UUID"
// @Param        body  body      dto.PricingRuleUpdateRequest  true  "Fields to update"
// @Success      204
// @Failure      400  {object}  dto.ErrorResponse
// @Failure      401  {object}  dto.ErrorResponse
// @Failure      404  {object}  dto.ErrorResponse
// @Failure      500  {object}  dto.ErrorResponse
// @Router       /pricing-rules/{id} [patch]
func (h *PricingRuleHandler) Update(ctx *gin.Context) {
	id, err := dto.IDParam(ctx)
	if err != nil {
		dto.FromError(ctx, err)

		return
	}

	data, err := dto.FromPricingRuleUpdateRequest(ctx)
	if err != nil {
		dto.MalformedJson(ctx)

		return
	}

	if err = h.svc.Update(ctx, id, data); err != nil {
		dto.FromError(ctx, err)

		return
	}

	dto.NoContent(ctx)
}

// Delete (PricingRule) godoc
// @Summary      Delete pricing rule
// @Tags         pricing-rules
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string  true  "Pricing rule UUID"
// @Success      204
// @Failure      400  {object}  dto.ErrorResponse
// @Failure      401  {object}  dto.ErrorResponse
// @Failure      404  {object}  dto.ErrorResponse
// @Failure      500  {object}  dto.ErrorResponse
// @Router       /pricing-rules/{id} [delete]
func (h *PricingRuleHandler) Delete(ctx *gin.Context) {
	id, err := dto.IDParam(ctx)
	if err != nil {
		dto.FromError(ctx, err)

		return
	}

	if err = h.svc.Delete(ctx, id); err != nil {
		dto.FromError(ctx, err)

		return
	}

	dto.NoContent(ctx)
}
