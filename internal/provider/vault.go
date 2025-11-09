package provider

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/vault/api"
	"github.com/mxmrykov/polonium-auth/internal/config"
)

type (
	IVault interface {
		Write(ctx context.Context, path string, data map[string]interface{}) error
		Read(ctx context.Context, path string) (map[string]interface{}, error)
	}

	vaultClient struct {
		client        *api.Client
		kv2           *api.KVv2
		connectionTtl time.Duration
	}
)

func NewVault(cfg *config.Vault) (IVault, error) {
	defConfig := api.DefaultConfig()
	defConfig.Address = cfg.Address

	client, err := api.NewClient(defConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create Vault client: %w", err)
	}

	client.SetToken(cfg.Token)

	return &vaultClient{
		client:        client,
		kv2:           client.KVv2(cfg.MountPath),
		connectionTtl: 15 * time.Second,
	}, nil
}

func (v *vaultClient) Write(ctx context.Context, path string, data map[string]interface{}) error {
	ctx, cancel := context.WithTimeout(ctx, v.connectionTtl)
	defer cancel()

	_, err := v.kv2.Put(ctx, path, data)
	return err
}

func (v *vaultClient) Read(ctx context.Context, path string) (map[string]interface{}, error) {
	ctx, cancel := context.WithTimeout(ctx, v.connectionTtl)
	defer cancel()

	secret, err := v.kv2.Get(ctx, path)
	if err != nil {
		return nil, err
	}

	if secret == nil {
		return nil, fmt.Errorf("secret not found")
	}

	return secret.Data, nil
}
