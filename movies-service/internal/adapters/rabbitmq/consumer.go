package rabbitmq

import (
	"context"
	"encoding/json"
	"log"

	"github.com/FranciscoHonorat/movies/movies-service/internal/core/domain"
	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQConsumer struct {
	channel   *amqp.Channel
	queueName string
}

func NewRabbitMQConsumer(url string, queueName string) (*RabbitMQConsumer, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}
	channel, err := conn.Channel()
	if err != nil {
		return nil, err
	}
	_, err = channel.QueueDeclare(
		queueName,
		true,  // durable
		false, // autoDelete
		false, // exclusive
		false, // noWait
		nil,   // args
	)
	if err != nil {
		return nil, err
	}

	return &RabbitMQConsumer{
		channel:   channel,
		queueName: queueName,
	}, nil
}

func (c *RabbitMQConsumer) Consume(ctx context.Context, handler func(domain.Movie) error) error {
	msgs, err := c.channel.Consume(
		c.queueName,
		"",    // consumer
		true,  // autoAck
		false, // exclusive
		false, // noLocal
		false, // noWait
		nil,   // args
	)
	if err != nil {
		return err
	}

	for msg := range msgs {
		var movieMessage domain.Movie
		if err := json.Unmarshal(msg.Body, &movieMessage); err != nil {
			log.Printf("Error unmarshaling message: %v", err)
			continue
		}
		if err := handler(movieMessage); err != nil {
			log.Printf("Error handling message: %v", err)
			continue
		}
	}

	return nil
}
