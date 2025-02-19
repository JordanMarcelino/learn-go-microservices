package provider

import (
	"github.com/jordanmarcelino/learn-go-microservices/mail-service/internal/config"
	"github.com/jordanmarcelino/learn-go-microservices/mail-service/internal/feign"
	"github.com/jordanmarcelino/learn-go-microservices/mail-service/internal/mq"
	pmq "github.com/jordanmarcelino/learn-go-microservices/pkg/mq"
)

func BootstrapAMQP(cfg *config.Config) []pmq.AMQPConsumer {
	orderClient := feign.NewOrderClient(cfg.Feign.OrderURL)

	return []pmq.AMQPConsumer{
		mq.NewSendVerificationConsumer(rabbitmq, mailer),
		mq.NewAccountVerifiedConsumer(rabbitmq, mailer),
		mq.NewPaymentReminderConsumer(rabbitmq, mailer, orderClient),
		mq.NewOrderCancelConsumer(rabbitmq, mailer, orderClient),
		mq.NewOrderSuccessConsumer(rabbitmq, mailer, orderClient),
	}
}
