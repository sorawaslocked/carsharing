package handler_test

import (
	"context"
	"errors"
	"testing"

	"github.com/sorawaslocked/car-rental-booking-service/internal/adapter/grpc/handler"
	"github.com/sorawaslocked/car-rental-booking-service/internal/adapter/grpc/handler/mocks"
	"github.com/sorawaslocked/car-rental-booking-service/internal/model"
	servicebookingpb "github.com/sorawaslocked/car-rental-protos/gen/service/booking"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
)

// --- CreatePricingRule ---

func TestPricingRuleHandler_CreatePricingRule_HappyPath(t *testing.T) {
	svc := mocks.NewMockPricingRuleService(t)
	svc.EXPECT().Create(context.Background(), model.PricingRuleCreate{
		Type:      "per_minute",
		RateTenge: 50,
	}).Return("r-1", nil)

	h := handler.NewPricingRuleHandler(discardLogger(), svc)
	resp, err := h.CreatePricingRule(context.Background(), &servicebookingpb.CreatePricingRuleRequest{
		Type:      "per_minute",
		RateTenge: 50,
	})

	require.NoError(t, err)
	assert.Equal(t, "r-1", resp.Id)
}

func TestPricingRuleHandler_CreatePricingRule_ServiceError(t *testing.T) {
	svc := mocks.NewMockPricingRuleService(t)
	svc.EXPECT().Create(context.Background(), model.PricingRuleCreate{}).
		Return("", errors.New("db error"))

	h := handler.NewPricingRuleHandler(discardLogger(), svc)
	_, err := h.CreatePricingRule(context.Background(), &servicebookingpb.CreatePricingRuleRequest{})

	assertCode(t, err, codes.Internal)
}

// --- GetPricingRule ---

func TestPricingRuleHandler_GetPricingRule_HappyPath(t *testing.T) {
	rule := model.PricingRule{ID: "r-1", Type: "per_minute", RateTenge: 50, IsActive: true}

	svc := mocks.NewMockPricingRuleService(t)
	svc.EXPECT().GetByID(context.Background(), "r-1").Return(rule, nil)

	h := handler.NewPricingRuleHandler(discardLogger(), svc)
	resp, err := h.GetPricingRule(context.Background(), &servicebookingpb.GetPricingRuleRequest{Id: "r-1"})

	require.NoError(t, err)
	assert.Equal(t, "r-1", resp.Rule.Id)
	assert.Equal(t, "per_minute", resp.Rule.Type)
	assert.Equal(t, int32(50), resp.Rule.RateTenge)
	assert.True(t, resp.Rule.IsActive)
}

func TestPricingRuleHandler_GetPricingRule_NotFound(t *testing.T) {
	svc := mocks.NewMockPricingRuleService(t)
	svc.EXPECT().GetByID(context.Background(), "missing").Return(model.PricingRule{}, model.ErrNotFound)

	h := handler.NewPricingRuleHandler(discardLogger(), svc)
	_, err := h.GetPricingRule(context.Background(), &servicebookingpb.GetPricingRuleRequest{Id: "missing"})

	assertCode(t, err, codes.NotFound)
}

// --- ListPricingRules ---

func TestPricingRuleHandler_ListPricingRules_HappyPath(t *testing.T) {
	rules := []model.PricingRule{{ID: "r-1"}, {ID: "r-2"}}

	svc := mocks.NewMockPricingRuleService(t)
	svc.EXPECT().List(context.Background(), model.PricingRuleListFilter{}).Return(rules, nil)

	h := handler.NewPricingRuleHandler(discardLogger(), svc)
	resp, err := h.ListPricingRules(context.Background(), &servicebookingpb.ListPricingRulesRequest{})

	require.NoError(t, err)
	assert.Len(t, resp.Rules, 2)
	assert.Equal(t, "r-1", resp.Rules[0].Id)
	assert.Equal(t, "r-2", resp.Rules[1].Id)
}

func TestPricingRuleHandler_ListPricingRules_ServiceError(t *testing.T) {
	svc := mocks.NewMockPricingRuleService(t)
	svc.EXPECT().List(context.Background(), model.PricingRuleListFilter{}).
		Return(nil, errors.New("db error"))

	h := handler.NewPricingRuleHandler(discardLogger(), svc)
	_, err := h.ListPricingRules(context.Background(), &servicebookingpb.ListPricingRulesRequest{})

	assertCode(t, err, codes.Internal)
}

// --- UpdatePricingRule ---

func TestPricingRuleHandler_UpdatePricingRule_HappyPath(t *testing.T) {
	svc := mocks.NewMockPricingRuleService(t)
	svc.EXPECT().Update(context.Background(), "r-1", model.PricingRuleUpdate{}).Return(nil)

	h := handler.NewPricingRuleHandler(discardLogger(), svc)
	resp, err := h.UpdatePricingRule(context.Background(), &servicebookingpb.UpdatePricingRuleRequest{Id: "r-1"})

	require.NoError(t, err)
	assert.NotNil(t, resp)
}

func TestPricingRuleHandler_UpdatePricingRule_NotFound(t *testing.T) {
	svc := mocks.NewMockPricingRuleService(t)
	svc.EXPECT().Update(context.Background(), "missing", model.PricingRuleUpdate{}).Return(model.ErrNotFound)

	h := handler.NewPricingRuleHandler(discardLogger(), svc)
	_, err := h.UpdatePricingRule(context.Background(), &servicebookingpb.UpdatePricingRuleRequest{Id: "missing"})

	assertCode(t, err, codes.NotFound)
}

// --- DeletePricingRule ---

func TestPricingRuleHandler_DeletePricingRule_HappyPath(t *testing.T) {
	svc := mocks.NewMockPricingRuleService(t)
	svc.EXPECT().Delete(context.Background(), "r-1").Return(nil)

	h := handler.NewPricingRuleHandler(discardLogger(), svc)
	resp, err := h.DeletePricingRule(context.Background(), &servicebookingpb.DeletePricingRuleRequest{Id: "r-1"})

	require.NoError(t, err)
	assert.NotNil(t, resp)
}

func TestPricingRuleHandler_DeletePricingRule_NotFound(t *testing.T) {
	svc := mocks.NewMockPricingRuleService(t)
	svc.EXPECT().Delete(context.Background(), "missing").Return(model.ErrNotFound)

	h := handler.NewPricingRuleHandler(discardLogger(), svc)
	_, err := h.DeletePricingRule(context.Background(), &servicebookingpb.DeletePricingRuleRequest{Id: "missing"})

	assertCode(t, err, codes.NotFound)
}

func TestPricingRuleHandler_DeletePricingRule_ServiceError(t *testing.T) {
	svc := mocks.NewMockPricingRuleService(t)
	svc.EXPECT().Delete(context.Background(), "r-1").Return(errors.New("constraint violation"))

	h := handler.NewPricingRuleHandler(discardLogger(), svc)
	_, err := h.DeletePricingRule(context.Background(), &servicebookingpb.DeletePricingRuleRequest{Id: "r-1"})

	assertCode(t, err, codes.Internal)
}
