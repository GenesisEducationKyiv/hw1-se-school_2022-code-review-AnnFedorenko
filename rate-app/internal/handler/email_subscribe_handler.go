package handler

import (
	"errors"
	"fmt"
	"net/http"
	"rate-api/internal/logger"
	"rate-api/internal/mailerror"
	"rate-api/internal/model"

	"github.com/gin-gonic/gin"
)

type EmailHandler struct {
	emailServ EmailServiceInterface
	log       logger.LoggerInterface
}

func NewEmailHandler(emailServ EmailServiceInterface, log logger.LoggerInterface) EmailHandler {
	return EmailHandler{emailServ, log}
}

func (h *EmailHandler) Subscribe(context *gin.Context) {
	h.log.Debug("EmailHandler started")
	var newEmail model.Email

	if err := context.Bind(&newEmail); err != nil {
		h.log.Error(fmt.Sprintf("Error occurred during adding subscription: %s", err.Error()))
		http.Error(context.Writer, err.Error(), http.StatusConflict)
		return
	}

	if err := h.emailServ.AddEmail(newEmail); err != nil {
		if errors.Is(err, mailerror.ErrEmailSubscribed) {
			h.log.Debug("Email already subscribed")
			http.Error(context.Writer, err.Error(), http.StatusConflict)
			return
		}
		if errors.Is(err, mailerror.ErrEmailNotValid) {
			h.log.Error("Not valid email provided")
			http.Error(context.Writer, err.Error(), http.StatusBadRequest)
			return
		}

		h.log.Error(fmt.Sprintf("Error occurred during adding subscription: %s", err.Error()))
		http.Error(context.Writer, err.Error(), http.StatusInternalServerError)
		return
	}

	context.Status(http.StatusCreated)
}
