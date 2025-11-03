package internal

import (
	"fmt"
	"net/smtp"
)

type EmailSender struct {
	host     string
	port     int
	username string
	password string
	from     string
	auth     smtp.Auth
}

func NewEmailSender(host string, port int, username, password, from string) *EmailSender {
	if host == "" || port == 0 || from == "" {
		return nil
	}

	var auth smtp.Auth
	if username != "" {
		auth = smtp.PlainAuth("", username, password, host)
	}

	return &EmailSender{
		host:     host,
		port:     port,
		username: username,
		password: password,
		from:     from,
		auth:     auth,
	}
}

func (s *EmailSender) Send(to, subject, body string) error {
	if s == nil {
		return fmt.Errorf("email sender not configured")
	}
	if to == "" {
		return fmt.Errorf("missing recipient email")
	}

	addr := fmt.Sprintf("%s:%d", s.host, s.port)
	msg := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\nMIME-Version: 1.0\r\nContent-Type: text/plain; charset=\"utf-8\"\r\n\r\n%s", s.from, to, subject, body)

	return smtp.SendMail(addr, s.auth, s.from, []string{to}, []byte(msg))
}
