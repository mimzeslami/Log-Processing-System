package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config holds application configuration
type Config struct {
	DBConnString string
	RabbitMQURL  string
	QueueName    string
}

func LoadConfig() *Config {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dbConn := os.Getenv("DB_CONN_STRING")
	rabbitURL := os.Getenv("RABBITMQ_URL")
	queueName := os.Getenv("QUEUE_NAME")

	if dbConn == "" || rabbitURL == "" || queueName == "" {
		log.Fatal("Required environment variables are missing")
	}

	return &Config{
		DBConnString: dbConn,
		RabbitMQURL:  rabbitURL,
		QueueName:    queueName,
	}
}
