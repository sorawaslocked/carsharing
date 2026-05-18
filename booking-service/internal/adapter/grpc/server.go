package grpc

import (
	"carsharing/booking-service/internal/adapter/grpc/handler"
	"carsharing/booking-service/internal/adapter/grpc/interceptor"
	servicebookingpb "github.com/sorawaslocked/car-rental-protos/gen/service/booking"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func NewServer(
	bookingHandler *handler.BookingHandler,
	ruleHandler *handler.PricingRuleHandler,
	healthHandler *handler.HealthHandler,
	base *interceptor.BaseInterceptor,
	logger *interceptor.LoggerInterceptor,
	auth *interceptor.AuthInterceptor,
) *grpc.Server {
	srv := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			base.Unary(),
			logger.Unary(),
			auth.Unary(),
		),
	)

	servicebookingpb.RegisterBookingServiceServer(srv, bookingHandler)
	servicebookingpb.RegisterPricingRuleServiceServer(srv, ruleHandler)
	servicebookingpb.RegisterHealthServiceServer(srv, healthHandler)

	reflection.Register(srv)

	return srv
}
