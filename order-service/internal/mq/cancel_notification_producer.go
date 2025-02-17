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

type CancelNotificationProducer struct {
	Channel  *amqp.Channel
	exchange string
}

func NewCancelNotificationProducer(conn *amqp.Connection) mq.AMQPProducer {
	exchange := constant.CancelNotificationExchange
	ch, err := conn.Channel()
	if err != nil {
		log.Logger.Fatalf("failed to open a channel: %s", err)
	}

	if err := ch.ExchangeDeclare(exchange, "topic", true, false, false, false, nil); err != nil {
		log.Logger.Fatalf("failed to declare an exchange: %s", err)
	}

	return &CancelNotificationProducer{
		Channel:  ch,
		exchange: exchange,
	}
}

func (p *CancelNotificationProducer) Send(ctx context.Context, event mq.AMQPEvent) error {
	bytes, err := sonic.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	for i := 1; i <= constant.AMQPRetryLimit; i++ {
		err = p.Channel.PublishWithContext(
			ctx,
			p.exchange,
			event.Key(),
			false,
			false,
			amqp.Publishing{
				ContentType: "application/json",
				Body:        bytes,
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

	log.Logger.Errorf("message is aborted: %s", err)
	return err
}

func (p *CancelNotificationProducer) Exchange() string {
	return p.exchange
}
