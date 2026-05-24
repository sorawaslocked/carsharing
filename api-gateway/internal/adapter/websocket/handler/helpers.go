package handler

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
)

const ctxTokenExpKey = "x-token-exp"

// tokenDeadlineCtx returns a context that expires when the bearer token does.
// The expiry is injected by the authentication middleware so no re-parsing is needed.
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
