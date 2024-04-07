package rabbitmq

import (
	"context"
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type MessageHandler interface {
	Handle(ctx context.Context, message amqp.Delivery) error
}

type Client struct {
	connStr string
	queue   string
	handler MessageHandler
}

func New(
	connStr string,
	queue string,
	handler MessageHandler,
) *Client {
	return &Client{
		connStr: connStr,
		queue:   queue,
		handler: handler,
	}
}

func (c *Client) Run(ctx context.Context) error {
	conn, err := amqp.Dial(c.connStr)
	if err != nil {
		return fmt.Errorf("rabbitmq connect: %w", err)
	}
	defer func() {
		if err := conn.Close(); err != nil {
			log.Printf("rabbitmq connection close: %v", err)
		}
	}()

	ch, err := conn.Channel()
	if err != nil {
		return fmt.Errorf("rabbitmq get channel: %w", err)
	}
	defer func() {
		if err := ch.Close(); err != nil {
			log.Printf("rabbitmq channel close: %v", err)
		}
	}()

	q, err := ch.QueueDeclare(
		c.queue,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("rabbitmq queue declare: %w", err)
	}

	messages, err := ch.ConsumeWithContext(
		ctx,
		q.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("rabbitmq consume: %w", err)
	}

	for message := range messages {
		if err := c.handler.Handle(ctx, message); err != nil {
			return fmt.Errorf("rabbitmq handler: %w", err)
		}

		if err := message.Ack(false); err != nil {
			return fmt.Errorf("rabbitmq ack message: %w", err)
		}
	}

	return nil
}
