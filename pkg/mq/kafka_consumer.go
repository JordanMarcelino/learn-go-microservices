package mq

import "context"

type KafkaHandler func(ctx context.Context, event KafkaEvent) error

type KafkaConsumer interface {
	Consume(ctx context.Context) error
	Handler() KafkaHandler
	Topic() string
	Close() error
}
