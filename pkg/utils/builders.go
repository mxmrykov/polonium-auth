package utils

import (
	"bytes"
	"fmt"
	"time"

	"github.com/mxmrykov/polonium-auth/internal/vars"
)

func BuildVerificationMsg(sender, code, to string) ([]byte, error) {
	now := time.Now().Format(time.RFC1123Z)
	htmlContent := fmt.Sprintf(vars.VerificationCode, code)

	var buffer bytes.Buffer

	buffer.WriteString(fmt.Sprintf("From: %s\r\n", sender))
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
