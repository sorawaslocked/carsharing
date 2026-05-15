package grpc

import (
	"log/slog"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	tripsvc "github.com/sorawaslocked/car-rental-protos/gen/service/trip"

	"github.com/sorawaslocked/car-rental-trip-service/internal/adapter/grpc/handler"
	"github.com/sorawaslocked/car-rental-trip-service/internal/adapter/grpc/interceptor"
)

func NewServer(
	log *slog.Logger,
	tripHandler *handler.TripHandler,
	streamHandler *handler.TripStreamHandler,
	healthHandler *handler.HealthHandler,
) *grpc.Server {
	srv := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			interceptor.AuthUnaryInterceptor,
			interceptor.LoggerUnaryInterceptor(log),
		),
		grpc.ChainStreamInterceptor(
			interceptor.AuthStreamInterceptor,
			interceptor.LoggerStreamInterceptor(log),
		),
	)

	tripsvc.RegisterTripServiceServer(srv, tripHandler)
	tripsvc.RegisterTripStreamServiceServer(srv, streamHandler)
	tripsvc.RegisterHealthServiceServer(srv, healthHandler)
	reflection.Register(srv)

	return srv
}
