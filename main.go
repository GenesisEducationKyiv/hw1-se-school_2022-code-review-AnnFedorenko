package main

import (
	"fmt"
	"log"
	"rate-api/config"
	"rate-api/repository"
	"rate-api/router"
	"rate-api/service"

	"github.com/gin-gonic/gin"
)

func main() {
	config.LoadConfig()
	serverAddr := fmt.Sprintf("%s:%s", config.Cfg.ServerHost, config.Cfg.ServerPort)

	rateRepo := repository.NewRateRepository()
	rateServ := service.NewRateService(*rateRepo)

	emailRepo := repository.NewEmailRepository(config.Cfg.EmailStorage)
	emailServ := service.NewEmailService(*emailRepo)

	emailSendServ := service.NewEmailSendService(*emailRepo, *rateServ)

	handler := router.NewHandler(rateServ, emailServ, emailSendServ)
	router := gin.Default()
	handler.RegisterRouter(router)
	log.Fatal(router.Run(serverAddr))
}
