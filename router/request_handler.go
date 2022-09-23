package router

import (
	"rate-api/model"
)

type RateServiceInterface interface {
	GetRate() (model.Rate, error)
	SetNext(next *RateServiceInterface)
}

type EmailServiceInterface interface {
	AddEmail(email model.Email) error
	GetAllEmails() []string
}

type EmailSendServiceInterface interface {
	SendEmails() error
}

type Handler struct {
	rateHandler      RateHandler
	emailHandler     EmailHandler
	emailSendHandler EmailSendHandler
}

func InitHandler(rateServ RateServiceInterface, emailServ EmailServiceInterface,
	emailSendServ EmailSendServiceInterface) *Handler {
	return &Handler{
		rateHandler:      NewRateHandler(rateServ),
		emailHandler:     NewEmailHandler(emailServ),
		emailSendHandler: NewEmailSendHandler(emailSendServ),
	}
}
