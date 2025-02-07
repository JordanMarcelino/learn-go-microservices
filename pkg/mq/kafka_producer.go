package mq

import (
	"context"
)

type KafkaMetadata struct {
	Retry int `json:"retry"`
}

type KafkaProducer interface {
	Send(ctx context.Context, event KafkaEvent) error
	Topic() string
}

type KafkaEvent interface {
	ID() string
}
