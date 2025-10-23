package grpc

import (
	"fmt"
	svc "github.com/sorawaslocked/car-rental-protos/gen/service"
	"github.com/sorawaslocked/car-rental-user-service/internal/adapter/grpc/handler"
	grpccfg "github.com/sorawaslocked/car-rental-user-service/internal/pkg/grpc"
	"github.com/sorawaslocked/car-rental-user-service/internal/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log/slog"
	"net"
)

type Server struct {
	s           *grpc.Server
	cfg         grpccfg.Config
	log         *slog.Logger
	authService *service.AuthService
}

func NewServer(
	cfg grpccfg.Config,
	log *slog.Logger,
	authService *service.AuthService,
) *Server {
	server := &Server{
		cfg:         cfg,
		log:         log,
		authService: authService,
	}

	server.register(authService)

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

func (s *Server) register(authService handler.AuthService) {
	s.s = grpc.NewServer()

	svc.RegisterAuthServiceServer(s.s, handler.NewAuthHandler(s.log, authService))

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
