package grpc

import (
	"log/slog"

	"carsharing/booking-service/internal/adapter/grpc/handler"
	"carsharing/booking-service/internal/adapter/grpc/interceptor"
	"carsharing/booking-service/internal/adapter/grpc/interceptor/auth"
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
