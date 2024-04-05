package main

import (
	"context"
	"fmt"
	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"time"
)

type ClickhouseClient struct {
	Connection driver.Conn
	Context    context.Context
}

func NewClickhouse() *ClickhouseClient {
	var ch = new(ClickhouseClient)
	return ch
}

func (c *ClickhouseClient) Connect(cfg Clickhouse) {
	var (
		ctx       = context.Background()
		conn, err = clickhouse.Open(&clickhouse.Options{
			Addr: []string{fmt.Sprintf("%s:%d", cfg.Hostname, cfg.Port)},
			Auth: clickhouse.Auth{
				Database: cfg.Database,
				Username: cfg.Username,
				Password: cfg.Password,
			},
			ClientInfo: clickhouse.ClientInfo{
				Products: []struct {
					Name    string
					Version string
				}{
					{Name: "an-example-go-client", Version: "0.1"},
				},
			},

			Debugf: func(format string, v ...interface{}) {
				fmt.Printf(format, v)
			},
		})
	)

	c.Connection = conn
	c.Context = ctx

	if err != nil {
		panic(err)
	}

	if err := conn.Ping(ctx); err != nil {
		if exception, ok := err.(*clickhouse.Exception); ok {
			fmt.Printf("Exception [%d] %s \n%s\n", exception.Code, exception.Message, exception.StackTrace)
		}
		panic(err)
	}
}

func (c *ClickhouseClient) Insert(temperature float32, humidity float32, time time.Time) error {
	return c.Connection.AsyncInsert(
		c.Context,
		`INSERT INTO measurements VALUES (?, ?, ?)`,
		false,
		temperature,
		humidity,
		time)
}

func (c *ClickhouseClient) Close() {
	_ = c.Connection.Close()
	_ = c.Context.Done()
}
