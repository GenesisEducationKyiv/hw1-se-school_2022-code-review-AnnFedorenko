package handler

type Handler struct {
	RateHandler      RateHandler
	EmailHandler     EmailHandler
	EmailSendHandler EmailSendHandler
}

func InitHandler(rateServ RateServiceInterface, emailServ EmailServiceInterface,
	emailSendServ EmailSendServiceInterface) *Handler {
	return &Handler{
		RateHandler:      NewRateHandler(rateServ),
		EmailHandler:     NewEmailHandler(emailServ),
		EmailSendHandler: NewEmailSendHandler(emailSendServ),
	}
}
