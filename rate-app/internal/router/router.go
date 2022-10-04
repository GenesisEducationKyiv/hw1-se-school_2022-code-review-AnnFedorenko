package router

import (
	"rate-api/internal/handler"

	"github.com/gin-gonic/gin"
)

func RegisterRouter(router *gin.Engine, h *handler.Handler) {
	router.GET("/rate", h.RateHandler.GetRate)
	router.POST("/subscribe", h.EmailHandler.Subscribe)
	router.POST("/sendEmails", h.EmailSendHandler.SendEmails)
}
