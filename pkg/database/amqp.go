package database

import (
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type AmqpOptions struct {
	Username string
	Password string
	Host     string
	VHost    string
	Port     int
}

func NewAMQP(opt *AmqpOptions) *amqp.Connection {
	url := fmt.Sprintf("amqp://%s:%s@%s:%d/%s",
		opt.Username,
		opt.Password,
		opt.Host,
		opt.Port,
		opt.VHost,
	)

	conn, err := amqp.Dial(url)
	if err != nil {
		log.Fatalf("failed to connect to rabbitmq: %v", err)
	}

	return conn
}
