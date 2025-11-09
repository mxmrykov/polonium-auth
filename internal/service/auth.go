package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/mxmrykov/polonium-auth/internal/model"
	"github.com/mxmrykov/polonium-auth/internal/repository"
	"github.com/mxmrykov/polonium-auth/internal/vars"
	"github.com/mxmrykov/polonium-auth/pkg/utils"
	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
)

type (
	IAuth interface {
		CanConfirmSignup(ctx context.Context, user string) error
		ConfirmEmail(user string) error
		ConfirmCode(user, code string) error
		SignupUnverified(ctx context.Context, user string, pwd string) error
		VerifyUser(ctx context.Context, user, pwd string) error
	}

	auth struct {
		authPg  repository.IAuthPostgres
		authRdb repository.IAuthRedis
		emailer repository.IEmailer
		vault   repository.IAuthVault
	}
)

func NewAuth(
	authPg repository.IAuthPostgres,
	authRdb repository.IAuthRedis,
	emailer repository.IEmailer,
	vault repository.IAuthVault,
) IAuth {
	return &auth{
		authPg:  authPg,
		authRdb: authRdb,
		emailer: emailer,
		vault:   vault,
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

	haveActiveSignupSession, err := a.authRdb.HasActiveECSession(user)

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

func (a *auth) ConfirmCode(user, code string) error {
	isAuthing, err := a.authRdb.HasActiveECSession(user)
	if err != nil {
		return fmt.Errorf("cannot check user auth state: %v", err)
	}

	if !isAuthing {
		return vars.ErrUserIsNotAuthing
	}

	actualAuthCode, err := a.authRdb.GetCode(user)
	if err != nil {
		return fmt.Errorf("cannot get auth code: %v", err)
	}

	if actualAuthCode != code {
		return vars.ErrInvalidAuthCode
	}

	return nil
}

func (a *auth) SignupUnverified(ctx context.Context, user string, pwd string) error {
	pwdHash, err := utils.Hash(pwd)

	if err != nil {
		return fmt.Errorf("cannot create user password: %v", err)
	}

	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      vars.TOTPIssuer,
		AccountName: user,
		Period:      30,
		SecretSize:  20,
		Digits:      otp.DigitsSix,
		Algorithm:   otp.AlgorithmSHA1,
	})

	if err != nil {
		return fmt.Errorf("cannot create TOTP for user: %v", err)
	}

	userModel := &model.User{
		Email:    user,
		Id:       uuid.New().String(),
		SshSign:  utils.NewCert(),
		Deployer: uuid.New().String(),
		Verified: false,
		Banned:   false,
	}

	if err := a.authPg.Signup(ctx, userModel); err != nil {
		return fmt.Errorf("cannot signup user in pg: %v", err)
	}

	if err := a.vault.PutNewUser(ctx, user, pwdHash, key.Secret()); err != nil {
		return fmt.Errorf("cannot signup user in vault: %v", err)
	}

	return nil
}

func (a *auth) VerifyUser(ctx context.Context, user, pwd string) error {
	exists, err := a.authPg.IsUserExists(ctx, user)
	if err != nil {
		return fmt.Errorf("cannot check user existance in db: %v", err)
	}

	if !exists {
		return vars.ErrUserNotFound
	}

	pwdHash, err := a.vault.GetPwdHash(ctx, user)

	if err != nil {
		return fmt.Errorf("cannot get user pwdHash: %v", err)
	}

	if !utils.CheckHash(pwd, pwdHash) {
		return vars.ErrIncorrectPwd
	}

	return nil
}
