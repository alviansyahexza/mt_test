package config

import (
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

func GetRabbitClient() (*amqp.Channel, error) {
	conn, err := amqp.Dial("amqp://user:password@localhost:5672/")
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %s", err)
	}
	fmt.Println("RabbitMQ connection established successfully")
	channel, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %s", err)
	}
	fmt.Println("RabbitMQ channel opened successfully")
	return channel, nil
}
