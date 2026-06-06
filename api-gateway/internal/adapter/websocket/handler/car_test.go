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

type mockCarSvc struct {
	fleetFn     func(ctx context.Context, filter model.CarFilter, send func([]model.SlimCar) error) error
	telemetryFn func(ctx context.Context, carID string, send func(model.CarTelemetryEvent) error) error
	statusFn    func(ctx context.Context, carID string, send func(model.CarStatusEvent) error) error
}

func (m *mockCarSvc) StreamCarsWithFilter(ctx context.Context, filter model.CarFilter, send func([]model.SlimCar) error) error {
	return m.fleetFn(ctx, filter, send)
}
func (m *mockCarSvc) StreamCarTelemetry(ctx context.Context, carID string, send func(model.CarTelemetryEvent) error) error {
	return m.telemetryFn(ctx, carID, send)
}
func (m *mockCarSvc) StreamCarStatusUpdates(ctx context.Context, carID string, send func(model.CarStatusEvent) error) error {
	return m.statusFn(ctx, carID, send)
}

// ---- test server helpers ---------------------------------------------------

func carTestServer(t *testing.T, svc wshandler.CarStreamService) *httptest.Server {
	t.Helper()
	gin.SetMode(gin.TestMode)
	r := gin.New()
	h := wshandler.NewCarWsHandler(svc, slog.New(slog.NewTextHandler(io.Discard, nil)))
	r.GET("/ws/fleet", func(c *gin.Context) {
		c.Set("x-token-exp", time.Now().Add(time.Hour))
		h.Fleet(c)
	})
	r.GET("/ws/cars/:id/telemetry", func(c *gin.Context) {
		c.Set("x-token-exp", time.Now().Add(time.Hour))
		h.Telemetry(c)
	})
	r.GET("/ws/cars/:id/status", func(c *gin.Context) {
		c.Set("x-token-exp", time.Now().Add(time.Hour))
		h.Status(c)
	})
	srv := httptest.NewServer(r)
	t.Cleanup(srv.Close)
	return srv
}

func carWsURL(srv *httptest.Server, path string) string {
	return "ws://" + srv.Listener.Addr().String() + path
}

// nopCarSvc returns a mock that panics on any unexpected method call.
func nopFleet(fn func(context.Context, model.CarFilter, func([]model.SlimCar) error) error) *mockCarSvc {
	return &mockCarSvc{fleetFn: fn}
}
func nopTelemetry(fn func(context.Context, string, func(model.CarTelemetryEvent) error) error) *mockCarSvc {
	return &mockCarSvc{telemetryFn: fn}
}
func nopStatus(fn func(context.Context, string, func(model.CarStatusEvent) error) error) *mockCarSvc {
	return &mockCarSvc{statusFn: fn}
}

// ---- Fleet tests -----------------------------------------------------------

func TestCarFleet_BatchDelivered(t *testing.T) {
	cars := []model.SlimCar{
		{ID: "car-1", LicensePlate: "AA001AA", Status: "available", FuelLevel: 0.8},
		{ID: "car-2", LicensePlate: "BB002BB", Status: "reserved", FuelLevel: 0.5},
	}

	svc := nopFleet(func(_ context.Context, _ model.CarFilter, send func([]model.SlimCar) error) error {
		return send(cars)
	})

	srv := carTestServer(t, svc)
	conn := dialWS(t, carWsURL(srv, "/ws/fleet"))

	var msg wsdto.CarFleetMessage
	require.NoError(t, wsjson.Read(context.Background(), conn, &msg))

	require.Len(t, msg.Cars, 2)
	assert.Equal(t, "car-1", msg.Cars[0].ID)
	assert.Equal(t, "AA001AA", msg.Cars[0].LicensePlate)
	assert.Equal(t, "car-2", msg.Cars[1].ID)
}

func TestCarFleet_ServiceErrorClosesWithInternalError(t *testing.T) {
	svc := nopFleet(func(_ context.Context, _ model.CarFilter, _ func([]model.SlimCar) error) error {
		return model.ErrInternalServerError
	})

	srv := carTestServer(t, svc)
	conn := dialWS(t, carWsURL(srv, "/ws/fleet"))

	_, _, readErr := conn.Read(context.Background())
	assert.Equal(t, websocket.StatusInternalError, websocket.CloseStatus(readErr))
}

func TestCarFleet_NormalClose(t *testing.T) {
	svc := nopFleet(func(_ context.Context, _ model.CarFilter, _ func([]model.SlimCar) error) error {
		return nil
	})

	srv := carTestServer(t, svc)
	conn := dialWS(t, carWsURL(srv, "/ws/fleet"))

	_, _, readErr := conn.Read(context.Background())
	assert.Equal(t, websocket.StatusNormalClosure, websocket.CloseStatus(readErr))
}

// ---- Telemetry tests -------------------------------------------------------

func TestCarTelemetry_EventDelivered(t *testing.T) {
	now := time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)
	event := model.CarTelemetryEvent{FuelLevel: 0.75, BatteryLevel: 0.9, MileageKM: 12345, RecordedAt: now}

	svc := nopTelemetry(func(_ context.Context, _ string, send func(model.CarTelemetryEvent) error) error {
		return send(event)
	})

	srv := carTestServer(t, svc)
	conn := dialWS(t, carWsURL(srv, "/ws/cars/car-1/telemetry"))

	var msg wsdto.CarTelemetryMessage
	require.NoError(t, wsjson.Read(context.Background(), conn, &msg))

	assert.InDelta(t, 0.75, msg.FuelLevel, 0.001)
	assert.InDelta(t, 0.9, msg.BatteryLevel, 0.001)
	assert.Equal(t, int64(12345), msg.MileageKM)
	assert.Equal(t, "2024-01-15T10:00:00Z", msg.RecordedAt)
}

func TestCarTelemetry_CarIDForwarded(t *testing.T) {
	var capturedID string

	svc := nopTelemetry(func(_ context.Context, carID string, _ func(model.CarTelemetryEvent) error) error {
		capturedID = carID
		return nil
	})

	srv := carTestServer(t, svc)
	conn := dialWS(t, carWsURL(srv, "/ws/cars/car-abc/telemetry"))

	_, _, _ = conn.Read(context.Background())
	assert.Equal(t, "car-abc", capturedID)
}

func TestCarTelemetry_ForbiddenClosesWithPolicyViolation(t *testing.T) {
	svc := nopTelemetry(func(_ context.Context, _ string, _ func(model.CarTelemetryEvent) error) error {
		return model.ErrForbidden
	})

	srv := carTestServer(t, svc)
	conn := dialWS(t, carWsURL(srv, "/ws/cars/car-1/telemetry"))

	_, _, readErr := conn.Read(context.Background())
	assert.Equal(t, websocket.StatusPolicyViolation, websocket.CloseStatus(readErr))
}

func TestCarTelemetry_ServiceErrorClosesWithInternalError(t *testing.T) {
	svc := nopTelemetry(func(_ context.Context, _ string, _ func(model.CarTelemetryEvent) error) error {
		return model.ErrInternalServerError
	})

	srv := carTestServer(t, svc)
	conn := dialWS(t, carWsURL(srv, "/ws/cars/car-1/telemetry"))

	_, _, readErr := conn.Read(context.Background())
	assert.Equal(t, websocket.StatusInternalError, websocket.CloseStatus(readErr))
}

func TestCarTelemetry_NormalClose(t *testing.T) {
	svc := nopTelemetry(func(_ context.Context, _ string, _ func(model.CarTelemetryEvent) error) error {
		return nil
	})

	srv := carTestServer(t, svc)
	conn := dialWS(t, carWsURL(srv, "/ws/cars/car-1/telemetry"))

	_, _, readErr := conn.Read(context.Background())
	assert.Equal(t, websocket.StatusNormalClosure, websocket.CloseStatus(readErr))
}

// ---- Status tests ----------------------------------------------------------

func TestCarStatus_EventDelivered(t *testing.T) {
	event := model.CarStatusEvent{FromStatus: "available", ToStatus: "reserved"}

	svc := nopStatus(func(_ context.Context, _ string, send func(model.CarStatusEvent) error) error {
		return send(event)
	})

	srv := carTestServer(t, svc)
	conn := dialWS(t, carWsURL(srv, "/ws/cars/car-1/status"))

	var msg wsdto.CarStatusMessage
	require.NoError(t, wsjson.Read(context.Background(), conn, &msg))

	assert.Equal(t, "car-1", msg.CarID)
	assert.Equal(t, "available", msg.FromStatus)
	assert.Equal(t, "reserved", msg.ToStatus)
}

func TestCarStatus_CarIDForwarded(t *testing.T) {
	var capturedID string

	svc := nopStatus(func(_ context.Context, carID string, _ func(model.CarStatusEvent) error) error {
		capturedID = carID
		return nil
	})

	srv := carTestServer(t, svc)
	conn := dialWS(t, carWsURL(srv, "/ws/cars/car-xyz/status"))

	_, _, _ = conn.Read(context.Background())
	assert.Equal(t, "car-xyz", capturedID)
}

func TestCarStatus_ForbiddenClosesWithPolicyViolation(t *testing.T) {
	svc := nopStatus(func(_ context.Context, _ string, _ func(model.CarStatusEvent) error) error {
		return model.ErrForbidden
	})

	srv := carTestServer(t, svc)
	conn := dialWS(t, carWsURL(srv, "/ws/cars/car-1/status"))

	_, _, readErr := conn.Read(context.Background())
	assert.Equal(t, websocket.StatusPolicyViolation, websocket.CloseStatus(readErr))
}

func TestCarStatus_ServiceErrorClosesWithInternalError(t *testing.T) {
	svc := nopStatus(func(_ context.Context, _ string, _ func(model.CarStatusEvent) error) error {
		return model.ErrInternalServerError
	})

	srv := carTestServer(t, svc)
	conn := dialWS(t, carWsURL(srv, "/ws/cars/car-1/status"))

	_, _, readErr := conn.Read(context.Background())
	assert.Equal(t, websocket.StatusInternalError, websocket.CloseStatus(readErr))
}

func TestCarStatus_NormalClose(t *testing.T) {
	svc := nopStatus(func(_ context.Context, _ string, _ func(model.CarStatusEvent) error) error {
		return nil
	})

	srv := carTestServer(t, svc)
	conn := dialWS(t, carWsURL(srv, "/ws/cars/car-1/status"))

	_, _, readErr := conn.Read(context.Background())
	assert.Equal(t, websocket.StatusNormalClosure, websocket.CloseStatus(readErr))
}
