package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type EmailSendHandler struct {
	emailSendServ EmailSendServiceInterface
}

func NewEmailSendHandler(emailSendServ EmailSendServiceInterface) EmailSendHandler {
	return EmailSendHandler{emailSendServ}
}

func (h *EmailSendHandler) sendEmails(context *gin.Context) {
	if err := h.emailSendServ.SendEmails(); err != nil {
		http.Error(context.Writer, "Failed to send emails", http.StatusInternalServerError)
		return
	}
	context.Status(http.StatusOK)
}
