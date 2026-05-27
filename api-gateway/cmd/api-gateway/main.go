// @title           Carsharing API Gateway
// @version         1.0
// @description     API Gateway for the Kazakhstan carsharing platform.
// @description     Provides REST and WebSocket endpoints for: authentication, user management,
// @description     fleet (cars, models, telemetry, status), bookings, trips, zones,
// @description     pricing rules, car insurance, and maintenance. Real-time updates are
// @description     available via WebSocket at /api/v1/ws.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.email  support@carrental.kz

// @license.name  Proprietary

// @host      localhost:4000
// @BasePath  /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description  Enter: Bearer <access_token>
package main

import (
	"log/slog"
	"os"

	"carsharing/api-gateway/docs"
	"carsharing/api-gateway/internal/app"
	"carsharing/api-gateway/internal/config"
	"carsharing/shared/pkg/log"
)

func main() {
	cfg := config.MustLoad()

	docs.SwaggerInfo.Host = cfg.HTTPServer.SwaggerHost

	logger := log.SetupLogger(cfg.Env)

	application, err := app.New(cfg, logger)
	if err != nil {
		logger.Error("initialising application", slog.String("error", err.Error()))
		os.Exit(1)
	}

	application.Run()
}
