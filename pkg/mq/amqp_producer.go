package mq

import "context"

type AMQPProducer interface {
	Send(ctx context.Context, event AMQPEvent) error
	Exchange() string
}

type AMQPEvent interface {
	Key() string
}
