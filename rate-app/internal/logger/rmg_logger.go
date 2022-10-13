package logger

import (
	"fmt"
)

type LoggerInterface interface {
	Debug(args ...interface{})
	Info(args ...interface{})
	Error(args ...interface{})
}

const (
	debugQueue = "debug"
	infoQueue  = "info"
	erroQueue  = "error"
)

type Logger struct {
	rabbitMqClient *RMQClient
}

func NewLogger(client *RMQClient) LoggerInterface {
	return &Logger{
		rabbitMqClient: client,
	}
}

func (log *Logger) Debug(args ...interface{}) {
	message := fmt.Sprint(args...)
	err := log.rabbitMqClient.Publish(debugQueue, message)
	if err != nil {
		fmt.Println("Failed to publish debug log")
	}
}

func (log *Logger) Info(args ...interface{}) {
	message := fmt.Sprint(args...)
	err := log.rabbitMqClient.Publish(infoQueue, message)
	if err != nil {
		fmt.Println("Failed to publish info log")
	}
}

func (log *Logger) Error(args ...interface{}) {
	message := fmt.Sprint(args...)
	err := log.rabbitMqClient.Publish(erroQueue, message)
	if err != nil {
		fmt.Println("Failed to publish error log")
	}
}
