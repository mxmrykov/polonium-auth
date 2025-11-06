package app

import (
	"context"
	"fmt"

	"github.com/mxmrykov/polonium-auth/internal/config"
	"github.com/mxmrykov/polonium-auth/internal/repository"
	"github.com/mxmrykov/polonium-auth/internal/server/httpHost"
)

type (
	Application struct {
		cfg        *config.PAuth
		httpServer httpHost.IServer
	}

	repositories struct {
		authPg  repository.IAuthPostgres
		authRdb repository.IAuthRedis
		emailer repository.IEmailer
	}
)

func New(cfg *config.PAuth) (*Application, error) {
	httpServer := httpHost.New(cfg)
	a := &Application{
		cfg:        cfg,
		httpServer: httpServer,
	}

	repos, err := a.initRepositories()

	if err != nil {
		return nil, fmt.Errorf("cannot init app: %v", err)
	}

	a.setupRoutesAPIV1(repos.authPg, repos.authRdb, repos.emailer)

	return a, nil
}

func (a *Application) Run() error {
	return a.httpServer.Start()
}

func (a *Application) Stop(ctx context.Context) error {
	return a.httpServer.Stop(ctx)
}
