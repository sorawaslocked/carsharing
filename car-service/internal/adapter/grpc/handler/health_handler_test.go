package handler

import (
	"context"
	"errors"
	"testing"

	"carsharing/car-service/internal/adapter/grpc/handler/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/protobuf/types/known/emptypb"
)

func TestHealthHandlerHealth(t *testing.T) {
	ctx := context.Background()

	t.Run("healthy when all deps ping successfully", func(t *testing.T) {
		pg := mocks.NewMockPinger(t)
		nc := mocks.NewMockPinger(t)
		h := NewHealthHandler(discardLogger(), map[string]Pinger{"postgres": pg, "nats": nc})

		pg.EXPECT().Ping(mock.Anything).Return(nil)
		nc.EXPECT().Ping(mock.Anything).Return(nil)

		resp, err := h.Health(ctx, &emptypb.Empty{})
		assert.NoError(t, err)
		assert.Equal(t, "healthy", resp.Status)
		assert.Equal(t, "car-service", resp.Name)
	})

	t.Run("unhealthy when one dep ping fails", func(t *testing.T) {
		pg := mocks.NewMockPinger(t)
		nc := mocks.NewMockPinger(t)
		h := NewHealthHandler(discardLogger(), map[string]Pinger{"postgres": pg, "nats": nc})

		pg.EXPECT().Ping(mock.Anything).Return(errors.New("connection refused"))
		nc.EXPECT().Ping(mock.Anything).Return(nil)

		resp, err := h.Health(ctx, &emptypb.Empty{})
		assert.NoError(t, err)
		assert.Equal(t, "unhealthy", resp.Status)
	})

	t.Run("unhealthy when all deps fail", func(t *testing.T) {
		pg := mocks.NewMockPinger(t)
		nc := mocks.NewMockPinger(t)
		h := NewHealthHandler(discardLogger(), map[string]Pinger{"postgres": pg, "nats": nc})

		pg.EXPECT().Ping(mock.Anything).Return(errors.New("connection refused"))
		nc.EXPECT().Ping(mock.Anything).Return(errors.New("disconnected"))

		resp, err := h.Health(ctx, &emptypb.Empty{})
		assert.NoError(t, err)
		assert.Equal(t, "unhealthy", resp.Status)
	})

	t.Run("healthy with no deps", func(t *testing.T) {
		h := NewHealthHandler(discardLogger(), map[string]Pinger{})

		resp, err := h.Health(ctx, &emptypb.Empty{})
		assert.NoError(t, err)
		assert.Equal(t, "healthy", resp.Status)
	})

	t.Run("uptime is non-negative", func(t *testing.T) {
		h := NewHealthHandler(discardLogger(), map[string]Pinger{})

		resp, err := h.Health(ctx, &emptypb.Empty{})
		assert.NoError(t, err)
		assert.GreaterOrEqual(t, resp.UptimeSeconds, uint64(0))
	})
}
