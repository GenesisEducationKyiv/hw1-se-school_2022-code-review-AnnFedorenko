package app

import (
	"fmt"
	"log"
	"rate-api/config"
	"rate-api/handler"
	"rate-api/mailclient"
	"rate-api/repository"
	route "rate-api/router"
	"rate-api/service"
	"rate-api/service/rate"
	"rate-api/service/rate/cache"
	"rate-api/service/rate/logger"

	"github.com/gin-gonic/gin"
)

var RateCreators = map[string]rate.RateFactory{
	"binance":  rate.NewBinanceRateCreator(),
	"coinbase": rate.NewCoinbaseRateCreator(),
	"coingate": rate.NewCoingateRateCreator(),
}

func Run() {
	serverAddr := fmt.Sprintf("%s:%s", config.Cfg.ServerHost, config.Cfg.ServerPort)

	rateServ := InitRateProvider()

	emailRepo := repository.NewEmailRepository(config.Cfg.EmailStorage)
	emailServ := service.NewEmailService(emailRepo)

	emailSendServ := service.NewEmailSendService(emailServ, rateServ, mailclient.NewSMTPClient(config.Cfg))

	handler := handler.InitHandler(rateServ, emailServ, emailSendServ)
	router := gin.Default()
	route.RegisterRouter(router, handler)
	log.Fatal(router.Run(serverAddr))
}

func InitRateProvider() handler.RateServiceInterface {
	defaultProviderName := config.Cfg.CryploCurrencyProvider
	defaultRateServ := RateCreators[defaultProviderName].GetRateService()

	prev := defaultRateServ
	for k, v := range RateCreators {
		if k != defaultProviderName {
			provider := v.GetRateService()
			prev.SetNext(&provider)
			prev = provider
		}
	}

	rateLogServ := logger.NewLogRateServiceDecorator(defaultRateServ)
	rateCacheServ := cache.NewCacheRateServiceProxy(rateLogServ)

	return rateCacheServ
}
