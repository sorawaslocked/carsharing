//go:build integration

package postgres_test

import (
	"context"
	"testing"
	"time"

	"carsharing/booking-service/internal/model"
	sharedmodel "carsharing/shared/model"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// --- Create ---

func TestBookingRepo_Create_ReturnsID(t *testing.T) {
	truncate(t)
	ctx := context.Background()
	ruleID := mustInsertPricingRule(t, testPricingRuleCreate())

	id, err := newBookingRepo().Create(ctx, testBookingCreate(ruleID), time.Now().Add(15*time.Minute))

	require.NoError(t, err)
	assert.NotEmpty(t, id)
}

func TestBookingRepo_Create_SnapshotsPricingFields(t *testing.T) {
	truncate(t)
	ctx := context.Background()
	r := newBookingRepo()

	ruleData := testPricingRuleCreate()
	ruleData.RateTenge = 250
	ruleData.RatePerKMTenge = ptr(int32(10))
	ruleID := mustInsertPricingRule(t, ruleData)

	id, err := r.Create(ctx, testBookingCreate(ruleID), time.Now().Add(15*time.Minute))
	require.NoError(t, err)

	booking, err := r.GetByID(ctx, id)
	require.NoError(t, err)
	assert.Equal(t, int32(250), booking.PricingSnapshot.RateTenge)
	assert.Equal(t, ptr(int32(10)), booking.PricingSnapshot.RatePerKMTenge)
}

func TestBookingRepo_Create_CreatesStatusHistoryEntry(t *testing.T) {
	truncate(t)
	ctx := context.Background()
	ruleID := mustInsertPricingRule(t, testPricingRuleCreate())

	id, err := newBookingRepo().Create(ctx, testBookingCreate(ruleID), time.Now().Add(15*time.Minute))
	require.NoError(t, err)

	history, err := newBookingRepo().GetStatusHistory(ctx, model.BookingStatusHistoryFilter{
		BookingID:  id,
		Pagination: sharedmodel.Pagination{Limit: 10, Offset: 0},
	})
	require.NoError(t, err)
	require.Len(t, history, 1)
	assert.Equal(t, "created", history[0].ToStatus)
	assert.Equal(t, "", history[0].FromStatus)
	assert.Equal(t, sharedmodel.ActorTypeSystem, history[0].ActorType)
}

func TestBookingRepo_Create_MissingPricingRule(t *testing.T) {
	truncate(t)
	ctx := context.Background()

	_, err := newBookingRepo().Create(ctx,
		testBookingCreate("00000000-0000-0000-0000-000000000000"),
		time.Now().Add(15*time.Minute),
	)

	assert.ErrorIs(t, err, model.ErrPricingRuleNotFound)
}

func TestBookingRepo_Create_InactivePricingRule(t *testing.T) {
	truncate(t)
	ctx := context.Background()

	ruleID := mustInsertPricingRule(t, testPricingRuleCreate())
	require.NoError(t, newPricingRuleRepo().Update(ctx, ruleID, model.PricingRuleUpdate{IsActive: ptr(false)}))

	_, err := newBookingRepo().Create(ctx, testBookingCreate(ruleID), time.Now().Add(15*time.Minute))

	assert.ErrorIs(t, err, model.ErrPricingRuleNotFound)
}

// --- GetByID ---

func TestBookingRepo_GetByID_Found(t *testing.T) {
	truncate(t)
	ctx := context.Background()
	r := newBookingRepo()

	ruleID := mustInsertPricingRule(t, testPricingRuleCreate())
	data := testBookingCreate(ruleID)
	id := mustInsertBooking(t, data)

	booking, err := r.GetByID(ctx, id)

	require.NoError(t, err)
	assert.Equal(t, id, booking.ID)
	assert.Equal(t, data.UserID, booking.UserID)
	assert.Equal(t, data.CarID, booking.CarID)
	assert.Equal(t, ruleID, booking.PricingRuleID)
	assert.Equal(t, model.BookingStatusCreated, booking.Status)
}

func TestBookingRepo_GetByID_NotFound(t *testing.T) {
	truncate(t)

	_, err := newBookingRepo().GetByID(context.Background(), "00000000-0000-0000-0000-000000000000")

	assert.ErrorIs(t, err, model.ErrBookingNotFound)
}

// --- List ---

func TestBookingRepo_List_Empty(t *testing.T) {
	truncate(t)

	bookings, err := newBookingRepo().List(context.Background(), model.BookingListFilter{
		Pagination: sharedmodel.Pagination{Limit: 10, Offset: 0},
	})

	require.NoError(t, err)
	assert.Empty(t, bookings)
}

func TestBookingRepo_List_ReturnsAll(t *testing.T) {
	truncate(t)
	ctx := context.Background()
	r := newBookingRepo()

	ruleID := mustInsertPricingRule(t, testPricingRuleCreate())
	mustInsertBooking(t, testBookingCreate(ruleID))
	mustInsertBooking(t, testBookingCreate(ruleID))

	bookings, err := r.List(ctx, model.BookingListFilter{Pagination: sharedmodel.Pagination{Limit: 10, Offset: 0}})

	require.NoError(t, err)
	assert.Len(t, bookings, 2)
}

func TestBookingRepo_List_FilterByUserID(t *testing.T) {
	truncate(t)
	ctx := context.Background()
	r := newBookingRepo()

	ruleID := mustInsertPricingRule(t, testPricingRuleCreate())
	id1 := mustInsertBooking(t, testBookingCreate(ruleID))

	other := testBookingCreate(ruleID)
	other.UserID = "00000000-0000-4000-8000-000000000099"
	mustInsertBooking(t, other)

	userID := "00000000-0000-4000-8000-000000000001"
	bookings, err := r.List(ctx, model.BookingListFilter{
		UserID:     &userID,
		Pagination: sharedmodel.Pagination{Limit: 10, Offset: 0},
	})

	require.NoError(t, err)
	require.Len(t, bookings, 1)
	assert.Equal(t, id1, bookings[0].ID)
}

func TestBookingRepo_List_FilterByCarID(t *testing.T) {
	truncate(t)
	ctx := context.Background()
	r := newBookingRepo()

	ruleID := mustInsertPricingRule(t, testPricingRuleCreate())
	id1 := mustInsertBooking(t, testBookingCreate(ruleID))

	other := testBookingCreate(ruleID)
	other.CarID = "00000000-0000-4000-8000-000000000099"
	mustInsertBooking(t, other)

	carID := "00000000-0000-4000-8000-000000000002"
	bookings, err := r.List(ctx, model.BookingListFilter{
		CarID:      &carID,
		Pagination: sharedmodel.Pagination{Limit: 10, Offset: 0},
	})

	require.NoError(t, err)
	require.Len(t, bookings, 1)
	assert.Equal(t, id1, bookings[0].ID)
}

func TestBookingRepo_List_FilterByStatus(t *testing.T) {
	truncate(t)
	ctx := context.Background()
	r := newBookingRepo()

	ruleID := mustInsertPricingRule(t, testPricingRuleCreate())
	id1 := mustInsertBooking(t, testBookingCreate(ruleID))
	id2 := mustInsertBooking(t, testBookingCreate(ruleID))
	require.NoError(t, r.UpdateStatus(ctx, id2, model.BookingStatusCancelled, sharedmodel.ActorTypeUser, nil, nil))

	status := string(model.BookingStatusCreated)
	bookings, err := r.List(ctx, model.BookingListFilter{
		Status:     &status,
		Pagination: sharedmodel.Pagination{Limit: 10, Offset: 0},
	})

	require.NoError(t, err)
	require.Len(t, bookings, 1)
	assert.Equal(t, id1, bookings[0].ID)
}

func TestBookingRepo_List_FilterByPricingRuleID(t *testing.T) {
	truncate(t)
	ctx := context.Background()
	r := newBookingRepo()

	ruleID1 := mustInsertPricingRule(t, testPricingRuleCreate())
	ruleID2 := mustInsertPricingRule(t, testPricingRuleCreate())
	id1 := mustInsertBooking(t, testBookingCreate(ruleID1))
	mustInsertBooking(t, testBookingCreate(ruleID2))

	bookings, err := r.List(ctx, model.BookingListFilter{
		PricingRuleID: &ruleID1,
		Pagination:    sharedmodel.Pagination{Limit: 10, Offset: 0},
	})

	require.NoError(t, err)
	require.Len(t, bookings, 1)
	assert.Equal(t, id1, bookings[0].ID)
}

func TestBookingRepo_List_Pagination(t *testing.T) {
	truncate(t)
	ctx := context.Background()
	r := newBookingRepo()

	ruleID := mustInsertPricingRule(t, testPricingRuleCreate())
	for range 4 {
		mustInsertBooking(t, testBookingCreate(ruleID))
	}

	page1, err := r.List(ctx, model.BookingListFilter{Pagination: sharedmodel.Pagination{Limit: 2, Offset: 0}})
	require.NoError(t, err)
	assert.Len(t, page1, 2)

	page2, err := r.List(ctx, model.BookingListFilter{Pagination: sharedmodel.Pagination{Limit: 2, Offset: 2}})
	require.NoError(t, err)
	assert.Len(t, page2, 2)

	assert.NotEqual(t, page1[0].ID, page2[0].ID)
}

// --- ListCreatedExpired ---

func TestBookingRepo_ListCreatedExpired_ReturnsExpired(t *testing.T) {
	truncate(t)
	ctx := context.Background()
	r := newBookingRepo()

	ruleID := mustInsertPricingRule(t, testPricingRuleCreate())
	id, err := r.Create(ctx, testBookingCreate(ruleID), time.Now().Add(-time.Minute))
	require.NoError(t, err)

	expired, err := r.ListCreatedExpired(ctx, time.Now())

	require.NoError(t, err)
	require.Len(t, expired, 1)
	assert.Equal(t, id, expired[0].ID)
}

func TestBookingRepo_ListCreatedExpired_IgnoresNonCreatedStatus(t *testing.T) {
	truncate(t)
	ctx := context.Background()
	r := newBookingRepo()

	ruleID := mustInsertPricingRule(t, testPricingRuleCreate())
	id, err := r.Create(ctx, testBookingCreate(ruleID), time.Now().Add(-time.Minute))
	require.NoError(t, err)
	require.NoError(t, r.UpdateStatus(ctx, id, model.BookingStatusCancelled, sharedmodel.ActorTypeSystem, nil, nil))

	expired, err := r.ListCreatedExpired(ctx, time.Now())

	require.NoError(t, err)
	assert.Empty(t, expired)
}

func TestBookingRepo_ListCreatedExpired_IgnoresNotYetExpired(t *testing.T) {
	truncate(t)
	ctx := context.Background()
	r := newBookingRepo()

	ruleID := mustInsertPricingRule(t, testPricingRuleCreate())
	_, err := r.Create(ctx, testBookingCreate(ruleID), time.Now().Add(time.Hour))
	require.NoError(t, err)

	expired, err := r.ListCreatedExpired(ctx, time.Now())

	require.NoError(t, err)
	assert.Empty(t, expired)
}

// --- UpdateStatus ---

func TestBookingRepo_UpdateStatus_UpdatesStatus(t *testing.T) {
	truncate(t)
	ctx := context.Background()
	r := newBookingRepo()

	ruleID := mustInsertPricingRule(t, testPricingRuleCreate())
	id := mustInsertBooking(t, testBookingCreate(ruleID))

	err := r.UpdateStatus(ctx, id, model.BookingStatusCancelled, sharedmodel.ActorTypeUser, nil, nil)
	require.NoError(t, err)

	booking, err := r.GetByID(ctx, id)
	require.NoError(t, err)
	assert.Equal(t, model.BookingStatusCancelled, booking.Status)
}

func TestBookingRepo_UpdateStatus_RecordsHistory(t *testing.T) {
	truncate(t)
	ctx := context.Background()
	r := newBookingRepo()

	ruleID := mustInsertPricingRule(t, testPricingRuleCreate())
	id := mustInsertBooking(t, testBookingCreate(ruleID))

	err := r.UpdateStatus(ctx, id, model.BookingStatusCancelled, sharedmodel.ActorTypeUser, nil, nil)
	require.NoError(t, err)

	history, err := r.GetStatusHistory(ctx, model.BookingStatusHistoryFilter{
		BookingID:  id,
		Pagination: sharedmodel.Pagination{Limit: 10, Offset: 0},
	})
	require.NoError(t, err)
	require.Len(t, history, 2) // initial 'created' + this transition
	last := history[len(history)-1]
	assert.Equal(t, string(model.BookingStatusCreated), last.FromStatus)
	assert.Equal(t, string(model.BookingStatusCancelled), last.ToStatus)
	assert.Equal(t, sharedmodel.ActorTypeUser, last.ActorType)
}

func TestBookingRepo_UpdateStatus_WithActorIDAndReason(t *testing.T) {
	truncate(t)
	ctx := context.Background()
	r := newBookingRepo()

	ruleID := mustInsertPricingRule(t, testPricingRuleCreate())
	id := mustInsertBooking(t, testBookingCreate(ruleID))

	actorID := "00000000-0000-4000-8000-000000000001"
	reason := "user requested cancellation"
	err := r.UpdateStatus(ctx, id, model.BookingStatusCancelled, sharedmodel.ActorTypeUser, &actorID, &reason)
	require.NoError(t, err)

	history, err := r.GetStatusHistory(ctx, model.BookingStatusHistoryFilter{
		BookingID:  id,
		Pagination: sharedmodel.Pagination{Limit: 10, Offset: 0},
	})
	require.NoError(t, err)
	last := history[len(history)-1]
	require.NotNil(t, last.ActorID)
	assert.Equal(t, actorID, *last.ActorID)
	require.NotNil(t, last.Reason)
	assert.Equal(t, reason, *last.Reason)
}

func TestBookingRepo_UpdateStatus_NotFound(t *testing.T) {
	truncate(t)

	err := newBookingRepo().UpdateStatus(context.Background(),
		"00000000-0000-0000-0000-000000000000",
		model.BookingStatusCancelled,
		sharedmodel.ActorTypeUser, nil, nil,
	)

	assert.ErrorIs(t, err, model.ErrBookingNotFound)
}

// --- GetStatusHistory ---

func TestBookingRepo_GetStatusHistory_ReturnsHistory(t *testing.T) {
	truncate(t)
	ctx := context.Background()
	r := newBookingRepo()

	ruleID := mustInsertPricingRule(t, testPricingRuleCreate())
	id := mustInsertBooking(t, testBookingCreate(ruleID))
	require.NoError(t, r.UpdateStatus(ctx, id, model.BookingStatusCompleted, sharedmodel.ActorTypeSystem, nil, nil))

	history, err := r.GetStatusHistory(ctx, model.BookingStatusHistoryFilter{
		BookingID:  id,
		Pagination: sharedmodel.Pagination{Limit: 10, Offset: 0},
	})

	require.NoError(t, err)
	require.Len(t, history, 2)
	assert.Equal(t, "created", history[0].ToStatus)
	assert.Equal(t, string(model.BookingStatusCompleted), history[1].ToStatus)
}

func TestBookingRepo_GetStatusHistory_FilterByTimeRange(t *testing.T) {
	truncate(t)
	ctx := context.Background()
	r := newBookingRepo()

	ruleID := mustInsertPricingRule(t, testPricingRuleCreate())
	id := mustInsertBooking(t, testBookingCreate(ruleID))
	require.NoError(t, r.UpdateStatus(ctx, id, model.BookingStatusCancelled, sharedmodel.ActorTypeUser, nil, nil))

	future := time.Now().Add(time.Hour)
	history, err := r.GetStatusHistory(ctx, model.BookingStatusHistoryFilter{
		BookingID: id,
		TimeRange: &sharedmodel.TimeRange{
			From: future,
			To:   future.Add(time.Hour),
		},
		Pagination: sharedmodel.Pagination{Limit: 10, Offset: 0},
	})

	require.NoError(t, err)
	assert.Empty(t, history)
}

func TestBookingRepo_GetStatusHistory_Pagination(t *testing.T) {
	truncate(t)
	ctx := context.Background()
	r := newBookingRepo()

	ruleID := mustInsertPricingRule(t, testPricingRuleCreate())
	id := mustInsertBooking(t, testBookingCreate(ruleID))
	require.NoError(t, r.UpdateStatus(ctx, id, model.BookingStatusCompleted, sharedmodel.ActorTypeSystem, nil, nil))

	page1, err := r.GetStatusHistory(ctx, model.BookingStatusHistoryFilter{
		BookingID:  id,
		Pagination: sharedmodel.Pagination{Limit: 1, Offset: 0},
	})
	require.NoError(t, err)
	assert.Len(t, page1, 1)

	page2, err := r.GetStatusHistory(ctx, model.BookingStatusHistoryFilter{
		BookingID:  id,
		Pagination: sharedmodel.Pagination{Limit: 1, Offset: 1},
	})
	require.NoError(t, err)
	assert.Len(t, page2, 1)

	assert.NotEqual(t, page1[0].ID, page2[0].ID)
}
