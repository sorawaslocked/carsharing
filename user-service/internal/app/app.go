package app

import (
	"fmt"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	pkggrpc "carsharing/shared/pkg/grpc"
	pkglog "carsharing/shared/pkg/log"
	pkgminio "carsharing/shared/pkg/minio"
	pkgnats "carsharing/shared/pkg/nats"
	pkgpostgres "carsharing/shared/pkg/postgres"
	pkgredis "carsharing/shared/pkg/redis"
	"carsharing/user-service/internal/adapter/brevo"
	grpcserver "carsharing/user-service/internal/adapter/grpc"
	grpcdocanalyzer "carsharing/user-service/internal/adapter/grpc/client"
	"carsharing/user-service/internal/adapter/grpc/handler"
	"carsharing/user-service/internal/adapter/grpc/interceptor"
	minioadapter "carsharing/user-service/internal/adapter/minio"
	natsadapter "carsharing/user-service/internal/adapter/nats"
	natshandler "carsharing/user-service/internal/adapter/nats/handler"
	"carsharing/user-service/internal/adapter/postgres"
	redisadapter "carsharing/user-service/internal/adapter/redis"
	"carsharing/user-service/internal/config"
	"carsharing/user-service/internal/service"
	"carsharing/user-service/internal/validation"

	"github.com/go-playground/validator/v10"
	"google.golang.org/grpc"
)

type App struct {
	log        *slog.Logger
	grpcServer *grpc.Server
	grpcAddr   string
	closer     closer
}

func New(log *slog.Logger, cfg config.Config) (*App, error) {
	var cl closer

	pool, err := pkgpostgres.NewPool(log, cfg.Postgres)
	if err != nil {
		return nil, fmt.Errorf("postgres: %w", err)
	}
	cl.add(pool.Close)

	validate := validator.New()
	if err := validation.RegisterCustomValidators(validate, log); err != nil {
		cl.closeAll()
		return nil, fmt.Errorf("validator: %w", err)
	}

	natsConn, err := pkgnats.NewPublisher(log, cfg.NATS)
	if err != nil {
		cl.closeAll()
		return nil, fmt.Errorf("nats: %w", err)
	}
	cl.add(natsConn.Close)

	baseClientInterceptor := interceptor.NewClientBaseInterceptor()
	analyzerConn, err := pkggrpc.NewClientConn(
		log,
		cfg.DocumentAnalyzer,
		grpc.WithChainUnaryInterceptor(baseClientInterceptor.Unary),
		grpc.WithChainStreamInterceptor(),
	)
	if err != nil {
		cl.closeAll()
		return nil, fmt.Errorf("document analyzer client: %w", err)
	}
	cl.add(func() { _ = analyzerConn.Close() })

	redisClient, err := pkgredis.NewClient(log, cfg.Redis)
	if err != nil {
		cl.closeAll()
		return nil, fmt.Errorf("redis: %w", err)
	}
	cl.add(func() { _ = redisClient.Close() })

	minioClient, err := pkgminio.NewClient(log, cfg.Minio)
	if err != nil {
		cl.closeAll()
		return nil, fmt.Errorf("minio: %w", err)
	}

	activationCodeCache := redisadapter.NewActivationCodeCache(log, redisClient)
	msMailer := brevo.New(log, cfg.Brevo)
	userRepo := postgres.NewUserRepository(log, pool)
	docRepo := postgres.NewDocumentRepository(log, pool)
	publisher := natsadapter.NewPublisher(log, natsConn)
	minioStorage := minioadapter.NewMinioObjectStorage(log, minioClient, cfg.Minio)
	analyzerClient := grpcdocanalyzer.NewDocumentAnalyzer(log, analyzerConn)

	userService := service.NewUserService(log, validate, userRepo, docRepo, minioStorage, analyzerClient, publisher, activationCodeCache, msMailer)

	docHandler := natshandler.NewDocumentHandler(log, natsConn, userService)
	if err := docHandler.Subscribe(); err != nil {
		cl.closeAll()
		return nil, fmt.Errorf("nats subscribe: %w", err)
	}

	healthHandler := handler.NewHealthHandler(log, map[string]handler.Pinger{
		"postgres":          postgres.NewPinger(log, pool),
		"redis":             redisadapter.NewPinger(log, redisClient),
		"nats":              natsadapter.NewPinger(log, natsConn),
		"minio":             minioadapter.NewPinger(log, minioClient),
		"document-analyzer": grpcdocanalyzer.NewPinger(log, analyzerConn),
	})

	grpcSrv, err := grpcserver.NewServer(log, cfg.GRPC, userService, healthHandler)
	if err != nil {
		cl.closeAll()
		return nil, fmt.Errorf("grpc server: %w", err)
	}

	return &App{
		log:        pkglog.WithComponent(log, "app"),
		grpcServer: grpcSrv,
		grpcAddr:   fmt.Sprintf("%s:%d", cfg.GRPC.Host, cfg.GRPC.Port),
		closer:     cl,
	}, nil
}

func (a *App) Run() error {
	lis, err := net.Listen("tcp", a.grpcAddr)
	if err != nil {
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
		return fmt.Errorf("gRPC server stopped unexpectedly: %w", err)
	}

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

	a.closer.closeAll()
}
