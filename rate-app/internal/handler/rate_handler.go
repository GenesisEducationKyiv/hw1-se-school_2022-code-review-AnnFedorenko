package handler

import (
	"net/http"
	"rate-api/internal/logger"
	"rate-api/internal/model"

	"github.com/gin-gonic/gin"
)

type RateServiceInterface interface {
	GetRate() (model.Rate, error)
	SetNext(next *RateServiceInterface)
}

type RateHandler struct {
	rateServ RateServiceInterface
	log      logger.LoggerInterface
}

func NewRateHandler(rateServ RateServiceInterface, log logger.LoggerInterface) RateHandler {
	return RateHandler{rateServ, log}
}

func (h *RateHandler) GetRate(context *gin.Context) {
	h.log.Debug("RateHandler started")
	rate, err := h.rateServ.GetRate()
	if err != nil {
		h.log.Error("Failed to get rate")
		http.Error(context.Writer, err.Error(), http.StatusInternalServerError)
		return
	}
	context.IndentedJSON(http.StatusOK, rate.Price)
}
