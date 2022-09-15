package app

import (
	"fmt"
	"log"
	"rate-api/cache"
	"rate-api/config"
	"rate-api/logger"
	"rate-api/repository"
	"rate-api/router"
	"rate-api/service"
	"rate-api/service/rate"

	"github.com/gin-gonic/gin"
)

func Run() {
	serverAddr := fmt.Sprintf("%s:%s", config.Cfg.ServerHost, config.Cfg.ServerPort)

	rate1Serv := rate.GetRateService("binance")
	rate2Serv := rate.GetRateService("coinbase")
	rate3Serv := rate.GetRateService("coingate")
	rate1Serv.SetNext(&rate2Serv)
	rate2Serv.SetNext(&rate3Serv)

	rateLogServ := logger.NewLogRateService(rate1Serv)
	rateCacheServ := cache.NewCacheRateService(rateLogServ)

	emailRepo := repository.NewEmailRepository(config.Cfg.EmailStorage)
	emailServ := service.NewEmailService(emailRepo)

	emailSendServ := service.NewEmailSendService(emailServ, rateCacheServ)

	handler := router.InitHandler(rateCacheServ, emailServ, emailSendServ)
	router := gin.Default()
	handler.RegisterRouter(router)
	log.Fatal(router.Run(serverAddr))
}
