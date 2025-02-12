package mq

import "context"

type KafkaHandler func(ctx context.Context, body []byte) error

type KafkaConsumer interface {
	Consume(ctx context.Context) error
	Handler() KafkaHandler
	Topic() string
	Close() error
}
