package repository

import (
	"context"
	_ "embed"
	"time"

	"github.com/mxmrykov/polonium-auth/internal/config"
	"github.com/mxmrykov/polonium-auth/internal/model"
	"github.com/mxmrykov/polonium-auth/internal/provider"
)

type (
	IAuthPostgres interface {
		IsUserExists(ctx context.Context, email string) (bool, error)
		Signup(ctx context.Context, user *model.User) error
		VerificateUser(ctx context.Context, user string) error
	}

	authPostgres struct {
		pg            *provider.PostgresProvider
		connectionTtl time.Duration
	}
)

var (
	//go:embed sql/isUserExists.sql
	isUserExistsQuery string

	//go:embed sql/signupUser.sql
	signupUserQuery string

	//go:embed sql/verificate.sql
	verificateQuery string
)

func NewAuthPostgres(cfg *config.Psql) (IAuthPostgres, error) {
	p, err := provider.NewPostgresPool(cfg)

	if err != nil {
		return nil, err
	}

	return &authPostgres{
		pg:            p.GetMaster(),
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

func (a *authPostgres) Signup(ctx context.Context, user *model.User) error {
	ctx, cancel := context.WithTimeout(ctx, a.connectionTtl)
	defer cancel()

	if _, err := a.pg.GetConnect().Exec(
		ctx, signupUserQuery,
		user.Email, user.Id, false, false, user.SshSign, user.Deployer,
	); err != nil {
		return err
	}

	return nil
}

func (a *authPostgres) VerificateUser(ctx context.Context, user string) error {
	ctx, cancel := context.WithTimeout(ctx, a.connectionTtl)
	defer cancel()

	if _, err := a.pg.GetConnect().Exec(
		ctx, verificateQuery,
		user,
	); err != nil {
		return err
	}

	return nil
}
