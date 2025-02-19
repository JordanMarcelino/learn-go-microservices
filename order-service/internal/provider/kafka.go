package provider

import (
	"github.com/IBM/sarama"
	"github.com/jordanmarcelino/learn-go-microservices/order-service/internal/config"
	"github.com/jordanmarcelino/learn-go-microservices/order-service/internal/constant"
	"github.com/jordanmarcelino/learn-go-microservices/order-service/internal/mq"
	"github.com/jordanmarcelino/learn-go-microservices/pkg/database"
	pmq "github.com/jordanmarcelino/learn-go-microservices/pkg/mq"
)

func BootstrapKafka(cfg *config.KafkaConfig) []pmq.KafkaConsumer {
	return []pmq.KafkaConsumer{
		mq.NewProductCreatedConsumer(
			database.NewKafkaConsumerGroup(&database.KafkaConsumerOptions{
				Brokers:       cfg.Brokers,
				ConsumerGroup: constant.ProductCreatedConsumerGroup,
				InitialOffset: sarama.OffsetOldest,
			}),
			dataStore,
		),
		mq.NewProductUpdatedConsumer(
			database.NewKafkaConsumerGroup(&database.KafkaConsumerOptions{
				Brokers:       cfg.Brokers,
				ConsumerGroup: constant.ProductUpdatedConsumerGroup,
				InitialOffset: sarama.OffsetOldest,
			}),
			dataStore,
		),
		mq.NewProductDeletedConsumer(
			database.NewKafkaConsumerGroup(&database.KafkaConsumerOptions{
				Brokers:       cfg.Brokers,
				ConsumerGroup: constant.ProductDeletedConsumerGroup,
				InitialOffset: sarama.OffsetOldest,
			}),
			dataStore,
		),
	}
}
