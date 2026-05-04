package util

import (
	"log"
	"os"
)

var ApiConfig *Config

type Config struct {
	PGUser            string
	PGPass            string
	PGDB              string
	PGHost            string
	JWTSecret         string
	TestDB            string
	RabbitMQ          string
	PGPort            string
	RABBITMQ_PORT     string
	RABBITMQ_HOST     string
	RABBITMQ_USER     string
	RABBITMQ_PASSWORD string
}

func LoadConfig() *Config {
	return &Config{
		PGUser:            os.Getenv("POSTGRES_USER"),
		PGPass:            os.Getenv("POSTGRES_PASSWORD"),
		PGDB:              os.Getenv("POSTGRES_DB"),
		PGHost:            os.Getenv("POSTGRES_HOST"),
		PGPort:            os.Getenv("POSTGRES_PORT"),
		JWTSecret:         os.Getenv("JWT_SECRET"),
		TestDB:            os.Getenv("TEST_DB_URL"),
		RABBITMQ_PORT:     os.Getenv("RABBITMQ_PORT"),
		RABBITMQ_HOST:     os.Getenv("RABBITMQ_HOST"),
		RABBITMQ_USER:     os.Getenv("RABBITMQ_USER"),
		RABBITMQ_PASSWORD: os.Getenv("RABBITMQ_PASSWORD"),
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
	if c.PGPort == "" {
		log.Fatal("Environment POSTGRES_DB not set")
	}
	if c.PGHost == "" {
		log.Fatal("Environment POSTGRES_Host not set")
	}
	if c.JWTSecret == "" {
		log.Fatal("Environment JWT_SECRET not set")
	}
	if c.TestDB == "" {
		log.Fatal("Environment TEST_DB_URL not set")
	}
	if c.RABBITMQ_HOST == "" {
		log.Fatal("Environment RABBITMQ_HOST not set")
	}
	if c.RABBITMQ_PORT == "" {
		log.Fatal("Environment RABBITMQ_PORT not set")
	}
	if c.RABBITMQ_USER == "" {
		log.Fatal("Environment RABBITMQ_USER not set")
	}
	if c.RABBITMQ_PASSWORD == "" {
		log.Fatal("Environment RABBITMQ_PASSWORD not set")
	}
}
