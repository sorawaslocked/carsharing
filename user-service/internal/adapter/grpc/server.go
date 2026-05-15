package grpc

import (
	"log/slog"

	usersvc "github.com/sorawaslocked/car-rental-protos/gen/service/user"
	"github.com/sorawaslocked/car-rental-user-service/internal/adapter/grpc/handler"
	"github.com/sorawaslocked/car-rental-user-service/internal/adapter/grpc/interceptor"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func NewServer(
	log *slog.Logger,
	userService handler.UserService,
	healthHandler *handler.HealthHandler,
) *grpc.Server {
	baseInterceptor := interceptor.NewBaseInterceptor()
	loggerInterceptor := interceptor.NewLoggerInterceptor(log)
	authInterceptor := interceptor.NewAuthInterceptor(log)

	s := grpc.NewServer(grpc.ChainUnaryInterceptor(
		baseInterceptor.Unary,
		loggerInterceptor.Unary,
		authInterceptor.Unary,
	))

	usersvc.RegisterUserServiceServer(s, handler.NewUserHandler(log, userService))
	usersvc.RegisterHealthServiceServer(s, healthHandler)

	reflection.Register(s)

	return s
}
