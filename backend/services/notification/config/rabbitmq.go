package config

import (
	"log"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitResources struct {
	Connection *amqp.Connection
	Channel    *amqp.Channel
	QueueName  string
}

func ConnectRabbitMQ() (*RabbitResources, error) {
	url := os.Getenv("RABBITMQ_URL")
	if url == "" {
		url = "amqp://guest:guest@rabbitmq:5672/"
	}

	queue := os.Getenv("RABBITMQ_QUEUE_NAME")
	if queue == "" {
		queue = "notifications"
	}

	var conn *amqp.Connection
	var err error

	maxAttempts := 10
	for attempt := 1; attempt <= maxAttempts; attempt++ {
		conn, err = amqp.Dial(url)
		if err == nil {
			break
		}
		log.Printf("RabbitMQ connection failed (attempt %d/%d): %v", attempt, maxAttempts, err)
		time.Sleep(time.Duration(attempt) * time.Second)
	}
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, err
	}

	if _, err := ch.QueueDeclare(
		queue,
		true,  // durable
		false, // auto-delete
		false, // exclusive
		false, // no-wait
		nil,
	); err != nil {
		ch.Close()
		conn.Close()
		return nil, err
	}

	log.Println("Connected to RabbitMQ")
	return &RabbitResources{
		Connection: conn,
		Channel:    ch,
		QueueName:  queue,
	}, nil
}
