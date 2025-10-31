package http

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	"github.com/sorawaslocked/car-rental-api-gateway/internal/adapter/http/handler"
	"github.com/sorawaslocked/car-rental-api-gateway/internal/adapter/http/middleware"
	"github.com/sorawaslocked/car-rental-api-gateway/internal/config"
	"log/slog"
	"net/http"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

type Server struct {
	router      *gin.Engine
	cfg         config.HTTPServer
	log         *slog.Logger
	authHandler *handler.Auth
	userHandler *handler.User
}

func New(
	env string,
	cfg config.HTTPServer,
	log *slog.Logger,
	authService handler.AuthService,
	userService handler.UserService,
) *Server {
	httpLog := log.With(
		slog.String("httpServerHost", cfg.Host),
		slog.Int("httpServerPort", cfg.Port),
	)

	gin.SetMode(cfg.GinMode)
	router := gin.New()

	// Middleware
	router.Use(gin.Recovery())
	if env == envLocal || env == envDev {
		router.Use(cors.Default())
	}
	router.Use(requestid.New())
	router.Use(middleware.Base())
	router.Use(middleware.Logger(httpLog))

	// Handlers
	authHandler := handler.NewAuth(authService)
	userHandler := handler.NewUser(userService)

	server := &Server{
		router:      router,
		cfg:         cfg,
		log:         httpLog,
		authHandler: authHandler,
		userHandler: userHandler,
	}

	server.setupRoutes()

	return server
}

func (s *Server) MustRun() {
	go func() {
		addr := fmt.Sprintf("%s:%d", s.cfg.Host, s.cfg.Port)

		s.log.Info("starting http server")
		err := s.router.Run(addr)
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			panic(err)
		}
	}()
}

func (s *Server) Stop(ctx context.Context) {
	s.log.Info("shutting down http server")
	<-ctx.Done()
}
