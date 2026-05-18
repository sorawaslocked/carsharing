package http

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"carsharing/api-gateway/internal/adapter/http/handler"
	wshandler "carsharing/api-gateway/internal/adapter/websocket/handler"
	"carsharing/api-gateway/internal/config"
	"github.com/gin-gonic/gin"
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
	pricingRuleHandler    *handler.PricingRuleHandler
	zoneHandler           *handler.ZoneHandler
	bookingHandler        *handler.BookingHandler
	tripHandler           *handler.TripHandler
	userWsHandler         *wshandler.UserWsHandler
	carWsHandler          *wshandler.CarWsHandler
	tripWsHandler         *wshandler.TripWsHandler
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
	pricingRuleService handler.PricingRuleService,
	zoneService handler.ZoneService,
	bookingService handler.BookingService,
	tripService handler.TripService,
	tokenManager TokenParser,
	userPermissionsCache UserPermissionsCache,
	userSessionCache UserSessionCache,
	carStreamService wshandler.CarStreamService,
	tripStreamService wshandler.TripStreamService,
	documentHub *wshandler.DocumentHub,
	carStatusHub *wshandler.CarStatusHub,
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
	pricingRuleHandler := handler.NewPricingRuleHandler(pricingRuleService)
	zoneHandler := handler.NewZoneHandler(zoneService)
	bookingHandler := handler.NewBookingHandler(bookingService)
	tripHandler := handler.NewTripHandler(tripService)

	// WebSocket handlers
	userWsHandler := wshandler.NewUserWsHandler(documentHub, log)
	carWsHandler := wshandler.NewCarWsHandler(carStreamService, carStatusHub, log)
	tripWsHandler := wshandler.NewTripWsHandler(tripStreamService, log)

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
		pricingRuleHandler:    pricingRuleHandler,
		zoneHandler:           zoneHandler,
		bookingHandler:        bookingHandler,
		tripHandler:           tripHandler,
		userWsHandler:         userWsHandler,
		carWsHandler:          carWsHandler,
		tripWsHandler:         tripWsHandler,
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
