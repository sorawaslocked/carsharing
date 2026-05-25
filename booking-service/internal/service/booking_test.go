package service_test

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"testing"

	"carsharing/booking-service/internal/model"
	"carsharing/booking-service/internal/service"
	"carsharing/booking-service/internal/service/mocks"
	"carsharing/booking-service/internal/validation"
	sharedmodel "carsharing/shared/model"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// Test UUIDs — all valid UUID4 format.
const (
	testUserID        = "00000000-0000-4000-8000-000000000001"
	testManagerID     = "00000000-0000-4000-8000-000000000002"
	testCarID         = "00000000-0000-4000-8000-000000000003"
	testPricingRuleID = "00000000-0000-4000-8000-000000000004"
	testBookingID     = "00000000-0000-4000-8000-000000000005"
	testBookingID2    = "00000000-0000-4000-8000-000000000006"
	testBookingID3    = "00000000-0000-4000-8000-000000000007"
	testMissingID     = "00000000-0000-4000-8000-000000000099"
)

func discardLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}

func ptr[T any](v T) *T { return &v }

func ctxAsBookingManager() context.Context {
	ctx := context.WithValue(context.Background(), "x-user-id", testManagerID)
	return context.WithValue(ctx, "x-user-roles", []sharedmodel.Role{sharedmodel.RoleBookingManager})
}

func ctxAsUser(userID string) context.Context {
	ctx := context.WithValue(context.Background(), "x-user-id", userID)
	return context.WithValue(ctx, "x-user-roles", []sharedmodel.Role{sharedmodel.RoleUser})
}

func newValidator() *validator.Validate {
	v := validator.New()
	validation.RegisterCustomValidators(v, discardLogger())
	return v
}

func newBookingSvc(
	repo *mocks.MockBookingRepository,
	rules *mocks.MockPricingRuleRepository,
	pub *mocks.MockEventPublisher,
	cars ...*mocks.MockCarChecker,
) *service.BookingService {
	var carChecker service.CarChecker
	if len(cars) > 0 {
		carChecker = cars[0]
	}
	return service.NewBookingService(discardLogger(), newValidator(), repo, rules, pub, carChecker)
}

// --- Create ---

func TestBookingService_Create_HappyPath(t *testing.T) {
	repo := mocks.NewMockBookingRepository(t)
	rules := mocks.NewMockPricingRuleRepository(t)
	pub := mocks.NewMockEventPublisher(t)
	cars := mocks.NewMockCarChecker(t)

	data := validation.BookingCreate{
		UserID:        testManagerID,
		CarID:         testCarID,
		PricingRuleID: testPricingRuleID,
	}
	booking := model.Booking{ID: testBookingID}

	cars.EXPECT().Exists(mock.Anything, testCarID).Return(true, nil)
	cars.EXPECT().GetStatus(mock.Anything, testCarID).Return(model.CarStatusAvailable, nil)
	rules.EXPECT().GetByID(mock.Anything, testPricingRuleID).Return(model.PricingRule{}, nil)
	repo.EXPECT().Create(mock.Anything, model.BookingCreate{
		UserID:        testManagerID,
		CarID:         testCarID,
		PricingRuleID: testPricingRuleID,
	}, mock.AnythingOfType("time.Time")).Return(testBookingID, nil)
	repo.EXPECT().GetByID(mock.Anything, testBookingID).Return(booking, nil)
	pub.EXPECT().PublishBookingCreated(mock.Anything, booking).Return(nil)

	id, err := newBookingSvc(repo, rules, pub, cars).Create(ctxAsBookingManager(), data)

	require.NoError(t, err)
	assert.Equal(t, testBookingID, id)
}

func TestBookingService_Create_RuleNotFound(t *testing.T) {
	repo := mocks.NewMockBookingRepository(t)
	rules := mocks.NewMockPricingRuleRepository(t)
	pub := mocks.NewMockEventPublisher(t)
	cars := mocks.NewMockCarChecker(t)

	cars.EXPECT().Exists(mock.Anything, testCarID).Return(true, nil)
	cars.EXPECT().GetStatus(mock.Anything, testCarID).Return(model.CarStatusAvailable, nil)
	rules.EXPECT().GetByID(mock.Anything, testMissingID).Return(model.PricingRule{}, model.ErrNotFound)

	_, err := newBookingSvc(repo, rules, pub, cars).Create(ctxAsBookingManager(), validation.BookingCreate{
		UserID:        testManagerID,
		CarID:         testCarID,
		PricingRuleID: testMissingID,
	})

	assert.ErrorIs(t, err, model.ErrPricingRuleNotFound)
}

func TestBookingService_Create_RepoError(t *testing.T) {
	repoErr := errors.New("db down")

	repo := mocks.NewMockBookingRepository(t)
	rules := mocks.NewMockPricingRuleRepository(t)
	pub := mocks.NewMockEventPublisher(t)
	cars := mocks.NewMockCarChecker(t)

	cars.EXPECT().Exists(mock.Anything, testCarID).Return(true, nil)
	cars.EXPECT().GetStatus(mock.Anything, testCarID).Return(model.CarStatusAvailable, nil)
	rules.EXPECT().GetByID(mock.Anything, testPricingRuleID).Return(model.PricingRule{}, nil)
	repo.EXPECT().Create(mock.Anything, mock.AnythingOfType("model.BookingCreate"), mock.AnythingOfType("time.Time")).Return("", repoErr)

	_, err := newBookingSvc(repo, rules, pub, cars).Create(ctxAsBookingManager(), validation.BookingCreate{
		UserID:        testManagerID,
		CarID:         testCarID,
		PricingRuleID: testPricingRuleID,
	})

	assert.ErrorIs(t, err, repoErr)
}

// If GetByID after create fails the service still returns the ID without error.
func TestBookingService_Create_FetchAfterCreateFails_StillReturnsID(t *testing.T) {
	repo := mocks.NewMockBookingRepository(t)
	rules := mocks.NewMockPricingRuleRepository(t)
	pub := mocks.NewMockEventPublisher(t)
	cars := mocks.NewMockCarChecker(t)

	cars.EXPECT().Exists(mock.Anything, testCarID).Return(true, nil)
	cars.EXPECT().GetStatus(mock.Anything, testCarID).Return(model.CarStatusAvailable, nil)
	rules.EXPECT().GetByID(mock.Anything, testPricingRuleID).Return(model.PricingRule{}, nil)
	repo.EXPECT().Create(mock.Anything, mock.AnythingOfType("model.BookingCreate"), mock.AnythingOfType("time.Time")).Return(testBookingID2, nil)
	repo.EXPECT().GetByID(mock.Anything, testBookingID2).Return(model.Booking{}, errors.New("transient error"))

	id, err := newBookingSvc(repo, rules, pub, cars).Create(ctxAsBookingManager(), validation.BookingCreate{
		UserID:        testManagerID,
		CarID:         testCarID,
		PricingRuleID: testPricingRuleID,
	})

	require.NoError(t, err)
	assert.Equal(t, testBookingID2, id)
}

// Publish failure must not surface as an error to the caller.
func TestBookingService_Create_PublishFailure_NoError(t *testing.T) {
	repo := mocks.NewMockBookingRepository(t)
	rules := mocks.NewMockPricingRuleRepository(t)
	pub := mocks.NewMockEventPublisher(t)
	cars := mocks.NewMockCarChecker(t)

	booking := model.Booking{ID: testBookingID3}

	cars.EXPECT().Exists(mock.Anything, testCarID).Return(true, nil)
	cars.EXPECT().GetStatus(mock.Anything, testCarID).Return(model.CarStatusAvailable, nil)
	rules.EXPECT().GetByID(mock.Anything, testPricingRuleID).Return(model.PricingRule{}, nil)
	repo.EXPECT().Create(mock.Anything, mock.AnythingOfType("model.BookingCreate"), mock.AnythingOfType("time.Time")).Return(testBookingID3, nil)
	repo.EXPECT().GetByID(mock.Anything, testBookingID3).Return(booking, nil)
	pub.EXPECT().PublishBookingCreated(mock.Anything, booking).Return(errors.New("nats unavailable"))

	id, err := newBookingSvc(repo, rules, pub, cars).Create(ctxAsBookingManager(), validation.BookingCreate{
		UserID:        testManagerID,
		CarID:         testCarID,
		PricingRuleID: testPricingRuleID,
	})

	require.NoError(t, err)
	assert.Equal(t, testBookingID3, id)
}

// --- GetByID ---

func TestBookingService_GetByID_HappyPath(t *testing.T) {
	want := model.Booking{ID: testBookingID, UserID: testUserID, Status: model.BookingStatusCreated}

	repo := mocks.NewMockBookingRepository(t)
	repo.EXPECT().GetByID(mock.Anything, testBookingID).Return(want, nil)

	got, err := newBookingSvc(repo, mocks.NewMockPricingRuleRepository(t), mocks.NewMockEventPublisher(t)).
		GetByID(ctxAsUser(testUserID), testBookingID)

	require.NoError(t, err)
	assert.Equal(t, want, got)
}

func TestBookingService_GetByID_NotFound(t *testing.T) {
	repo := mocks.NewMockBookingRepository(t)
	repo.EXPECT().GetByID(mock.Anything, testMissingID).Return(model.Booking{}, model.ErrBookingNotFound)

	_, err := newBookingSvc(repo, mocks.NewMockPricingRuleRepository(t), mocks.NewMockEventPublisher(t)).
		GetByID(ctxAsBookingManager(), testMissingID)

	assert.ErrorIs(t, err, model.ErrBookingNotFound)
}

// --- List ---

func TestBookingService_List_HappyPath(t *testing.T) {
	want := []model.Booking{{ID: testBookingID}, {ID: testBookingID2}}

	repo := mocks.NewMockBookingRepository(t)
	repo.EXPECT().List(mock.Anything, mock.AnythingOfType("model.BookingListFilter")).Return(want, nil)

	got, err := newBookingSvc(repo, mocks.NewMockPricingRuleRepository(t), mocks.NewMockEventPublisher(t)).
		List(ctxAsBookingManager(), validation.BookingListFilter{})

	require.NoError(t, err)
	assert.Equal(t, want, got)
}

func TestBookingService_List_RepoError(t *testing.T) {
	repoErr := errors.New("query failed")

	repo := mocks.NewMockBookingRepository(t)
	repo.EXPECT().List(mock.Anything, mock.AnythingOfType("model.BookingListFilter")).Return(nil, repoErr)

	_, err := newBookingSvc(repo, mocks.NewMockPricingRuleRepository(t), mocks.NewMockEventPublisher(t)).
		List(ctxAsBookingManager(), validation.BookingListFilter{})

	assert.ErrorIs(t, err, repoErr)
}

// --- Cancel ---

func TestBookingService_Cancel_HappyPath(t *testing.T) {
	booking := model.Booking{ID: testBookingID, UserID: testUserID, Status: model.BookingStatusCreated}
	reason := "no longer needed"

	repo := mocks.NewMockBookingRepository(t)
	pub := mocks.NewMockEventPublisher(t)

	repo.EXPECT().GetByID(mock.Anything, testBookingID).Return(booking, nil)
	repo.EXPECT().UpdateStatus(mock.Anything, testBookingID, model.BookingStatusCancelled, sharedmodel.ActorTypeUser, ptr(testUserID), &reason).Return(nil)
	pub.EXPECT().PublishBookingCancelled(mock.Anything, booking, reason).Return(nil)

	err := newBookingSvc(repo, mocks.NewMockPricingRuleRepository(t), pub).
		Cancel(ctxAsUser(testUserID), testBookingID, &reason)

	require.NoError(t, err)
}

func TestBookingService_Cancel_NotFound(t *testing.T) {
	repo := mocks.NewMockBookingRepository(t)
	repo.EXPECT().GetByID(mock.Anything, testMissingID).Return(model.Booking{}, model.ErrBookingNotFound)

	err := newBookingSvc(repo, mocks.NewMockPricingRuleRepository(t), mocks.NewMockEventPublisher(t)).
		Cancel(ctxAsBookingManager(), testMissingID, nil)

	assert.ErrorIs(t, err, model.ErrBookingNotFound)
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
			repo.EXPECT().GetByID(mock.Anything, testBookingID).Return(model.Booking{UserID: testUserID, Status: status}, nil)

			err := newBookingSvc(repo, mocks.NewMockPricingRuleRepository(t), mocks.NewMockEventPublisher(t)).
				Cancel(ctxAsUser(testUserID), testBookingID, nil)

			assert.ErrorIs(t, err, model.ErrInvalidBookingStatusTransition)
		})
	}
}

func TestBookingService_Cancel_UpdateStatusError(t *testing.T) {
	repoErr := errors.New("update failed")
	booking := model.Booking{ID: testBookingID, UserID: testUserID, Status: model.BookingStatusCreated}

	repo := mocks.NewMockBookingRepository(t)
	repo.EXPECT().GetByID(mock.Anything, testBookingID).Return(booking, nil)
	repo.EXPECT().UpdateStatus(mock.Anything, testBookingID, model.BookingStatusCancelled, sharedmodel.ActorTypeUser, ptr(testUserID), (*string)(nil)).Return(repoErr)

	err := newBookingSvc(repo, mocks.NewMockPricingRuleRepository(t), mocks.NewMockEventPublisher(t)).
		Cancel(ctxAsUser(testUserID), testBookingID, nil)

	assert.ErrorIs(t, err, repoErr)
}

// Publish failure on cancel must not surface as an error.
func TestBookingService_Cancel_PublishFailure_NoError(t *testing.T) {
	booking := model.Booking{ID: testBookingID, UserID: testUserID, Status: model.BookingStatusCreated}

	repo := mocks.NewMockBookingRepository(t)
	pub := mocks.NewMockEventPublisher(t)

	repo.EXPECT().GetByID(mock.Anything, testBookingID).Return(booking, nil)
	repo.EXPECT().UpdateStatus(mock.Anything, testBookingID, model.BookingStatusCancelled, sharedmodel.ActorTypeUser, ptr(testUserID), (*string)(nil)).Return(nil)
	pub.EXPECT().PublishBookingCancelled(mock.Anything, booking, "").Return(errors.New("nats unavailable"))

	err := newBookingSvc(repo, mocks.NewMockPricingRuleRepository(t), pub).
		Cancel(ctxAsUser(testUserID), testBookingID, nil)

	require.NoError(t, err)
}

// --- Complete ---

func TestBookingService_Complete_HappyPath(t *testing.T) {
	booking := model.Booking{ID: testBookingID, Status: model.BookingStatusCreated}

	repo := mocks.NewMockBookingRepository(t)
	pub := mocks.NewMockEventPublisher(t)

	repo.EXPECT().GetByID(mock.Anything, testBookingID).Return(booking, nil)
	repo.EXPECT().UpdateStatus(mock.Anything, testBookingID, model.BookingStatusCompleted, sharedmodel.ActorTypeSystem, (*string)(nil), (*string)(nil)).Return(nil)
	pub.EXPECT().PublishBookingCompleted(mock.Anything, booking).Return(nil)

	require.NoError(t, newBookingSvc(repo, mocks.NewMockPricingRuleRepository(t), pub).Complete(context.Background(), testBookingID))
}

func TestBookingService_Complete_NotFound(t *testing.T) {
	repo := mocks.NewMockBookingRepository(t)
	repo.EXPECT().GetByID(mock.Anything, testMissingID).Return(model.Booking{}, model.ErrBookingNotFound)

	err := newBookingSvc(repo, mocks.NewMockPricingRuleRepository(t), mocks.NewMockEventPublisher(t)).
		Complete(context.Background(), testMissingID)

	assert.ErrorIs(t, err, model.ErrBookingNotFound)
}

func TestBookingService_Complete_InvalidTransition(t *testing.T) {
	repo := mocks.NewMockBookingRepository(t)
	repo.EXPECT().GetByID(mock.Anything, testBookingID).Return(model.Booking{Status: model.BookingStatusCompleted}, nil)

	err := newBookingSvc(repo, mocks.NewMockPricingRuleRepository(t), mocks.NewMockEventPublisher(t)).
		Complete(context.Background(), testBookingID)

	assert.ErrorIs(t, err, model.ErrInvalidBookingStatusTransition)
}

// --- UpdateStatus ---

func TestBookingService_UpdateStatus_HappyPath(t *testing.T) {
	booking := model.Booking{ID: testBookingID, Status: model.BookingStatusCreated}

	repo := mocks.NewMockBookingRepository(t)
	repo.EXPECT().GetByID(mock.Anything, testBookingID).Return(booking, nil)
	repo.EXPECT().UpdateStatus(mock.Anything, testBookingID, model.BookingStatusCancelled, sharedmodel.ActorTypeUser, ptr(testManagerID), (*string)(nil)).Return(nil)

	err := newBookingSvc(repo, mocks.NewMockPricingRuleRepository(t), mocks.NewMockEventPublisher(t)).
		UpdateStatus(ctxAsBookingManager(), testBookingID, validation.BookingStatusUpdate{Status: "cancelled"})

	require.NoError(t, err)
}

func TestBookingService_UpdateStatus_InvalidStatusString(t *testing.T) {
	err := newBookingSvc(
		mocks.NewMockBookingRepository(t),
		mocks.NewMockPricingRuleRepository(t),
		mocks.NewMockEventPublisher(t),
	).UpdateStatus(ctxAsBookingManager(), testBookingID, validation.BookingStatusUpdate{Status: "INVALID"})

	var ve validation.Errors
	assert.ErrorAs(t, err, &ve)
}

func TestBookingService_UpdateStatus_NotFound(t *testing.T) {
	repo := mocks.NewMockBookingRepository(t)
	repo.EXPECT().GetByID(mock.Anything, testMissingID).Return(model.Booking{}, model.ErrBookingNotFound)

	err := newBookingSvc(repo, mocks.NewMockPricingRuleRepository(t), mocks.NewMockEventPublisher(t)).
		UpdateStatus(ctxAsBookingManager(), testMissingID, validation.BookingStatusUpdate{Status: "cancelled"})

	assert.ErrorIs(t, err, model.ErrBookingNotFound)
}

func TestBookingService_UpdateStatus_InvalidTransition(t *testing.T) {
	repo := mocks.NewMockBookingRepository(t)
	repo.EXPECT().GetByID(mock.Anything, testBookingID).Return(model.Booking{Status: model.BookingStatusExpired}, nil)

	err := newBookingSvc(repo, mocks.NewMockPricingRuleRepository(t), mocks.NewMockEventPublisher(t)).
		UpdateStatus(ctxAsBookingManager(), testBookingID, validation.BookingStatusUpdate{Status: "cancelled"})

	assert.ErrorIs(t, err, model.ErrInvalidBookingStatusTransition)
}

// --- GetStatusHistory ---

func TestBookingService_GetStatusHistory_HappyPath(t *testing.T) {
	want := []model.BookingStatusReading{{ID: "h-1"}, {ID: "h-2"}}

	repo := mocks.NewMockBookingRepository(t)
	repo.EXPECT().GetByID(mock.Anything, testBookingID).Return(model.Booking{ID: testBookingID, UserID: testUserID}, nil)
	repo.EXPECT().GetStatusHistory(mock.Anything, mock.AnythingOfType("model.BookingStatusHistoryFilter")).Return(want, nil)

	got, err := newBookingSvc(repo, mocks.NewMockPricingRuleRepository(t), mocks.NewMockEventPublisher(t)).
		GetStatusHistory(ctxAsUser(testUserID), validation.BookingStatusHistoryFilter{BookingID: testBookingID})

	require.NoError(t, err)
	assert.Equal(t, want, got)
}

func TestBookingService_GetStatusHistory_RepoError(t *testing.T) {
	repoErr := errors.New("query failed")

	repo := mocks.NewMockBookingRepository(t)
	repo.EXPECT().GetByID(mock.Anything, testBookingID).Return(model.Booking{ID: testBookingID, UserID: testUserID}, nil)
	repo.EXPECT().GetStatusHistory(mock.Anything, mock.AnythingOfType("model.BookingStatusHistoryFilter")).Return(nil, repoErr)

	_, err := newBookingSvc(repo, mocks.NewMockPricingRuleRepository(t), mocks.NewMockEventPublisher(t)).
		GetStatusHistory(ctxAsUser(testUserID), validation.BookingStatusHistoryFilter{BookingID: testBookingID})

	assert.ErrorIs(t, err, repoErr)
}
