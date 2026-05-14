package grpc

import (
	"log/slog"
	"net"

	"github.com/sorawaslocked/car-rental-booking-service/internal/adapter/grpc/handler"
	"github.com/sorawaslocked/car-rental-booking-service/internal/adapter/grpc/interceptor"
	"github.com/sorawaslocked/car-rental-booking-service/internal/config"
	pkglog "github.com/sorawaslocked/car-rental-booking-service/internal/pkg/log"
	servicebookingpb "github.com/sorawaslocked/car-rental-protos/gen/service/booking"
	"google.golang.org/grpc"
)

type Server struct {
	log    *slog.Logger
	server *grpc.Server
	addr   string
}

func NewServer(
	log *slog.Logger,
	cfg config.GRPCServer,
	bookingHandler *handler.BookingHandler,
	ruleHandler *handler.PricingRuleHandler,
	healthHandler *handler.HealthHandler,
	base *interceptor.BaseInterceptor,
	logger *interceptor.LoggerInterceptor,
	auth *interceptor.AuthInterceptor,
) *Server {
	grpcSrv := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			base.Unary(),
			logger.Unary(),
			auth.Unary(),
		),
	)

	servicebookingpb.RegisterBookingServiceServer(grpcSrv, bookingHandler)
	servicebookingpb.RegisterPricingRuleServiceServer(grpcSrv, ruleHandler)
	servicebookingpb.RegisterHealthServiceServer(grpcSrv, healthHandler)

	return &Server{
		log:    pkglog.WithComponent(log, "grpc.Server"),
		server: grpcSrv,
		addr:   cfg.Address,
	}
}

func (s *Server) Run() error {
	lis, err := net.Listen("tcp", s.addr)
	if err != nil {
		return err
	}

	s.log.Info("grpc server listening", slog.String("addr", s.addr))

	return s.server.Serve(lis)
}

func (s *Server) Stop() {
	s.server.GracefulStop()
}
