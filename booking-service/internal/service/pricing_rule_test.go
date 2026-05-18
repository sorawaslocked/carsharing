package service_test

import (
	"context"
	"errors"
	"testing"

	"carsharing/booking-service/internal/model"
	"carsharing/booking-service/internal/service"
	"carsharing/booking-service/internal/service/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newRuleSvc(repo *mocks.MockPricingRuleRepository) *service.PricingRuleService {
	return service.NewPricingRuleService(discardLogger(), repo)
}

// --- Create ---

func TestPricingRuleService_Create_HappyPath(t *testing.T) {
	const ruleID = "rule-1"

	repo := mocks.NewMockPricingRuleRepository(t)
	repo.EXPECT().Create(context.Background(), model.PricingRuleCreate{}).Return(ruleID, nil)

	id, err := newRuleSvc(repo).Create(context.Background(), model.PricingRuleCreate{})

	require.NoError(t, err)
	assert.Equal(t, ruleID, id)
}

func TestPricingRuleService_Create_RepoError(t *testing.T) {
	repoErr := errors.New("db error")

	repo := mocks.NewMockPricingRuleRepository(t)
	repo.EXPECT().Create(context.Background(), model.PricingRuleCreate{}).Return("", repoErr)

	_, err := newRuleSvc(repo).Create(context.Background(), model.PricingRuleCreate{})

	assert.ErrorIs(t, err, repoErr)
}

// --- GetByID ---

func TestPricingRuleService_GetByID_HappyPath(t *testing.T) {
	want := model.PricingRule{ID: "rule-1", RateTenge: 500}

	repo := mocks.NewMockPricingRuleRepository(t)
	repo.EXPECT().GetByID(context.Background(), "rule-1").Return(want, nil)

	got, err := newRuleSvc(repo).GetByID(context.Background(), "rule-1")

	require.NoError(t, err)
	assert.Equal(t, want, got)
}

func TestPricingRuleService_GetByID_NotFound(t *testing.T) {
	repo := mocks.NewMockPricingRuleRepository(t)
	repo.EXPECT().GetByID(context.Background(), "missing").Return(model.PricingRule{}, model.ErrNotFound)

	_, err := newRuleSvc(repo).GetByID(context.Background(), "missing")

	assert.ErrorIs(t, err, model.ErrNotFound)
}

// --- List ---

func TestPricingRuleService_List_HappyPath(t *testing.T) {
	want := []model.PricingRule{{ID: "r-1"}, {ID: "r-2"}}
	filter := model.PricingRuleListFilter{}

	repo := mocks.NewMockPricingRuleRepository(t)
	repo.EXPECT().List(context.Background(), filter).Return(want, nil)

	got, err := newRuleSvc(repo).List(context.Background(), filter)

	require.NoError(t, err)
	assert.Equal(t, want, got)
}

func TestPricingRuleService_List_RepoError(t *testing.T) {
	repoErr := errors.New("query failed")

	repo := mocks.NewMockPricingRuleRepository(t)
	repo.EXPECT().List(context.Background(), model.PricingRuleListFilter{}).Return(nil, repoErr)

	_, err := newRuleSvc(repo).List(context.Background(), model.PricingRuleListFilter{})

	assert.ErrorIs(t, err, repoErr)
}

// --- Update ---

func TestPricingRuleService_Update_HappyPath(t *testing.T) {
	repo := mocks.NewMockPricingRuleRepository(t)
	repo.EXPECT().Update(context.Background(), "rule-1", model.PricingRuleUpdate{}).Return(nil)

	require.NoError(t, newRuleSvc(repo).Update(context.Background(), "rule-1", model.PricingRuleUpdate{}))
}

func TestPricingRuleService_Update_NotFound(t *testing.T) {
	repo := mocks.NewMockPricingRuleRepository(t)
	repo.EXPECT().Update(context.Background(), "missing", model.PricingRuleUpdate{}).Return(model.ErrNotFound)

	assert.ErrorIs(t, newRuleSvc(repo).Update(context.Background(), "missing", model.PricingRuleUpdate{}), model.ErrNotFound)
}

// --- Delete ---

func TestPricingRuleService_Delete_HappyPath(t *testing.T) {
	repo := mocks.NewMockPricingRuleRepository(t)
	repo.EXPECT().Delete(context.Background(), "rule-1").Return(nil)

	require.NoError(t, newRuleSvc(repo).Delete(context.Background(), "rule-1"))
}

func TestPricingRuleService_Delete_NotFound(t *testing.T) {
	repo := mocks.NewMockPricingRuleRepository(t)
	repo.EXPECT().Delete(context.Background(), "missing").Return(model.ErrNotFound)

	assert.ErrorIs(t, newRuleSvc(repo).Delete(context.Background(), "missing"), model.ErrNotFound)
}

func TestPricingRuleService_Delete_RepoError(t *testing.T) {
	repoErr := errors.New("constraint violation")

	repo := mocks.NewMockPricingRuleRepository(t)
	repo.EXPECT().Delete(context.Background(), "rule-1").Return(repoErr)

	assert.ErrorIs(t, newRuleSvc(repo).Delete(context.Background(), "rule-1"), repoErr)
}
