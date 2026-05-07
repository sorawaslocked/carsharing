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

	userServiceLog := log.With(
		slog.String("grpcURL", cfg.GRPCServer.Client.UserServiceURL),
	)
	userServiceLog.Info("connecting to grpc server")

	userServiceGrpcConn, err := grpcconn.Connect(
		cfg.GRPCServer.Client.UserServiceURL,
		cfg.GRPCServer.Client,
	)
	if err != nil {
		userServiceLog.Error("connecting to grpc server", pkglog.Err(err))

		return nil
	}

	//if err = grpcconn.PingServer(userServiceGrpcConn); err != nil {
	//	userServiceLog.Error("pinging grpc server", pkglog.Err(err))
	//}

	userServiceGrpcClient := usersvc.NewUserServiceClient(userServiceGrpcConn)
	userServiceGrpcHandler := grpchandler.NewUserHandler(userServiceGrpcClient, log)

	jwtManager := pkgjwt.NewManager(cfg.JWT, log)

	lazy := &lazySessionCache{}
	userService := service.NewUserService(userServiceGrpcHandler, jwtManager, lazy)

	carModelServiceGrpcHandler := grpchandler.NewCarModelHandler()
	carModelService := service.NewCarModelService(carModelServiceGrpcHandler)

	carServiceGrpcHandler := grpchandler.NewCarHandler()
	carService := service.NewCarService(carServiceGrpcHandler)

	carInsuranceServiceGrpcHandler := grpchandler.NewInsuranceHandler()
	carInsuranceService := service.NewCarInsuranceService(carInsuranceServiceGrpcHandler)

	carMaintenanceServiceGrpcHandler := grpchandler.NewCarMaintenanceHandler()
	carMaintenanceService := service.NewCarMaintenanceService(carMaintenanceServiceGrpcHandler)

	zoneServiceGrpcHandler := grpchandler.NewZoneHandler()
	zoneService := service.NewZoneService(zoneServiceGrpcHandler)

	rdb, err := pkgredis.NewClient(context.Background(), &cfg.Redis, log)
	if err != nil {
		// error already logged inside NewClient
		return nil
	}

	userCache := redisadapter.NewUserCache(rdb, userService, cfg.Redis, log)
	lazy.real = userCache

	natsConn, err := pkgnats.Connect(cfg.NATS, log)
	if err != nil {
		// error already logged inside Connect
		return nil
	}

	natsUserSub := natshandler.NewUserSubscriber(natsConn, userCache, log)
	if err = natsUserSub.Subscribe(); err != nil {
		// error already logged inside Subscribe
		return nil
	}

	healthCheckers := []httphandler.HealthChecker{userServiceGrpcHandler}

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
		zoneService,
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
