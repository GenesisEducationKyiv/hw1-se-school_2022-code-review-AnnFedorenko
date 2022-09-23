package service

import (
	"rate-api/mailclient"
	"rate-api/router"
)

type EmailSendService struct {
	es     router.EmailServiceInterface
	rs     router.RateServiceInterface
	client mailclient.EmailClient
}

func NewEmailSendService(es router.EmailServiceInterface,
	rs router.RateServiceInterface, client mailclient.EmailClient) router.EmailSendServiceInterface {
	return &EmailSendService{es, rs, client}
}

func (s *EmailSendService) SendEmails() error {
	//	var cfg = config.Cfg
	sender := "BTC rate app"

	receiver := s.es.GetAllEmails()
	if receiver == nil {
		return nil
	}

	msg, err := s.createEmailMessage()
	if err != nil {
		return err
	}

	err = s.client.Send(sender, receiver, msg)
	if err != nil {
		return err
	}
	return nil
}

func (s *EmailSendService) createEmailMessage() ([]byte, error) {
	rate, err := s.rs.GetRate()
	if err != nil {
		return nil, err
	}

	msg := []byte("From: Bitcoin rate helper\r\n" +
		"Subject: BTCUAH Rate\r\n\r\n" +
		rate.Price +
		"\r\n")

	return msg, nil
}
