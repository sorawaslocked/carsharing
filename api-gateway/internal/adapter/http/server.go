package http

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	"github.com/sorawaslocked/car-rental-api-gateway/internal/adapter/http/handler"
	"github.com/sorawaslocked/car-rental-api-gateway/internal/adapter/http/middleware"
	"github.com/sorawaslocked/car-rental-api-gateway/internal/config"
	"log/slog"
	"net/http"
)

type Server struct {
	router      *gin.Engine
	httpCfg     config.HTTPServer
	log         *slog.Logger
	authHandler *handler.Auth
	userHandler *handler.User
}

func New(
	httpCfg config.HTTPServer,
	cookieCfg config.Cookie,
	log *slog.Logger,
	authService handler.AuthService,
	userService handler.UserService,
) *Server {
	httpLog := log.With(
		slog.String("httpServerHost", httpCfg.Host),
		slog.Int("httpServerPort", httpCfg.Port),
	)

	gin.SetMode(httpCfg.GinMode)
	router := gin.New()

	// Middleware
	router.Use(gin.Recovery())
	router.Use(middleware.Cors())
	router.Use(requestid.New())
	router.Use(middleware.Base())
	router.Use(middleware.Logger(httpLog))

	// Handlers
	authHandler := handler.NewAuth(authService, cookieCfg)
	userHandler := handler.NewUser(userService)

	server := &Server{
		router:      router,
		httpCfg:     httpCfg,
		log:         httpLog,
		authHandler: authHandler,
		userHandler: userHandler,
	}

	server.setupRoutes()

	return server
}

func (s *Server) MustRun() {
	go func() {
		addr := fmt.Sprintf("%s:%d", s.httpCfg.Host, s.httpCfg.Port)

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
