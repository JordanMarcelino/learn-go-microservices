package mq

import (
	"context"
	"fmt"
	"math"
	"sync"
	"time"

	"github.com/bytedance/sonic"
	"github.com/jordanmarcelino/learn-go-microservices/mail-service/internal/constant"
	"github.com/jordanmarcelino/learn-go-microservices/mail-service/internal/dto"
	"github.com/jordanmarcelino/learn-go-microservices/mail-service/internal/log"
	"github.com/jordanmarcelino/learn-go-microservices/pkg/mq"
	"github.com/jordanmarcelino/learn-go-microservices/pkg/utils/smtputils"
	amqp "github.com/rabbitmq/amqp091-go"
)

type AccountVerifiedConsumer struct {
	Channel *amqp.Channel
	Mailer  smtputils.Mailer
	queue   string
	wg      *sync.WaitGroup
}

func NewAccountVerifiedConsumer(conn *amqp.Connection, mailer smtputils.Mailer) mq.AMQPConsumer {
	queue := constant.AccountVerifiedQueue
	exchange := constant.AccountVerifiedExchange
	key := constant.AccountVerifiedKey

	ch, err := conn.Channel()
	if err != nil {
		log.Logger.Fatalf("failed to open a channel: %s", err)
	}

	if _, err := ch.QueueDeclare(queue, true, false, false, false, nil); err != nil {
		log.Logger.Fatalf("failed to declare a queue: %s", err)
	}

	if err := ch.QueueBind(queue, key, exchange, false, nil); err != nil {
		log.Logger.Fatalf("failed to bind a queue: %s", err)
	}

	return &AccountVerifiedConsumer{
		Channel: ch,
		Mailer:  mailer,
		queue:   queue,
		wg:      &sync.WaitGroup{},
	}
}

func (c *AccountVerifiedConsumer) Consume(ctx context.Context, nWorker int) error {
	for i := 1; i <= nWorker; i++ {
		c.wg.Add(1)
		go c.Start(ctx, i)
	}
	return nil
}

func (c *AccountVerifiedConsumer) Start(ctx context.Context, workerID int) {
	defer c.wg.Done()

	msgs, err := c.Channel.ConsumeWithContext(ctx, c.Queue(), fmt.Sprintf("%v-%v", c.Queue(), workerID), false, false, false, false, nil)
	if err != nil {
		log.Logger.Errorf("failed to register %v-%v: %s", c.Queue(), workerID, err)
	}

	for {
		select {
		case <-ctx.Done():
			log.Logger.Infof("%v-%v shutting down...", c.Queue(), workerID)
			return
		case msg, ok := <-msgs:
			if !ok {
				log.Logger.Infof("%v-%v: message channel closed", c.Queue(), workerID)
				return
			}

			log.Logger.Infof("%v-%v: received a message %v", c.Queue(), workerID, string(msg.Body))
			for i := 1; i <= constant.AMQPRetryLimit; i++ {
				if err := c.Handler()(ctx, msg.Body); err != nil {
					log.Logger.Errorf("failed to consume message: %s", err)

					if i == constant.AMQPRetryLimit {
						log.Logger.Errorf("failed to consume message after %d retries: %s", constant.AMQPRetryLimit, err)
					} else {
						delay := math.Pow(constant.AMQPRetryDelay, float64(i))
						time.Sleep(time.Duration(delay) * constant.AMQPRetryDelay)
						log.Logger.Infof("retrying to consume message, attempt %d", i)
					}
				} else {
					_ = msg.Ack(false)
					break
				}
			}
		}
	}
}

func (c *AccountVerifiedConsumer) Handler() mq.AMQPHandler {
	return func(ctx context.Context, body []byte) error {
		var event dto.AccountVerifiedEvent
		if err := sonic.Unmarshal(body, &event); err != nil {
			log.Logger.Errorf("failed to unmarshal message: %s", err)
			return err
		}

		return c.Mailer.SendMail(ctx, event.Email, constant.AccountVerifiedSubject, constant.AccountVerifiedTemplate)
	}
}

func (c *AccountVerifiedConsumer) Queue() string {
	return c.queue
}

func (c *AccountVerifiedConsumer) Close() error {
	log.Logger.Infof("Closing consumer for queue: %s", c.Queue())
	c.wg.Wait()
	return c.Channel.Close()
}
