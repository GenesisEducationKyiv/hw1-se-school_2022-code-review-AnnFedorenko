package mailclient

import (
	"fmt"
	"net/smtp"
	"rate-api/config"
)

type SMTPClient struct {
	user     string
	password string
	host     string
	port     int
}

type EmailClient interface {
	Send(from string, to []string, msg []byte) error
}

func NewSMTPClient(cfg config.Config) EmailClient {
	return &SMTPClient{
		user:     cfg.SMTPUsername,
		password: cfg.SMTPPassword,
		host:     cfg.SMTPHost,
		port:     cfg.SMTPPort,
	}
}

func (c *SMTPClient) Send(from string, to []string, msg []byte) error {
	auth := smtp.PlainAuth("", c.user, c.password, c.host)
	addr := fmt.Sprintf("%s:%d", c.host, c.port)

	if err := smtp.SendMail(addr, auth, from, to, msg); err != nil {
		return err
	}

	return nil
}
