package http

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sorawaslocked/car-rental-api-gateway/internal/adapter/http/handler"
	"github.com/sorawaslocked/car-rental-api-gateway/internal/config"
	"log"
	"net/http"
)

type Server struct {
	s           *gin.Engine
	cfg         config.HTTPServer
	authHandler *handler.Auth
}

func New(cfg config.HTTPServer, authService handler.AuthService) *Server {
	gin.SetMode(cfg.GinMode)
	s := gin.New()

	// Middleware
	s.Use(gin.Recovery())

	// Handlers
	authHandler := handler.NewAuth(authService)

	server := &Server{
		s:           s,
		cfg:         cfg,
		authHandler: authHandler,
	}

	server.setupRoutes()

	return server
}

func (s *Server) setupRoutes() {
	v1 := s.s.Group("/api/v1")
	{
		auth := v1.Group("/auth")
		{
			auth.POST("/register", s.authHandler.Register)
			auth.POST("/login", s.authHandler.Login)
			auth.POST("/refresh-token", s.authHandler.RefreshToken)
		}
	}
}

func (s *Server) MustRun() {
	go func() {
		addr := fmt.Sprintf("%s:%d", s.cfg.Host, s.cfg.Port)

		log.Println("starting http server at", addr)
		err := s.s.Run(addr)
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			panic(err)
		}
	}()
}

func (s *Server) Stop(ctx context.Context) {
	<-ctx.Done()
}
