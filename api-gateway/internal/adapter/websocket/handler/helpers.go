package handler

import (
	"bufio"
	"context"
	"fmt"
	"net"
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

// wsResponseWriter wraps gin's ResponseWriter and overrides Hijack() to bypass
// the Written() guard introduced in gin v1.11. coder/websocket calls
// WriteHeaderNow() before Hijack(), which sets Written=true; gin v1.11 then
// rejects Hijack(). Our Hijack() skips that check and delegates directly to
// the underlying net/http hijacker, while gin retains full ownership of headers
// and its post-handler WriteHeaderNow() becomes a no-op (already Written=true).
type wsResponseWriter struct {
	gin.ResponseWriter
}

func (w wsResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	type ginUnwrapper interface {
		Unwrap() http.ResponseWriter
	}
	if uw, ok := w.ResponseWriter.(ginUnwrapper); ok {
		if hj, ok := uw.Unwrap().(http.Hijacker); ok {
			return hj.Hijack()
		}
	}
	return nil, nil, fmt.Errorf("websocket: underlying http.ResponseWriter does not support hijacking")
}

// acceptWebSocket upgrades the HTTP connection to a WebSocket using a wrapper
// that fixes the gin v1.11 / coder/websocket incompatibility.
// InsecureSkipVerify disables coder/websocket's built-in origin check; origin
// validation is already enforced by the Gin CORS middleware earlier in the stack.
func acceptWebSocket(c *gin.Context, opts *websocket.AcceptOptions) (*websocket.Conn, error) {
	if opts == nil {
		opts = &websocket.AcceptOptions{}
	}
	opts.InsecureSkipVerify = true
	return websocket.Accept(wsResponseWriter{c.Writer}, c.Request, opts)
}
