package app

import (
	"context"
	"log/slog"

	grpcserver "github.com/sorawaslocked/car-rental-booking-service/internal/adapter/grpc"
	"github.com/sorawaslocked/car-rental-booking-service/internal/adapter/grpc/handler"
	"github.com/sorawaslocked/car-rental-booking-service/internal/adapter/grpc/interceptor"
	natsadapter "github.com/sorawaslocked/car-rental-booking-service/internal/adapter/nats"
	"github.com/sorawaslocked/car-rental-booking-service/internal/adapter/postgres"
	"github.com/sorawaslocked/car-rental-booking-service/internal/config"
	pkgnats "github.com/sorawaslocked/car-rental-booking-service/internal/pkg/nats"
	pkgpostgres "github.com/sorawaslocked/car-rental-booking-service/internal/pkg/postgres"
	"github.com/sorawaslocked/car-rental-booking-service/internal/service"
)

type App struct {
	GRPCServer *grpcserver.Server
}

func New(ctx context.Context, log *slog.Logger, cfg config.Config) (*App, error) {
	db, err := pkgpostgres.New(cfg.Postgres)
	if err != nil {
		return nil, err
	}

	nc, err := pkgnats.New(cfg.NATS)
	if err != nil {
		return nil, err
	}

	bookingRepo := postgres.NewBookingRepo(log, db)
	ruleRepo := postgres.NewPricingRuleRepo(log, db)

	publisher := natsadapter.NewPublisher(log, nc)

	bookingSvc := service.NewBookingService(log, bookingRepo, ruleRepo, publisher)
	ruleSvc := service.NewPricingRuleService(log, ruleRepo)

	consumer := natsadapter.NewConsumer(log, nc, bookingSvc)
	if err := consumer.Subscribe(ctx); err != nil {
		return nil, err
	}

	go bookingSvc.StartExpiryWatcher(ctx)

	bookingHandler := handler.NewBookingHandler(log, bookingSvc)
	ruleHandler := handler.NewPricingRuleHandler(log, ruleSvc)
	healthHandler := handler.NewHealthHandler(log, db, nc)

	baseInterceptor := interceptor.NewBaseInterceptor()
	loggerInterceptor := interceptor.NewLoggerInterceptor(log)
	authInterceptor := interceptor.NewAuthInterceptor()

	grpcSrv := grpcserver.NewServer(
		log, cfg.GRPCServer,
		bookingHandler, ruleHandler, healthHandler,
		baseInterceptor, loggerInterceptor, authInterceptor,
	)

	return &App{GRPCServer: grpcSrv}, nil
}
