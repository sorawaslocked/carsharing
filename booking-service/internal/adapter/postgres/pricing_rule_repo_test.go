//go:build integration

package postgres_test

import (
	"context"
	"testing"

	"carsharing/booking-service/internal/model"
	sharedmodel "carsharing/shared/model"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// --- Create ---

func TestPricingRuleRepo_Create_ReturnsID(t *testing.T) {
	truncate(t)
	ctx := context.Background()

	id, err := newPricingRuleRepo().Create(ctx, testPricingRuleCreate())

	require.NoError(t, err)
	assert.NotEmpty(t, id)
}

func TestPricingRuleRepo_Create_DefaultIsActive(t *testing.T) {
	truncate(t)
	ctx := context.Background()
	r := newPricingRuleRepo()

	id, err := r.Create(ctx, testPricingRuleCreate())
	require.NoError(t, err)

	rule, err := r.GetByID(ctx, id)
	require.NoError(t, err)
	assert.True(t, rule.IsActive)
}

func TestPricingRuleRepo_Create_WithOptionalFields(t *testing.T) {
	truncate(t)
	ctx := context.Background()
	r := newPricingRuleRepo()

	data := testPricingRuleCreate()
	data.RatePerKMTenge = ptr(int32(50))
	data.FreeMinutes = ptr(int32(10))
	data.MinChargeTenge = ptr(int32(200))
	data.OvertimePolicy = ptr("charge_double")
	data.OvertimeRateTenge = ptr(int32(300))

	id, err := r.Create(ctx, data)
	require.NoError(t, err)

	rule, err := r.GetByID(ctx, id)
	require.NoError(t, err)
	assert.Equal(t, data.RatePerKMTenge, rule.RatePerKMTenge)
	assert.Equal(t, data.FreeMinutes, rule.FreeMinutes)
	assert.Equal(t, data.MinChargeTenge, rule.MinChargeTenge)
	assert.Equal(t, data.OvertimePolicy, rule.OvertimePolicy)
	assert.Equal(t, data.OvertimeRateTenge, rule.OvertimeRateTenge)
}

// --- GetByID ---

func TestPricingRuleRepo_GetByID_Found(t *testing.T) {
	truncate(t)
	ctx := context.Background()
	r := newPricingRuleRepo()

	data := testPricingRuleCreate()
	id := mustInsertPricingRule(t, data)

	rule, err := r.GetByID(ctx, id)

	require.NoError(t, err)
	assert.Equal(t, id, rule.ID)
	assert.Equal(t, data.Type, rule.Type)
	assert.Equal(t, data.RateTenge, rule.RateTenge)
	assert.True(t, rule.IsActive)
}

func TestPricingRuleRepo_GetByID_NotFound(t *testing.T) {
	truncate(t)

	_, err := newPricingRuleRepo().GetByID(context.Background(), "00000000-0000-0000-0000-000000000000")

	assert.ErrorIs(t, err, model.ErrPricingRuleNotFound)
}

// --- List ---

func TestPricingRuleRepo_List_Empty(t *testing.T) {
	truncate(t)

	rules, err := newPricingRuleRepo().List(context.Background(), model.PricingRuleListFilter{
		Pagination: sharedmodel.Pagination{Limit: 10, Offset: 0},
	})

	require.NoError(t, err)
	assert.Empty(t, rules)
}

func TestPricingRuleRepo_List_ReturnsAll(t *testing.T) {
	truncate(t)
	ctx := context.Background()
	r := newPricingRuleRepo()

	mustInsertPricingRule(t, testPricingRuleCreate())
	mustInsertPricingRule(t, testPricingRuleCreate())

	rules, err := r.List(ctx, model.PricingRuleListFilter{
		Pagination: sharedmodel.Pagination{Limit: 10, Offset: 0},
	})

	require.NoError(t, err)
	assert.Len(t, rules, 2)
}

func TestPricingRuleRepo_List_FilterByIsActive(t *testing.T) {
	truncate(t)
	ctx := context.Background()
	r := newPricingRuleRepo()

	activeID := mustInsertPricingRule(t, testPricingRuleCreate())
	inactiveID := mustInsertPricingRule(t, testPricingRuleCreate())
	require.NoError(t, r.Update(ctx, inactiveID, model.PricingRuleUpdate{IsActive: ptr(false)}))

	active, err := r.List(ctx, model.PricingRuleListFilter{
		IsActive:   ptr(true),
		Pagination: sharedmodel.Pagination{Limit: 10, Offset: 0},
	})
	require.NoError(t, err)
	require.Len(t, active, 1)
	assert.Equal(t, activeID, active[0].ID)
}

func TestPricingRuleRepo_List_Pagination(t *testing.T) {
	truncate(t)
	ctx := context.Background()
	r := newPricingRuleRepo()

	for range 4 {
		mustInsertPricingRule(t, testPricingRuleCreate())
	}

	page1, err := r.List(ctx, model.PricingRuleListFilter{Pagination: sharedmodel.Pagination{Limit: 2, Offset: 0}})
	require.NoError(t, err)
	assert.Len(t, page1, 2)

	page2, err := r.List(ctx, model.PricingRuleListFilter{Pagination: sharedmodel.Pagination{Limit: 2, Offset: 2}})
	require.NoError(t, err)
	assert.Len(t, page2, 2)

	assert.NotEqual(t, page1[0].ID, page2[0].ID)
}

// --- Update ---

func TestPricingRuleRepo_Update_RateTenge(t *testing.T) {
	truncate(t)
	ctx := context.Background()
	r := newPricingRuleRepo()

	id := mustInsertPricingRule(t, testPricingRuleCreate())

	err := r.Update(ctx, id, model.PricingRuleUpdate{RateTenge: ptr(int32(999))})
	require.NoError(t, err)

	rule, err := r.GetByID(ctx, id)
	require.NoError(t, err)
	assert.Equal(t, int32(999), rule.RateTenge)
}

func TestPricingRuleRepo_Update_IsActive(t *testing.T) {
	truncate(t)
	ctx := context.Background()
	r := newPricingRuleRepo()

	id := mustInsertPricingRule(t, testPricingRuleCreate())

	err := r.Update(ctx, id, model.PricingRuleUpdate{IsActive: ptr(false)})
	require.NoError(t, err)

	rule, err := r.GetByID(ctx, id)
	require.NoError(t, err)
	assert.False(t, rule.IsActive)
}

func TestPricingRuleRepo_Update_NotFound(t *testing.T) {
	truncate(t)

	err := newPricingRuleRepo().Update(context.Background(), "00000000-0000-0000-0000-000000000000", model.PricingRuleUpdate{
		RateTenge: ptr(int32(200)),
	})

	assert.ErrorIs(t, err, model.ErrPricingRuleNotFound)
}

// --- Delete ---

func TestPricingRuleRepo_Delete_Removes(t *testing.T) {
	truncate(t)
	ctx := context.Background()
	r := newPricingRuleRepo()

	id := mustInsertPricingRule(t, testPricingRuleCreate())

	err := r.Delete(ctx, id)
	require.NoError(t, err)

	_, err = r.GetByID(ctx, id)
	assert.ErrorIs(t, err, model.ErrPricingRuleNotFound)
}

func TestPricingRuleRepo_Delete_NotFound(t *testing.T) {
	truncate(t)

	err := newPricingRuleRepo().Delete(context.Background(), "00000000-0000-0000-0000-000000000000")

	assert.ErrorIs(t, err, model.ErrPricingRuleNotFound)
}
