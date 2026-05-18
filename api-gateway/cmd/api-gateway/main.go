// @title           Car Rental API Gateway
// @version         1.0
// @description     API Gateway for the Kazakhstan carsharing platform. Handles auth, users, fleet, zones, insurance, and maintenance.
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
	"carsharing/api-gateway/internal/app"
	"carsharing/api-gateway/internal/config"
	"carsharing/shared/pkg/log"

	_ "carsharing/api-gateway/docs"
)

func main() {
	cfg := config.MustLoad()

	logger := log.SetupLogger(cfg.Env)

	application := app.New(cfg, logger)

	if application != nil {
		application.Run()
	}
}
