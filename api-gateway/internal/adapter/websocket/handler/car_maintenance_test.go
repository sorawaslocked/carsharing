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

type mockMaintenanceSvc struct {
	fn func(ctx context.Context, send func(model.CarMaintenanceEvent) error) error
}

func (m *mockMaintenanceSvc) StreamMaintenanceEvents(ctx context.Context, send func(model.CarMaintenanceEvent) error) error {
	return m.fn(ctx, send)
}

// ---- test server helpers ---------------------------------------------------

func maintenanceTestServer(t *testing.T, svc wshandler.CarMaintenanceStreamService) *httptest.Server {
	t.Helper()
	gin.SetMode(gin.TestMode)
	r := gin.New()
	h := wshandler.NewCarMaintenanceWsHandler(svc, slog.New(slog.NewTextHandler(io.Discard, nil)))
	r.GET("/ws", func(c *gin.Context) {
		c.Set("x-token-exp", time.Now().Add(time.Hour))
		h.MaintenanceEvents(c)
	})
	srv := httptest.NewServer(r)
	t.Cleanup(srv.Close)
	return srv
}

func maintenanceWsURL(srv *httptest.Server) string {
	return "ws://" + srv.Listener.Addr().String() + "/ws"
}

// ---- tests -----------------------------------------------------------------

func TestMaintenanceEvents_EventDelivered(t *testing.T) {
	now := time.Date(2024, 3, 10, 8, 30, 0, 0, time.UTC)
	event := model.CarMaintenanceEvent{
		CarID:      "car-1",
		TemplateID: "tmpl-42",
		RecordID:   "rec-99",
		EventType:  "warn",
		OccurredAt: now,
	}

	svc := &mockMaintenanceSvc{fn: func(_ context.Context, send func(model.CarMaintenanceEvent) error) error {
		return send(event)
	}}

	srv := maintenanceTestServer(t, svc)
	conn := dialWS(t, maintenanceWsURL(srv))

	var msg wsdto.CarMaintenanceEventMessage
	require.NoError(t, wsjson.Read(context.Background(), conn, &msg))

	assert.Equal(t, "car-1", msg.CarID)
	assert.Equal(t, "tmpl-42", msg.TemplateID)
	assert.Equal(t, "rec-99", msg.RecordID)
	assert.Equal(t, "warn", msg.EventType)
	assert.Equal(t, "2024-03-10T08:30:00Z", msg.OccurredAt)
}

func TestMaintenanceEvents_ForbiddenClosesWithPolicyViolation(t *testing.T) {
	svc := &mockMaintenanceSvc{fn: func(_ context.Context, _ func(model.CarMaintenanceEvent) error) error {
		return model.ErrForbidden
	}}

	srv := maintenanceTestServer(t, svc)
	conn := dialWS(t, maintenanceWsURL(srv))

	_, _, readErr := conn.Read(context.Background())
	assert.Equal(t, websocket.StatusPolicyViolation, websocket.CloseStatus(readErr))
}

func TestMaintenanceEvents_UnauthorizedClosesWithPolicyViolation(t *testing.T) {
	svc := &mockMaintenanceSvc{fn: func(_ context.Context, _ func(model.CarMaintenanceEvent) error) error {
		return model.ErrUnauthorized
	}}

	srv := maintenanceTestServer(t, svc)
	conn := dialWS(t, maintenanceWsURL(srv))

	_, _, readErr := conn.Read(context.Background())
	assert.Equal(t, websocket.StatusPolicyViolation, websocket.CloseStatus(readErr))
}

func TestMaintenanceEvents_ServiceErrorClosesWithInternalError(t *testing.T) {
	svc := &mockMaintenanceSvc{fn: func(_ context.Context, _ func(model.CarMaintenanceEvent) error) error {
		return model.ErrInternalServerError
	}}

	srv := maintenanceTestServer(t, svc)
	conn := dialWS(t, maintenanceWsURL(srv))

	_, _, readErr := conn.Read(context.Background())
	assert.Equal(t, websocket.StatusInternalError, websocket.CloseStatus(readErr))
}

func TestMaintenanceEvents_NormalClose(t *testing.T) {
	svc := &mockMaintenanceSvc{fn: func(_ context.Context, _ func(model.CarMaintenanceEvent) error) error {
		return nil
	}}

	srv := maintenanceTestServer(t, svc)
	conn := dialWS(t, maintenanceWsURL(srv))

	_, _, readErr := conn.Read(context.Background())
	assert.Equal(t, websocket.StatusNormalClosure, websocket.CloseStatus(readErr))
}

func TestMaintenanceEvents_MultipleEventsDeliveredInOrder(t *testing.T) {
	events := []model.CarMaintenanceEvent{
		{CarID: "car-1", EventType: "warn", OccurredAt: time.Now()},
		{CarID: "car-2", EventType: "pull", OccurredAt: time.Now()},
	}

	svc := &mockMaintenanceSvc{fn: func(_ context.Context, send func(model.CarMaintenanceEvent) error) error {
		for _, e := range events {
			if err := send(e); err != nil {
				return err
			}
		}
		return nil
	}}

	srv := maintenanceTestServer(t, svc)
	conn := dialWS(t, maintenanceWsURL(srv))

	var first, second wsdto.CarMaintenanceEventMessage
	require.NoError(t, wsjson.Read(context.Background(), conn, &first))
	require.NoError(t, wsjson.Read(context.Background(), conn, &second))

	assert.Equal(t, "car-1", first.CarID)
	assert.Equal(t, "warn", first.EventType)
	assert.Equal(t, "car-2", second.CarID)
	assert.Equal(t, "pull", second.EventType)
}
