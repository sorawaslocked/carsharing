package grpc

import (
	grpccfg "car-rental-user-service/internal/pkg/grpc"
	"car-rental-user-service/internal/service"
	"fmt"
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

	server.register()

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

func (s *Server) register() {
	s.s = grpc.NewServer()

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
