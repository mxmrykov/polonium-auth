package provider

import (
	"context"
	"errors"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/mxmrykov/polonium-auth/internal/config"
)

type (
	IRedis interface {
		Get(key string) (string, error)
		IsExists(key string) (bool, error)
		Set(key, value string, ttl time.Duration) error
		Drop(key string) error
	}

	rdb struct {
		db *redis.Client
	}
)

func NewRedis(cfg *config.Redis) IRedis {
	r := redis.NewClient(&redis.Options{
		Addr:     cfg.Host,
		Password: cfg.Password,
		DB:       cfg.Db,
	})

	return &rdb{
		db: r,
	}
}

func (r *rdb) Get(key string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	val, err := r.db.Get(ctx, key).Result()
	if err != nil {
		return "", err
	}

	return val, nil
}

func (r *rdb) IsExists(key string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err := r.db.Get(ctx, key).Result()

	if err != nil {
		if errors.Is(err, redis.Nil) {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

func (r *rdb) Set(key, value string, ttl time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return r.db.Set(ctx, key, value, ttl).Err()
}

func (r *rdb) Drop(key string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return r.db.Del(ctx, key).Err()
}
