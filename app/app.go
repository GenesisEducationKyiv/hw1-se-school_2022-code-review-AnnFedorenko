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

var RateCreators = map[string]rate.RateFactory{
	"binance":  &rate.BinanceRateCreator{},
	"coinbase": &rate.CoinbaseRateCreator{},
	"coingate": &rate.CoingateRateCreator{},
}

func Run() {
	serverAddr := fmt.Sprintf("%s:%s", config.Cfg.ServerHost, config.Cfg.ServerPort)

	rateServ := InitRateProvider()

	emailRepo := repository.NewEmailRepository(config.Cfg.EmailStorage)
	emailServ := service.NewEmailService(emailRepo)

	emailSendServ := service.NewEmailSendService(emailServ, rateServ)

	handler := router.InitHandler(rateServ, emailServ, emailSendServ)
	router := gin.Default()
	handler.RegisterRouter(router)
	log.Fatal(router.Run(serverAddr))
}

func InitRateProvider() router.RateServiceInterface {
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

	rateLogServ := logger.NewLogRateService(defaultRateServ)
	rateCacheServ := cache.NewCacheRateService(rateLogServ)

	return rateCacheServ
}
