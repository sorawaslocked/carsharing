package handler_test

import (
	"context"
	"io"
	"log/slog"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	wsdto "carsharing/api-gateway/internal/adapter/websocket/dto"
	wshandler "carsharing/api-gateway/internal/adapter/websocket/handler"
	"carsharing/api-gateway/internal/model"
)

// ---- mock service ----------------------------------------------------------

type mockTripSvc struct {
	fn func(ctx context.Context, tripID string, send func(model.TripLiveFeed) error) error
}

func (m *mockTripSvc) StreamTripLiveFeed(ctx context.Context, tripID string, send func(model.TripLiveFeed) error) error {
	return m.fn(ctx, tripID, send)
}

// ---- test server helpers ---------------------------------------------------

func tripTestServer(t *testing.T, svc wshandler.TripStreamService) *httptest.Server {
	t.Helper()
	gin.SetMode(gin.TestMode)
	r := gin.New()
	h := wshandler.NewTripWsHandler(svc, slog.New(slog.NewTextHandler(io.Discard, nil)))
	r.GET("/ws/:id", func(c *gin.Context) {
		c.Set("x-token-exp", time.Now().Add(time.Hour))
		h.LiveFeed(c)
	})
	srv := httptest.NewServer(r)
	t.Cleanup(srv.Close)
	return srv
}

func tripWsURL(srv *httptest.Server, tripID string) string {
	return "ws://" + srv.Listener.Addr().String() + "/ws/" + tripID
}

// ---- tests -----------------------------------------------------------------

func TestTripLiveFeed_EventDelivered(t *testing.T) {
	feed := model.TripLiveFeed{ElapsedSeconds: 120, CurrentCostTenge: 500, DistanceTraveledKM: 3.5}

	svc := &mockTripSvc{fn: func(_ context.Context, _ string, send func(model.TripLiveFeed) error) error {
		return send(feed)
	}}

	srv := tripTestServer(t, svc)
	conn := dialWS(t, tripWsURL(srv, "trip-1"))

	var msg wsdto.TripLiveFeedMessage
	require.NoError(t, wsjson.Read(context.Background(), conn, &msg))

	assert.Equal(t, int64(120), msg.ElapsedSeconds)
	assert.Equal(t, int32(500), msg.CurrentCostTenge)
	assert.Equal(t, float64(3.5), msg.DistanceTraveledKM)
}

func TestTripLiveFeed_TripIDForwarded(t *testing.T) {
	var capturedID string

	svc := &mockTripSvc{fn: func(_ context.Context, tripID string, _ func(model.TripLiveFeed) error) error {
		capturedID = tripID
		return nil
	}}

	srv := tripTestServer(t, svc)
	conn := dialWS(t, tripWsURL(srv, "trip-xyz"))

	_, _, _ = conn.Read(context.Background())
	assert.Equal(t, "trip-xyz", capturedID)
}

func TestTripLiveFeed_ForbiddenClosesWithPolicyViolation(t *testing.T) {
	svc := &mockTripSvc{fn: func(_ context.Context, _ string, _ func(model.TripLiveFeed) error) error {
		return model.ErrForbidden
	}}

	srv := tripTestServer(t, svc)
	conn := dialWS(t, tripWsURL(srv, "trip-1"))

	_, _, readErr := conn.Read(context.Background())
	assert.Equal(t, websocket.StatusPolicyViolation, websocket.CloseStatus(readErr))
}

func TestTripLiveFeed_UnauthorizedClosesWithPolicyViolation(t *testing.T) {
	svc := &mockTripSvc{fn: func(_ context.Context, _ string, _ func(model.TripLiveFeed) error) error {
		return model.ErrUnauthorized
	}}

	srv := tripTestServer(t, svc)
	conn := dialWS(t, tripWsURL(srv, "trip-1"))

	_, _, readErr := conn.Read(context.Background())
	assert.Equal(t, websocket.StatusPolicyViolation, websocket.CloseStatus(readErr))
}

func TestTripLiveFeed_ServiceErrorClosesWithInternalError(t *testing.T) {
	svc := &mockTripSvc{fn: func(_ context.Context, _ string, _ func(model.TripLiveFeed) error) error {
		return model.ErrInternalServerError
	}}

	srv := tripTestServer(t, svc)
	conn := dialWS(t, tripWsURL(srv, "trip-1"))

	_, _, readErr := conn.Read(context.Background())
	assert.Equal(t, websocket.StatusInternalError, websocket.CloseStatus(readErr))
}

func TestTripLiveFeed_NormalClose(t *testing.T) {
	svc := &mockTripSvc{fn: func(_ context.Context, _ string, _ func(model.TripLiveFeed) error) error {
		return nil
	}}

	srv := tripTestServer(t, svc)
	conn := dialWS(t, tripWsURL(srv, "trip-1"))

	_, _, readErr := conn.Read(context.Background())
	assert.Equal(t, websocket.StatusNormalClosure, websocket.CloseStatus(readErr))
}

func TestTripLiveFeed_MultipleEventsDeliveredInOrder(t *testing.T) {
	feeds := []model.TripLiveFeed{
		{ElapsedSeconds: 10, CurrentCostTenge: 100, DistanceTraveledKM: 0.5},
		{ElapsedSeconds: 20, CurrentCostTenge: 200, DistanceTraveledKM: 1.0},
	}

	svc := &mockTripSvc{fn: func(_ context.Context, _ string, send func(model.TripLiveFeed) error) error {
		for _, f := range feeds {
			if err := send(f); err != nil {
				return err
			}
		}
		return nil
	}}

	srv := tripTestServer(t, svc)
	conn := dialWS(t, tripWsURL(srv, "trip-1"))

	var first, second wsdto.TripLiveFeedMessage
	require.NoError(t, wsjson.Read(context.Background(), conn, &first))
	require.NoError(t, wsjson.Read(context.Background(), conn, &second))

	assert.Equal(t, int64(10), first.ElapsedSeconds)
	assert.Equal(t, int64(20), second.ElapsedSeconds)
}
