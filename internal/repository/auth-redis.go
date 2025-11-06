package repository

import (
	"fmt"
	"time"

	"github.com/mxmrykov/polonium-auth/internal/config"
	"github.com/mxmrykov/polonium-auth/internal/provider"
)

type (
	IAuthRedis interface {
		HasActiveSignupSession(user string) (bool, error)
		SetCode(user, code string) error
		DropCode(user string) error
	}

	authRedis struct {
		rdb provider.IRedis
	}
)

func NewAuthRedis(cfg *config.Redis) IAuthRedis {
	return &authRedis{
		rdb: provider.NewRedis(cfg),
	}
}

func (a *authRedis) HasActiveSignupSession(user string) (bool, error) {
	key := fmt.Sprintf("codes/signup/email-confirmation/%s", user)

	exists, err := a.rdb.IsExists(key)

	if err != nil {
		return false, fmt.Errorf("cannot check user existance: %v", err)
	}

	return exists, nil
}

func (a *authRedis) SetCode(user, code string) error {
	key := fmt.Sprintf("codes/signup/email-confirmation/%s", user)
	return a.rdb.Set(key, code, 2*time.Minute)
}

func (a *authRedis) DropCode(user string) error {
	key := fmt.Sprintf("codes/signup/email-confirmation/%s", user)
	return a.rdb.Drop(key)
}
