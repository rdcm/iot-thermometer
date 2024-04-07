package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"

	"consumer/clickhouse"
	"consumer/config"
	"consumer/handler"
	"consumer/rabbitmq"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	if err := appMain(ctx); err != nil {
		log.Fatal(err)
	}
}

func appMain(ctx context.Context) error {
	cfg, err := config.ReadCfg()
	if err != nil {
		return fmt.Errorf("read config: %w", err)
	}

	clickhouseClient := clickhouse.New(
		cfg.Clickhouse.Hostname,
		cfg.Clickhouse.Port,
		cfg.Clickhouse.Database,
		cfg.Clickhouse.Username,
		cfg.Clickhouse.Password,
	)
	if err := clickhouseClient.Connect(ctx); err != nil {
		return fmt.Errorf("clickhouse client connect: %w", err)
	}
	defer func() {
		if err := clickhouseClient.Close(); err != nil {
			log.Printf("clickhouse client close: %v", err)
		}
	}()

	measurementHandler := handler.New(clickhouseClient)

	rabbitmqClient := rabbitmq.New(
		cfg.Rabbit.GetConnStr(),
		cfg.Rabbit.Queue,
		measurementHandler,
	)

	if err := rabbitmqClient.Run(ctx); err != nil {
		return fmt.Errorf("run rabbitmq client: %w", err)
	}

	return nil
}
