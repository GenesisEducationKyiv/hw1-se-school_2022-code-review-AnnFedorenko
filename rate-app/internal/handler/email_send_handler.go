package handler

import (
	"net/http"
	"rate-api/internal/logger"
	"rate-api/internal/model"

	"github.com/gin-gonic/gin"
)

type EmailSendServiceInterface interface {
	SendEmails() error
}

type EmailServiceInterface interface {
	AddEmail(email model.Email) error
	GetAllEmails() []string
}

type EmailSendHandler struct {
	emailSendServ EmailSendServiceInterface
	log           logger.LoggerInterface
}

func NewEmailSendHandler(emailSendServ EmailSendServiceInterface, log logger.LoggerInterface) EmailSendHandler {
	return EmailSendHandler{emailSendServ, log}
}

func (h *EmailSendHandler) SendEmails(context *gin.Context) {
	h.log.Debug("EmailSendHandler started")
	if err := h.emailSendServ.SendEmails(); err != nil {
		h.log.Error("Failed to send emails")
		http.Error(context.Writer, "Failed to send emails", http.StatusInternalServerError)
		return
	}
	context.Status(http.StatusOK)
}
