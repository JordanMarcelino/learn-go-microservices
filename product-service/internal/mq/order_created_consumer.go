package mq

import (
	"context"
	"errors"
	"math"
	"sync"
	"time"

	"github.com/IBM/sarama"
	"github.com/bytedance/sonic"
	"github.com/jordanmarcelino/learn-go-microservices/pkg/mq"
	"github.com/jordanmarcelino/learn-go-microservices/product-service/internal/constant"
	"github.com/jordanmarcelino/learn-go-microservices/product-service/internal/dto"
	"github.com/jordanmarcelino/learn-go-microservices/product-service/internal/log"
	"github.com/jordanmarcelino/learn-go-microservices/product-service/internal/repository"
)

type OrderCreatedConsumer struct {
	Consumer  sarama.ConsumerGroup
	Datastore repository.DataStore
	topic     string
	wg        *sync.WaitGroup
}

func NewOrderCreatedConsumer(consumer sarama.ConsumerGroup, dataStore repository.DataStore) mq.KafkaConsumer {
	return &OrderCreatedConsumer{
		Consumer:  consumer,
		Datastore: dataStore,
		topic:     constant.OrderCreatedTopic,
		wg:        &sync.WaitGroup{},
	}
}

func (c *OrderCreatedConsumer) Consume(ctx context.Context) error {
	c.wg.Add(1)
	return c.Consumer.Consume(ctx, []string{c.topic}, c)
}

func (c *OrderCreatedConsumer) Handler() mq.KafkaHandler {
	return func(ctx context.Context, body []byte) error {
		event := new(dto.OrderCreatedEvent)
		if err := sonic.Unmarshal(body, event); err != nil {
			return err
		}

		return c.Datastore.Atomic(ctx, func(ds repository.DataStore) error {
			productRepository := ds.ProductRepository()

			productIds := []int64{}
			for _, item := range event.Items {
				productIds = append(productIds, item.ProductID)
			}

			products, err := productRepository.FindAllByIDForUpdate(ctx, productIds)
			if err != nil {
				return err
			}
			if len(products) != len(productIds) {
				return errors.New("product not found")
			}

			for i, product := range products {
				product.Quantity -= event.Items[i].Quantity
			}

			return productRepository.UpdateAllQuantity(ctx, products)
		})
	}
}

func (c *OrderCreatedConsumer) Topic() string {
	return c.topic
}

func (c *OrderCreatedConsumer) Close() error {
	log.Logger.Infof("Closing consumer for topic: %s", c.Topic())
	c.wg.Wait()
	return c.Consumer.Close()
}

func (c *OrderCreatedConsumer) Setup(_ sarama.ConsumerGroupSession) error { return nil }

func (c *OrderCreatedConsumer) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }

func (c *OrderCreatedConsumer) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	defer c.wg.Done()

	for message := range claim.Messages() {
		log.Logger.Infof("topic %v received a message %v", c.Topic(), string(message.Value))

		for i := 1; i < constant.KafkaConsumerRetryLimit+1; i++ {
			if err := c.Handler()(sess.Context(), message.Value); err != nil {
				log.Logger.Errorf("failed to consume message: %s", err)

				if i > constant.KafkaConsumerRetryLimit {
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
