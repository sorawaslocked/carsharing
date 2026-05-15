package service_test

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"testing"

	"github.com/sorawaslocked/car-rental-booking-service/internal/model"
	"github.com/sorawaslocked/car-rental-booking-service/internal/pkg/utils"
	"github.com/sorawaslocked/car-rental-booking-service/internal/service"
	"github.com/sorawaslocked/car-rental-booking-service/internal/service/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func discardLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}

func ptr[T any](v T) *T { return &v }

func ctxAsBookingManager() context.Context {
	return utils.SetMetadata(context.Background(), "", "", "manager-1", "booking_manager")
}

func ctxAsUser(userID string) context.Context {
	return utils.SetMetadata(context.Background(), "", "", userID, "user")
}

func newBookingSvc(repo *mocks.MockBookingRepository, rules *mocks.MockPricingRuleRepository, pub *mocks.MockEventPublisher) *service.BookingService {
	return service.NewBookingService(discardLogger(), repo, rules, pub)
}

// --- Create ---

func TestBookingService_Create_HappyPath(t *testing.T) {
	const bookingID = "booking-1"

	repo := mocks.NewMockBookingRepository(t)
	rules := mocks.NewMockPricingRuleRepository(t)
	pub := mocks.NewMockEventPublisher(t)

	rules.EXPECT().GetByID(mock.Anything, "rule-1").Return(model.PricingRule{}, nil)
	repo.EXPECT().Create(mock.Anything, model.BookingCreate{PricingRuleID: "rule-1"}, mock.AnythingOfType("time.Time")).Return(bookingID, nil)
	repo.EXPECT().GetByID(mock.Anything, bookingID).Return(model.Booking{ID: bookingID}, nil)
	pub.EXPECT().PublishBookingCreated(mock.Anything, model.Booking{ID: bookingID}).Return(nil)

	id, err := newBookingSvc(repo, rules, pub).Create(ctxAsBookingManager(), model.BookingCreate{PricingRuleID: "rule-1"})

	require.NoError(t, err)
	assert.Equal(t, bookingID, id)
}

func TestBookingService_Create_RuleNotFound(t *testing.T) {
	repo := mocks.NewMockBookingRepository(t)
	rules := mocks.NewMockPricingRuleRepository(t)
	pub := mocks.NewMockEventPublisher(t)

	rules.EXPECT().GetByID(mock.Anything, "missing").Return(model.PricingRule{}, model.ErrNotFound)

	_, err := newBookingSvc(repo, rules, pub).Create(ctxAsBookingManager(), model.BookingCreate{PricingRuleID: "missing"})

	assert.ErrorIs(t, err, model.ErrNotFound)
}

func TestBookingService_Create_RepoError(t *testing.T) {
	repoErr := errors.New("db down")

	repo := mocks.NewMockBookingRepository(t)
	rules := mocks.NewMockPricingRuleRepository(t)
	pub := mocks.NewMockEventPublisher(t)

	rules.EXPECT().GetByID(mock.Anything, "rule-1").Return(model.PricingRule{}, nil)
	repo.EXPECT().Create(mock.Anything, model.BookingCreate{PricingRuleID: "rule-1"}, mock.AnythingOfType("time.Time")).Return("", repoErr)

	_, err := newBookingSvc(repo, rules, pub).Create(ctxAsBookingManager(), model.BookingCreate{PricingRuleID: "rule-1"})

	assert.ErrorIs(t, err, repoErr)
}

// If GetByID after create fails the service still returns the ID without error.
func TestBookingService_Create_FetchAfterCreateFails_StillReturnsID(t *testing.T) {
	const bookingID = "booking-2"

	repo := mocks.NewMockBookingRepository(t)
	rules := mocks.NewMockPricingRuleRepository(t)
	pub := mocks.NewMockEventPublisher(t)

	rules.EXPECT().GetByID(mock.Anything, "rule-1").Return(model.PricingRule{}, nil)
	repo.EXPECT().Create(mock.Anything, model.BookingCreate{PricingRuleID: "rule-1"}, mock.AnythingOfType("time.Time")).Return(bookingID, nil)
	repo.EXPECT().GetByID(mock.Anything, bookingID).Return(model.Booking{}, errors.New("transient error"))

	id, err := newBookingSvc(repo, rules, pub).Create(ctxAsBookingManager(), model.BookingCreate{PricingRuleID: "rule-1"})

	require.NoError(t, err)
	assert.Equal(t, bookingID, id)
}

// Publish failure must not surface as an error to the caller.
func TestBookingService_Create_PublishFailure_NoError(t *testing.T) {
	const bookingID = "booking-3"

	repo := mocks.NewMockBookingRepository(t)
	rules := mocks.NewMockPricingRuleRepository(t)
	pub := mocks.NewMockEventPublisher(t)

	rules.EXPECT().GetByID(mock.Anything, "rule-1").Return(model.PricingRule{}, nil)
	repo.EXPECT().Create(mock.Anything, model.BookingCreate{PricingRuleID: "rule-1"}, mock.AnythingOfType("time.Time")).Return(bookingID, nil)
	repo.EXPECT().GetByID(mock.Anything, bookingID).Return(model.Booking{ID: bookingID}, nil)
	pub.EXPECT().PublishBookingCreated(mock.Anything, model.Booking{ID: bookingID}).Return(errors.New("nats unavailable"))

	id, err := newBookingSvc(repo, rules, pub).Create(ctxAsBookingManager(), model.BookingCreate{PricingRuleID: "rule-1"})

	require.NoError(t, err)
	assert.Equal(t, bookingID, id)
}

// --- GetByID ---

func TestBookingService_GetByID_HappyPath(t *testing.T) {
	want := model.Booking{ID: "b-1", UserID: "u-1", Status: model.BookingStatusCreated}

	repo := mocks.NewMockBookingRepository(t)
	repo.EXPECT().GetByID(mock.Anything, "b-1").Return(want, nil)

	got, err := newBookingSvc(repo, mocks.NewMockPricingRuleRepository(t), mocks.NewMockEventPublisher(t)).
		GetByID(ctxAsUser("u-1"), "b-1")

	require.NoError(t, err)
	assert.Equal(t, want, got)
}

func TestBookingService_GetByID_NotFound(t *testing.T) {
	repo := mocks.NewMockBookingRepository(t)
	repo.EXPECT().GetByID(mock.Anything, "missing").Return(model.Booking{}, model.ErrNotFound)

	_, err := newBookingSvc(repo, mocks.NewMockPricingRuleRepository(t), mocks.NewMockEventPublisher(t)).
		GetByID(context.Background(), "missing")

	assert.ErrorIs(t, err, model.ErrNotFound)
}

// --- List ---

func TestBookingService_List_HappyPath(t *testing.T) {
	want := []model.Booking{{ID: "b-1"}, {ID: "b-2"}}
	filter := model.BookingListFilter{}

	repo := mocks.NewMockBookingRepository(t)
	repo.EXPECT().List(mock.Anything, filter).Return(want, nil)

	got, err := newBookingSvc(repo, mocks.NewMockPricingRuleRepository(t), mocks.NewMockEventPublisher(t)).
		List(context.Background(), filter)

	require.NoError(t, err)
	assert.Equal(t, want, got)
}

func TestBookingService_List_RepoError(t *testing.T) {
	repoErr := errors.New("query failed")

	repo := mocks.NewMockBookingRepository(t)
	repo.EXPECT().List(mock.Anything, model.BookingListFilter{}).Return(nil, repoErr)

	_, err := newBookingSvc(repo, mocks.NewMockPricingRuleRepository(t), mocks.NewMockEventPublisher(t)).
		List(context.Background(), model.BookingListFilter{})

	assert.ErrorIs(t, err, repoErr)
}

// --- Cancel ---

func TestBookingService_Cancel_HappyPath(t *testing.T) {
	booking := model.Booking{ID: "b-1", UserID: "u-1", Status: model.BookingStatusCreated}
	reason := "no longer needed"

	repo := mocks.NewMockBookingRepository(t)
	pub := mocks.NewMockEventPublisher(t)

	repo.EXPECT().GetByID(mock.Anything, "b-1").Return(booking, nil)
	repo.EXPECT().UpdateStatus(mock.Anything, "b-1", "cancelled", "user", ptr("u-1"), &reason).Return(nil)
	pub.EXPECT().PublishBookingCancelled(mock.Anything, booking, reason).Return(nil)

	err := newBookingSvc(repo, mocks.NewMockPricingRuleRepository(t), pub).
		Cancel(ctxAsUser("u-1"), "b-1", &reason)

	require.NoError(t, err)
}

func TestBookingService_Cancel_NotFound(t *testing.T) {
	repo := mocks.NewMockBookingRepository(t)
	repo.EXPECT().GetByID(mock.Anything, "missing").Return(model.Booking{}, model.ErrNotFound)

	err := newBookingSvc(repo, mocks.NewMockPricingRuleRepository(t), mocks.NewMockEventPublisher(t)).
		Cancel(context.Background(), "missing", nil)

	assert.ErrorIs(t, err, model.ErrNotFound)
}

func TestBookingService_Cancel_InvalidTransition(t *testing.T) {
	terminalStatuses := []model.BookingStatus{
		model.BookingStatusExpired,
		model.BookingStatusCompleted,
		model.BookingStatusCancelled,
	}

	for _, status := range terminalStatuses {
		t.Run(string(status), func(t *testing.T) {
			repo := mocks.NewMockBookingRepository(t)
			repo.EXPECT().GetByID(mock.Anything, "b-1").Return(model.Booking{UserID: "u-1", Status: status}, nil)

			err := newBookingSvc(repo, mocks.NewMockPricingRuleRepository(t), mocks.NewMockEventPublisher(t)).
				Cancel(ctxAsUser("u-1"), "b-1", nil)

			assert.ErrorIs(t, err, model.ErrInvalidTransition)
		})
	}
}

func TestBookingService_Cancel_UpdateStatusError(t *testing.T) {
	repoErr := errors.New("update failed")
	booking := model.Booking{ID: "b-1", UserID: "u-1", Status: model.BookingStatusCreated}

	repo := mocks.NewMockBookingRepository(t)
	repo.EXPECT().GetByID(mock.Anything, "b-1").Return(booking, nil)
	repo.EXPECT().UpdateStatus(mock.Anything, "b-1", "cancelled", "user", ptr("u-1"), (*string)(nil)).Return(repoErr)

	err := newBookingSvc(repo, mocks.NewMockPricingRuleRepository(t), mocks.NewMockEventPublisher(t)).
		Cancel(ctxAsUser("u-1"), "b-1", nil)

	assert.ErrorIs(t, err, repoErr)
}

// Publish failure on cancel must not surface as an error.
func TestBookingService_Cancel_PublishFailure_NoError(t *testing.T) {
	booking := model.Booking{ID: "b-1", UserID: "u-1", Status: model.BookingStatusCreated}

	repo := mocks.NewMockBookingRepository(t)
	pub := mocks.NewMockEventPublisher(t)

	repo.EXPECT().GetByID(mock.Anything, "b-1").Return(booking, nil)
	repo.EXPECT().UpdateStatus(mock.Anything, "b-1", "cancelled", "user", ptr("u-1"), (*string)(nil)).Return(nil)
	pub.EXPECT().PublishBookingCancelled(mock.Anything, booking, "").Return(errors.New("nats unavailable"))

	err := newBookingSvc(repo, mocks.NewMockPricingRuleRepository(t), pub).
		Cancel(ctxAsUser("u-1"), "b-1", nil)

	require.NoError(t, err)
}

// --- Complete ---

func TestBookingService_Complete_HappyPath(t *testing.T) {
	booking := model.Booking{ID: "b-1", Status: model.BookingStatusCreated}

	repo := mocks.NewMockBookingRepository(t)
	pub := mocks.NewMockEventPublisher(t)

	repo.EXPECT().GetByID(mock.Anything, "b-1").Return(booking, nil)
	repo.EXPECT().UpdateStatus(mock.Anything, "b-1", "completed", "system", (*string)(nil), (*string)(nil)).Return(nil)
	pub.EXPECT().PublishBookingCompleted(mock.Anything, booking).Return(nil)

	require.NoError(t, newBookingSvc(repo, mocks.NewMockPricingRuleRepository(t), pub).Complete(context.Background(), "b-1"))
}

func TestBookingService_Complete_NotFound(t *testing.T) {
	repo := mocks.NewMockBookingRepository(t)
	repo.EXPECT().GetByID(mock.Anything, "missing").Return(model.Booking{}, model.ErrNotFound)

	err := newBookingSvc(repo, mocks.NewMockPricingRuleRepository(t), mocks.NewMockEventPublisher(t)).
		Complete(context.Background(), "missing")

	assert.ErrorIs(t, err, model.ErrNotFound)
}

func TestBookingService_Complete_InvalidTransition(t *testing.T) {
	repo := mocks.NewMockBookingRepository(t)
	repo.EXPECT().GetByID(mock.Anything, "b-1").Return(model.Booking{Status: model.BookingStatusCompleted}, nil)

	err := newBookingSvc(repo, mocks.NewMockPricingRuleRepository(t), mocks.NewMockEventPublisher(t)).
		Complete(context.Background(), "b-1")

	assert.ErrorIs(t, err, model.ErrInvalidTransition)
}

// --- UpdateStatus ---

func TestBookingService_UpdateStatus_HappyPath(t *testing.T) {
	booking := model.Booking{ID: "b-1", Status: model.BookingStatusCreated}

	repo := mocks.NewMockBookingRepository(t)
	repo.EXPECT().GetByID(mock.Anything, "b-1").Return(booking, nil)
	repo.EXPECT().UpdateStatus(mock.Anything, "b-1", "cancelled", "system", (*string)(nil), (*string)(nil)).Return(nil)

	err := newBookingSvc(repo, mocks.NewMockPricingRuleRepository(t), mocks.NewMockEventPublisher(t)).
		UpdateStatus(context.Background(), "b-1", "cancelled", nil)

	require.NoError(t, err)
}

func TestBookingService_UpdateStatus_InvalidStatusString(t *testing.T) {
	err := newBookingSvc(
		mocks.NewMockBookingRepository(t),
		mocks.NewMockPricingRuleRepository(t),
		mocks.NewMockEventPublisher(t),
	).UpdateStatus(context.Background(), "b-1", "INVALID", nil)

	assert.ErrorIs(t, err, model.ErrInvalidStatus)
}

func TestBookingService_UpdateStatus_NotFound(t *testing.T) {
	repo := mocks.NewMockBookingRepository(t)
	repo.EXPECT().GetByID(mock.Anything, "missing").Return(model.Booking{}, model.ErrNotFound)

	err := newBookingSvc(repo, mocks.NewMockPricingRuleRepository(t), mocks.NewMockEventPublisher(t)).
		UpdateStatus(context.Background(), "missing", "cancelled", nil)

	assert.ErrorIs(t, err, model.ErrNotFound)
}

func TestBookingService_UpdateStatus_InvalidTransition(t *testing.T) {
	repo := mocks.NewMockBookingRepository(t)
	repo.EXPECT().GetByID(mock.Anything, "b-1").Return(model.Booking{Status: model.BookingStatusExpired}, nil)

	err := newBookingSvc(repo, mocks.NewMockPricingRuleRepository(t), mocks.NewMockEventPublisher(t)).
		UpdateStatus(context.Background(), "b-1", "cancelled", nil)

	assert.ErrorIs(t, err, model.ErrInvalidTransition)
}

// --- GetStatusHistory ---

func TestBookingService_GetStatusHistory_HappyPath(t *testing.T) {
	want := []model.BookingStatusReading{{ID: "h-1"}, {ID: "h-2"}}
	filter := model.BookingStatusHistoryFilter{BookingID: "b-1"}

	repo := mocks.NewMockBookingRepository(t)
	repo.EXPECT().GetStatusHistory(mock.Anything, filter).Return(want, nil)

	got, err := newBookingSvc(repo, mocks.NewMockPricingRuleRepository(t), mocks.NewMockEventPublisher(t)).
		GetStatusHistory(context.Background(), filter)

	require.NoError(t, err)
	assert.Equal(t, want, got)
}

func TestBookingService_GetStatusHistory_RepoError(t *testing.T) {
	repoErr := errors.New("query failed")
	filter := model.BookingStatusHistoryFilter{}

	repo := mocks.NewMockBookingRepository(t)
	repo.EXPECT().GetStatusHistory(mock.Anything, filter).Return(nil, repoErr)

	_, err := newBookingSvc(repo, mocks.NewMockPricingRuleRepository(t), mocks.NewMockEventPublisher(t)).
		GetStatusHistory(context.Background(), filter)

	assert.ErrorIs(t, err, repoErr)
}
