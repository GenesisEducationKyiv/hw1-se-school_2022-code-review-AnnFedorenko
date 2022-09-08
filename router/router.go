package router

import (
	"github.com/gin-gonic/gin"
)

func (h *Handler) RegisterRouter(router *gin.Engine) {
	router.GET("/rate", h.getRate)
	router.POST("/subscribe", h.subscribe)
	router.POST("/sendEmails", h.sendEmails)
}
