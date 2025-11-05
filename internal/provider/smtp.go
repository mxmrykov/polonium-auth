package provider

import (
	"bytes"
	"fmt"
	"net/mail"
	"net/smtp"
	"time"

	"github.com/mxmrykov/polonium-auth/internal/config"
	"github.com/mxmrykov/polonium-auth/internal/vars"
)

type (
	ISMTP interface {
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

func (s *SMPT) SendVerificationCode(code, to string) error {
	if !valid(to) {
		return fmt.Errorf("invalid to address: %s", to)
	}

	msg, err := s.buildMsg(code, to)
	if err != nil {
		return fmt.Errorf("cannot build verification msg: %v", err)
	}

	return s.send(to, msg)
}

func (s *SMPT) send(to string, msg []byte) error {
	return smtp.SendMail(s.host, s.auth, s.sender, []string{to}, msg)
}

func (s *SMPT) buildMsg(code, to string) ([]byte, error) {
	now := time.Now().Format(time.RFC1123Z)
	htmlContent := fmt.Sprintf(vars.VerificationCode, code)

	var buffer bytes.Buffer

	buffer.WriteString(fmt.Sprintf("From: %s\r\n", s.sender))
	buffer.WriteString(fmt.Sprintf("To: %s\r\n", to))
	buffer.WriteString("Subject: Verify Your Email Address\r\n")
	buffer.WriteString(fmt.Sprintf("Date: %s\r\n", now))
	buffer.WriteString("MIME-Version: 1.0\r\n")
	buffer.WriteString("Content-Type: text/html; charset=\"utf-8\"\r\n")
	buffer.WriteString("Content-Transfer-Encoding: 7bit\r\n")

	buffer.WriteString(htmlContent)
	buffer.WriteString("\r\n")

	return buffer.Bytes(), nil
}

func valid(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}
