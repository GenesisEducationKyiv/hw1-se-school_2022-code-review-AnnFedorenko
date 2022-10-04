package handler

import (
	"net/http"
	"rate-api/internal/model"

	"github.com/gin-gonic/gin"
)

type RateServiceInterface interface {
	GetRate() (model.Rate, error)
	SetNext(next *RateServiceInterface)
}

type RateHandler struct {
	rateServ RateServiceInterface
}

func NewRateHandler(rateServ RateServiceInterface) RateHandler {
	return RateHandler{rateServ}
}

func (h *RateHandler) GetRate(context *gin.Context) {
	rate, err := h.rateServ.GetRate()
	if err != nil {
		http.Error(context.Writer, err.Error(), http.StatusInternalServerError)
		return
	}
	context.IndentedJSON(http.StatusOK, rate.Price)
}
