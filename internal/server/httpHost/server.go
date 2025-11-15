package httpHost

import (
	"context"
	"crypto/tls"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mxmrykov/polonium-auth/internal/config"
)

type (
	IServer interface {
		Start() error
		Stop(ctx context.Context) error
		Router() *gin.Engine
	}

	Server struct {
		cfg    *config.PAuth
		server *http.Server
		router *gin.Engine
	}
)

func New(cfg *config.PAuth, TLS ...*tls.Config) IServer {
	router := gin.New()
	server := &http.Server{
		Addr:    cfg.PublicServer.Port,
		Handler: router,
	}

	s := &Server{
		cfg:    cfg,
		server: server,
		router: router,
	}

	if len(TLS) > 0 && TLS[0] != nil {
		s.server.TLSConfig = TLS[0]
	}

	return s
}

func (s *Server) Start() error {
	return s.server.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

func (s *Server) Router() *gin.Engine {
	return s.router
}
