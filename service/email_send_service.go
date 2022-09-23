package service

import "rate-api/handler"

type EmailClient interface {
	Send(from string, to []string, msg []byte) error
}

type EmailSendService struct {
	es     handler.EmailServiceInterface
	rs     handler.RateServiceInterface
	client EmailClient
}

func NewEmailSendService(es handler.EmailServiceInterface,
	rs handler.RateServiceInterface, client EmailClient) handler.EmailSendServiceInterface {
	return &EmailSendService{es, rs, client}
}

func (s *EmailSendService) SendEmails() error {
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
