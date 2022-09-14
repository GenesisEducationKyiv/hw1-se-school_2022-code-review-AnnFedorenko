package router

import (
	"github.com/gin-gonic/gin"
)

func (h *Handler) RegisterRouter(router *gin.Engine) {
	router.GET("/rate", h.rateHandler.getRate)
	router.POST("/subscribe", h.emailHandler.subscribe)
	router.POST("/sendEmails", h.emailSendHandler.sendEmails)
}
