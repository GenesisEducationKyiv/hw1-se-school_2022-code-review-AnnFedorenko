package mailclient

import (
	"fmt"
	"net/smtp"
	"rate-api/config"
	"rate-api/service"
)

type SMTPClient struct {
	user     string
	password string
	host     string
	port     int
}

func NewSMTPClient(cfg config.Config) service.EmailClient {
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
