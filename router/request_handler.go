package router

import (
	"errors"
	"net/http"
	"rate-api/model"
	"rate-api/service"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	rateServ      *service.RateService
	emailServ     *service.EmailService
	emailSendServ *service.EmailSendService
}

func NewHandler(rateServ *service.RateService, emailServ *service.EmailService,
	emailSendServ *service.EmailSendService) *Handler {
	return &Handler{
		rateServ:      rateServ,
		emailServ:     emailServ,
		emailSendServ: emailSendServ,
	}
}

func (h *Handler) getRate(context *gin.Context) {
	rate, err := h.rateServ.GetRate()
	if err != nil {
		http.Error(context.Writer, err.Error(), http.StatusInternalServerError)
		return
	}
	context.IndentedJSON(http.StatusOK, rate.Price)
}

func (h *Handler) subscribe(context *gin.Context) {
	var newEmail model.Email

	if err := context.Bind(&newEmail); err != nil {
		http.Error(context.Writer, err.Error(), http.StatusConflict)
		return
	}

	if err := h.emailServ.AddEmail(newEmail); err != nil {
		if errors.Is(err, service.ErrEmailSubscribed) {
			http.Error(context.Writer, err.Error(), http.StatusConflict)
			return
		}
		if errors.Is(err, service.ErrEmailNotValid) {
			http.Error(context.Writer, err.Error(), http.StatusBadRequest)
			return
		}

		http.Error(context.Writer, err.Error(), http.StatusInternalServerError)
		return
	}

	context.Status(http.StatusCreated)
}

func (h *Handler) sendEmails(context *gin.Context) {
	if err := h.emailSendServ.SendEmails(); err != nil {
		http.Error(context.Writer, "Failed to send emails", http.StatusInternalServerError)
		return
	}
	context.Status(http.StatusOK)
}
