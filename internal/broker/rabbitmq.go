package broker

import (
	"context"
	"encoding/json"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

type EventPublisher interface {
	Publish(ctx context.Context, queueName string, message interface{}) error
}

type RabbitMq struct {
	Conn    *amqp.Connection
	Channel *amqp.Channel
}

func NewRabbitMq(rabbitmqUrl string) (*RabbitMq, error) {
	conn, err := amqp.Dial(rabbitmqUrl)

	if err != nil {
		return nil, fmt.Errorf("Error while connecting rabbitmq: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("Error while creating rabbitmq channel: %w", err)
	}

	return &RabbitMq{
		Conn:    conn,
		Channel: ch,
	}, nil
}

func (rq *RabbitMq) Close() error {
	err := rq.Channel.Close()
	if err != nil {
		return fmt.Errorf("Error while closing channel: %w", err)
	}
	err = rq.Conn.Close()
	if err != nil {
		return fmt.Errorf("Error while closing rabbitmq connection: %w", err)
	}
	return nil
}

func (rq *RabbitMq) DeclareQueue(name string) (amqp.Queue, error) {
	queue, err := rq.Channel.QueueDeclare(name, true, false, false, false, nil)
	if err != nil {
		return amqp.Queue{}, fmt.Errorf("Error while creating queue: %w", err)
	}

	return queue, nil
}

func (rq *RabbitMq) Publish(ctx context.Context, queueName string, message interface{}) error {

	data, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("Error while encoding message to json: %w", err)
	}

	msg := amqp.Publishing{
		ContentType: "application/json",
		Body:        data,
	}

	err = rq.Channel.PublishWithContext(ctx, "", queueName, false, false, msg)
	if err != nil {
		return fmt.Errorf("Error while publishing messages: %w", err)
	}

	return nil
}
