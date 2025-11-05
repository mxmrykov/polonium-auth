package provider

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mxmrykov/polonium-auth/internal/config"
)

type PostgresProvider struct {
	pool *pgxpool.Pool
	name string
}

func NewPostgresProvider(conf *config.Psql) (*PostgresProvider, error) {
	const defaultConnect = 50

	if conf.MaxConn < 1 {
		conf.MaxConn = defaultConnect
	}

	connStr := fmt.Sprintf("postgres://%s:%s@%s/%s?pool_max_conns=%d&pool_max_conn_idle_time=3s", conf.User, conf.Pass, conf.Host, conf.DB, conf.MaxConn)

	pgxPoolConfig, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		return nil, fmt.Errorf("fail postgres parse config with err %s", err.Error())
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), pgxPoolConfig)
	if err != nil {
		return nil, fmt.Errorf("fail postgres connect with err %s", err.Error())
	}
	if err = pool.Ping(context.Background()); err != nil {
		return nil, fmt.Errorf("fail postgres ping with err %s", err.Error())
	}

	return &PostgresProvider{
		pool: pool,
	}, nil
}

func (p *PostgresProvider) GetConnect() *pgxpool.Pool {
	return p.pool
}

func (p *PostgresProvider) Name() string {
	return p.name
}

func (p *PostgresProvider) Close() {
	p.pool.Close()
}

type PostgresPool struct {
	master  *PostgresProvider
	replica *PostgresProvider
}

func NewPostgresPool(poolCfg *config.Psql) (*PostgresPool, error) {
	master, err := NewPostgresProvider(poolCfg)
	if err != nil {
		return nil, err
	}

	p := &PostgresPool{
		master: master,
	}

	if err != nil {
		return nil, fmt.Errorf("failed to run migrations, %w", err)
	}

	return p, nil
}

func (p *PostgresPool) GetMaster() *PostgresProvider {
	return p.master
}

func (p *PostgresPool) Close() {
	p.master.Close()
}
