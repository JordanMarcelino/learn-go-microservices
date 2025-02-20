package provider

import (
	"github.com/IBM/sarama"
	"github.com/jordanmarcelino/learn-go-microservices/pkg/database"
	pmq "github.com/jordanmarcelino/learn-go-microservices/pkg/mq"
	"github.com/jordanmarcelino/learn-go-microservices/product-service/internal/config"
	"github.com/jordanmarcelino/learn-go-microservices/product-service/internal/constant"
	"github.com/jordanmarcelino/learn-go-microservices/product-service/internal/mq"
)

func BootstrapKafka(cfg *config.KafkaConfig) []pmq.KafkaConsumer {
	return []pmq.KafkaConsumer{
		mq.NewOrderCreatedConsumer(
			database.NewKafkaConsumerGroup(&database.KafkaConsumerOptions{
				Brokers:       cfg.Brokers,
				ConsumerGroup: constant.OrderCreatedConsumerGroup,
				InitialOffset: sarama.OffsetOldest,
			}),
			dataStore,
		),
		mq.NewOrderCancelledConsumer(
			database.NewKafkaConsumerGroup(&database.KafkaConsumerOptions{
				Brokers:       cfg.Brokers,
				ConsumerGroup: constant.OrderCancelledConsumerGroup,
				InitialOffset: sarama.OffsetOldest,
			}),
			dataStore,
		),
	}
}
