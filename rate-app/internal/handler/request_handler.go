package handler

import "rate-api/internal/logger"

type Handler struct {
	RateHandler      RateHandler
	EmailHandler     EmailHandler
	EmailSendHandler EmailSendHandler
	log              logger.LoggerInterface
}

func InitHandler(rateServ RateServiceInterface, emailServ EmailServiceInterface,
	emailSendServ EmailSendServiceInterface, log logger.LoggerInterface) *Handler {
	return &Handler{
		RateHandler:      NewRateHandler(rateServ, log),
		EmailHandler:     NewEmailHandler(emailServ, log),
		EmailSendHandler: NewEmailSendHandler(emailSendServ, log),
	}
}
