package main

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Rabbit     RabbitMq   `yaml:"rabbitmq"`
	Clickhouse Clickhouse `yaml:"clickhouse"`
}

type RabbitMq struct {
	Port     string `yaml:"port"     env:"MQTT_PORT"`
	Host     string `yaml:"host"     env:"MQTT_HOST"`
	User     string `yaml:"user"     env:"MQTT_USER"`
	Password string `yaml:"password" env:"MQTT_PASSWORD"`
	Queue    string `yaml:"queue"    env:"MQTT_QUEUE"`
}

type Clickhouse struct {
	Hostname string `yaml:"hostname" env:"CH_HOSTNAME"`
	Port     int    `yaml:"port"     env:"CH_PORT"`
	Password string `yaml:"password" env:"CH_PASSWORD"`
	Database string `yaml:"database" env:"CH_DATABASE"`
	Username string `yaml:"username" env:"CH_USERNAME"`
}

func (c *RabbitMq) getConnStr() string {
	return fmt.Sprintf("amqp://%s:%s@%s:%s/", c.User, c.Password, c.Host, c.Port)
}

func readCfg() Config {
	var cfg Config

	err := cleanenv.ReadConfig("config.yml", &cfg)
	if err != nil {
		panic(err)
	}

	return cfg
}
