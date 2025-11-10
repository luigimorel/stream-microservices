package services

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"sync"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Client struct {
	m          *sync.Mutex
	queueName  string
	infoLog    *slog.Logger
	errorLog   *slog.Logger
	isReady    bool
	done       chan bool
	connection *amqp.Connection
	channel    *amqp.Channel
}

func NewProducer(queueName, addr string) (*Client, error) {
	conn, err := amqp.Dial(addr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed to open a channel: %w", err)
	}

	_, err = ch.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to declare queue: %w", err)
	}

	client := &Client{
		m:          &sync.Mutex{},
		queueName:  queueName,
		infoLog:    slog.Default(),
		errorLog:   slog.Default(),
		isReady:    true,
		done:       make(chan bool),
		connection: conn,
		channel:    ch,
	}

	client.infoLog.Info("RabbitMQ producer ready", "queue", queueName)
	return client, nil
}

// Publish sends a message to the queue
func (c *Client) Publish(body string) error {
	c.m.Lock()
	defer c.m.Unlock()

	if !c.isReady {
		return fmt.Errorf("connection not ready")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := c.channel.PublishWithContext(ctx,
		"",
		c.queueName,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        []byte(body),
		})
	if err != nil {
		return fmt.Errorf("failed to publish message: %w", err)
	}

	c.infoLog.Info("Message sent", "body", body)
	return nil
}

// Close cleans up the connection and channel
func (c *Client) Close() {
	c.infoLog.Info("Closing RabbitMQ connection")
	c.channel.Close()
	c.connection.Close()
	close(c.done)
}

func Connect() (*amqp.Channel, error) {
	conn, err := amqp.Dial(os.Getenv("RABBITMQ_URL"))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to rabbitmq: %w", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed to open the Channel: %w", err)
	}
	defer ch.Close()

	return ch, nil
}
