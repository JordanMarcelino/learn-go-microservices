package mq

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/bytedance/sonic"
	"github.com/jordanmarcelino/learn-go-microservices/order-service/internal/constant"
	"github.com/jordanmarcelino/learn-go-microservices/order-service/internal/log"
	"github.com/jordanmarcelino/learn-go-microservices/pkg/mq"
	amqp "github.com/rabbitmq/amqp091-go"
)

type PaymentReminderProducer struct {
	Channel  *amqp.Channel
	exchange string
}

func NewPaymentReminderProducer(conn *amqp.Connection) mq.AMQPProducer {
	exchange := constant.PaymentReminderExchange
	ch, err := conn.Channel()
	if err != nil {
		log.Logger.Fatalf("failed to open a channel: %s", err)
	}

	args := amqp.Table{
		"x-delayed-type": "topic",
	}
	if err := ch.ExchangeDeclare(exchange, "x-delayed-message", true, false, false, false, args); err != nil {
		if amqpErr, ok := err.(*amqp.Error); ok && amqpErr.Code != amqp.PreconditionFailed {
			log.Logger.Fatalf("failed to declare an exchange: %s", err)
		}
	}

	return &PaymentReminderProducer{
		Channel:  ch,
		exchange: exchange,
	}
}

func (p *PaymentReminderProducer) Send(ctx context.Context, event mq.AMQPEvent) error {
	bytes, err := sonic.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	for _, delay := range constant.PaymentNotificationDelays {
		for i := 1; i <= constant.AMQPRetryLimit; i++ {
			err = p.Channel.PublishWithContext(
				ctx,
				p.exchange,
				event.Key(),
				false,
				false,
				amqp.Publishing{
					ContentType: "application/json",
					Headers: amqp.Table{
						"x-delay": delay.Milliseconds(),
					},
					Body: bytes,
				},
			)
			if err == nil {
				log.Logger.Infof("message published: %s", string(bytes))
				return nil
			}

			log.Logger.Errorf("failed to publish message: %s", err)
			delay := math.Pow(constant.AMQPRetryDelay, float64(i))
			time.Sleep(time.Duration(delay) * constant.AMQPRetryDelay)
		}
	}

	log.Logger.Errorf("message is aborted: %s", err)
	return err
}

func (p *PaymentReminderProducer) Exchange() string {
	return p.exchange
}
