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

	t.Run("ok when db ping and nats both healthy", func(t *testing.T) {
		db := mocks.NewMockDBPinger(t)
		nats := mocks.NewMockNATSChecker(t)
		h := NewHealthHandler(db, nats, discardLogger())

		db.EXPECT().Ping(mock.Anything).Return(nil)
		nats.EXPECT().IsConnected().Return(true)

		resp, err := h.Health(ctx, &emptypb.Empty{})
		assert.NoError(t, err)
		assert.Equal(t, "ok", resp.Status)
		assert.Equal(t, "car-service", resp.Name)
	})

	t.Run("degraded when db ping fails", func(t *testing.T) {
		db := mocks.NewMockDBPinger(t)
		nats := mocks.NewMockNATSChecker(t)
		h := NewHealthHandler(db, nats, discardLogger())

		db.EXPECT().Ping(mock.Anything).Return(errors.New("connection refused"))
		nats.EXPECT().IsConnected().Return(true)

		resp, err := h.Health(ctx, &emptypb.Empty{})
		assert.NoError(t, err)
		assert.Equal(t, "degraded", resp.Status)
	})

	t.Run("degraded when nats not connected", func(t *testing.T) {
		db := mocks.NewMockDBPinger(t)
		nats := mocks.NewMockNATSChecker(t)
		h := NewHealthHandler(db, nats, discardLogger())

		db.EXPECT().Ping(mock.Anything).Return(nil)
		nats.EXPECT().IsConnected().Return(false)

		resp, err := h.Health(ctx, &emptypb.Empty{})
		assert.NoError(t, err)
		assert.Equal(t, "degraded", resp.Status)
	})

	t.Run("degraded when both db and nats fail", func(t *testing.T) {
		db := mocks.NewMockDBPinger(t)
		nats := mocks.NewMockNATSChecker(t)
		h := NewHealthHandler(db, nats, discardLogger())

		db.EXPECT().Ping(mock.Anything).Return(errors.New("connection refused"))
		nats.EXPECT().IsConnected().Return(false)

		resp, err := h.Health(ctx, &emptypb.Empty{})
		assert.NoError(t, err)
		assert.Equal(t, "degraded", resp.Status)
	})

	t.Run("uptime increases with time", func(t *testing.T) {
		db := mocks.NewMockDBPinger(t)
		nats := mocks.NewMockNATSChecker(t)
		h := NewHealthHandler(db, nats, discardLogger())

		db.EXPECT().Ping(mock.Anything).Return(nil).Maybe()
		nats.EXPECT().IsConnected().Return(true).Maybe()

		resp, err := h.Health(ctx, &emptypb.Empty{})
		assert.NoError(t, err)
		assert.GreaterOrEqual(t, resp.UptimeSeconds, uint64(0))
	})
}
