package app

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-playground/validator/v10"
	natsio "github.com/nats-io/nats.go"
	goredis "github.com/redis/go-redis/v9"
	grpcserver "github.com/sorawaslocked/car-rental-user-service/internal/adapter/grpc"
	grpcdocanalyzer "github.com/sorawaslocked/car-rental-user-service/internal/adapter/grpc/client"
	"github.com/sorawaslocked/car-rental-user-service/internal/adapter/grpc/handler"
	"github.com/sorawaslocked/car-rental-user-service/internal/adapter/mailer"
	minioadapter "github.com/sorawaslocked/car-rental-user-service/internal/adapter/minio"
	natsadapter "github.com/sorawaslocked/car-rental-user-service/internal/adapter/nats"
	natshandler "github.com/sorawaslocked/car-rental-user-service/internal/adapter/nats/handler"
	"github.com/sorawaslocked/car-rental-user-service/internal/adapter/postgres"
	redisadapter "github.com/sorawaslocked/car-rental-user-service/internal/adapter/redis"
	"github.com/sorawaslocked/car-rental-user-service/internal/config"
	pkggrpc "github.com/sorawaslocked/car-rental-user-service/internal/pkg/grpc"
	pkglog "github.com/sorawaslocked/car-rental-user-service/internal/pkg/log"
	miniocfg "github.com/sorawaslocked/car-rental-user-service/internal/pkg/minio"
	natscfg "github.com/sorawaslocked/car-rental-user-service/internal/pkg/nats"
	postgrescfg "github.com/sorawaslocked/car-rental-user-service/internal/pkg/postgres"
	rediscfg "github.com/sorawaslocked/car-rental-user-service/internal/pkg/redis"
	validatecfg "github.com/sorawaslocked/car-rental-user-service/internal/pkg/validate"
	"github.com/sorawaslocked/car-rental-user-service/internal/service"
	"google.golang.org/grpc"
)

type App struct {
	log          *slog.Logger
	grpcServer   *grpc.Server
	grpcAddr     string
	db           *sql.DB
	natsConn     *natsio.Conn
	analyzerConn *grpc.ClientConn
	redisClient  *goredis.Client
}

func New(log *slog.Logger, cfg config.Config) (*App, error) {
	db, err := postgrescfg.NewDB(cfg.Postgres)
	if err != nil {
		return nil, fmt.Errorf("postgres: %w", err)
	}

	validate := validator.New()
	if err := validate.RegisterValidation("min_age", validatecfg.MinAge); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("validator: %w", err)
	}
	if err := validate.RegisterValidation("complex_password", validatecfg.ComplexPassword); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("validator: %w", err)
	}

	natsConn, err := natscfg.NewConn(cfg.NATS)
	if err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("nats: %w", err)
	}

	analyzerConn, err := pkggrpc.NewClientConn(cfg.DocumentAnalyzer)
	if err != nil {
		_ = db.Close()
		natsConn.Close()
		return nil, fmt.Errorf("document analyzer client: %w", err)
	}

	redisClient := rediscfg.Client(cfg.Redis)

	minioClient, err := miniocfg.NewClient(cfg.Minio)
	if err != nil {
		_ = db.Close()
		natsConn.Close()
		_ = analyzerConn.Close()
		_ = redisClient.Close()
		return nil, fmt.Errorf("minio: %w", err)
	}
	if err := miniocfg.EnsureBucket(context.Background(), minioClient, cfg.Minio); err != nil {
		_ = db.Close()
		natsConn.Close()
		_ = analyzerConn.Close()
		_ = redisClient.Close()
		return nil, fmt.Errorf("minio bucket: %w", err)
	}

	activationCodeCache := redisadapter.NewActivationCodeRedisCache(redisClient)
	msMailer := mailer.New(log, cfg.Mailtrap)
	userRepo := postgres.NewUserRepository(log, db)
	docRepo := postgres.NewDocumentRepository(log, db)
	publisher := natsadapter.NewPublisher(log, natsConn)
	minioStorage := minioadapter.NewMinioObjectStorage(log, minioClient, cfg.Minio)
	analyzerClient := grpcdocanalyzer.NewDocumentAnalyzer(log, analyzerConn)

	userService := service.NewUserService(log, validate, userRepo, docRepo, minioStorage, analyzerClient, publisher, activationCodeCache, msMailer)

	docHandler := natshandler.NewDocumentHandler(log, natsConn, userService)
	if err := docHandler.Subscribe(); err != nil {
		_ = db.Close()
		natsConn.Close()
		_ = analyzerConn.Close()
		_ = redisClient.Close()
		return nil, fmt.Errorf("nats subscribe: %w", err)
	}

	healthHandler := handler.NewHealthHandler(log, map[string]handler.Pinger{
		"postgres":          postgres.NewChecker(db),
		"redis":             redisadapter.NewChecker(redisClient),
		"nats":              natsadapter.NewChecker(natsConn),
		"minio":             minioadapter.NewChecker(minioClient, cfg.Minio.BucketName),
		"document-analyzer": grpcdocanalyzer.NewChecker(analyzerConn),
	})

	grpcSrv := grpcserver.NewServer(log, userService, healthHandler)

	return &App{
		log:          pkglog.WithComponent(log, "app"),
		grpcServer:   grpcSrv,
		grpcAddr:     cfg.GRPC.Addr,
		db:           db,
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

	_ = a.db.Close()
	a.natsConn.Close()
	_ = a.analyzerConn.Close()
	_ = a.redisClient.Close()
}
