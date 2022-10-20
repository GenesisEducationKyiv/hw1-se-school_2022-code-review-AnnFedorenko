package router

import (
	"customers/internal/handler"

	"github.com/dtm-labs/dtm/dtmutil"
	"github.com/gin-gonic/gin"
)

func InitRoutes(router *gin.Engine, h *handler.Handler) {
	router.POST("/register-customer", dtmutil.WrapHandler2(h.CustomerHandler.RegisterCustomer))
	router.POST("/register-customer-compensate",
		dtmutil.WrapHandler2(h.CustomerHandler.RegisterCustomerCompensate))

	router.POST("/withdraw-money", dtmutil.WrapHandler2(h.CustomerTransactionHandler.WithdrawMoney))
	router.POST("/withdraw-money-compensate",
		dtmutil.WrapHandler2(h.CustomerTransactionHandler.WithdrawMoneyCompensate))
}
