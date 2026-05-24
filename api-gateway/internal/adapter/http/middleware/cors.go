package middleware

import (
	"carsharing/api-gateway/internal/config"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Cors(cfg config.CorsConfig) gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowOrigins:     cfg.AllowOrigins,
		AllowMethods:     cfg.AllowMethods,
		AllowHeaders:     cfg.AllowHeaders,
		ExposeHeaders:    cfg.ExposeHeaders,
		AllowCredentials: cfg.AllowCredentials,
		MaxAge:           cfg.MaxAge,
	})
}
