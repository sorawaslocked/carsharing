package app

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	grpcserver "carsharing/car-service/internal/adapter/grpc"
	"carsharing/car-service/internal/adapter/grpc/handler"
	"carsharing/car-service/internal/adapter/grpc/interceptor"
	minioadapter "carsharing/car-service/internal/adapter/minio"
	natsadapter "carsharing/car-service/internal/adapter/nats"
	"carsharing/car-service/internal/adapter/postgres"
	"carsharing/car-service/internal/config"
	"carsharing/car-service/internal/service"
	"carsharing/car-service/internal/validation"
	pkggrpc "carsharing/shared/pkg/grpc"
	pkglog "carsharing/shared/pkg/log"
	pkgminio "carsharing/shared/pkg/minio"
	pkgnats "carsharing/shared/pkg/nats"
	pkgpostgres "carsharing/shared/pkg/postgres"

	"github.com/go-playground/validator/v10"
	"google.golang.org/grpc"
)

type App struct {
	log               *slog.Logger
	grpcServer        *grpc.Server
	grpcAddr          string
	closer            closer
	telematicsService *service.TelematicsService
}

func New(log *slog.Logger, cfg config.Config) (*App, error) {
	var cl closer

	pool, err := pkgpostgres.NewPool(log, cfg.PG)
	if err != nil {
		return nil, fmt.Errorf("postgres: %w", err)
	}
	cl.add(pool.Close)

	minioClient, err := pkgminio.NewClient(log, cfg.MinIO)
	if err != nil {
		cl.closeAll()
		return nil, fmt.Errorf("minio: %w", err)
	}

	ncPub, err := pkgnats.NewPublisher(log, cfg.NATSPublisher)
	if err != nil {
		cl.closeAll()
		return nil, fmt.Errorf("nats publisher: %w", err)
	}
	cl.add(ncPub.Close)

	ncSub, err := pkgnats.NewSubscriber(log, cfg.NATSSubscriber)
	if err != nil {
		cl.closeAll()
		return nil, fmt.Errorf("nats subscriber: %w", err)
	}
	cl.add(ncSub.Close)

	baseClientInterceptor := interceptor.NewClientBaseInterceptor()
	telematicsConn, err := pkggrpc.NewClientConn(
		log,
		cfg.TelematicsStream,
		grpc.WithChainUnaryInterceptor(baseClientInterceptor.Unary),
		grpc.WithChainStreamInterceptor(),
	)
	if err != nil {
		cl.closeAll()
		return nil, fmt.Errorf("telematics stream client: %w", err)
	}
	cl.add(func() { _ = telematicsConn.Close() })

	validate := validator.New()
	if err = validation.RegisterCustomValidators(validate, log); err != nil {
		cl.closeAll()
		return nil, fmt.Errorf("register validators: %w", err)
	}

	objectStorage := minioadapter.NewObjectStorage(minioClient, cfg.MinIO.Bucket)
	natsPublisher := natsadapter.NewPublisher(ncPub)

	carModelRepo := postgres.NewCarModelRepository(pool, log)
	carRepo := postgres.NewCarRepository(pool, log)
	carInsuranceRepo := postgres.NewCarInsuranceRepository(pool, log)
	templateRepo := postgres.NewCarMaintenanceTemplateRepository(pool, log)
	recordRepo := postgres.NewCarMaintenanceRecordRepository(pool, log)
	serviceStateRepo := postgres.NewCarServiceStateRepository(pool, log)
	statusLogRepo := postgres.NewCarStatusLogRepository(pool, log)
	telematicsRepo := postgres.NewTelematicsRepository(pool, log)
	zoneRepo := postgres.NewZoneRepository(pool, log)

	telematicsStreamClient := grpcserver.NewTelematicsStreamClient(telematicsConn, log)
	telematicsService := service.NewTelematicsService(telematicsStreamClient, telematicsRepo, carRepo, log)

	carModelService := service.NewCarModelService(carModelRepo, objectStorage, validate, log)
	carService := service.NewCarService(carModelRepo, carRepo, statusLogRepo, telematicsRepo, objectStorage, natsPublisher, validate, log)
	carService.SetCarCreatedNotifier(telematicsService)

	carInsuranceService := service.NewCarInsuranceService(carInsuranceRepo, objectStorage, validate, log)
	carMaintenanceService := service.NewCarMaintenanceService(templateRepo, recordRepo, serviceStateRepo, carRepo, carService, objectStorage, validate, log)
	zoneService := service.NewZoneService(zoneRepo, validate, log)

	natsSubscriber := natsadapter.NewSubscriber(ncSub, carService, log)
	if err = natsSubscriber.Subscribe(); err != nil {
		cl.closeAll()
		return nil, fmt.Errorf("nats subscribe: %w", err)
	}

	healthHandler := handler.NewHealthHandler(log, map[string]handler.Pinger{
		"postgres":        postgres.NewPinger(log, pool),
		"nats-publisher":  natsadapter.NewPinger(log, ncPub),
		"nats-subscriber": natsadapter.NewPinger(log, ncSub),
		"minio":           minioadapter.NewPinger(log, minioClient),
	})

	grpcSrv, err := grpcserver.NewServer(
		log,
		cfg.GRPC,
		carModelService, carService, carInsuranceService, carMaintenanceService, zoneService,
		telematicsService,
		healthHandler,
	)
	if err != nil {
		cl.closeAll()
		return nil, fmt.Errorf("grpc server: %w", err)
	}

	return &App{
		log:               pkglog.WithComponent(log, "app"),
		grpcServer:        grpcSrv,
		grpcAddr:          fmt.Sprintf("%s:%d", cfg.GRPC.Host, cfg.GRPC.Port),
		closer:            cl,
		telematicsService: telematicsService,
	}, nil
}

func (a *App) Run() error {
	ctx, cancel := context.WithCancel(context.Background())

	if err := a.telematicsService.Start(ctx); err != nil {
		cancel()
		return fmt.Errorf("telematics service start: %w", err)
	}

	lis, err := net.Listen("tcp", a.grpcAddr)
	if err != nil {
		cancel()
		return fmt.Errorf("gRPC listen: %w", err)
	}
	a.log.Info("gRPC server listening", slog.String("addr", a.grpcAddr))

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	errCh := make(chan error, 1)
	go func() {
		errCh <- a.grpcServer.Serve(lis)
	}()

	select {
	case sig := <-quit:
		a.log.Info("received signal, shutting down", slog.String("signal", sig.String()))
	case err := <-errCh:
		cancel()
		return fmt.Errorf("gRPC server stopped unexpectedly: %w", err)
	}

	cancel()
	a.stop()

	return nil
}

func (a *App) stop() {
	a.log.Info("shutting down")

	stopped := make(chan struct{})
	go func() {
		a.grpcServer.GracefulStop()
		close(stopped)
	}()

	select {
	case <-stopped:
		a.log.Info("gRPC server stopped gracefully")
	case <-time.After(15 * time.Second):
		a.log.Warn("graceful stop timed out, forcing stop")
		a.grpcServer.Stop()
	}

	a.telematicsService.Stop()
	a.closer.closeAll()
}
