package grpc

import (
	"fmt"
	"log/slog"
	"net"

	"github.com/sorawaslocked/car-rental-car-service/internal/adapter/grpc/handler"
	"github.com/sorawaslocked/car-rental-car-service/internal/adapter/grpc/interceptor"
	authinterceptor "github.com/sorawaslocked/car-rental-car-service/internal/adapter/grpc/interceptor/auth"
	grpccfg "github.com/sorawaslocked/car-rental-car-service/internal/pkg/grpc"
	carsvc "github.com/sorawaslocked/car-rental-protos/gen/service/car"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type Server struct {
	s   *grpc.Server
	cfg grpccfg.Config
	log *slog.Logger
}

func NewServer(
	cfg grpccfg.Config,
	log *slog.Logger,
	carModelService handler.CarModelService,
	carService handler.CarService,
) *Server {
	server := &Server{
		cfg: cfg,
		log: log,
	}

	server.register(carModelService, carService, log)

	return server
}

func (s *Server) MustRun() {
	go func() {
		if err := s.run(); err != nil {
			panic(err)
		}
	}()
}

func (s *Server) Stop() {
	s.log.Info(
		"stopping grpc server",
		slog.String("addr", fmt.Sprintf("%s:%d", s.cfg.Host, s.cfg.Port)),
	)

	s.s.GracefulStop()
}

func (s *Server) register(
	carModelService handler.CarModelService,
	carService handler.CarService,
	log *slog.Logger,
) {
	baseInterceptor := interceptor.NewBaseInterceptor()
	authInterceptor := authinterceptor.NewInterceptor()

	s.s = grpc.NewServer(grpc.ChainUnaryInterceptor(
		baseInterceptor.Unary,
		authInterceptor.Unary,
	))

	carsvc.RegisterCarModelServiceServer(s.s, handler.NewCarModelHandler(carModelService, log))
	carsvc.RegisterCarServiceServer(s.s, handler.NewCarHandler(carService, log))

	reflection.Register(s.s)
}

func (s *Server) run() error {
	addr := fmt.Sprintf("%s:%d", s.cfg.Host, s.cfg.Port)

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	s.log.Info("starting grpc server", slog.String("addr", addr))
	if err := s.s.Serve(listener); err != nil {
		return err
	}

	return nil
}
