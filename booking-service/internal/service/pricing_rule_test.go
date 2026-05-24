package service_test

import (
	"context"
	"errors"
	"testing"

	"carsharing/booking-service/internal/model"
	"carsharing/booking-service/internal/service"
	"carsharing/booking-service/internal/service/mocks"
	"carsharing/booking-service/internal/validation"
	sharedmodel "carsharing/shared/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newRuleSvc(repo *mocks.MockPricingRuleRepository) *service.PricingRuleService {
	return service.NewPricingRuleService(discardLogger(), newValidator(), repo, nil, nil)
}

var validPricingRuleCreate = validation.PricingRuleCreate{Type: "by_minute", RateTenge: 100}
var mappedPricingRuleCreate = model.PricingRuleCreate{Type: "by_minute", RateTenge: 100}
var defaultPricingRuleListFilter = model.PricingRuleListFilter{
	Pagination: sharedmodel.Pagination{Limit: sharedmodel.DefaultPaginationLimit, Offset: sharedmodel.DefaultPaginationOffset},
}

// --- Create ---

func TestPricingRuleService_Create_HappyPath(t *testing.T) {
	repo := mocks.NewMockPricingRuleRepository(t)
	repo.EXPECT().Create(context.Background(), mappedPricingRuleCreate).Return(testPricingRuleID, nil)

	id, err := newRuleSvc(repo).Create(context.Background(), validPricingRuleCreate)

	require.NoError(t, err)
	assert.Equal(t, testPricingRuleID, id)
}

func TestPricingRuleService_Create_RepoError(t *testing.T) {
	repoErr := errors.New("db error")

	repo := mocks.NewMockPricingRuleRepository(t)
	repo.EXPECT().Create(context.Background(), mappedPricingRuleCreate).Return("", repoErr)

	_, err := newRuleSvc(repo).Create(context.Background(), validPricingRuleCreate)

	assert.ErrorIs(t, err, repoErr)
}

// --- GetByID ---

func TestPricingRuleService_GetByID_HappyPath(t *testing.T) {
	want := model.PricingRule{ID: testPricingRuleID, RateTenge: 500}

	repo := mocks.NewMockPricingRuleRepository(t)
	repo.EXPECT().GetByID(context.Background(), testPricingRuleID).Return(want, nil)

	got, err := newRuleSvc(repo).GetByID(context.Background(), testPricingRuleID)

	require.NoError(t, err)
	assert.Equal(t, want, got)
}

func TestPricingRuleService_GetByID_NotFound(t *testing.T) {
	repo := mocks.NewMockPricingRuleRepository(t)
	repo.EXPECT().GetByID(context.Background(), testMissingID).Return(model.PricingRule{}, model.ErrPricingRuleNotFound)

	_, err := newRuleSvc(repo).GetByID(context.Background(), testMissingID)

	assert.ErrorIs(t, err, model.ErrPricingRuleNotFound)
}

// --- List ---

func TestPricingRuleService_List_HappyPath(t *testing.T) {
	want := []model.PricingRule{{ID: testPricingRuleID}, {ID: testBookingID}}

	repo := mocks.NewMockPricingRuleRepository(t)
	repo.EXPECT().List(context.Background(), defaultPricingRuleListFilter).Return(want, nil)

	got, err := newRuleSvc(repo).List(context.Background(), validation.PricingRuleListFilter{})

	require.NoError(t, err)
	assert.Equal(t, want, got)
}

func TestPricingRuleService_List_RepoError(t *testing.T) {
	repoErr := errors.New("query failed")

	repo := mocks.NewMockPricingRuleRepository(t)
	repo.EXPECT().List(context.Background(), defaultPricingRuleListFilter).Return(nil, repoErr)

	_, err := newRuleSvc(repo).List(context.Background(), validation.PricingRuleListFilter{})

	assert.ErrorIs(t, err, repoErr)
}

// --- Update ---

func TestPricingRuleService_Update_HappyPath(t *testing.T) {
	repo := mocks.NewMockPricingRuleRepository(t)
	repo.EXPECT().Update(context.Background(), testPricingRuleID, model.PricingRuleUpdate{}).Return(nil)

	require.NoError(t, newRuleSvc(repo).Update(context.Background(), testPricingRuleID, validation.PricingRuleUpdate{}))
}

func TestPricingRuleService_Update_NotFound(t *testing.T) {
	repo := mocks.NewMockPricingRuleRepository(t)
	repo.EXPECT().Update(context.Background(), testMissingID, model.PricingRuleUpdate{}).Return(model.ErrPricingRuleNotFound)

	assert.ErrorIs(t, newRuleSvc(repo).Update(context.Background(), testMissingID, validation.PricingRuleUpdate{}), model.ErrPricingRuleNotFound)
}

// --- Delete ---

func TestPricingRuleService_Delete_HappyPath(t *testing.T) {
	repo := mocks.NewMockPricingRuleRepository(t)
	repo.EXPECT().Delete(context.Background(), testPricingRuleID).Return(nil)

	require.NoError(t, newRuleSvc(repo).Delete(context.Background(), testPricingRuleID))
}

func TestPricingRuleService_Delete_NotFound(t *testing.T) {
	repo := mocks.NewMockPricingRuleRepository(t)
	repo.EXPECT().Delete(context.Background(), testMissingID).Return(model.ErrPricingRuleNotFound)

	assert.ErrorIs(t, newRuleSvc(repo).Delete(context.Background(), testMissingID), model.ErrPricingRuleNotFound)
}

func TestPricingRuleService_Delete_RepoError(t *testing.T) {
	repoErr := errors.New("constraint violation")

	repo := mocks.NewMockPricingRuleRepository(t)
	repo.EXPECT().Delete(context.Background(), testPricingRuleID).Return(repoErr)

	assert.ErrorIs(t, newRuleSvc(repo).Delete(context.Background(), testPricingRuleID), repoErr)
}
