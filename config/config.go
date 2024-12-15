package config

import (
	"os"
	"strconv"
)

type Config struct {
	PostgresUser string
	PostgresPass string
	PostgresDB   string
	PostgresHost string
	PostgresPort int

	RabbitUser  string
	RabbitPass  string
	RabbitHost  string
	RabbitPort  int
	RabbitQueue string
}

func LoadConfig() (*Config, error) {
	pPort, _ := strconv.Atoi(os.Getenv("POSTGRES_PORT"))
	rPort, _ := strconv.Atoi(os.Getenv("RABBITMQ_PORT"))

	return &Config{
		PostgresUser: os.Getenv("POSTGRES_USER"),
		PostgresPass: os.Getenv("POSTGRES_PASSWORD"),
		PostgresDB:   os.Getenv("POSTGRES_DB"),
		PostgresHost: os.Getenv("POSTGRES_HOST"),
		PostgresPort: pPort,

		RabbitUser:  os.Getenv("RABBITMQ_USER"),
		RabbitPass:  os.Getenv("RABBITMQ_PASSWORD"),
		RabbitHost:  os.Getenv("RABBITMQ_HOST"),
		RabbitPort:  rPort,
		RabbitQueue: os.Getenv("RABBITMQ_QUEUE"),
	}, nil
}
