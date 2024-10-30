package config

import (
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"log"
)

type Config struct {
	Port              string `envconfig:"PORT" default:":4000"`
	DatabaseName      string `envconfig:"DATABASE_NAME" default:"goth.db"`
	SessionCookieName string `envconfig:"SESSION_COOKIE_NAME" default:"session"`
	StripeSecretKey   string `envconfig:"STRIPE_SECRET_KEY" default:"StripeSecretKey"`
}

func loadConfig() (*Config, error) {
	// Load env variables from .env file

	err := godotenv.Load()
	if err != nil {
		log.Printf("Error loading .env file")
	}

	var cfg Config
	err = envconfig.Process("", &cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}

func MustLoadConfig() *Config {
	cfg, err := loadConfig()
	if err != nil {
		panic(err)
	}
	return cfg
}
