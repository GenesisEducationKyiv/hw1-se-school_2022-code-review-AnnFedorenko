package handler

import (
	"net/http"
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
}

func NewEmailSendHandler(emailSendServ EmailSendServiceInterface) EmailSendHandler {
	return EmailSendHandler{emailSendServ}
}

func (h *EmailSendHandler) SendEmails(context *gin.Context) {
	if err := h.emailSendServ.SendEmails(); err != nil {
		http.Error(context.Writer, "Failed to send emails", http.StatusInternalServerError)
		return
	}
	context.Status(http.StatusOK)
}
