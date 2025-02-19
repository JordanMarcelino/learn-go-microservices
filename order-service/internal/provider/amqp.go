package provider

import (
	"github.com/jordanmarcelino/learn-go-microservices/order-service/internal/config"
	"github.com/jordanmarcelino/learn-go-microservices/order-service/internal/mq"
	pmq "github.com/jordanmarcelino/learn-go-microservices/pkg/mq"
)

func BootstrapAMQP(cfg *config.Config) []pmq.AMQPConsumer {
	cancelNotificationProducer := mq.NewCancelNotificationProducer(rabbitmq)

	return []pmq.AMQPConsumer{
		mq.NewAutoCancelConsumer(rabbitmq, dataStore, cancelNotificationProducer),
	}
}
