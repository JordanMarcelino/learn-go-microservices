package provider

import (
	"github.com/IBM/sarama"
	"github.com/gin-gonic/gin"
	"github.com/jordanmarcelino/learn-go-microservices/order-service/internal/config"
	"github.com/jordanmarcelino/learn-go-microservices/order-service/internal/constant"
	"github.com/jordanmarcelino/learn-go-microservices/order-service/internal/controller"
	"github.com/jordanmarcelino/learn-go-microservices/order-service/internal/mq"
	"github.com/jordanmarcelino/learn-go-microservices/order-service/internal/repository"
	"github.com/jordanmarcelino/learn-go-microservices/order-service/internal/usecase"
	"github.com/jordanmarcelino/learn-go-microservices/pkg/database"
)

func BootstrapOrder(cfg *config.Config, router *gin.Engine) {
	kafkaAdmin.CreateTopic(constant.OrderCreatedTopic, constant.OrderCreatedTopicNumPartitions, constant.OrderCreatedTopicReplicationFactor)
	kafkaAdmin.CreateTopic(constant.OrderCancelledTopic, constant.OrderCancelledTopicNumPartitions, constant.OrderCancelledTopicReplicationFactor)

	orderCreatedProducer := mq.NewOrderCreatedProducer(database.NewKafkaAsyncProducer(&database.KafkaProducerOptions{
		Brokers:        cfg.Kafka.Brokers,
		RetryMax:       cfg.Kafka.RetryMax,
		FlushFrequency: cfg.Kafka.FlushFrequency,
		ReturnSuccess:  cfg.Kafka.ReturnSuccess,
		Topic:          constant.OrderCreatedTopic,
		Acks:           sarama.WaitForLocal,
	}))
	orderCancelledProducer := mq.NewOrderCreatedProducer(database.NewKafkaAsyncProducer(&database.KafkaProducerOptions{
		Brokers:        cfg.Kafka.Brokers,
		RetryMax:       cfg.Kafka.RetryMax,
		FlushFrequency: cfg.Kafka.FlushFrequency,
		ReturnSuccess:  cfg.Kafka.ReturnSuccess,
		Topic:          constant.OrderCancelledTopic,
		Acks:           sarama.WaitForLocal,
	}))
	paymentReminderProducer := mq.NewPaymentReminderProducer(rabbitmq)
	autoCancelProducer := mq.NewAutoCancelProducer(rabbitmq)
	orderSuccessProducer := mq.NewOrderSuccessProducer(rabbitmq)
	cancelNotificationProducer := mq.NewCancelNotificationProducer(rabbitmq)
	lockRepository := repository.NewLockRepository(rdl)

	orderUseCase := usecase.NewOrderUseCase(
		dataStore,
		lockRepository,
		orderCreatedProducer,
		orderCancelledProducer,
		cancelNotificationProducer,
		paymentReminderProducer,
		autoCancelProducer,
		orderSuccessProducer,
	)
	orderController := controller.NewOrderController(orderUseCase)

	orderController.Route(router)
}
