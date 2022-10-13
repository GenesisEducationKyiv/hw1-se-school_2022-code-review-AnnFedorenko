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
	rmqClient, err := rmq.NewRMQClient(rmq.RMQConfig{
		Host:     os.Getenv("RMQ_HOST"),
		Port:     os.Getenv("RMQ_PORT"),
		User:     os.Getenv("RMQ_USER"),
		Password: os.Getenv("RMQ_PASSWORD"),
		Exchange: os.Getenv("LOG_EXCHANGE"),
	})
	failOnError(err, "Failed to declare rabbitMQ consumer")
	defer func() {
		if err := rmqClient.CloseConnection(); err != nil {
			panic(err)
		}
	}()

	logLevel := os.Getenv("LOG_LEVEL")

	err = rmqClient.QueueDeclare(logLevel)
	failOnError(err, "Failed to declare a queue")

	msgs, err := rmqClient.Consume(logLevel)
	failOnError(err, "Failed to register a consumer")

	var forever chan struct{}

	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
		}
	}()

	log.Printf(" [*] Waiting for messages.")
	<-forever
}
