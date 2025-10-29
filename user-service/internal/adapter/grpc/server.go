package grpc

import (
	"fmt"
	authsvc "github.com/sorawaslocked/car-rental-protos/gen/service/auth"
	usersvc "github.com/sorawaslocked/car-rental-protos/gen/service/user"
	"github.com/sorawaslocked/car-rental-user-service/internal/adapter/grpc/handler"
	grpccfg "github.com/sorawaslocked/car-rental-user-service/internal/pkg/grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log/slog"
	"net"
)

type Server struct {
	s   *grpc.Server
	cfg grpccfg.Config
	log *slog.Logger
}

func NewServer(
	cfg grpccfg.Config,
	log *slog.Logger,
	authService handler.AuthService,
	userService handler.UserService,
) *Server {
	server := &Server{
		cfg: cfg,
		log: log,
	}

	server.register(authService, userService)

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

func (s *Server) register(authService handler.AuthService, userService handler.UserService) {
	s.s = grpc.NewServer()

	authsvc.RegisterAuthServiceServer(s.s, handler.NewAuthHandler(s.log, authService))
	usersvc.RegisterUserServiceServer(s.s, handler.NewUserHandler(s.log, userService))

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
