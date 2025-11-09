package service

import (
	"context"
	"fmt"
	"slices"

	"github.com/mxmrykov/polonium-auth/internal/repository"
	"github.com/mxmrykov/polonium-auth/internal/vars"
	"github.com/mxmrykov/polonium-auth/pkg/utils"
	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
)

type (
	ITOTP interface {
		CreateUserQR(ctx context.Context, user string) ([]byte, error)
		IsCodeCorrect(ctx context.Context, user, code string) (bool, error)
	}

	ttp struct {
		vault repository.IAuthVault
	}
)

func NewTOTP(vault repository.IAuthVault) ITOTP {
	return &ttp{
		vault: vault,
	}
}

func (t *ttp) CreateUserQR(ctx context.Context, user string) ([]byte, error) {
	secret, err := t.vault.GetTOTPSecret(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("cannot get user TOTP secret: %v", err)
	}

	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      vars.TOTPIssuer,
		AccountName: user,
		Period:      30,
		SecretSize:  20,
		Digits:      otp.DigitsSix,
		Algorithm:   otp.AlgorithmSHA1,
		Secret:      []byte(secret),
	})

	if err != nil {
		return nil, fmt.Errorf("cannot create TOTP for user: %v", err)
	}

	return utils.GenerateQR(key.URL())
}

func (t *ttp) IsCodeCorrect(ctx context.Context, user, code string) (bool, error) {
	secret, err := t.vault.GetTOTPSecret(ctx, user)
	if err != nil {
		return false, fmt.Errorf("cannot get user TOTP secret: %v", err)
	}

	codes := utils.GetTOTPCodes(secret, 30)
	return slices.Contains(codes, code), nil
}
