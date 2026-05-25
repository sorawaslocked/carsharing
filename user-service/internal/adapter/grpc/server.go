package grpc

import (
	"log/slog"

	pkggrpc "carsharing/shared/pkg/grpc"
	"carsharing/user-service/internal/adapter/grpc/handler"
	"carsharing/user-service/internal/adapter/grpc/interceptor"
	"carsharing/user-service/internal/adapter/grpc/interceptor/auth"

	usersvc "carsharing/protos/gen/service/user"
	"google.golang.org/grpc"
)

func NewServer(
	log *slog.Logger,
	cfg pkggrpc.ServerConfig,
	userService handler.UserService,
	documentSubscriber handler.DocumentAnalyzedSubscriber,
	healthHandler *handler.HealthHandler,
) (*grpc.Server, error) {
	baseInterceptor := interceptor.NewBaseInterceptor()
	loggerInterceptor := interceptor.NewLoggerInterceptor(log)
	authInterceptor := auth.NewInterceptor(log)

	s, err := pkggrpc.NewServer(
		log,
		cfg,
		grpc.ChainUnaryInterceptor(
			baseInterceptor.Unary,
			loggerInterceptor.Unary,
			authInterceptor.Unary,
		),
		grpc.ChainStreamInterceptor(
			baseInterceptor.Stream,
			loggerInterceptor.Stream,
			authInterceptor.Stream,
		),
	)
	if err != nil {
		return nil, err
	}

	usersvc.RegisterUserServiceServer(s, handler.NewUserHandler(log, userService, documentSubscriber))
	usersvc.RegisterHealthServiceServer(s, healthHandler)

	return s, nil
}
