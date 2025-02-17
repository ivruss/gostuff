package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"time"
)

type Config struct {
	HTTPServer      HTTPServer      `yaml:"http_server"`
	Logger          Logger          `yaml:"logger"`
	Application     Application     `yaml:"application"`
	MessagingClient MessagingClient `yaml:"messaging_client"`
}

type MessagingClient struct {
	Endpoint string `yaml:"endpoint"`
}

type Application struct {
	Timeout      time.Duration `yaml:"timeout" env-default:"5s"`
	ForceTimeout time.Duration `yaml:"force_timeout" env-default:"60s"`
}

type Logger struct {
	Level string `yaml:"level" env-default:"debug"`
	Env   string `yaml:"environment" env-default:"prod"`
}

type HTTPServer struct {
	Address string `yaml:"address" env-default:"0.0.0.0"`
	Port    int    `yaml:"port" env-default:"8080"`
}

func MustLoad(confpath string) *Config {
	var cfg Config
	err := cleanenv.ReadConfig(confpath, &cfg)
	if err != nil {
		log.Fatalf("error reading config file: %s", err)
	}

	return &cfg
}
