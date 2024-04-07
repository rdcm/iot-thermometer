package clickhouse

import (
	"context"
	"fmt"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
)

type Client struct {
	connection driver.Conn
	hostname   string
	port       int
	database   string
	username   string
	password   string
}

func New(
	hostname string,
	port int,
	database string,
	username string,
	password string,
) *Client {
	return &Client{
		hostname: hostname,
		port:     port,
		database: database,
		username: username,
		password: password,
	}
}

func (c *Client) Connect(ctx context.Context) error {
	opt := &clickhouse.Options{
		Addr: []string{fmt.Sprintf("%s:%d", c.hostname, c.port)},
		Auth: clickhouse.Auth{
			Database: c.database,
			Username: c.username,
			Password: c.password,
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
	}

	conn, err := clickhouse.Open(opt)
	if err != nil {
		return fmt.Errorf("clickhouse open: %w", err)
	}

	c.connection = conn

	if err := conn.Ping(ctx); err != nil {
		return fmt.Errorf("clickhouse ping: %w", err)
	}

	return nil
}

func (c *Client) Close() error {
	if err := c.connection.Close(); err != nil {
		return fmt.Errorf("clickhouse close: %w", err)
	}

	return nil
}

func (c *Client) Store(
	ctx context.Context,
	temperature float32,
	humidity float32,
	time time.Time,
) error {
	return c.connection.AsyncInsert(
		ctx,
		`INSERT INTO measurements VALUES (?, ?, ?)`,
		false,
		temperature,
		humidity,
		time,
	)
}
