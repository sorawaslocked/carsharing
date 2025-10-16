package http

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sorawaslocked/car-rental-api-gateway/internal/config"
	"net/http"
)

type Server struct {
	s   *gin.Engine
	cfg config.HTTPServer
}

func New(cfg config.HTTPServer) *Server {
	gin.SetMode(cfg.GinMode)
	s := gin.New()

	// Middleware
	s.Use(gin.Recovery())

	server := &Server{
		s:   s,
		cfg: cfg,
	}

	return server
}

func (s *Server) MustRun() {
	go func() {
		addr := fmt.Sprintf("%s:%d", s.cfg.Host, s.cfg.Port)

		err := s.s.Run(addr)
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			panic(err)
		}
	}()
}

func (s *Server) Stop(ctx context.Context) {
	<-ctx.Done()
}
