package repository

import (
	"fmt"

	"github.com/mxmrykov/polonium-auth/internal/config"
	"github.com/mxmrykov/polonium-auth/internal/provider"
	"github.com/mxmrykov/polonium-auth/internal/vars"
	"github.com/mxmrykov/polonium-auth/pkg/utils"
)

type (
	IEmailer interface {
		SendVerificationCode(code, to string) error
	}

	emailer struct {
		smtp provider.ISMTP
	}
)

func NewEmailer(cfg *config.Smtp) IEmailer {
	return &emailer{
		smtp: provider.NewSmtp(cfg),
	}
}

func (e *emailer) SendVerificationCode(code, to string) error {
	if !utils.IsEmailValid(to) {
		return vars.ErrInvalidEmail
	}

	msg, err := utils.BuildVerificationMsg(e.smtp.SenderGetter(), code, to)
	if err != nil {
		return fmt.Errorf("cannot build verification msg: %v", err)
	}

	return e.smtp.Send(to, msg)
}
