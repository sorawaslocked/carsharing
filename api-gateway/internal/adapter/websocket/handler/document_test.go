package handler_test

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
	"github.com/gin-gonic/gin"

	wsdto "carsharing/api-gateway/internal/adapter/websocket/dto"
	wshandler "carsharing/api-gateway/internal/adapter/websocket/handler"
	"carsharing/api-gateway/internal/model"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ---- mock service ----------------------------------------------------------

type mockDocSvc struct {
	fn func(ctx context.Context, userID *string, passed *bool, send func(model.DocumentAnalyzedEvent) error) error
}

func (m *mockDocSvc) StreamDocumentAnalyzed(ctx context.Context, userID *string, passed *bool, send func(model.DocumentAnalyzedEvent) error) error {
	return m.fn(ctx, userID, passed, send)
}

// ---- test server helpers ---------------------------------------------------

func wsTestServer(t *testing.T, svc wshandler.DocumentStreamService) *httptest.Server {
	t.Helper()
	gin.SetMode(gin.TestMode)

	r := gin.New()
	h := wshandler.NewUserWsHandler(svc, slog.New(slog.NewTextHandler(io.Discard, nil)))

	r.GET("/ws", func(c *gin.Context) {
		c.Set("x-token-exp", time.Now().Add(time.Hour))
		h.DocumentUpdates(c)
	})

	srv := httptest.NewServer(r)
	t.Cleanup(srv.Close)
	return srv
}

func wsURL(srv *httptest.Server, query string) string {
	addr := "ws://" + srv.Listener.Addr().String() + "/ws"
	if query != "" {
		addr += "?" + query
	}
	return addr
}

func dialWS(t *testing.T, url string) *websocket.Conn {
	t.Helper()
	conn, _, err := websocket.Dial(context.Background(), url, &websocket.DialOptions{
		CompressionMode: websocket.CompressionDisabled,
	})
	require.NoError(t, err)
	t.Cleanup(func() { conn.CloseNow() })
	return conn
}

// ---- tests -----------------------------------------------------------------

func TestDocumentUpdates_EventDelivered(t *testing.T) {
	event := model.DocumentAnalyzedEvent{
		DocumentID: "doc-123",
		UserID:     "user-456",
		Passed:     true,
		Defects:    nil,
	}

	svc := &mockDocSvc{fn: func(_ context.Context, _ *string, _ *bool, send func(model.DocumentAnalyzedEvent) error) error {
		return send(event)
	}}

	srv := wsTestServer(t, svc)
	conn := dialWS(t, wsURL(srv, ""))

	var msg wsdto.DocumentAnalyzedMessage
	err := wsjson.Read(context.Background(), conn, &msg)
	require.NoError(t, err)

	assert.Equal(t, "doc-123", msg.DocumentID)
	assert.Equal(t, "user-456", msg.UserID)
	assert.True(t, msg.Passed)
	assert.Empty(t, msg.Defects)
}

func TestDocumentUpdates_EventWithDefectsDelivered(t *testing.T) {
	event := model.DocumentAnalyzedEvent{
		DocumentID: "doc-789",
		UserID:     "user-1",
		Passed:     false,
		Defects: []model.DocumentDefect{
			{Type: "blur", Description: "too blurry"},
			{Type: "crop", Description: "edges cropped"},
		},
	}

	svc := &mockDocSvc{fn: func(_ context.Context, _ *string, _ *bool, send func(model.DocumentAnalyzedEvent) error) error {
		return send(event)
	}}

	srv := wsTestServer(t, svc)
	conn := dialWS(t, wsURL(srv, ""))

	var msg wsdto.DocumentAnalyzedMessage
	require.NoError(t, wsjson.Read(context.Background(), conn, &msg))

	assert.False(t, msg.Passed)
	require.Len(t, msg.Defects, 2)
	assert.Equal(t, "blur", msg.Defects[0].Type)
	assert.Equal(t, "crop", msg.Defects[1].Type)
}

func TestDocumentUpdates_PermissionDenied_ClosesWithPolicyViolation(t *testing.T) {
	svc := &mockDocSvc{fn: func(_ context.Context, _ *string, _ *bool, _ func(model.DocumentAnalyzedEvent) error) error {
		return model.ErrForbidden
	}}

	srv := wsTestServer(t, svc)
	conn := dialWS(t, wsURL(srv, ""))

	// drain until close frame
	_, _, readErr := conn.Read(context.Background())
	status := websocket.CloseStatus(readErr)
	assert.Equal(t, websocket.StatusPolicyViolation, status)
}

func TestDocumentUpdates_Unauthorized_ClosesWithPolicyViolation(t *testing.T) {
	svc := &mockDocSvc{fn: func(_ context.Context, _ *string, _ *bool, _ func(model.DocumentAnalyzedEvent) error) error {
		return model.ErrUnauthorized
	}}

	srv := wsTestServer(t, svc)
	conn := dialWS(t, wsURL(srv, ""))

	_, _, readErr := conn.Read(context.Background())
	assert.Equal(t, websocket.StatusPolicyViolation, websocket.CloseStatus(readErr))
}

func TestDocumentUpdates_ServiceError_ClosesWithInternalError(t *testing.T) {
	svc := &mockDocSvc{fn: func(_ context.Context, _ *string, _ *bool, _ func(model.DocumentAnalyzedEvent) error) error {
		return model.ErrInternalServerError
	}}

	srv := wsTestServer(t, svc)
	conn := dialWS(t, wsURL(srv, ""))

	_, _, readErr := conn.Read(context.Background())
	assert.Equal(t, websocket.StatusInternalError, websocket.CloseStatus(readErr))
}

func TestDocumentUpdates_NormalClose(t *testing.T) {
	svc := &mockDocSvc{fn: func(_ context.Context, _ *string, _ *bool, _ func(model.DocumentAnalyzedEvent) error) error {
		return nil
	}}

	srv := wsTestServer(t, svc)
	conn := dialWS(t, wsURL(srv, ""))

	_, _, readErr := conn.Read(context.Background())
	assert.Equal(t, websocket.StatusNormalClosure, websocket.CloseStatus(readErr))
}

func TestDocumentUpdates_PassedQueryParamForwarded(t *testing.T) {
	var capturedPassed *bool

	svc := &mockDocSvc{fn: func(_ context.Context, _ *string, passed *bool, _ func(model.DocumentAnalyzedEvent) error) error {
		capturedPassed = passed
		return nil
	}}

	srv := wsTestServer(t, svc)
	conn := dialWS(t, wsURL(srv, "passed=true"))

	_, _, _ = conn.Read(context.Background()) // wait for close
	require.NotNil(t, capturedPassed)
	assert.True(t, *capturedPassed)
}

func TestDocumentUpdates_PassedFalseForwarded(t *testing.T) {
	var capturedPassed *bool

	svc := &mockDocSvc{fn: func(_ context.Context, _ *string, passed *bool, _ func(model.DocumentAnalyzedEvent) error) error {
		capturedPassed = passed
		return nil
	}}

	srv := wsTestServer(t, svc)
	conn := dialWS(t, wsURL(srv, "passed=false"))

	_, _, _ = conn.Read(context.Background())
	require.NotNil(t, capturedPassed)
	assert.False(t, *capturedPassed)
}

func TestDocumentUpdates_NoPassed_FilterIsNil(t *testing.T) {
	var capturedPassed *bool = ptr(true) // pre-set to non-nil to confirm it's actually nil after

	svc := &mockDocSvc{fn: func(_ context.Context, _ *string, passed *bool, _ func(model.DocumentAnalyzedEvent) error) error {
		capturedPassed = passed
		return nil
	}}

	srv := wsTestServer(t, svc)
	conn := dialWS(t, wsURL(srv, ""))

	_, _, _ = conn.Read(context.Background())
	assert.Nil(t, capturedPassed, "no passed query param should produce a nil filter")
}

func TestDocumentUpdates_UserIDQueryParamForwarded(t *testing.T) {
	var capturedUserID *string

	svc := &mockDocSvc{fn: func(_ context.Context, userID *string, _ *bool, _ func(model.DocumentAnalyzedEvent) error) error {
		capturedUserID = userID
		return nil
	}}

	srv := wsTestServer(t, svc)
	conn := dialWS(t, wsURL(srv, "userID=some-user-id"))

	_, _, _ = conn.Read(context.Background())
	require.NotNil(t, capturedUserID)
	assert.Equal(t, "some-user-id", *capturedUserID)
}

func TestDocumentUpdates_MessageIsValidJSON(t *testing.T) {
	event := model.DocumentAnalyzedEvent{DocumentID: "doc-1", UserID: "u-1", Passed: true}

	svc := &mockDocSvc{fn: func(_ context.Context, _ *string, _ *bool, send func(model.DocumentAnalyzedEvent) error) error {
		return send(event)
	}}

	srv := wsTestServer(t, svc)
	conn := dialWS(t, wsURL(srv, ""))

	_, raw, err := conn.Read(context.Background())
	require.NoError(t, err)

	var m map[string]interface{}
	require.NoError(t, json.Unmarshal(raw, &m))
	assert.Equal(t, "doc-1", m["documentID"])
}

// ptr is a generic helper for pointer literals, duplicated here to keep
// the test file self-contained.
func ptr[T any](v T) *T { return &v }
