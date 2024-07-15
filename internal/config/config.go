package config

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	ExporterPort           string `envconfig:"EXPORTER_PORT" default:"9100"`
	LiteServerAddr         string `envconfig:"LITE_SERVER_ADDR"`
	LiteServerKey          string `envconfig:"LITE_SERVER_KEY"`
	GlobalLiteClientConfig string `envconfig:"GLOBAL_CONFIG_URL" default:"https://ton.org/global.config.json"`
}

func LoadConfig() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found")
	}

	var config Config
	err = envconfig.Process("", &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
