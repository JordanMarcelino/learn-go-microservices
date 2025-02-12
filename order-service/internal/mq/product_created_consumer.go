package mq

import (
	"context"
	"math"
	"sync"
	"time"

	"github.com/IBM/sarama"
	"github.com/bytedance/sonic"
	"github.com/jordanmarcelino/learn-go-microservices/order-service/internal/constant"
	"github.com/jordanmarcelino/learn-go-microservices/order-service/internal/dto"
	"github.com/jordanmarcelino/learn-go-microservices/order-service/internal/entity"
	"github.com/jordanmarcelino/learn-go-microservices/order-service/internal/log"
	"github.com/jordanmarcelino/learn-go-microservices/order-service/internal/repository"
	"github.com/jordanmarcelino/learn-go-microservices/pkg/mq"
)

type ProductCreatedConsumer struct {
	Consumer  sarama.ConsumerGroup
	Datastore repository.DataStore
	topic     string
	wg        *sync.WaitGroup
}

func NewProductCreatedConsumer(consumer sarama.ConsumerGroup, dataStore repository.DataStore) mq.KafkaConsumer {
	return &ProductCreatedConsumer{
		Consumer:  consumer,
		Datastore: dataStore,
		topic:     constant.ProductCreatedTopic,
		wg:        &sync.WaitGroup{},
	}
}

func (c *ProductCreatedConsumer) Consume(ctx context.Context) error {
	c.wg.Add(1)
	return c.Consumer.Consume(ctx, []string{c.topic}, c)
}

func (c *ProductCreatedConsumer) Handler() mq.KafkaHandler {
	return func(ctx context.Context, body []byte) error {
		event := new(dto.ProductCreatedEvent)
		if err := sonic.Unmarshal(body, event); err != nil {
			return err
		}

		return c.Datastore.Atomic(ctx, func(ds repository.DataStore) error {
			productRepository := ds.ProductRepository()

			product := &entity.Product{Name: event.Name, Description: event.Description, Price: event.Price, Quantity: event.Quantity}
			return productRepository.Save(ctx, product)
		})
	}
}

func (c *ProductCreatedConsumer) Topic() string {
	return c.topic
}

func (c *ProductCreatedConsumer) Close() error {
	log.Logger.Infof("Closing consumer for topic: %s", c.Topic())
	c.wg.Wait()
	return c.Consumer.Close()
}

func (c *ProductCreatedConsumer) Setup(_ sarama.ConsumerGroupSession) error { return nil }

func (c *ProductCreatedConsumer) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }

func (c *ProductCreatedConsumer) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	defer c.wg.Done()

	for message := range claim.Messages() {
		log.Logger.Infof("topic %v received a message %v", c.Topic(), string(message.Value))

		for i := 1; i <= constant.KafkaRetryLimit; i++ {
			if err := c.Handler()(sess.Context(), message.Value); err != nil {
				log.Logger.Errorf("failed to consume message: %s", err)

				if i == constant.KafkaRetryLimit {
					log.Logger.Errorf("failed to consume message after %d retries: %s", constant.KafkaRetryLimit, err)
				} else {
					delay := math.Pow(constant.KafkaRetryDelay, float64(i))
					time.Sleep(time.Duration(delay))
					log.Logger.Infof("retrying to consume message, attempt %d", i)
				}
			} else {
				sess.MarkMessage(message, "")
				break
			}
		}

	}

	return nil
}
