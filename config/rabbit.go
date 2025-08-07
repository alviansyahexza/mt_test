package config

import (
	"fmt"
	"log"
	"os"

	"github.com/streadway/amqp"
)

func GetRabbitClient() (*amqp.Channel, error) {
	conn, err := amqp.Dial(os.Getenv("RABBITMQ_URL"))
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
