package provider

import (
	"github.com/IBM/sarama"
	"github.com/gin-gonic/gin"
	"github.com/jordanmarcelino/learn-go-microservices/pkg/database"
	"github.com/jordanmarcelino/learn-go-microservices/product-service/internal/config"
	"github.com/jordanmarcelino/learn-go-microservices/product-service/internal/constant"
	"github.com/jordanmarcelino/learn-go-microservices/product-service/internal/controller"
	"github.com/jordanmarcelino/learn-go-microservices/product-service/internal/mq"
	"github.com/jordanmarcelino/learn-go-microservices/product-service/internal/usecase"
)

func BootstrapProduct(cfg *config.Config, router *gin.Engine) {
	kafkaAdmin.CreateTopic(constant.ProductCreatedTopic, constant.ProductCreatedTopicNumPartitions, constant.ProductCreatedTopicReplicationFactor)
	kafkaAdmin.CreateTopic(constant.ProductDeletedTopic, constant.ProductDeletedTopicNumPartitions, constant.ProductDeletedTopicReplicationFactor)
	kafkaAdmin.CreateTopic(constant.ProductUpdatedTopic, constant.ProductUpdatedTopicNumPartitions, constant.ProductUpdatedTopicReplicationFactor)

	productCreatedProducer := mq.NewProductCreatedProducer(database.NewKafkaAsyncProducer(&database.KafkaProducerOptions{
		Brokers:        cfg.Kafka.Brokers,
		RetryMax:       cfg.Kafka.RetryMax,
		FlushFrequency: cfg.Kafka.FlushFrequency,
		ReturnSuccess:  cfg.Kafka.ReturnSuccess,
		Topic:          constant.ProductCreatedTopic,
		Acks:           sarama.WaitForLocal,
	}))
	productUpdatedProducer := mq.NewProductUpdatedProducer(database.NewKafkaSyncProducer(&database.KafkaProducerOptions{
		Brokers:        cfg.Kafka.Brokers,
		RetryMax:       cfg.Kafka.RetryMax,
		FlushFrequency: cfg.Kafka.FlushFrequency,
		ReturnSuccess:  cfg.Kafka.ReturnSuccess,
		Topic:          constant.ProductUpdatedTopic,
		Acks:           sarama.WaitForLocal,
	}))
	productDeletedProducer := mq.NewProductDeletedProducer(database.NewKafkaAsyncProducer(&database.KafkaProducerOptions{
		Brokers:        cfg.Kafka.Brokers,
		RetryMax:       cfg.Kafka.RetryMax,
		FlushFrequency: cfg.Kafka.FlushFrequency,
		ReturnSuccess:  cfg.Kafka.ReturnSuccess,
		Topic:          constant.ProductDeletedTopic,
		Acks:           sarama.WaitForLocal,
	}))
	productUseCase := usecase.NewProductUseCase(dataStore, productCreatedProducer, productUpdatedProducer, productDeletedProducer)
	productController := controller.NewProductController(productUseCase)

	productController.Route(router)
}
