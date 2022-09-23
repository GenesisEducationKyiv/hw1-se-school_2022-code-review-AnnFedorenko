package router

import (
	"errors"
	"net/http"
	"rate-api/mailerror"
	"rate-api/model"

	"github.com/gin-gonic/gin"
)

type EmailHandler struct {
	emailServ EmailServiceInterface
}

func NewEmailHandler(emailServ EmailServiceInterface) EmailHandler {
	return EmailHandler{emailServ}
}

func (h *EmailHandler) subscribe(context *gin.Context) {
	var newEmail model.Email

	if err := context.Bind(&newEmail); err != nil {
		http.Error(context.Writer, err.Error(), http.StatusConflict)
		return
	}

	if err := h.emailServ.AddEmail(newEmail); err != nil {
		if errors.Is(err, mailerror.ErrEmailSubscribed) {
			http.Error(context.Writer, err.Error(), http.StatusConflict)
			return
		}
		if errors.Is(err, mailerror.ErrEmailNotValid) {
			http.Error(context.Writer, err.Error(), http.StatusBadRequest)
			return
		}

		http.Error(context.Writer, err.Error(), http.StatusInternalServerError)
		return
	}

	context.Status(http.StatusCreated)
}
