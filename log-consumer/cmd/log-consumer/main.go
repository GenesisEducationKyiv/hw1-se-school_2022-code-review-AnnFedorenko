package main

import (
	"log"
	"os"

	"log-consumer/pkg/rmq"

	"github.com/joho/godotenv"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func main() {
	err := godotenv.Load(".env")
	failOnError(err, "Failed to load env")
	rmqClient, _ := rmq.NewRMQClient(rmq.RMQConfig{
		Host:     os.Getenv("RMQ_HOST"),
		Port:     os.Getenv("RMQ_PORT"),
		User:     os.Getenv("RMQ_USER"),
		Password: os.Getenv("RMQ_PASSWORD"),
		Exchange: os.Getenv("LOG_EXCHANGE"),
	})

	logLevel := os.Getenv("LOG_LEVEL")

	err = rmqClient.QueueDeclare(logLevel)
	failOnError(err, "Failed to declare a queue")
	defer rmqClient.CloseConnection()

	msgs, err := rmqClient.Consume(logLevel)
	failOnError(err, "Failed to register a consumer")
	defer rmqClient.CloseConnection()

	var forever chan struct{}

	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
		}
	}()

	log.Printf(" [*] Waiting for messages.")
	<-forever
}
