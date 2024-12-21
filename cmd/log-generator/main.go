package main

import (
	"fmt"
	"interview/config"
	"log"
	"math/rand"
	"time"

	"github.com/streadway/amqp"
)

// Simulated log data
var statusCodes = []int{200, 403, 404, 500}
var endpoints = []string{"/api/login", "/api/upload", "/api/profile", "/api/logout"}
var messages = []string{
	"User logged in",
	"File uploaded",
	"Access denied",
	"Server error occurred",
}
var ips = []string{"192.168.1.10", "192.168.1.11", "192.168.1.12", "192.168.1.13"}

func main() {
	config := config.LoadConfig()
	var rabbitMQURL = config.RabbitMQURL
	var queueName = config.QueueName
	// Connect to RabbitMQ
	conn, err := amqp.Dial(rabbitMQURL)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}
	defer ch.Close()

	_, err = ch.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to declare queue: %v", err)
	}

	log.Println("Log generator started. Sending structured logs to RabbitMQ...")
	for {
		logMessage := generateLog()
		err := ch.Publish(
			"",        // exchange
			queueName, // routing key
			false,     // mandatory
			false,     // immediate
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        []byte(logMessage),
			},
		)
		if err != nil {
			log.Printf("Failed to publish log: %v", err)
		} else {
			log.Printf("Published log: %s", logMessage)
		}

		time.Sleep(time.Duration(rand.Intn(5)+1) * time.Second) // Random delay
	}
}

func generateLog() string {
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("%d,%s,%s,%s,%s",
		statusCodes[rand.Intn(len(statusCodes))],
		endpoints[rand.Intn(len(endpoints))],
		messages[rand.Intn(len(messages))],
		time.Now().UTC().Format(time.RFC3339),
		ips[rand.Intn(len(ips))],
	)
}
