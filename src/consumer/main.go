package main

import (
	"encoding/json"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"time"
)

func main() {
	cfg := readCfg()
	ch := NewClickhouse()
	ch.Connect(cfg.Clickhouse)

	defer func() {
		ch.Close()
	}()

	rabbitMq := NewRabbitMq(cfg.Rabbit)
	rabbitMq.Connect()
	rabbitMq.Consume(func(msg amqp.Delivery) {
		data := Measurement{}
		parseErr := json.Unmarshal(msg.Body[:], &data)
		if parseErr != nil {
			log.Printf("Failed umarshall data. Error: %s", parseErr)
			return
		}

		insertErr := ch.Insert(data.Temperature, data.Humidity, time.Unix(data.Timestamp, 0))
		if insertErr != nil {
			log.Printf("Failed insert data. Error: %s", insertErr)
			return
		}
	})

	defer func() {
		rabbitMq.Close()
	}()

	var forever chan struct{}
	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
