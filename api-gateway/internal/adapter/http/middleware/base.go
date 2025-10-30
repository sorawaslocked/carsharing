package middleware

import (
	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
)

func Base() gin.HandlerFunc {
	return func(c *gin.Context) {
		if authHeader := c.GetHeader("Authorization"); authHeader != "" {
			c.Set("authorization", authHeader)
		}

		c.Set("x-request-id", requestid.Get(c))
		c.Set("x-client-ip", c.ClientIP())

		c.Next()
	}
}
