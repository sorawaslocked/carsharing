package grpc

import (
	"fmt"
	"log/slog"
	"net"

	usersvc "github.com/sorawaslocked/car-rental-protos/gen/service/user"
	"github.com/sorawaslocked/car-rental-user-service/internal/adapter/grpc/handler"
	"github.com/sorawaslocked/car-rental-user-service/internal/adapter/grpc/interceptor"
	grpccfg "github.com/sorawaslocked/car-rental-user-service/internal/pkg/grpc"
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
	userService handler.UserService,
) *Server {
	server := &Server{
		cfg: cfg,
		log: log,
	}

	server.register(userService)

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

func (s *Server) register(userService handler.UserService) {
	baseInterceptor := interceptor.NewBaseInterceptor()
	loggerInterceptor := interceptor.NewLoggerInterceptor(s.log)

	s.s = grpc.NewServer(grpc.ChainUnaryInterceptor(
		baseInterceptor.Unary,
		loggerInterceptor.Unary,
	))

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
