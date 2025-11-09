package repository

import (
	"fmt"
	"time"

	"github.com/mxmrykov/polonium-auth/internal/config"
	"github.com/mxmrykov/polonium-auth/internal/provider"
	"github.com/mxmrykov/polonium-auth/internal/vars"
)

type (
	IAuthRedis interface {
		HasActiveECSession(user string) (bool, error)
		SetCode(user, code string) error
		GetCode(user string) (string, error)
		DropCode(user string) error
		NewAuthSession(user, session string) error
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

func (a *authRedis) HasActiveECSession(user string) (bool, error) {
	key := fmt.Sprintf(vars.CodesSignupEmailConfirmation, user)

	exists, err := a.rdb.IsExists(key)

	if err != nil {
		return false, fmt.Errorf("cannot check user existance: %v", err)
	}

	return exists, nil
}

func (a *authRedis) SetCode(user, code string) error {
	key := fmt.Sprintf(vars.CodesSignupEmailConfirmation, user)
	return a.rdb.Set(key, code, 2*time.Minute)
}

func (a *authRedis) GetCode(user string) (string, error) {
	key := fmt.Sprintf(vars.CodesSignupEmailConfirmation, user)
	return a.rdb.Get(key)
}

func (a *authRedis) DropCode(user string) error {
	key := fmt.Sprintf(vars.CodesSignupEmailConfirmation, user)
	return a.rdb.Drop(key)
}

func (a *authRedis) NewAuthSession(user, session string) error {
	key := fmt.Sprintf(vars.AuthSessionsUsers, user)
	return a.rdb.Set(key, session, time.Hour)
}
