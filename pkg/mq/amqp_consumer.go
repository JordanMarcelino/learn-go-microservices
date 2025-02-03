package mq

import "context"

type AMQPHandler func(ctx context.Context, body []byte) error

type AMQPConsumer interface {
	Consume(ctx context.Context, nWorker int) error
	Handler() AMQPHandler
	Queue() string
	Close() error
}
