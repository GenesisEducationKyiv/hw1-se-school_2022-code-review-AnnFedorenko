package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"rate-api/config"
	"rate-api/service"

	"github.com/gin-gonic/gin"
)

func getRate(context *gin.Context) {
	rate, err := service.GetRateFromBinance()
	if err != nil {
		http.Error(context.Writer, err.Error(), http.StatusInternalServerError)
		return
	}
	context.IndentedJSON(http.StatusOK, rate.Price)
}

func subscribe(context *gin.Context) {
	var newEmail service.Email

	if err := context.Bind(&newEmail); err != nil {
		http.Error(context.Writer, err.Error(), http.StatusConflict)
		return
	}

	if err := service.AddEmail(newEmail); err != nil {
		if errors.Is(err, service.ErrEmailSubscribed) {
			http.Error(context.Writer, err.Error(), http.StatusConflict)
			return
		}
		if errors.Is(err, service.ErrEmailNotValid) {
			http.Error(context.Writer, err.Error(), http.StatusBadRequest)
			return
		}

		http.Error(context.Writer, err.Error(), http.StatusInternalServerError)
		return
	}

	context.Status(http.StatusCreated)
}

func sendEmails(context *gin.Context) {
	if err := service.SendEmails(); err != nil {
		http.Error(context.Writer, "Failed to send emails", http.StatusInternalServerError)
		return
	}
	context.Status(http.StatusOK)
}

func main() {
	config.LoadConfig()
	serverAddr := fmt.Sprintf("%s:%s", config.Cfg.ServerHost, config.Cfg.ServerPort)

	router := gin.Default()
	router.GET("/rate", getRate)
	router.POST("/subscribe", subscribe)
	router.POST("/sendEmails", sendEmails)
	log.Fatal(router.Run(serverAddr))
}
