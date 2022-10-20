package router

import (
	"orders/internal/handler"

	"github.com/dtm-labs/dtm/dtmutil"
	"github.com/gin-gonic/gin"
)

func InitRoutes(router *gin.Engine, h *handler.Handler) {
	router.POST("/create", h.OrderHandler.CreateOrder)
	router.POST("/register-order",
		dtmutil.WrapHandler2(h.OrderHandler.RegisterOrder))

	router.POST("/register-order-compensate", dtmutil.WrapHandler2(h.OrderHandler.RegisterOrderCompensate))
}
