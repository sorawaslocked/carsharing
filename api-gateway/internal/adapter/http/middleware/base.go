package middleware

import (
	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
)

const (
	ctxRequestIDKey = "x-request-id"
	ctxClientIPKey  = "x-client-ip"
	ctxDeviceIDKey  = "x-device-id"
)

func Base() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(ctxRequestIDKey, requestid.Get(c))
		c.Set(ctxClientIPKey, c.ClientIP())

		deviceID := c.GetHeader("x-device-id")
		if deviceID != "" {
			c.Set(ctxDeviceIDKey, deviceID)
		}

		c.Next()
	}
}
