package grpc

import (
	"log/slog"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	tripsvc "carsharing/protos/gen/service/trip"
	pkggrpc "carsharing/shared/pkg/grpc"

	"carsharing/trip-service/internal/adapter/grpc/handler"
	"carsharing/trip-service/internal/adapter/grpc/interceptor"
	"carsharing/trip-service/internal/adapter/grpc/interceptor/auth"
)

func NewServer(
	log *slog.Logger,
	cfg pkggrpc.ServerConfig,
	tripHandler *handler.TripHandler,
	streamHandler *handler.TripStreamHandler,
	healthHandler *handler.HealthHandler,
) (*grpc.Server, error) {
	baseInterceptor := interceptor.NewBaseInterceptor()
	loggerInterceptor := interceptor.NewLoggerInterceptor(log)
	authInterceptor := auth.NewInterceptor(log)

	srv, err := pkggrpc.NewServer(
		log,
		cfg,
		grpc.ChainUnaryInterceptor(
			baseInterceptor.Unary,
			loggerInterceptor.Unary,
			authInterceptor.Unary,
		),
		grpc.ChainStreamInterceptor(
			baseInterceptor.Stream,
			loggerInterceptor.Stream,
			authInterceptor.Stream,
		),
	)
	if err != nil {
		return nil, err
	}

	tripsvc.RegisterTripServiceServer(srv, tripHandler)
	tripsvc.RegisterTripStreamServiceServer(srv, streamHandler)
	tripsvc.RegisterHealthServiceServer(srv, healthHandler)
	reflection.Register(srv)

	return srv, nil
}
