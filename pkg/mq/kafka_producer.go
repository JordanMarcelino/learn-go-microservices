package mq

import "context"

type KafkaProducer interface {
	Send(ctx context.Context, event KafkaEvent) error
	Topic() string
}

type KafkaEvent interface {
	ID() string
}
