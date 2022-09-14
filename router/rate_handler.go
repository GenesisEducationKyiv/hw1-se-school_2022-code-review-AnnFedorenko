package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type RateHandler struct {
	rateServ RateServiceInterface
}

func NewRateHandler(rateServ RateServiceInterface) RateHandler {
	return RateHandler{rateServ}
}

func (h *RateHandler) getRate(context *gin.Context) {
	rate, err := h.rateServ.GetRate()
	if err != nil {
		http.Error(context.Writer, err.Error(), http.StatusInternalServerError)
		return
	}
	context.IndentedJSON(http.StatusOK, rate.Price)
}
