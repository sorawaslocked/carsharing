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

	pkggrpc "carsharing/shared/pkg/grpc"
	pkglog "carsharing/shared/pkg/log"
	pkgminio "carsharing/shared/pkg/minio"
	pkgnats "carsharing/shared/pkg/nats"
	pkgpostgres "carsharing/shared/pkg/postgres"
	pkgredis "carsharing/shared/pkg/redis"
	grpcserver "carsharing/user-service/internal/adapter/grpc"
	grpcdocanalyzer "carsharing/user-service/internal/adapter/grpc/client"
	"carsharing/user-service/internal/adapter/grpc/handler"
	"carsharing/user-service/internal/adapter/grpc/interceptor"
	"carsharing/user-service/internal/adapter/mailer"
	minioadapter "carsharing/user-service/internal/adapter/minio"
	natsadapter "carsharing/user-service/internal/adapter/nats"
	natshandler "carsharing/user-service/internal/adapter/nats/handler"
	"carsharing/user-service/internal/adapter/postgres"
	redisadapter "carsharing/user-service/internal/adapter/redis"
	"carsharing/user-service/internal/config"
	validatecfg "carsharing/user-service/internal/pkg/validate"
	"carsharing/user-service/internal/service"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgxpool"
	natsio "github.com/nats-io/nats.go"
	goredis "github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
)

type App struct {
	log          *slog.Logger
	grpcServer   *grpc.Server
	grpcAddr     string
	pool         *pgxpool.Pool
	natsConn     *natsio.Conn
	analyzerConn *grpc.ClientConn
	redisClient  *goredis.Client
}

func New(log *slog.Logger, cfg config.Config) (*App, error) {
	pool, err := pkgpostgres.NewPool(log, cfg.Postgres)
	if err != nil {
		return nil, fmt.Errorf("postgres: %w", err)
	}

	validate := validator.New()
	if err := validate.RegisterValidation("min_age", validatecfg.MinAge); err != nil {
		pool.Close()
		return nil, fmt.Errorf("validator: %w", err)
	}
	if err := validate.RegisterValidation("complex_password", validatecfg.ComplexPassword); err != nil {
		pool.Close()
		return nil, fmt.Errorf("validator: %w", err)
	}

	natsConn, err := pkgnats.NewPublisher(log, cfg.NATS)
	if err != nil {
		pool.Close()
		return nil, fmt.Errorf("nats: %w", err)
	}

	baseClientInterceptor := interceptor.NewClientBaseInterceptor()
	analyzerConn, err := pkggrpc.NewClientConn(
		log,
		cfg.DocumentAnalyzer,
		grpc.WithChainUnaryInterceptor(baseClientInterceptor.Unary),
		grpc.WithChainStreamInterceptor(),
	)
	if err != nil {
		pool.Close()
		natsConn.Close()
		return nil, fmt.Errorf("document analyzer client: %w", err)
	}

	redisClient, err := pkgredis.NewClient(log, cfg.Redis)
	if err != nil {
		pool.Close()
		natsConn.Close()
		_ = analyzerConn.Close()
		return nil, fmt.Errorf("redis: %w", err)
	}

	minioClient, err := pkgminio.NewClient(log, cfg.Minio)
	if err != nil {
		pool.Close()
		natsConn.Close()
		_ = analyzerConn.Close()
		_ = redisClient.Close()
		return nil, fmt.Errorf("minio: %w", err)
	}
	if err := pkgminio.EnsureBucket(context.Background(), log, minioClient, cfg.Minio); err != nil {
		pool.Close()
		natsConn.Close()
		_ = analyzerConn.Close()
		_ = redisClient.Close()
		return nil, fmt.Errorf("minio bucket: %w", err)
	}

	activationCodeCache := redisadapter.NewActivationCodeRedisCache(redisClient)
	msMailer := mailer.New(log, cfg.Brevo)
	userRepo := postgres.NewUserRepository(log, pool)
	docRepo := postgres.NewDocumentRepository(log, pool)
	publisher := natsadapter.NewPublisher(log, natsConn)
	minioStorage := minioadapter.NewMinioObjectStorage(log, minioClient, cfg.Minio)
	analyzerClient := grpcdocanalyzer.NewDocumentAnalyzer(log, analyzerConn)

	userService := service.NewUserService(log, validate, userRepo, docRepo, minioStorage, analyzerClient, publisher, activationCodeCache, msMailer)

	docHandler := natshandler.NewDocumentHandler(log, natsConn, userService)
	if err := docHandler.Subscribe(); err != nil {
		pool.Close()
		natsConn.Close()
		_ = analyzerConn.Close()
		_ = redisClient.Close()
		return nil, fmt.Errorf("nats subscribe: %w", err)
	}

	healthHandler := handler.NewHealthHandler(log, map[string]handler.Pinger{
		"postgres":          postgres.NewChecker(log, pool),
		"redis":             redisadapter.NewChecker(redisClient),
		"nats":              natsadapter.NewChecker(natsConn),
		"minio":             minioadapter.NewChecker(minioClient, cfg.Minio.Bucket),
		"document-analyzer": grpcdocanalyzer.NewChecker(analyzerConn),
	})

	grpcSrv, err := grpcserver.NewServer(log, cfg.GRPC, userService, healthHandler)
	if err != nil {
		pool.Close()
		natsConn.Close()
		_ = analyzerConn.Close()
		_ = redisClient.Close()
		return nil, fmt.Errorf("grpc server: %w", err)
	}

	return &App{
		log:          pkglog.WithComponent(log, "app"),
		grpcServer:   grpcSrv,
		grpcAddr:     fmt.Sprintf("%s:%d", cfg.GRPC.Host, cfg.GRPC.Port),
		pool:         pool,
		natsConn:     natsConn,
		analyzerConn: analyzerConn,
		redisClient:  redisClient,
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

	a.pool.Close()
	a.natsConn.Close()
	_ = a.analyzerConn.Close()
	_ = a.redisClient.Close()
}
