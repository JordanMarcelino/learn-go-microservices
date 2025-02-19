package mq

import (
	"context"
	"errors"
	"fmt"
	"math"
	"sync"
	"time"

	"github.com/bytedance/sonic"
	"github.com/jordanmarcelino/learn-go-microservices/order-service/internal/constant"
	"github.com/jordanmarcelino/learn-go-microservices/order-service/internal/dto"
	"github.com/jordanmarcelino/learn-go-microservices/order-service/internal/log"
	"github.com/jordanmarcelino/learn-go-microservices/order-service/internal/repository"
	"github.com/jordanmarcelino/learn-go-microservices/pkg/mq"
	amqp "github.com/rabbitmq/amqp091-go"
)

type AutoCancelConsumer struct {
	Channel                    *amqp.Channel
	DataStore                  repository.DataStore
	CancelNotificationProducer mq.AMQPProducer
	queue                      string
	wg                         *sync.WaitGroup
}

func NewAutoCancelConsumer(
	conn *amqp.Connection,
	dataStore repository.DataStore,
	cancelNotificationProducer mq.AMQPProducer,
) mq.AMQPConsumer {
	queue := constant.AutoCancelQueue
	exchange := constant.AutoCancelExchange
	key := constant.AutoCancelKey

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

	return &AutoCancelConsumer{
		Channel:                    ch,
		DataStore:                  dataStore,
		CancelNotificationProducer: cancelNotificationProducer,
		queue:                      queue,
		wg:                         &sync.WaitGroup{},
	}
}

func (c *AutoCancelConsumer) Consume(ctx context.Context, nWorker int) error {
	for i := 1; i <= nWorker; i++ {
		c.wg.Add(1)
		go c.Start(ctx, i)
	}
	return nil
}

func (c *AutoCancelConsumer) Start(ctx context.Context, workerID int) {
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

func (c *AutoCancelConsumer) Handler() mq.AMQPHandler {
	return func(ctx context.Context, body []byte) error {
		var event dto.AutoCancelEvent
		if err := sonic.Unmarshal(body, &event); err != nil {
			log.Logger.Errorf("failed to unmarshal message: %s", err)
			return err
		}

		return c.DataStore.Atomic(ctx, func(ds repository.DataStore) error {
			orderRepository := ds.OrderRepository()

			order, err := orderRepository.FindByID(ctx, event.OrderID)
			if err != nil {
				return err
			}
			if order == nil {
				return errors.New("order not found")
			}

			if order.Status == constant.ORDER_SUCCESS || order.Status == constant.ORDER_CANCELLED {
				return nil
			}

			order.Status = constant.ORDER_CANCELLED
			if err := orderRepository.UpdateStatus(ctx, order); err != nil {
				return err
			}

			return c.CancelNotificationProducer.Send(ctx, &dto.CancelNotificationEvent{OrderID: event.OrderID, UserID: event.UserID, Email: event.Email})
		})

	}
}

func (c *AutoCancelConsumer) Queue() string {
	return c.queue
}

func (c *AutoCancelConsumer) Close() error {
	log.Logger.Infof("Closing consumer for queue: %s", c.Queue())
	c.wg.Wait()
	return c.Channel.Close()
}
