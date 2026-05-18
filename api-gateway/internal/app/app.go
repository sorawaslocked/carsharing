package app

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	grpchandler "carsharing/api-gateway/internal/adapter/grpc/handler"
	"carsharing/api-gateway/internal/adapter/grpc/interceptor"
	httpserver "carsharing/api-gateway/internal/adapter/http"
	httphandler "carsharing/api-gateway/internal/adapter/http/handler"
	natshandler "carsharing/api-gateway/internal/adapter/nats/handler"
	redisadapter "carsharing/api-gateway/internal/adapter/redis"
	wshandler "carsharing/api-gateway/internal/adapter/websocket/handler"
	"carsharing/api-gateway/internal/config"
	"carsharing/api-gateway/internal/service"
	pkggrpc "carsharing/shared/pkg/grpc"
	pkgjwt "carsharing/shared/pkg/jwt"
	pkglog "carsharing/shared/pkg/log"
	pkgnats "carsharing/shared/pkg/nats"
	pkgredis "carsharing/shared/pkg/redis"
	bookingsvc "github.com/sorawaslocked/car-rental-protos/gen/service/booking"
	carsvc "github.com/sorawaslocked/car-rental-protos/gen/service/car"
	tripsvc "github.com/sorawaslocked/car-rental-protos/gen/service/trip"
	usersvc "github.com/sorawaslocked/car-rental-protos/gen/service/user"
	"google.golang.org/grpc"
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
	cfg                    config.Config
	log                    *slog.Logger
	httpServer             *httpserver.Server
	natsUserSubscriber     *natshandler.UserSubscriber
	natsDocumentSubscriber *natshandler.DocumentSubscriber
	natsCarSubscriber      *natshandler.CarSubscriber
}

func New(cfg config.Config, log *slog.Logger) *App {
	log = log.With(slog.String("appID", "api-gateway"))

	baseInterceptor := interceptor.NewBaseInterceptor()
	unaryOpt := grpc.WithUnaryInterceptor(baseInterceptor.Unary)
	streamOpt := grpc.WithStreamInterceptor(baseInterceptor.Stream)

	// User service gRPC connection
	userServiceLog := log.With(slog.String("grpcService", "user-service"))
	userServiceLog.Info("connecting to grpc server")

	userServiceGrpcConn, err := pkggrpc.NewClientConn(userServiceLog, cfg.GRPCServer.UserService, unaryOpt, streamOpt)
	if err != nil {
		userServiceLog.Error("connecting to grpc server", pkglog.Err(err))
		return nil
	}

	// Car service gRPC connection
	carServiceLog := log.With(slog.String("grpcService", "car-service"))
	carServiceLog.Info("connecting to grpc server")

	carServiceGrpcConn, err := pkggrpc.NewClientConn(carServiceLog, cfg.GRPCServer.CarService, unaryOpt, streamOpt)
	if err != nil {
		carServiceLog.Error("connecting to grpc server", pkglog.Err(err))
		return nil
	}

	// Booking service gRPC connection
	bookingServiceLog := log.With(slog.String("grpcService", "booking-service"))
	bookingServiceLog.Info("connecting to grpc server")

	bookingServiceGrpcConn, err := pkggrpc.NewClientConn(bookingServiceLog, cfg.GRPCServer.BookingService, unaryOpt, streamOpt)
	if err != nil {
		bookingServiceLog.Error("connecting to grpc server", pkglog.Err(err))
		return nil
	}

	// Trip service gRPC connection
	tripServiceLog := log.With(slog.String("grpcService", "trip-service"))
	tripServiceLog.Info("connecting to grpc server")

	tripServiceGrpcConn, err := pkggrpc.NewClientConn(tripServiceLog, cfg.GRPCServer.TripService, unaryOpt, streamOpt)
	if err != nil {
		tripServiceLog.Error("connecting to grpc server", pkglog.Err(err))
		return nil
	}

	// gRPC clients
	userGrpcClient := usersvc.NewUserServiceClient(userServiceGrpcConn)
	userHealthGrpcClient := usersvc.NewHealthServiceClient(userServiceGrpcConn)
	carGrpcClient := carsvc.NewCarServiceClient(carServiceGrpcConn)
	carStreamGrpcClient := carsvc.NewCarStreamServiceClient(carServiceGrpcConn)
	carHealthGrpcClient := carsvc.NewHealthServiceClient(carServiceGrpcConn)
	carModelGrpcClient := carsvc.NewCarModelServiceClient(carServiceGrpcConn)
	carInsuranceGrpcClient := carsvc.NewCarInsuranceServiceClient(carServiceGrpcConn)
	carMaintenanceGrpcClient := carsvc.NewCarMaintenanceServiceClient(carServiceGrpcConn)
	zoneGrpcClient := carsvc.NewZoneServiceClient(carServiceGrpcConn)
	pricingRuleGrpcClient := bookingsvc.NewPricingRuleServiceClient(bookingServiceGrpcConn)
	bookingGrpcClient := bookingsvc.NewBookingServiceClient(bookingServiceGrpcConn)
	bookingHealthGrpcClient := bookingsvc.NewHealthServiceClient(bookingServiceGrpcConn)
	tripGrpcClient := tripsvc.NewTripServiceClient(tripServiceGrpcConn)
	tripStreamGrpcClient := tripsvc.NewTripStreamServiceClient(tripServiceGrpcConn)
	tripHealthGrpcClient := tripsvc.NewHealthServiceClient(tripServiceGrpcConn)

	// gRPC handlers
	userServiceGrpcHandler := grpchandler.NewUserHandler(userGrpcClient, log)
	userHealthGrpcHandler := grpchandler.NewHealthHandler(userHealthGrpcClient, log)
	carGrpcHandler := grpchandler.NewCarHandler(carGrpcClient, carStreamGrpcClient, log)
	carHealthGrpcHandler := grpchandler.NewHealthHandler(carHealthGrpcClient, log)
	carModelGrpcHandler := grpchandler.NewCarModelHandler(carModelGrpcClient, log)
	carInsuranceGrpcHandler := grpchandler.NewCarInsuranceHandler(carInsuranceGrpcClient, log)
	carMaintenanceGrpcHandler := grpchandler.NewCarMaintenanceHandler(carMaintenanceGrpcClient, log)
	zoneGrpcHandler := grpchandler.NewZoneHandler(zoneGrpcClient, log)
	pricingRuleGrpcHandler := grpchandler.NewPricingRuleHandler(pricingRuleGrpcClient, log)
	bookingGrpcHandler := grpchandler.NewBookingHandler(bookingGrpcClient, log)
	bookingHealthGrpcHandler := grpchandler.NewHealthHandler(bookingHealthGrpcClient, log)
	tripGrpcHandler := grpchandler.NewTripHandler(tripGrpcClient, tripStreamGrpcClient, log)
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
	rdb, err := pkgredis.NewClient(log, cfg.Redis)
	if err != nil {
		return nil
	}

	userCache := redisadapter.NewUserCache(rdb, userService, cfg.Cache, log)
	lazy.real = userCache

	// NATS
	natsConn, err := pkgnats.NewSubscriber(log, cfg.NATS)
	if err != nil {
		return nil
	}

	natsUserSub := natshandler.NewUserSubscriber(natsConn, userCache, log)
	if err = natsUserSub.Subscribe(); err != nil {
		return nil
	}

	documentHub := wshandler.NewDocumentHub()
	carStatusHub := wshandler.NewCarStatusHub()

	natsDocumentSub := natshandler.NewDocumentSubscriber(natsConn, documentHub, log)
	if err = natsDocumentSub.Subscribe(); err != nil {
		return nil
	}

	natsCarSub := natshandler.NewCarSubscriber(natsConn, carStatusHub, log)
	if err = natsCarSub.Subscribe(); err != nil {
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
		carService,
		tripService,
		documentHub,
		carStatusHub,
	)

	return &App{
		cfg:                    cfg,
		log:                    log,
		httpServer:             httpServer,
		natsUserSubscriber:     natsUserSub,
		natsDocumentSubscriber: natsDocumentSub,
		natsCarSubscriber:      natsCarSub,
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
	a.natsUserSubscriber.Close()
	a.natsDocumentSubscriber.Close()
	a.natsCarSubscriber.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	a.httpServer.Stop(ctx)
}
