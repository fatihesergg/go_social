package util

import (
	"log"
	"os"
)

var ApiConfig *Config

type Config struct {
	PGUser    string
	PGPass    string
	PGDB      string
	JWTSecret string
	TestDB    string
	RabbitMQ  string
}

func LoadConfig() *Config {
	return &Config{
		PGUser:    os.Getenv("POSTGRES_USER"),
		PGPass:    os.Getenv("POSTGRES_PASSWORD"),
		PGDB:      os.Getenv("POSTGRES_DB"),
		JWTSecret: os.Getenv("JWT_SECRET"),
		TestDB:    os.Getenv("TEST_DB_URL"),
		RabbitMQ:  os.Getenv("RABBITMQ_URL"),
	}
}

func (c *Config) Validate() {
	if c.PGUser == "" {
		log.Fatal("Environment POSTGRES_USER not set")
	}
	if c.PGPass == "" {
		log.Fatal("Environment POSTGRES_PASSWORD not set")
	}
	if c.PGDB == "" {
		log.Fatal("Environment POSTGRES_DB not set")
	}
	if c.JWTSecret == "" {
		log.Fatal("Environment JWT_SECRET not set")
	}
	if c.TestDB == "" {
		log.Fatal("Environment TEST_DB_URL not set")
	}
	if c.RabbitMQ == "" {
		log.Fatal("Environment RABBITMQ_URL not set")
	}
}
