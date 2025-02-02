package provider

import (
	"github.com/jordanmarcelino/learn-go-microservices/mail-service/internal/config"
	"github.com/jordanmarcelino/learn-go-microservices/pkg/database"
	"github.com/jordanmarcelino/learn-go-microservices/pkg/utils/smtputils"
	amqp "github.com/rabbitmq/amqp091-go"
)

var (
	rabbitmq *amqp.Connection
	mailer   smtputils.Mailer
)

func BootstrapGlobal(cfg *config.Config) {
	rabbitmq = database.NewAMQP((*database.AmqpOptions)(cfg.Amqp))
	mailer = smtputils.NewMailer()
}
