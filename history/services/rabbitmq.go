package services

import (
	"log/slog"
	"os"

	amqp "github.com/rabbitmq/amqp091-go"
)

func Connect() (*amqp.Channel, error) {
	rabbitmqURL := os.Getenv("RABBITMQ_URL")
	conn, err := amqp.Dial(string(rabbitmqURL))
	if err != nil {
		slog.Error("Failed to connect to RabbitMQ", "error", err)
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		slog.Error("Failed to open a channel", "error", err)
		return nil, err
	}

	return ch, nil
}
