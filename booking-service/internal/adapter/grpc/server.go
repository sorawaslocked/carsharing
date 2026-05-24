package grpc

import (
	"log/slog"

	"carsharing/booking-service/internal/adapter/grpc/handler"
	"carsharing/booking-service/internal/adapter/grpc/interceptor"
	servicebookingpb "carsharing/protos/gen/service/booking"
	pkggrpc "carsharing/shared/pkg/grpc"
	"google.golang.org/grpc"
)

func NewServer(
	log *slog.Logger,
	cfg pkggrpc.ServerConfig,
	bookingHandler *handler.BookingHandler,
	ruleHandler *handler.PricingRuleHandler,
	healthHandler *handler.HealthHandler,
	base *interceptor.BaseInterceptor,
	logger *interceptor.LoggerInterceptor,
	auth *interceptor.AuthInterceptor,
) (*grpc.Server, error) {
	srv, err := pkggrpc.NewServer(
		log,
		cfg,
		grpc.ChainUnaryInterceptor(
			base.Unary(),
			logger.Unary(),
			auth.Unary(),
		),
		grpc.ChainStreamInterceptor(),
	)
	if err != nil {
		return nil, err
	}

	servicebookingpb.RegisterBookingServiceServer(srv, bookingHandler)
	servicebookingpb.RegisterPricingRuleServiceServer(srv, ruleHandler)
	servicebookingpb.RegisterHealthServiceServer(srv, healthHandler)

	return srv, nil
}
