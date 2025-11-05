package httpHost

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mxmrykov/polonium-auth/internal/config"
)

type (
	IServer interface {
		Start() error
		Stop(ctx context.Context) error
	}

	Server struct {
		cfg    *config.PAuth
		server *http.Server
		router *gin.Engine
	}
)

func New(cfg *config.PAuth) IServer {
	router := gin.New()
	server := &http.Server{
		Addr:    cfg.Server.Port,
		Handler: router,
	}

	return &Server{
		cfg:    cfg,
		server: server,
		router: router,
	}
}

func (s *Server) Start() error {
	return s.server.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
