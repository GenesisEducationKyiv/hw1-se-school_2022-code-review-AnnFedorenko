package logger

import (
	"context"
	"fmt"
	"net"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RMQConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Exchange string
}

type RMQClient struct {
	config     RMQConfig
	connection *amqp.Connection
	channel    *amqp.Channel
}

func NewRMQClient(cfg RMQConfig) (*RMQClient, error) {
	url := fmt.Sprintf("amqp://%s:%s@%s/",
		cfg.User, cfg.Password, net.JoinHostPort(cfg.Host, cfg.Port))

	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}

	channel, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	err = channel.ExchangeDeclare(
		cfg.Exchange, // name
		"direct",     // kind
		true,         // durable
		false,        // autoDelete
		false,        // internal
		false,        // noWait
		nil,          // args
	)
	if err != nil {
		return nil, err
	}

	client := &RMQClient{
		config:     cfg,
		connection: conn,
		channel:    channel,
	}

	return client, err
}

func (c *RMQClient) Publish(queueName string, messageBody string) error {
	err := queueBind(queueName, c)
	message := amqp.Publishing{
		DeliveryMode: amqp.Persistent,
		Timestamp:    time.Now(),
		ContentType:  "text/plain",
		Body:         []byte(messageBody),
	}
	if err != nil {
		return err
	}

	err = c.channel.PublishWithContext(
		context.Background(),
		c.config.Exchange, // exchange
		queueName,         // routing key
		true,              // mandatory
		false,             // immediate
		message)

	return err
}

func (c *RMQClient) CloseConnection() error {
	if err := c.connection.Close(); err != nil {
		return err
	}
	if err := c.channel.Close(); err != nil {
		return err
	}
	return nil
}

func queueBind(queueName string, client *RMQClient) error {
	queue, err := client.channel.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // autoDelete
		false,     // exclusive
		false,     // noWait
		nil,       // args
	)
	if err != nil {
		return err
	}

	err = client.channel.QueueBind(
		queue.Name,             // queue name
		queue.Name,             // routing key
		client.config.Exchange, // exchange
		false,                  // noWait
		nil,                    // args
	)
	return err
}
