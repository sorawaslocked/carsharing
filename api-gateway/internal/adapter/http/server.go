package http

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sorawaslocked/car-rental-api-gateway/internal/adapter/http/handler"
	"github.com/sorawaslocked/car-rental-api-gateway/internal/config"
)

type Server struct {
	router                *gin.Engine
	httpCfg               config.HTTPServer
	log                   *slog.Logger
	healthHandler         *handler.HealthHandler
	userHandler           *handler.UserHandler
	carModelHandler       *handler.CarModelHandler
	carHandler            *handler.CarHandler
	carInsuranceHandler   *handler.CarInsuranceHandler
	carMaintenanceHandler *handler.CarMaintenanceHandler
	zoneHandler           *handler.ZoneHandler
}

func New(
	httpCfg config.HTTPServer,
	cookieCfg config.Cookie,
	log *slog.Logger,
	healthCheckers []handler.HealthChecker,
	userService handler.UserService,
	carModelService handler.CarModelService,
	carService handler.CarService,
	carInsuranceService handler.CarInsuranceService,
	carMaintenanceService handler.CarMaintenanceService,
	zoneService handler.ZoneService,
	tokenManager TokenParser,
	userPermissionsCache UserPermissionsCache,
	userSessionCache UserSessionCache,
) *Server {
	httpLog := log.With(
		slog.String("httpServerHost", httpCfg.Host),
		slog.Int("httpServerPort", httpCfg.Port),
	)

	gin.SetMode(httpCfg.GinMode)
	router := gin.New()

	// Handlers
	userHandler := handler.NewUser(userService, cookieCfg)
	healthHandler := handler.NewHealthHandler(healthCheckers)
	carModelHandler := handler.NewCarModelHandler(carModelService)
	carHandler := handler.NewCarHandler(carService)
	carInsuranceHandler := handler.NewCarInsuranceHandler(carInsuranceService)
	carMaintenanceHandler := handler.NewCarMaintenanceHandler(carMaintenanceService)
	zoneHandler := handler.NewZoneHandler(zoneService)

	server := &Server{
		router:                router,
		httpCfg:               httpCfg,
		log:                   httpLog,
		userHandler:           userHandler,
		healthHandler:         healthHandler,
		carModelHandler:       carModelHandler,
		carHandler:            carHandler,
		carInsuranceHandler:   carInsuranceHandler,
		carMaintenanceHandler: carMaintenanceHandler,
		zoneHandler:           zoneHandler,
	}

	server.setupMiddleware()
	server.setupRoutes(tokenManager, userPermissionsCache, userSessionCache)

	return server
}

func (s *Server) MustRun() {
	go func() {
		addr := fmt.Sprintf("%s:%d", s.httpCfg.Host, s.httpCfg.Port)

		s.log.Info("starting http server")
		err := s.router.Run(addr)
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			panic(err)
		}
	}()
}

func (s *Server) Stop(ctx context.Context) {
	s.log.Info("shutting down http server")
	<-ctx.Done()
}
