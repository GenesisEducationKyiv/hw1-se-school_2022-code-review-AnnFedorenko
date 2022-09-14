package app

import (
	"fmt"
	"log"
	"rate-api/config"
	"rate-api/repository"
	"rate-api/router"
	"rate-api/service"

	"github.com/gin-gonic/gin"
)

func Run() {
	serverAddr := fmt.Sprintf("%s:%s", config.Cfg.ServerHost, config.Cfg.ServerPort)

	rateServ := service.NewRateService()

	emailRepo := repository.NewEmailRepository(config.Cfg.EmailStorage)
	emailServ := service.NewEmailService(emailRepo)

	emailSendServ := service.NewEmailSendService(emailServ, rateServ)

	handler := router.InitHandler(rateServ, emailServ, emailSendServ)
	router := gin.Default()
	handler.RegisterRouter(router)
	log.Fatal(router.Run(serverAddr))
}
