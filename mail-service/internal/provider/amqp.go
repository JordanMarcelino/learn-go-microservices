package provider

import (
	"github.com/jordanmarcelino/learn-go-microservices/mail-service/internal/mq"
	pmq "github.com/jordanmarcelino/learn-go-microservices/pkg/mq"
)

func BootstrapAMQP() []pmq.AMQPConsumer {
	return []pmq.AMQPConsumer{
		mq.NewSendVerificationConsumer(rabbitmq, mailer),
	}
}
