package service

import (
	"context"
	"testing"
	"time"

	"carsharing/car-service/internal/model"
	"carsharing/car-service/internal/service/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// mockStreamClient is a minimal hand-written mock for TelemetryStreamClient.
// TelemetryStreamClient is not in the generated mock set, so we keep it here.
type mockStreamClient struct {
	subscribeFunc func(ctx context.Context, car model.Car) (<-chan model.TelemetryUpdate, error)
}

func (m *mockStreamClient) Subscribe(ctx context.Context, car model.Car) (<-chan model.TelemetryUpdate, error) {
	return m.subscribeFunc(ctx, car)
}

// blockingStreamClient returns a channel that closes when the provided context
// is cancelled, allowing goroutines to exit cleanly in tests without the
// reconnect delay.
func blockingStreamClient() *mockStreamClient {
	return &mockStreamClient{
		subscribeFunc: func(ctx context.Context, _ model.Car) (<-chan model.TelemetryUpdate, error) {
			ch := make(chan model.TelemetryUpdate)
			go func() {
				<-ctx.Done()
				close(ch)
			}()
			return ch, nil
		},
	}
}

func newTelemetrySvc(t *testing.T, client TelemetryStreamClient, telemetryRepo TelemetryReadingRepository, carRepo CarRepository) *TelemetryService {
	t.Helper()
	return NewTelemetryService(discardLogger(), newTestValidator(t), client, telemetryRepo, carRepo, 2*time.Minute)
}

// --- Ping ---

func TestTelemetryService_Ping(t *testing.T) {
	ctx := context.Background()

	t.Run("healthy when no streams configured", func(t *testing.T) {
		svc := newTelemetrySvc(t, nil, nil, nil)
		assert.NoError(t, svc.Ping(ctx))
	})

	t.Run("ErrTelemetryNoStreamsConnected when streams configured but none ever connected", func(t *testing.T) {
		svc := newTelemetrySvc(t, nil, nil, nil)
		svc.totalStreams.Store(1)

		assert.ErrorIs(t, svc.Ping(ctx), model.ErrTelemetryNoStreamsConnected)
	})

	t.Run("ErrTelemetryAllStreamsDisconnected when all streams previously active but now inactive", func(t *testing.T) {
		svc := newTelemetrySvc(t, nil, nil, nil)
		svc.totalStreams.Store(1)
		past := time.Now().Add(-3 * time.Minute)
		svc.lastSeenAt.Store(&past)

		var err model.ErrTelemetryAllStreamsDisconnected
		assert.ErrorAs(t, svc.Ping(ctx), &err)
	})

	t.Run("ErrTelemetryStreamStale when active but no updates within threshold", func(t *testing.T) {
		svc := NewTelemetryService(discardLogger(), newTestValidator(t), nil, nil, nil, 1*time.Minute)
		svc.totalStreams.Store(1)
		svc.activeStreams.Store(1)
		stale := time.Now().Add(-2 * time.Minute)
		svc.lastSeenAt.Store(&stale)

		var err model.ErrTelemetryStreamStale
		assert.ErrorAs(t, svc.Ping(ctx), &err)
	})

	t.Run("healthy when active streams have recent data", func(t *testing.T) {
		svc := newTelemetrySvc(t, nil, nil, nil)
		svc.totalStreams.Store(1)
		svc.activeStreams.Store(1)
		now := time.Now()
		svc.lastSeenAt.Store(&now)

		assert.NoError(t, svc.Ping(ctx))
	})
}

// --- Start ---

func TestTelemetryService_Start(t *testing.T) {
	t.Run("returns error when car listing fails", func(t *testing.T) {
		carRepo := mocks.NewMockCarRepository(t)
		svc := newTelemetrySvc(t, nil, nil, carRepo)

		carRepo.EXPECT().Find(mock.Anything, model.CarFilter{}).Return(nil, model.ErrSql)

		err := svc.Start(context.Background())
		assert.ErrorIs(t, err, model.ErrSql)
	})

	t.Run("starts one goroutine per car", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())

		carRepo := mocks.NewMockCarRepository(t)
		svc := newTelemetrySvc(t, blockingStreamClient(), nil, carRepo)

		carRepo.EXPECT().Find(mock.Anything, model.CarFilter{}).
			Return([]model.Car{{ID: "c-1"}, {ID: "c-2"}}, nil)

		err := svc.Start(ctx)
		require.NoError(t, err)
		assert.Equal(t, int32(2), svc.totalStreams.Load())

		cancel()
		svc.Stop()
	})
}

// --- OnCarCreated ---

func TestTelemetryService_OnCarCreated(t *testing.T) {
	t.Run("warns and skips when called before Start", func(t *testing.T) {
		svc := newTelemetrySvc(t, nil, nil, nil)

		svc.OnCarCreated(model.Car{ID: "c-1"})

		assert.Equal(t, int32(0), svc.totalStreams.Load())
	})

	t.Run("starts goroutine after Start context is set", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())

		svc := newTelemetrySvc(t, blockingStreamClient(), nil, nil)
		svc.mu.Lock()
		svc.ctx = ctx
		svc.mu.Unlock()

		svc.OnCarCreated(model.Car{ID: "c-1"})

		assert.Equal(t, int32(1), svc.totalStreams.Load())

		cancel()
		svc.Stop()
	})
}

// --- SubscribeUpdates ---

func TestTelemetryService_SubscribeUpdates(t *testing.T) {
	carID := "c-1"

	t.Run("subscriber receives fanned-out update", func(t *testing.T) {
		svc := newTelemetrySvc(t, nil, nil, nil)

		ch, unsub := svc.SubscribeUpdates(carID)
		defer unsub()

		update := model.TelemetryUpdate{CarID: carID, MileageKM: 1000}
		svc.fanOut(carID, update)

		select {
		case got := <-ch:
			assert.Equal(t, update, got)
		case <-time.After(time.Second):
			t.Fatal("timed out waiting for update")
		}
	})

	t.Run("unsubscribe stops delivery", func(t *testing.T) {
		svc := newTelemetrySvc(t, nil, nil, nil)

		ch, unsub := svc.SubscribeUpdates(carID)
		unsub()

		svc.fanOut(carID, model.TelemetryUpdate{CarID: carID, MileageKM: 1000})

		select {
		case <-ch:
			t.Fatal("received update after unsubscribe")
		case <-time.After(20 * time.Millisecond):
		}
	})

	t.Run("multiple subscribers each receive the update", func(t *testing.T) {
		svc := newTelemetrySvc(t, nil, nil, nil)

		ch1, unsub1 := svc.SubscribeUpdates(carID)
		ch2, unsub2 := svc.SubscribeUpdates(carID)
		defer unsub1()
		defer unsub2()

		update := model.TelemetryUpdate{CarID: carID, MileageKM: 2000}
		svc.fanOut(carID, update)

		for _, ch := range []<-chan model.TelemetryUpdate{ch1, ch2} {
			select {
			case got := <-ch:
				assert.Equal(t, update, got)
			case <-time.After(time.Second):
				t.Fatal("timed out waiting for update")
			}
		}
	})
}

// --- applyUpdate ---

func TestTelemetryService_applyUpdate(t *testing.T) {
	ctx := context.Background()
	carID := "c0000000-0000-4000-8000-000000000001"

	t.Run("inserts telemetry reading and updates car on valid update", func(t *testing.T) {
		carRepo := mocks.NewMockCarRepository(t)
		telemetryRepo := mocks.NewMockTelemetryReadingRepository(t)
		svc := newTelemetrySvc(t, nil, telemetryRepo, carRepo)

		update := model.TelemetryUpdate{CarID: carID, MileageKM: 10_000}

		carRepo.EXPECT().FindByID(ctx, carID).Return(model.Car{ID: carID, MileageKM: 9_000}, nil)
		telemetryRepo.EXPECT().Insert(ctx, mock.MatchedBy(func(r model.TelemetryReading) bool {
			return r.CarID == carID && r.MileageKM != nil && *r.MileageKM == 10_000
		})).Return(nil)
		carRepo.EXPECT().Update(ctx, carID, mock.MatchedBy(func(u model.CarUpdate) bool {
			return u.MileageKM != nil && *u.MileageKM == 10_000
		})).Return(nil)

		err := svc.applyUpdate(ctx, discardLogger(), update)
		assert.NoError(t, err)
		assert.NotNil(t, svc.lastSeenAt.Load())
	})

	t.Run("rejects mileage regression", func(t *testing.T) {
		carRepo := mocks.NewMockCarRepository(t)
		svc := newTelemetrySvc(t, nil, nil, carRepo)

		update := model.TelemetryUpdate{CarID: carID, MileageKM: 5_000}
		carRepo.EXPECT().FindByID(ctx, carID).Return(model.Car{ID: carID, MileageKM: 9_000}, nil)

		err := svc.applyUpdate(ctx, discardLogger(), update)
		assert.ErrorIs(t, err, model.ErrMileageRegression)
	})

	t.Run("propagates car repo find error", func(t *testing.T) {
		carRepo := mocks.NewMockCarRepository(t)
		svc := newTelemetrySvc(t, nil, nil, carRepo)

		update := model.TelemetryUpdate{CarID: carID, MileageKM: 10_000}
		carRepo.EXPECT().FindByID(ctx, carID).Return(model.Car{}, model.ErrCarNotFound)

		err := svc.applyUpdate(ctx, discardLogger(), update)
		assert.ErrorIs(t, err, model.ErrCarNotFound)
	})

	t.Run("propagates telemetry reading insert error", func(t *testing.T) {
		carRepo := mocks.NewMockCarRepository(t)
		telemetryRepo := mocks.NewMockTelemetryReadingRepository(t)
		svc := newTelemetrySvc(t, nil, telemetryRepo, carRepo)

		update := model.TelemetryUpdate{CarID: carID, MileageKM: 10_000}
		carRepo.EXPECT().FindByID(ctx, carID).Return(model.Car{ID: carID, MileageKM: 9_000}, nil)
		telemetryRepo.EXPECT().Insert(ctx, mock.Anything).Return(model.ErrSql)

		err := svc.applyUpdate(ctx, discardLogger(), update)
		assert.ErrorIs(t, err, model.ErrSql)
	})

	t.Run("propagates car update error", func(t *testing.T) {
		carRepo := mocks.NewMockCarRepository(t)
		telemetryRepo := mocks.NewMockTelemetryReadingRepository(t)
		svc := newTelemetrySvc(t, nil, telemetryRepo, carRepo)

		update := model.TelemetryUpdate{CarID: carID, MileageKM: 10_000}
		carRepo.EXPECT().FindByID(ctx, carID).Return(model.Car{ID: carID, MileageKM: 9_000}, nil)
		telemetryRepo.EXPECT().Insert(ctx, mock.Anything).Return(nil)
		carRepo.EXPECT().Update(ctx, carID, mock.Anything).Return(model.ErrSql)

		err := svc.applyUpdate(ctx, discardLogger(), update)
		assert.ErrorIs(t, err, model.ErrSql)
	})
}
