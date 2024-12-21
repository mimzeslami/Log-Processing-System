package main

import (
	"fmt"
	"interview/config"
	"interview/internal/db"
	"interview/internal/models"
	"interview/internal/mq"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"
)

func main() {
	// Database and RabbitMQ setup
	config := config.LoadConfig()
	dbConnString := config.DBConnString
	rabbitMQURL := config.RabbitMQURL
	queueName := config.QueueName

	database := db.NewDB(dbConnString)
	rabbitMQ := mq.NewRabbitMQ(rabbitMQURL, queueName)
	defer rabbitMQ.Close()

	// Start log processor
	var wg sync.WaitGroup
	const numWorkers = 3
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go processLogs(database, rabbitMQ, &wg)
	}

	wg.Wait()
	log.Println("All logs processed successfully")
}

func processLogs(database *db.DB, rabbitMQ *mq.RabbitMQ, wg *sync.WaitGroup) {
	defer wg.Done()

	messages, err := rabbitMQ.ConsumeMessages()
	if err != nil {
		log.Fatalf("Failed to consume messages: %v", err)
	}

	for msg := range messages {
		logMessage := string(msg.Body)
		log.Printf("Received log: %s\n", logMessage)

		parsedLog, err := parseLog(logMessage)
		if err != nil {
			log.Printf("Invalid log format: %v", err)
			continue
		}

		err = database.InsertStructuredLog(parsedLog)
		if err != nil {
			log.Printf("Failed to insert log into database: %v", err)
		} else {
			log.Printf("Successfully processed log: %+v", parsedLog)
		}
	}
}

func parseLog(logMessage string) (*models.StructuredLog, error) {
	parts := strings.Split(logMessage, ",")
	if len(parts) != 5 {
		return nil, fmt.Errorf("log message must have 5 fields")
	}

	statusCode, err := strconv.Atoi(parts[0])
	if err != nil {
		return nil, fmt.Errorf("invalid status code: %v", err)
	}

	timestamp, err := time.Parse(time.RFC3339, parts[3])
	if err != nil {
		return nil, fmt.Errorf("invalid timestamp: %v", err)
	}

	return &models.StructuredLog{
		StatusCode: statusCode,
		API:        parts[1],
		Message:    parts[2],
		Timestamp:  timestamp,
		IPAddress:  parts[4],
	}, nil
}
