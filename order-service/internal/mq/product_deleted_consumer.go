package mq

import (
	"context"
	"errors"
	"math"
	"sync"
	"time"

	"github.com/IBM/sarama"
	"github.com/bytedance/sonic"
	"github.com/jordanmarcelino/learn-go-microservices/order-service/internal/constant"
	"github.com/jordanmarcelino/learn-go-microservices/order-service/internal/dto"
	"github.com/jordanmarcelino/learn-go-microservices/order-service/internal/log"
	"github.com/jordanmarcelino/learn-go-microservices/order-service/internal/repository"
	"github.com/jordanmarcelino/learn-go-microservices/pkg/mq"
)

type ProductDeletedConsumer struct {
	Consumer  sarama.ConsumerGroup
	Datastore repository.DataStore
	topic     string
	wg        *sync.WaitGroup
}

func NewProductDeletedConsumer(consumer sarama.ConsumerGroup, dataStore repository.DataStore) mq.KafkaConsumer {
	return &ProductDeletedConsumer{
		Consumer:  consumer,
		Datastore: dataStore,
		topic:     constant.ProductDeletedTopic,
		wg:        &sync.WaitGroup{},
	}
}

func (c *ProductDeletedConsumer) Consume(ctx context.Context) error {
	c.wg.Add(1)
	return c.Consumer.Consume(ctx, []string{c.topic}, c)
}

func (c *ProductDeletedConsumer) Handler() mq.KafkaHandler {
	return func(ctx context.Context, body []byte) error {
		event := new(dto.ProductUpdatedEvent)
		if err := sonic.Unmarshal(body, event); err != nil {
			return err
		}

		return c.Datastore.Atomic(ctx, func(ds repository.DataStore) error {
			productRepository := ds.ProductRepository()

			product, err := productRepository.FindByID(ctx, event.Id)
			if err != nil {
				return err
			}
			if product == nil {
				return errors.New("product not found")
			}

			return productRepository.DeleteByID(ctx, product.ID)
		})
	}
}

func (c *ProductDeletedConsumer) Topic() string {
	return c.topic
}

func (c *ProductDeletedConsumer) Close() error {
	log.Logger.Infof("Closing consumer for topic: %s", c.Topic())
	c.wg.Wait()
	return c.Consumer.Close()
}

func (c *ProductDeletedConsumer) Setup(_ sarama.ConsumerGroupSession) error { return nil }

func (c *ProductDeletedConsumer) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }

func (c *ProductDeletedConsumer) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	defer c.wg.Done()

	for message := range claim.Messages() {
		log.Logger.Infof("topic %v received a message %v", c.Topic(), string(message.Value))

		for i := 1; i <= constant.KafkaConsumerRetryLimit+1; i++ {
			if err := c.Handler()(sess.Context(), message.Value); err != nil {
				log.Logger.Errorf("failed to consume message: %s", err)

				if i == constant.KafkaConsumerRetryLimit {
					log.Logger.Errorf("failed to consume message after %d retries: %s", constant.KafkaConsumerRetryLimit, err)
				} else {
					delay := math.Pow(constant.KafkaConsumerRetryDelay, float64(i))
					time.Sleep(time.Duration(delay) * time.Second)
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
