package app

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	grpchandler "github.com/sorawaslocked/car-rental-api-gateway/internal/adapter/grpc/handler"
	httpserver "github.com/sorawaslocked/car-rental-api-gateway/internal/adapter/http"
	httphandler "github.com/sorawaslocked/car-rental-api-gateway/internal/adapter/http/handler"
	natshandler "github.com/sorawaslocked/car-rental-api-gateway/internal/adapter/nats/handler"
	redisadapter "github.com/sorawaslocked/car-rental-api-gateway/internal/adapter/redis"
	"github.com/sorawaslocked/car-rental-api-gateway/internal/config"
	grpcconn "github.com/sorawaslocked/car-rental-api-gateway/internal/pkg/grpc"
	pkgjwt "github.com/sorawaslocked/car-rental-api-gateway/internal/pkg/jwt"
	pkglog "github.com/sorawaslocked/car-rental-api-gateway/internal/pkg/log"
	pkgnats "github.com/sorawaslocked/car-rental-api-gateway/internal/pkg/nats"
	pkgredis "github.com/sorawaslocked/car-rental-api-gateway/internal/pkg/redis"
	"github.com/sorawaslocked/car-rental-api-gateway/internal/service"
	bookingsvc "github.com/sorawaslocked/car-rental-protos/gen/service/booking"
	carsvc "github.com/sorawaslocked/car-rental-protos/gen/service/car"
	tripsvc "github.com/sorawaslocked/car-rental-protos/gen/service/trip"
	usersvc "github.com/sorawaslocked/car-rental-protos/gen/service/user"
)

// lazySessionCache breaks the init cycle between UserService and UserCache:
// UserService needs UserCache as its session cache, and UserCache needs UserService
// as its user provider. Both are constructed first; real is wired after.
type lazySessionCache struct {
	real service.UserSessionCache
}

func (l *lazySessionCache) IsSignedIn(ctx context.Context, userID, deviceID string) (bool, error) {
	return l.real.IsSignedIn(ctx, userID, deviceID)
}

func (l *lazySessionCache) SetSignedIn(ctx context.Context, userID, deviceID string, v bool) error {
	return l.real.SetSignedIn(ctx, userID, deviceID, v)
}

type App struct {
	cfg            config.Config
	log            *slog.Logger
	httpServer     *httpserver.Server
	natsSubscriber *natshandler.UserSubscriber
}

func New(cfg config.Config, log *slog.Logger) *App {
	log = log.With(slog.String("appID", "api-gateway"))

	// User service gRPC connection
	userServiceLog := log.With(slog.String("grpcURL", cfg.GRPCServer.Client.UserServiceURL))
	userServiceLog.Info("connecting to grpc server")

	userServiceGrpcConn, err := grpcconn.Connect(cfg.GRPCServer.Client.UserServiceURL, cfg.GRPCServer.Client)
	if err != nil {
		userServiceLog.Error("connecting to grpc server", pkglog.Err(err))
		return nil
	}

	// Car service gRPC connection
	carServiceLog := log.With(slog.String("grpcURL", cfg.GRPCServer.Client.CarServiceURL))
	carServiceLog.Info("connecting to grpc server")

	carServiceGrpcConn, err := grpcconn.Connect(cfg.GRPCServer.Client.CarServiceURL, cfg.GRPCServer.Client)
	if err != nil {
		carServiceLog.Error("connecting to grpc server", pkglog.Err(err))
		return nil
	}

	// Booking service gRPC connection
	bookingServiceLog := log.With(slog.String("grpcURL", cfg.GRPCServer.Client.BookingServiceURL))
	bookingServiceLog.Info("connecting to grpc server")

	bookingServiceGrpcConn, err := grpcconn.Connect(cfg.GRPCServer.Client.BookingServiceURL, cfg.GRPCServer.Client)
	if err != nil {
		bookingServiceLog.Error("connecting to grpc server", pkglog.Err(err))
		return nil
	}

	// Trip service gRPC connection
	tripServiceLog := log.With(slog.String("grpcURL", cfg.GRPCServer.Client.TripServiceURL))
	tripServiceLog.Info("connecting to grpc server")

	tripServiceGrpcConn, err := grpcconn.Connect(cfg.GRPCServer.Client.TripServiceURL, cfg.GRPCServer.Client)
	if err != nil {
		tripServiceLog.Error("connecting to grpc server", pkglog.Err(err))
		return nil
	}

	// gRPC clients
	userGrpcClient := usersvc.NewUserServiceClient(userServiceGrpcConn)
	userHealthGrpcClient := usersvc.NewHealthServiceClient(userServiceGrpcConn)
	carGrpcClient := carsvc.NewCarServiceClient(carServiceGrpcConn)
	carHealthGrpcClient := carsvc.NewHealthServiceClient(carServiceGrpcConn)
	carModelGrpcClient := carsvc.NewCarModelServiceClient(carServiceGrpcConn)
	carInsuranceGrpcClient := carsvc.NewCarInsuranceServiceClient(carServiceGrpcConn)
	carMaintenanceGrpcClient := carsvc.NewCarMaintenanceServiceClient(carServiceGrpcConn)
	zoneGrpcClient := carsvc.NewZoneServiceClient(carServiceGrpcConn)
	pricingRuleGrpcClient := bookingsvc.NewPricingRuleServiceClient(bookingServiceGrpcConn)
	bookingGrpcClient := bookingsvc.NewBookingServiceClient(bookingServiceGrpcConn)
	bookingHealthGrpcClient := bookingsvc.NewHealthServiceClient(bookingServiceGrpcConn)
	tripGrpcClient := tripsvc.NewTripServiceClient(tripServiceGrpcConn)
	tripHealthGrpcClient := tripsvc.NewHealthServiceClient(tripServiceGrpcConn)

	// gRPC handlers
	userServiceGrpcHandler := grpchandler.NewUserHandler(userGrpcClient, log)
	userHealthGrpcHandler := grpchandler.NewHealthHandler(userHealthGrpcClient, log)
	carGrpcHandler := grpchandler.NewCarHandler(carGrpcClient, log)
	carHealthGrpcHandler := grpchandler.NewHealthHandler(carHealthGrpcClient, log)
	carModelGrpcHandler := grpchandler.NewCarModelHandler(carModelGrpcClient, log)
	carInsuranceGrpcHandler := grpchandler.NewCarInsuranceHandler(carInsuranceGrpcClient, log)
	carMaintenanceGrpcHandler := grpchandler.NewCarMaintenanceHandler(carMaintenanceGrpcClient, log)
	zoneGrpcHandler := grpchandler.NewZoneHandler(zoneGrpcClient, log)
	pricingRuleGrpcHandler := grpchandler.NewPricingRuleHandler(pricingRuleGrpcClient, log)
	bookingGrpcHandler := grpchandler.NewBookingHandler(bookingGrpcClient, log)
	bookingHealthGrpcHandler := grpchandler.NewHealthHandler(bookingHealthGrpcClient, log)
	tripGrpcHandler := grpchandler.NewTripHandler(tripGrpcClient, log)
	tripHealthGrpcHandler := grpchandler.NewHealthHandler(tripHealthGrpcClient, log)

	// JWT
	jwtManager := pkgjwt.NewManager(cfg.JWT, log)

	// Services
	lazy := &lazySessionCache{}
	userService := service.NewUserService(userServiceGrpcHandler, jwtManager, lazy)
	carService := service.NewCarService(carGrpcHandler)
	carModelService := service.NewCarModelService(carModelGrpcHandler)
	carInsuranceService := service.NewCarInsuranceService(carInsuranceGrpcHandler)
	carMaintenanceService := service.NewCarMaintenanceService(carMaintenanceGrpcHandler)
	zoneService := service.NewZoneService(zoneGrpcHandler)
	pricingRuleService := service.NewPricingRuleService(pricingRuleGrpcHandler)
	bookingService := service.NewBookingService(bookingGrpcHandler)
	tripService := service.NewTripService(tripGrpcHandler)

	// Redis
	rdb, err := pkgredis.NewClient(context.Background(), &cfg.Redis, log)
	if err != nil {
		return nil
	}

	userCache := redisadapter.NewUserCache(rdb, userService, cfg.Redis, log)
	lazy.real = userCache

	// NATS
	natsConn, err := pkgnats.Connect(cfg.NATS, log)
	if err != nil {
		return nil
	}

	natsUserSub := natshandler.NewUserSubscriber(natsConn, userCache, log)
	if err = natsUserSub.Subscribe(); err != nil {
		return nil
	}

	healthCheckers := []httphandler.HealthChecker{
		userHealthGrpcHandler,
		carHealthGrpcHandler,
		bookingHealthGrpcHandler,
		tripHealthGrpcHandler,
	}

	httpServer := httpserver.New(
		cfg.HTTPServer,
		cfg.HTTPServer.Cookie,
		log,
		healthCheckers,
		userService,
		carModelService,
		carService,
		carInsuranceService,
		carMaintenanceService,
		pricingRuleService,
		zoneService,
		bookingService,
		tripService,
		jwtManager,
		userCache,
		userCache,
	)

	return &App{
		cfg:            cfg,
		log:            log,
		httpServer:     httpServer,
		natsSubscriber: natsUserSub,
	}
}

func (a *App) Run() {
	a.httpServer.MustRun()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	s := <-shutdown

	a.log.Info("received shutdown signal", slog.String("signal", s.String()))

	a.Stop()

	a.log.Info("graceful shutdown complete")
}

func (a *App) Stop() {
	a.natsSubscriber.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	a.httpServer.Stop(ctx)
}
