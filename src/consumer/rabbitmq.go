package main

import amqp "github.com/rabbitmq/amqp091-go"

type MessageHandler func(message amqp.Delivery)

type RabbitMqClient struct {
	ConnectionString string
	Queue            string
	Connection       *amqp.Connection
	Channel          *amqp.Channel
}

func NewRabbitMq(cfg RabbitMq) *RabbitMqClient {
	client := new(RabbitMqClient)
	client.Queue = cfg.Queue
	client.ConnectionString = cfg.getConnStr()
	return client
}

func (c *RabbitMqClient) Connect() {
	conn, err := amqp.Dial(c.ConnectionString)
	if err != nil {
		panic(err)
	}

	c.Connection = conn
}

func (c *RabbitMqClient) Consume(handler MessageHandler) {

	ch, err := c.Connection.Channel()
	if err != nil {
		panic(err)
	}

	c.Channel = ch

	q, err := c.Channel.QueueDeclare(
		c.Queue,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		panic(err)
	}

	messages, err := c.Channel.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		panic(err)
	}

	go func() {
		for message := range messages {
			handler(message)
		}
	}()
}

func (c *RabbitMqClient) Close() {
	_ = c.Connection.Close()
	_ = c.Channel.Close()
}
