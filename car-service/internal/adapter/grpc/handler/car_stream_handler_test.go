package handler

import (
	"context"
	"testing"

	"carsharing/car-service/internal/adapter/grpc/handler/mocks"
	"carsharing/car-service/internal/model"
	"carsharing/car-service/internal/validation"
	carsvc "carsharing/protos/gen/service/car"
	"google.golang.org/grpc/metadata"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/codes"
)

// streamServer is a minimal in-process implementation of grpc.ServerStreamingServer[T].
type streamServer[T any] struct {
	ctx  context.Context
	msgs []*T
}

func newStream[T any](ctx context.Context) *streamServer[T] {
	return &streamServer[T]{ctx: ctx}
}

func (s *streamServer[T]) Send(msg *T) error            { s.msgs = append(s.msgs, msg); return nil }
func (s *streamServer[T]) Context() context.Context     { return s.ctx }
func (s *streamServer[T]) SetHeader(metadata.MD) error  { return nil }
func (s *streamServer[T]) SendHeader(metadata.MD) error { return nil }
func (s *streamServer[T]) SetTrailer(metadata.MD)       {}
func (s *streamServer[T]) RecvMsg(any) error            { return nil }
func (s *streamServer[T]) SendMsg(any) error            { return nil }

func TestStreamCarsWithFilter(t *testing.T) {
	t.Run("sends cars then exits on context cancel", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())

		svc := mocks.NewMockCarService(t)
		h := NewCarStreamHandler(discardLogger(), svc, nil)
		stream := newStream[carsvc.StreamCarsWithFilterResponse](ctx)

		svc.EXPECT().
			List(mock.Anything, mock.Anything).
			RunAndReturn(func(c context.Context, f validation.CarFilter) ([]model.Car, error) {
				cancel()
				return []model.Car{{ID: "c-1", Status: model.CarStatusAvailable}}, nil
			})

		err := h.StreamCarsWithFilter(&carsvc.StreamCarsWithFilterRequest{}, stream)
		assert.NoError(t, err)
		assert.Len(t, stream.msgs, 1)
		assert.Equal(t, "c-1", stream.msgs[0].Car[0].Id)
	})

	t.Run("service error is returned", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		svc := mocks.NewMockCarService(t)
		h := NewCarStreamHandler(discardLogger(), svc, nil)
		stream := newStream[carsvc.StreamCarsWithFilterResponse](ctx)

		svc.EXPECT().List(mock.Anything, mock.Anything).Return(nil, errInternal)

		err := h.StreamCarsWithFilter(&carsvc.StreamCarsWithFilterRequest{}, stream)
		assert.Equal(t, codes.Internal, grpcCode(err))
	})

	t.Run("min_fuel_level filters out low-fuel cars", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())

		svc := mocks.NewMockCarService(t)
		h := NewCarStreamHandler(discardLogger(), svc, nil)
		stream := newStream[carsvc.StreamCarsWithFilterResponse](ctx)

		highFuel := float32(80)
		lowFuel := float32(10)
		minFuel := float32(50)

		svc.EXPECT().
			List(mock.Anything, mock.Anything).
			RunAndReturn(func(c context.Context, f validation.CarFilter) ([]model.Car, error) {
				cancel()
				return []model.Car{
					{ID: "c-high", FuelLevel: &highFuel},
					{ID: "c-low", FuelLevel: &lowFuel},
				}, nil
			})

		err := h.StreamCarsWithFilter(&carsvc.StreamCarsWithFilterRequest{MinFuelLevel: &minFuel}, stream)
		assert.NoError(t, err)
		assert.Len(t, stream.msgs, 1)
		assert.Equal(t, "c-high", stream.msgs[0].Car[0].Id)
	})
}

func TestStreamCarTelemetry(t *testing.T) {
	carID := "c0000000-0000-4000-8000-000000000001"

	t.Run("sends initial state then updates from channel", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		svc := mocks.NewMockCarService(t)
		sub := mocks.NewMockTelemetrySubscriber(t)
		h := NewCarStreamHandler(discardLogger(), svc, sub)
		stream := newStream[carsvc.StreamCarTelemetryResponse](ctx)

		mileage := int64(50_000)
		svc.EXPECT().Get(mock.Anything, carID).Return(model.Car{
			ID: carID, MileageKM: mileage,
		}, nil)

		ch := make(chan model.TelemetryUpdate, 1)
		updatedMileage := int64(51_000)
		ch <- model.TelemetryUpdate{MileageKM: updatedMileage}
		close(ch)

		sub.EXPECT().Subscribe(mock.Anything, mock.Anything).Return((<-chan model.TelemetryUpdate)(ch), nil)

		err := h.StreamCarTelemetry(&carsvc.StreamCarTelemetryRequest{CarId: carID}, stream)
		assert.NoError(t, err)
		assert.Len(t, stream.msgs, 2)
		assert.Equal(t, mileage, stream.msgs[0].MileageKm)
		assert.Equal(t, updatedMileage, stream.msgs[1].MileageKm)
	})

	t.Run("car not found maps to NotFound", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		svc := mocks.NewMockCarService(t)
		h := NewCarStreamHandler(discardLogger(), svc, nil)
		stream := newStream[carsvc.StreamCarTelemetryResponse](ctx)

		svc.EXPECT().Get(mock.Anything, carID).Return(model.Car{}, model.ErrNotFound)

		err := h.StreamCarTelemetry(&carsvc.StreamCarTelemetryRequest{CarId: carID}, stream)
		assert.Equal(t, codes.NotFound, grpcCode(err))
	})

	t.Run("subscribe error is returned", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		svc := mocks.NewMockCarService(t)
		sub := mocks.NewMockTelemetrySubscriber(t)
		h := NewCarStreamHandler(discardLogger(), svc, sub)
		stream := newStream[carsvc.StreamCarTelemetryResponse](ctx)

		svc.EXPECT().Get(mock.Anything, carID).Return(model.Car{ID: carID}, nil)
		sub.EXPECT().Subscribe(mock.Anything, mock.Anything).Return(nil, errInternal)

		err := h.StreamCarTelemetry(&carsvc.StreamCarTelemetryRequest{CarId: carID}, stream)
		assert.Equal(t, codes.Internal, grpcCode(err))
	})

	t.Run("context cancellation stops the stream", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())

		svc := mocks.NewMockCarService(t)
		sub := mocks.NewMockTelemetrySubscriber(t)
		h := NewCarStreamHandler(discardLogger(), svc, sub)
		stream := newStream[carsvc.StreamCarTelemetryResponse](ctx)

		svc.EXPECT().Get(mock.Anything, carID).Return(model.Car{ID: carID}, nil)

		ch := make(chan model.TelemetryUpdate)
		sub.EXPECT().Subscribe(mock.Anything, mock.Anything).Return((<-chan model.TelemetryUpdate)(ch), nil)

		cancel()
		err := h.StreamCarTelemetry(&carsvc.StreamCarTelemetryRequest{CarId: carID}, stream)
		assert.NoError(t, err)
	})
}
