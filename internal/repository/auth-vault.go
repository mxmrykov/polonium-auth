package repository

import (
	"context"
	"fmt"

	"github.com/mxmrykov/polonium-auth/internal/config"
	"github.com/mxmrykov/polonium-auth/internal/provider"
	"github.com/mxmrykov/polonium-auth/internal/vars"
)

type (
	IAuthVault interface {
		PutNewUser(ctx context.Context, user, pwd, totpSecret string) error
		GetPwdHash(ctx context.Context, user string) (string, error)
		GetTOTPSecret(ctx context.Context, user string) (string, error)
	}

	authVault struct {
		vault provider.IVault
	}
)

func NewAuthVault(cfg *config.Vault) (IAuthVault, error) {
	vaultProvider, err := provider.NewVault(cfg)

	if err != nil {
		return nil, err
	}

	return &authVault{vault: vaultProvider}, nil
}

func (a *authVault) PutNewUser(ctx context.Context, user, pwd, totpSecret string) error {
	if err := a.vault.Write(
		ctx,
		fmt.Sprintf(vars.UsersGlobalLoginPwd, user),
		map[string]interface{}{
			"val": pwd,
		},
	); err != nil {
		return err
	}

	return a.vault.Write(
		ctx,
		fmt.Sprintf(vars.UsersTOTPCodes, user),
		map[string]interface{}{
			"val": totpSecret,
		},
	)
}

func (a *authVault) GetPwdHash(ctx context.Context, user string) (string, error) {
	route := fmt.Sprintf(vars.UsersGlobalLoginPwd, user)

	pwd, err := a.vault.Read(ctx, route)

	if err != nil {
		return "", err
	}

	val, ok := pwd["val"]

	if !ok {
		return "", vars.ErrNoSuchVariableInVault
	}

	return val.(string), nil
}

func (a *authVault) GetTOTPSecret(ctx context.Context, user string) (string, error) {
	route := fmt.Sprintf(vars.UsersTOTPCodes, user)

	secret, err := a.vault.Read(ctx, route)

	if err != nil {
		return "", err
	}

	val, ok := secret["val"]

	if !ok {
		return "", vars.ErrNoSuchVariableInVault
	}

	return val.(string), nil
}
