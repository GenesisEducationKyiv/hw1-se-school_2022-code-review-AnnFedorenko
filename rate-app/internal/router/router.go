package router

import (
	handle "rate-api/internal/handler"

	"github.com/dtm-labs/dtm/dtmutil"
	"github.com/gin-gonic/gin"
)

func RegisterRouter(router *gin.Engine, handler *handle.Handler) {
	router.GET("/rate", handler.RateHandler.GetRate)
	router.POST("/subscribe", handler.EmailHandler.EmailPlusCustomerSubscribe)
	router.POST("/sendEmails", handler.EmailSendHandler.SendEmails)

	router.POST("/register-subscriber",
		dtmutil.WrapHandler2(handler.EmailHandler.Subscribe))
	router.POST("/register-subscriber-compensate",
		dtmutil.WrapHandler2(handler.EmailHandler.SubscribeEmailCompensate))
}
