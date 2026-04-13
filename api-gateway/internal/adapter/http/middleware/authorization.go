package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/sorawaslocked/car-rental-api-gateway/internal/adapter/http/dto"
	"github.com/sorawaslocked/car-rental-api-gateway/internal/model"
)

func VerificationChecker() gin.HandlerFunc {
	return func(c *gin.Context) {
		verified, exists := c.Get(ctxUserVerifiedKey)
		if !exists {
			dto.FromError(c, model.ErrInternalServerError)

			return
		}

		if !verified.(bool) {
			dto.FromError(c, model.ErrForbidden)

			return
		}

		c.Next()
	}
}

func SuspensionChecker() gin.HandlerFunc {
	return func(c *gin.Context) {
		suspended, exists := c.Get(ctxUserSuspendedKey)
		if !exists {
			dto.FromError(c, model.ErrInternalServerError)

			return
		}

		if suspended.(bool) {
			dto.FromError(c, model.ErrForbidden)

			return
		}

		c.Next()
	}
}
