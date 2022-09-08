package service

import (
	"fmt"
	"net/smtp"
	"rate-api/config"
	"rate-api/repository"
)

type EmailSendService struct {
	repo repository.EmailRepository
	rs   RateService
}

func NewEmailSendService(repo repository.EmailRepository, rs RateService) *EmailSendService {
	return &EmailSendService{repo: repo,
		rs: rs}
}

func (s *EmailSendService) SendEmails() error {
	var cfg = config.Cfg
	user := cfg.SMTPUsername
	password := cfg.SMTPPassword
	addr := fmt.Sprintf("%s:%d", cfg.SMTPHost, cfg.SMTPPort)
	host := cfg.SMTPHost
	sender := "BTC rate app"

	receiver := s.repo.GetAllEmails()
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
