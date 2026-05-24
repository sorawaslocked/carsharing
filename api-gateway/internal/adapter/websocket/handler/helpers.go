package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/coder/websocket"
	"github.com/gin-gonic/gin"
)

const ctxTokenExpKey = "x-token-exp"

// tokenDeadlineCtx returns a context that expires when the bearer token does.
// The expiry is injected by the authentication middleware so no re-parsing is needed.
// Uses c (gin.Context) as the parent so gRPC interceptors can read gin-stored values
// (user ID, roles, etc.) via ctx.Value.
func tokenDeadlineCtx(c *gin.Context) (context.Context, context.CancelFunc) {
	exp, exists := c.Get(ctxTokenExpKey)
	if !exists {
		return context.WithCancel(c)
	}

	expTime, ok := exp.(time.Time)
	if !ok {
		return context.WithCancel(c)
	}

	return context.WithDeadline(c, expTime)
}

// acceptWebSocket upgrades the HTTP connection to a WebSocket.
// It works around a gin v1.11+ incompatibility with coder/websocket: gin's
// Hijack() now rejects calls when Written() is true, but coder/websocket
// calls WriteHeaderNow() (which sets Written) before hijacking. Passing the
// unwrapped http.ResponseWriter bypasses gin's guard while keeping gin's
// status tracking intact for the request logger.
func acceptWebSocket(c *gin.Context, opts *websocket.AcceptOptions) (*websocket.Conn, error) {
	type ginUnwrapper interface {
		Unwrap() http.ResponseWriter
	}
	w := http.ResponseWriter(c.Writer)
	if uw, ok := c.Writer.(ginUnwrapper); ok {
		w = uw.Unwrap()
	}
	c.Status(http.StatusSwitchingProtocols)
	return websocket.Accept(w, c.Request, opts)
}
