package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/sorawaslocked/car-rental-api-gateway/internal/adapter/http/dto"
	"github.com/sorawaslocked/car-rental-api-gateway/internal/model"
)

func DocumentVerificationChecker() gin.HandlerFunc {
	return func(c *gin.Context) {
		isVerified, exists := c.Get(ctxUserDocumentVerifiedKey)
		if !exists {
			dto.FromError(c, model.ErrInternalServerError)
			c.Abort()

			return
		}

		if !isVerified.(bool) {
			dto.FromError(c, model.ErrForbidden)
			c.Abort()

			return
		}

		c.Next()
	}
}

func EmailVerificationChecker() gin.HandlerFunc {
	return func(c *gin.Context) {
		isVerified, exists := c.Get(ctxUserEmailVerifiedKey)
		if !exists {
			dto.FromError(c, model.ErrInternalServerError)
			c.Abort()

			return
		}

		if !isVerified.(bool) {
			dto.FromError(c, model.ErrForbidden)
			c.Abort()

			return
		}

		c.Next()
	}
}

func SuspensionChecker() gin.HandlerFunc {
	return func(c *gin.Context) {
		isSuspended, exists := c.Get(ctxUserSuspendedKey)
		if !exists {
			dto.FromError(c, model.ErrInternalServerError)
			c.Abort()

			return
		}

		if isSuspended.(bool) {
			dto.FromError(c, model.ErrForbidden)
			c.Abort()

			return
		}

		c.Next()
	}
}
