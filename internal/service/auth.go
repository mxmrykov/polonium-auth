package service

import (
	"context"
	"fmt"

	"github.com/mxmrykov/polonium-auth/internal/repository"
	"github.com/mxmrykov/polonium-auth/internal/vars"
	"github.com/mxmrykov/polonium-auth/pkg/utils"
)

type (
	IAuth interface {
		CanConfirmSignup(ctx context.Context, user string) error
		ConfirmEmail(user string) error
	}

	auth struct {
		authPg  repository.IAuthPostgres
		authRdb repository.IAuthRedis
		emailer repository.IEmailer
	}
)

func NewAuth(authPg repository.IAuthPostgres,
	authRdb repository.IAuthRedis,
	emailer repository.IEmailer,
) IAuth {
	return &auth{
		authPg:  authPg,
		authRdb: authRdb,
		emailer: emailer,
	}
}

func (a *auth) CanConfirmSignup(ctx context.Context, user string) error {
	exists, err := a.authPg.IsUserExists(ctx, user)
	if err != nil {
		return fmt.Errorf("cannot check user existance in db: %v", err)
	}

	if exists {
		return vars.ErrUserAlreadyExists
	}

	haveActiveSignupSession, err := a.authRdb.HasActiveSignupSession(user)

	if err != nil {
		return fmt.Errorf("cannot check auth session existance: %v", err)
	}

	if haveActiveSignupSession {
		return vars.ErrUserAlreadyConfirmingSignup
	}

	return nil
}

func (a *auth) ConfirmEmail(user string) error {
	code := utils.RandVerificationCode()

	if err := a.authRdb.SetCode(user, code); err != nil {
		return fmt.Errorf("cannot set confirmation code in redis: %v", err)
	}

	if err := a.emailer.SendVerificationCode(code, user); err != nil {
		_ = a.authRdb.DropCode(user)
		return fmt.Errorf("cannot send confirmation code: %v", err)
	}

	return nil
}
