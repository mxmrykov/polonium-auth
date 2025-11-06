package provider

import (
	"net/smtp"

	"github.com/mxmrykov/polonium-auth/internal/config"
)

type (
	ISMTP interface {
		Send(to string, msg []byte) error
		SenderGetter() string
	}

	SMPT struct {
		auth         smtp.Auth
		host, sender string
	}
)

func NewSmtp(cfg *config.Smtp) ISMTP {
	auth := smtp.PlainAuth("", cfg.Sender, cfg.Password, cfg.Host)

	return &SMPT{
		auth:   auth,
		host:   cfg.Host + ":" + cfg.Port,
		sender: cfg.Sender,
	}
}

func (s *SMPT) Send(to string, msg []byte) error {
	return smtp.SendMail(s.host, s.auth, s.sender, []string{to}, msg)
}

func (s *SMPT) SenderGetter() string {
	return s.sender
}
