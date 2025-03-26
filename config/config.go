package config

import (
	"time"
)

type Storage struct {
	DBName   string `yaml:"database"`
	Host     string `yaml:"host" env-default:"db"`
	Port     int    `yaml:"port" env-default:"5432"`
	User     string `yaml:"username"`
	Password string `env:"DB_PASS"`
	SSLMode  string `yaml:"ssl_mode" env-default:"disable"`
}

type Application struct {
	Timeout      time.Duration `yaml:"timeout" env-default:"5s"`
	ForceTimeout time.Duration `yaml:"force_timeout" env-default:"10s"`
}

type Logger struct {
	Env string `yaml:"environment" env-default:"prod"`
}

type HTTPServer struct {
	Address string `yaml:"address" env-default:"0.0.0.0"`
	Port    int    `yaml:"port" env-default:"8080"`
}

type GRPCServer struct {
	Address string `yaml:"address" env-default:"0.0.0.0"`
	Port    int    `yaml:"port" env-default:"7070"`
}

type GRPCConn struct {
	Endpoint string `yaml:"endpoint"`
}

type RabbitConsumer struct {
	BindingKey   string `yaml:"binding_key"`
	QueueName    string `yaml:"queue_name"`
	ExchangeName string `yaml:"exchange_name"`
}

type RabbitProducer struct {
	ExchangeName string `yaml:"exchange_name"`
	RoutingKey   string `yaml:"routing_key"`
}

type RabbitConn struct {
	Host     string `yaml:"host" env-default:"localhost"`
	Port     int    `yaml:"port" env-default:"5672"`
	User     string `env:"RABBITMQ_USER" env-default:"guest"`
	Password string `env:"RABBITMQ_PASS" env-default:"guest"`
	VHost    string `yaml:"vhost" env-default:""`
}

type Redis struct {
	Address  string `yaml:"address" env-default:"localhost:6379"`
	Password string `env:"REDIS_PASSWORD" env-default:""`
	DB       int    `yaml:"db" env-default:"0"`
}
