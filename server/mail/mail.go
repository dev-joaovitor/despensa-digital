package mail

import (
	"fmt"
	"net/smtp"
)

type MailService struct {
	host string
	port string
	fromEmail string
	auth smtp.Auth
}

func NewEmailService(host, port, username, password, fromEmail string) *MailService {
	var auth smtp.Auth = nil

	if username != "" && password != "" {
		auth = smtp.PlainAuth(fromEmail, username, password, host)
	}

	return &MailService{
		host: host,
		port: port,
		fromEmail: fromEmail,
		auth: auth,
	}
}

func (s *MailService) SendPasswordReset(toEmail string, recoveryCode string) error {
	to := []string{toEmail}

	subject := "Subject: Redefinir Senha"
	mime := "MIME-version: 1.0;\nContent-Type: text/plain; charset=\"UTF-8\";\r\n"

	body := fmt.Sprintf(
		"To: %s\r\n"+
		"%s"+
		"%s\r\n\r\n"+
		"Seu código de verificação é: %s\r\n\r\n"+
		"Este código irá expirar em 5 minutos",
		toEmail, mime, subject, recoveryCode,
	)

	addr := fmt.Sprintf("%s:%s", s.host, s.port)

	err := smtp.SendMail(addr, s.auth, s.fromEmail, to, []byte(body))
	if err != nil {
		return fmt.Errorf("SMTP transmission failed: %w", err)
	}
	return nil
}
