package services

import (
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

const (
	reconnectDelay = 5 * time.Second
	reInitDelay    = 2 * time.Second
	resendDelay    = 5 * time.Second
)

func NewClient(queueName, addr string) *Client {
	client := &Client{
		m: &sync.Mutex{},
		infoLog: slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		),
		errorLog: slog.New(
			slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError}),
		),
		done:      make(chan bool),
		queueName: queueName,
	}

	go client.handleReconnect(addr)

	return client
}

func (client *Client) handleReconnect(addr string) {
	for {
		client.m.Lock()
		client.isReady = false
		client.m.Unlock()

		client.infoLog.Info("attempting to connect")

		conn, err := client.connect(addr)
		if err != nil {
			client.errlog.Println("failed to connect. Retrying...")

			select {
			case <-client.done:
				return
			case <-time.After(reconnectDelay):
			}
			continue
		}

		if done := client.handleReInit(conn); done {
			break
		}
	}
}

func (client *Client) connect(addr string) (*amqp.Connection, error) {
	conn, err := amqp.Dial(addr)
	if err != nil {
		return nil, err
	}

	client.changeConnection(conn)
	client.infolog.Println("connected")
	return conn, nil
}

func (client *Client) changeConnection(connection *amqp.Connection) {
	client.connection = connection
	client.notifyConnClose = make(chan *amqp.Error, 1)
	client.connection.NotifyClose(client.notifyConnClose)
}
