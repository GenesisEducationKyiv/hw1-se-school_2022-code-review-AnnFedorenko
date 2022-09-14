package service

import (
	"fmt"
	"net/smtp"
	"rate-api/config"
	"rate-api/router"
)

type EmailSendService struct {
	es router.EmailServiceInterface
	rs router.RateServiceInterface
}

func NewEmailSendService(es router.EmailServiceInterface,
	rs router.RateServiceInterface) router.EmailSendServiceInterface {
	return &EmailSendService{es, rs}
}

func (s *EmailSendService) SendEmails() error {
	var cfg = config.Cfg
	user := cfg.SMTPUsername
	password := cfg.SMTPPassword
	addr := fmt.Sprintf("%s:%d", cfg.SMTPHost, cfg.SMTPPort)
	host := cfg.SMTPHost
	sender := "BTC rate app"

	receiver := s.es.GetAllEmails()
	if receiver == nil {
		return nil
	}

	rate, err := s.rs.GetRate()
	if err != nil {
		return err
	}

	msg := []byte("From: Bitcoin rate helper\r\n" +
		"Subject: BTCUAH Rate\r\n\r\n" +
		rate.Price +
		"\r\n")

	auth := smtp.PlainAuth("", user, password, host)

	if err = smtp.SendMail(addr, auth, sender, receiver, msg); err != nil {
		return err
	}

	return nil
}
