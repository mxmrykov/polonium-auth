package repository

import (
	"context"
	_ "embed"
	"time"

	"github.com/mxmrykov/polonium-auth/internal/config"
	"github.com/mxmrykov/polonium-auth/internal/provider"
)

type (
	IAuthPostgres interface {
		IsUserExists(ctx context.Context, email string) (bool, error)
	}

	authPostgres struct {
		pg            *provider.PostgresProvider
		connectionTtl time.Duration
	}
)

var (
	//go:embed sql/isUserExists.sql
	isUserExistsQuery string
)

func NewAuthPostgres(cfg *config.Psql) (IAuthPostgres, error) {
	p, err := provider.NewPostgresProvider(cfg)

	if err != nil {
		return nil, err
	}

	return &authPostgres{
		pg:            p,
		connectionTtl: 15 * time.Second,
	}, nil
}

func (a *authPostgres) IsUserExists(ctx context.Context, email string) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, a.connectionTtl)
	defer cancel()

	var exists bool
	if err := a.pg.GetConnect().QueryRow(ctx, isUserExistsQuery, email).Scan(&exists); err != nil {
		return false, err
	}

	return exists, nil
}
